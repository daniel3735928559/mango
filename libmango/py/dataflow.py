from error import *

class m_dataflow:
    def __init__(self,interface,transport,serialiser,dispatch_cb,error_cb):
        self.version = "0.1"
        self.transport = transport
        self.interface = interface
        self.serialiser = serialiser
        self.dispatch_cb = dispatch_cb
        self.error_cb = error_cb

    def send(self,header,msg):
        try:
            data = self.serialiser.serialise(header, msg)
            self.transport.tx(data)
        except m_error as exc:
            self.error_cb(exc)
            return None
    
    def recv(self):
        data = self.transport.rx()
        try:
            header,args = self.serialiser.deserialise(data)
            self.interface.validate(header['command'],args)
            result = self.dispatch_cb(header,args)
        except m_error as exc:
            self.error_cb(header['src_node'],str(exc))
            return None

