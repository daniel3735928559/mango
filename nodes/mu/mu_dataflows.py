import io, re, socketserver, socket, zmq, subprocess, traceback
from dataflow import *
from libmango import *
from serialiser import *
from transport import *
from mu_transport import *

class mu_dataflow():
    def __init__(self,transport,serialiser,dispatch_cb,error_cb):
        self.transport = transport
        self.serialiser = serialiser
        self.dispatch_cb = dispatch_cb
        self.error_cb = error_cb
    
    def recv(self):
        data = self.transport.rx()
        if data == b"tx" or data == b"rx":
            print("just a hello")
            return
        try:
            header,args = self.serialiser.deserialise(data)
            print("HA",header,args)
            result = self.dispatch_cb(header,args)
        except m_error as exc:
            traceback.print_exc()
            self.error_cb(header['src_node'],str(exc))
            return None

    def send(self,header,msg):
        data = self.serialiser.serialise(header, msg)
        self.transport.tx(data)
