import zmq, socket

class m_transport:
    def __init__(self,target):
        self.target = target
    
    def tx(self,payload):
        pass

    def rx(self):
        pass
    
    def die(self):
        pass


class m_ZMQ_transport(m_transport):
    def __init__(self,owner,target,context,poller,server=False):
        super().__init__(target)
        self.owner = owner
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
        if nid is None: nid = bytes(self.owner.node_id,"UTF-8")
        else: nid = bytes(nid,"UTF-8")
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


class m_srv_sock:
    def __init__(self,port):
        self.socket = socket.socket()
        self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        print("Binding ",port)
        self.socket.bind(('',port))
        self.socket.listen(1)
        self.socket.setblocking(0)

    def die(self):
        self.socket.close()

    def tx(self,data):
        pass

    def rx(self):
        c,a = self.socket.accept()
        return c,a

class m_client_sock:
    def __init__(self,sock):
        self.socket = sock

    def die(self):
        self.socket.close()

    def rx(self):
        data = self.socket.recv(4096)

        if(len(data) == 0):
            return None

        return data

    def tx(self,data):
        txd = self.socket.send(data)
        sent = txd
        while(sent < len(data) and txd > 0):
            txd = self.socket.send(data)
            sent += txd

        return sent
