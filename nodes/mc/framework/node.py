import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from . transform import *
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node

STARTING = 0
RUNNING = 1
REAP = 2

class Node: 
    def __init__(self, node_id, group, node_type, key, dataflow, route, iface):
        # self.dataflow is the socket (or whatever) that you can use
        # to talk to this node.  It will usually be set by mc to
        # self.connections[0], and to send on it you can just use
        # route=bytearray(self.node_id,'utf-8')
        self.node_type = node_type
        self.dataflow = dataflow
        self.group = group
        self.key = key
        self.node_id = node_id
        self.name = self.node_id
        self.interface = iface
        self.route = route
        self.status = STARTING
        self.hb_thread = None
        self.hb_stopper = None
        self.last_heartbeat_time = time.time()
        self.last_alive_time = self.last_heartbeat_time
        self.routes = {}

    def get_id(self):
        return str(self)
        
    def __repr__(self):
        return "{}/{}".format(self.group, self.node_id)

    # a message came in for this node.  Validate it and then pass it
    # on to the node through the dataflow
            
    def handle(self, header, args, route=None):
        if route is None:
            route = self.route
        try:
            args = self.interface.validate(header['name'],args)
            self.dataflow.send(header,args,route)
        except Exception as exc:
            print('OOPS',exc)
            exc_type, exc_value, exc_traceback = sys.exc_info()
            traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
            traceback.print_exception(exc_type, exc_value, exc_traceback,file=sys.stdout)
