import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, subprocess, shlex, json
from route_parser import route_parser
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from mangolib import m_node
from lxml import etree

class NodeType: 
    def __init__(self,name,wd,runner):
        self.name = name
        self.wd = wd
        self.runner = runner

class Node: 
    def __init__(self,node_id,key,dataflow,master=None,local=True):
        # This is the socket (or whatever) that you can use to talk to
        # this node.  It will usually be set by mc to
        # self.connections[0], and to send on it you can just use
        # route=bytearray(self.node_id,'utf-8')
        self.dataflow = dataflow
        self.key = key
        self.node_id = node_id
        self.flags = 0
        self.ports = {"stdio":Port("stdio",self)}
        self.local = local
        if (not master is None) and local:
            self.ports["mc"] = Port("mc",self)
            self.ports["mc"].add_route(Route(self.ports["mc"],master))
    def __repr__(self):
        return self.node_id

class Route:
    def __init__(self,source,endpoint):
        self.source = source
        self.endpoint = endpoint
        self.transmogrifiers = [] # a list of "modifier type",whatever pairs:

    def to_string(self):
        return ""
        
    def apply(self,message,header,args):
        new_header = header
        new_args = args
        new_message = None
        print("T",self.transmogrifiers)
        for t,o in self.transmogrifiers:
            if(t == "add"):
                for k in o:
                    if(k in header):
                        new_header[k] = o[k] if o[k][0] != '$' else (args[o[k][1:]] if o[k][1] != '_' else message.decode('ASCII'))
                    else:
                        new_args[k] = o[k] if o[k][0] != '$' else (args[o[k][1:]] if o[k][1] != '_' else message.decode('ASCII'))
            elif(t == "addn"):
                for k in o:
                    if not k in args:
                        new_args[k] = o[k]
            elif(t == "raw"):
                new_message = new_args[o]
                break
            elif(t == "del"):
                for k in o:
                    if k in new_args:
                        del new_args[k]
            elif(t == "sub"):
                new_args = {}
                for k in o:
                    new_args[k] = o[k] if o[k][0] != '$' else (args[o[k][1:]] if o[k][1] != '_' else message.decode('ASCII'))
            elif(t == "filter"):
                return None

            # Syntax for a bash transmogrifier is: 
            # ...> sh $key [cmd] > ...
            # Which is equivalent to dict[$key] = $(echo $key | cmd)
            elif(t == "sh"):
                k,cmd = o.split(' ',1)
                if(not k in new_args):
                    print("Bad key",k)
                    return None
                args[k] = subprocess.check_output('printf {} | {}'.format(args[k],cmd),shell=True).decode()

        #new_header['port'] = self.endpoint.name
        return new_header,new_args,new_message

    def send(self,message,header,args):
        h,a,m = self.apply(message,header,args)
        if m is None:
            #self.endpoint.owner.conn.send(a,self.source.name,dest=self.endpoint.get_id(),route=bytearray(self.endpoint.owner.node_id,'utf-8'))
            print("ROUTE send",h,a,m,self.endpoint.get_id())
            src=h['source']
            #if not self.endpoint.owner.local: src = self.endpoint.owner.master.srv_addr
            print("R SEND",self.endpoint.owner,self.endpoint.owner.dataflow)
            self.endpoint.owner.dataflow.send(a,header['source'],dport=self.endpoint.get_id(),mid=header['mid'])

        else:
            print("R SEND RAW",self.endpoint.owner,self.endpoint.owner.dataflow)
            self.endpoint.owner.dataflow.send_raw(m,route=bytearray(self.endpoint.owner.node_id,'utf-8'))

    def __repr__(self):
        return str(self.source) + " > " + "".join([str(t)+" > " for t in self.transmogrifiers]) + str(self.endpoint)

    # def reply(self,port,header,reply,raw,route=None):
    #     print("MC REPLY",port,header,reply,raw,route)
    #     self.source.owner.conn.send(a,port=self.source.name,dest=self.endpoint.get_id(),route=bytearray(self.source.owner.node_id,'utf-8'),source_node=header['source'],mid=header['mid'])


class Port:
    def __init__(self,name,owner):
        self.name = name
        self.owner = owner
        self.routes = {} # dictionary of endpoint : [Route object]

    def get_id(self):
        return self.owner.node_id + "/" + self.name

    def add_route(self,r):
        self.routes[r.endpoint] = r
        print("ROUTE ADD",str(r))
        return True

    def del_route(self,r):
        if(r.endpoint in self.routes): 
            del self.routes[r.endpoint]
            return True
        return False

    # pass the message through the transmogrifiers for each route and
    # return a list of target nodes/ports and corresponding
    # transmogrified messages
    def send(self,message,header,args):
        print("Port sending",str(self))
        for r in self.routes:
            print("SENDING ON",str(r))
            self.routes[r].send(message,header,args)
            
    def __repr__(self):
        return self.owner.node_id + "/" + self.name
            
class Remote:
    def __init__(self,host,port):
        self.host = host
        self.port = port
