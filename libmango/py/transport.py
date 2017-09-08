import zmq, socket

class m_ZMQ_transport():
    def __init__(self,target,context,poller,nid=None,server=False):
        self.target = target
        if server:
            self.socket = context.socket(zmq.ROUTER)
            self.socket.bind(target)
        else:
            self.socket = context.socket(zmq.DEALER)
            if nid: self.socket.setsockopt_string(zmq.IDENTITY,nid)
            self.socket.connect(target)
        
    def set_id(self,nid=None):
        nid = bytes(nid,"UTF-8")
        self.socket.disconnect(self.target)
        self.socket.setsockopt(zmq.IDENTITY,nid)
        self.socket.connect(self.target)

    def die(self):
        self.socket.disconnect(self.target)

    def tx(self,payload):
        self.socket.send(payload)

    def rx(self):
        return self.socket.recv()
