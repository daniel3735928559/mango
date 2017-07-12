import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from transform import *
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node

class Node: 
    def __init__(self,node_id,key,dataflow,route,iface,flags={},local=True):
        # self.dataflow is the socket (or whatever) that you can use
        # to talk to this node.  It will usually be set by mc to
        # self.connections[0], and to send on it you can just use
        # route=bytearray(self.node_id,'utf-8')
        self.dataflow = dataflow
        self.key = key
        self.node_id = node_id
        self.flags = flags
        self.interface = iface
        self.last_heartbeat_time = time.time()
        self.last_alive_time = time.time()
        self.local = local
        self.route = route
        self.hb_thread = None
        self.hb_stopper = None
        self.routes = {}
                
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
    def emit(self,message,header,args):
        print("Node emitting",str(self))
        for r in self.routes:
            print("SENDING ON",str(r))
            self.routes[r].send(message,header,args)

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


            
class Remote:
    def __init__(self,host,port):
        self.host = host
        self.port = port
