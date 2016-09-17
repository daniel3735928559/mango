import time
import zmq

def mc_heartbeat_worker(node,url,timeout,context):
    socket = context.socket(zmq.DEALER)
    socket.setsockopt(zmq.IDENTITY, bytes(node,'utf-8'))
    socket.connect(url)
    while True:
        time.sleep(timeout)
        socket.send(bytes(node,'utf-8'))
