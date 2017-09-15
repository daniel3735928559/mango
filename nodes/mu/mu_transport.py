# Many thanks to https://gist.github.com/rich20bb/4190781

import zmq

class mu_transport:
    def __init__(self, context, poller, bind='*', port=12000):
        self.socket = context.socket(zmq.ROUTER)
        poller.register(self.socket,zmq.POLLIN)
        if port is None:
            self.port = self.socket.bind_to_random_port("tcp://{}".format(bind))
        else:
            self.socket.bind("tcp://{}:{}".format(bind,port))
            self.port = port
        self.clients = set()

    def tx(self, data):
        print("TX",data,self.clients)
        for c in self.clients:
            self.socket.send_multipart([c,data])

    def rx(self):
        client = self.socket.recv()
        data = self.socket.recv()
        print("MUT RX",client,data)
        if not client in self.clients and data == b"rx":
            self.clients.add(client)
        return data
