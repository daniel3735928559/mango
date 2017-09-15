import zmq
import random
import sys
import time
import json

port = "12000"

context = zmq.Context()
poller = zmq.Poller()
socket = context.socket(zmq.ROUTER)
poller.register(socket,zmq.POLLIN)
socket.bind("tcp://*:%s" % port)
ctr = 0
clients = {}
while True:
    ctr += 1
    print("looping")
    socks = dict(poller.poll(1000))
    #self.local_gateway.socket.send()
    for s in socks:
        #self.debug_print(s,socks[s])
        if socks[s] == zmq.POLLIN:
            client = s.recv()
            data = s.recv()
            if client in clients:
                print("old buddy", client)
            else: 
                print("new friend", client)
                clients[client] = True
        elif(socks[s] & zmq.POLLERR != 0):
            print("rip",s)
    if ctr % 3 == 0:
        first_data_element = random.randrange(2,20)
        second_data_element = random.randrange(0,360)
        message = json.dumps({'First Data':first_data_element, 'Second Data':second_data_element})
        print(message)
        for c in clients:
            socket.send_multipart([c,bytes(message,'ascii')])
    time.sleep(0.5)
