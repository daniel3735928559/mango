import zmq, socket

class m_ZMQ_transport():
    def __init__(self,target,context,poller,server=False):
        self.target = target
        if server:
            self.socket = context.socket(zmq.ROUTER)
            print("binding to" + str(target))
            self.socket.bind(target)
        else:
            self.socket = context.socket(zmq.DEALER)
            print("connecting to" + str(target))
            #self.set_id()
            self.socket.connect(target)
        
    def set_id(self,nid=None):
        nid = bytes(nid,"UTF-8")
        print('Setting id: '+str(nid))
        self.socket.disconnect(self.target)
        self.socket.setsockopt(zmq.IDENTITY,nid)
        self.socket.connect(self.target)

    def die(self):
        self.socket.disconnect(self.target)

    def tx(self,payload):
        # print(payload)
        self.socket.send(payload)
        return "sent"

    def rx(self):
        return self.socket.recv()
