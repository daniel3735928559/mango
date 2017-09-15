import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, json, traceback, yaml, subprocess, shlex
from libmango import *
from mu_dataflows import *
from mu_transport import *

class mu(m_node):    
    def __init__(self):
        super().__init__(debug=True)
        self.mu_transport = mu_transport(self.context, self.poller, "0.0.0.0")
        
        self.mu_dataflow = mu_dataflow(self.mu_transport, self.serialiser, self.ui_to_world, self.handle_error)
        self.poller.register(self.mu_transport.socket,zmq.POLLIN)
        self.dataflows[self.mu_transport.socket] = self.mu_dataflow
        self.interface.default_handler = self.world_to_ui
        subprocess.Popen(["python", "server.py"], cwd=os.path.dirname(os.path.realpath(__file__)), env={"MU_HTTP_PORT":os.getenv("MU_HTTP_PORT"),"MANGO_SIDECHANNEL_PORT":str(self.mu_transport.port),"MU_ROOT_DIR":os.getcwd()})
        self.run()

    def ui_to_world(self, header, args):
        print("UI2W",header,args)
        self.m_send(header['name'], args)

    def world_to_ui(self, header, args):
        print("W2UI", header, args)
        self.mu_dataflow.send(header,args)
        
mu()
