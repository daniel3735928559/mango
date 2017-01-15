import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from route_parser import route_parser
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node
from lxml import etree

class NodeType: 
    def __init__(self,name,wd,runner):
        self.name = name
        self.wd = wd
        self.runner = runner

class Node: 
    def __init__(self,node_id,key,dataflow,route,iface,master=None,ports=[],flags={},local=True):
        # self.dataflow is the socket (or whatever) that you can use
        # to talk to this node.  It will usually be set by mc to
        # self.connections[0], and to send on it you can just use
        # route=bytearray(self.node_id,'utf-8')
        self.dataflow = dataflow
        self.key = key
        self.node_id = node_id
        self.flags = flags
        self.interface = iface
        self.ports = {"stdio":Port("stdio",self)}
        self.alive_time = time.time()
        print("PPP",self.ports)
        self.local = local
        self.route = route
        if (not master is None) and local:
            self.ports["mc"] = Port("mc",self)
            self.ports["mc"].add_route(Route(self.ports["mc"],master))
        else:
            self.ports["stderr"] = Port("stderr",self)
        for x in ports:
            if x != "mc" and x != "stdio":
                self.ports[x] = Port(x,self)
            
    def send(self, header, args, route=None):
        if route is None:
            route = self.route
        try:
            print("HH",header,self.interface.interface,type(self.interface.interface))
            if self.flags.get("strict",True): args = self.interface.validate(header['name'],args)
            print("Sending",self.node_id,route,header,args)
            self.dataflow.send(header,args,route)
        except Exception as exc:
            print('OOPS',exc)
            exc_type, exc_value, exc_traceback = sys.exc_info()
            traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
            traceback.print_exception(exc_type, exc_value, exc_traceback,file=sys.stdout)

    # def __repr__(self):
    #     return self.node_id

class Route:
    def __init__(self,source,endpoint):
        self.source = source
        self.endpoint = endpoint
        self.transmogrifiers = [] # a list of "modifier type",whatever pairs:

    def to_string(self):
        if len(self.transmogrifiers) > 0:
            return "{} > {} > {}".format(str(self.source), str(self.transmogrifiers), str(self.endpoint))
        return "{} > {}".format(str(self.source), str(self.endpoint))
        
    def apply(self,message,header,args):
        new_header = header
        new_args = args
        new_message = None
        print("T",self.transmogrifiers)
        for t,o in self.transmogrifiers:
            if(t == "add"):
                for k in o:
                    #if(k in header):
                    #    new_header[k] = o[k] if o[k][0] != '$' else (args[o[k][1:]] if o[k][1] != '_' else message.decode('ASCII'))
                    #else:
                    #    new_args[k] = o[k] if o[k][0] != '$' else (args[o[k][1:]] if o[k][1] != '_' else message.decode('ASCII'))
                    if o[k] == "$_":
                        new_args[k] = message.decode('ASCII')
                    elif o[k][0] != '$':
                        new_args[k] = o[k]
                    elif o[k][1:] in new_args:
                        new_args[k] = new_args[o[k][1:]]
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
            elif(t == "comm"):
                new_header['name'] = o
            elif(t == "sub"):
                n_args = {}
                for k in o:
                    if o[k] == "$_":
                        new_args[k] = message.decode('ASCII')
                    elif o[k][0] != '$':
                        n_args[k] = o[k]
                    elif o[k][1:] in new_args:
                        n_args[k] = new_args[o[k][1:]]
                new_args = n_args
            elif(t == "filter"):
                if header['name'] != o:
                    print("FILTER BLOCK",header['name'],'is not',o)
                    return None,None,None
                print("FILTER PASS",header['name'],o)

            # Syntax for a bash transmogrifier is: 
            # ...> sh $key [cmd] > ...
            # Which is equivalent to dict[$key] = $(echo $key | cmd)
            elif(t == "sh"):
                print(new_args)
                k,cmd = o.split(' ',1)
                if(not k in new_args):
                    print("Bad key",k)
                    return None,None,None
                new_args[k] = subprocess.check_output('printf \'{}\' | {}'.format(new_args[k],cmd),shell=True).decode()

        #new_header['port'] = self.endpoint.name
        return new_header,new_args,new_message

    def send(self,message,header,args):
        h,a,m = self.apply(message,header,args)
        if h is None:
            print("BLOCKED",message,header,args)
            return
        if m is None:
            #self.endpoint.owner.conn.send(a,self.source.name,dest=self.endpoint.get_id(),route=bytearray(self.endpoint.owner.node_id,'utf-8'))
            print("ROUTE send",h,a,m,self.endpoint.get_id())
            src=h['src_node']
            #if not self.endpoint.owner.local: src = self.endpoint.owner.master.srv_addr
            h['port'] = self.endpoint.name
            print("R SEND",self.endpoint.owner,self.endpoint.owner.dataflow,h,a)

            # Special case so that mc can reply directly
            
            if str(self.endpoint.owner.node_id) == "mc":
                self.endpoint.owner.send(h,a,self.source.owner.route)
            else:
                self.endpoint.owner.send(h,a)

        else:
            print("R SEND RAW",self.endpoint.owner,self.endpoint.owner.dataflow)
            self.endpoint.owner.dataflow.send_raw(m,bytearray(self.endpoint.owner.node_id,'utf-8'))

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

    def del_route_to(self,p):
        if(p in self.routes): 
            del self.routes[p]
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
