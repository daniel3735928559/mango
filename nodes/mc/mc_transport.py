import zmq, socket

class mc_ZMQ_transport():
    def __init__(self,socket,route):
        self.socket = socket
        self.route = route

    def die(self):
        pass
        #self.socket.disconnect(self.target)

    def tx(self,payload):
        # print(payload)
        self.socket.send(self.route,zmq.SNDMORE)
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
