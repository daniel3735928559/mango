from error import *

class m_dataflow:
    def __init__(self,interface,transport,serialiser,dispatch_cb,reply_cb,error_cb):
        self.version = "0.1"
        self.transport = transport
        self.interface = interface
        self.serialiser = serialiser
        self.dispatch_cb = dispatch_cb
        self.reply_cb = reply_cb
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
            if header['command'] == 'call':
                args = self.interface.validate(self.interface[header['function']]['args'],args)
                result = self.dispatch_cb(header,args)
                result = self.interface.validate(self.interface[header['function']]['returns'],result)
                self.send(result,header['port'])
            elif header['command'] == 'reply':
                self.reply_cb(header,args,data,self)
        except m_error as exc:
            self.error_cb(exc)
            return None

