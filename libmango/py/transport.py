import zmq, socket

class m_ZMQ_transport():
    def __init__(self,target,context,poller,nid):
        self.target = target
        self.socket = context.socket(zmq.DEALER)
        self.socket.setsockopt_string(zmq.IDENTITY,nid)
        self.socket.connect(target)

    def die(self):
        self.socket.disconnect(self.target)

    def tx(self,payload):
        self.socket.send_string(payload)

    def rx(self):
        return self.socket.recv()
