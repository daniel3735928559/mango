import time
import zmq
import threading

def mc_heartbeat_worker(node,url,timeout,stopper,context):
    socket = context.socket(zmq.DEALER)
    socket.setsockopt(zmq.IDENTITY, bytes(node,'utf-8'))
    socket.connect(url)
    while True:
        time.sleep(timeout)
        if stopper.is_set():
            break
        socket.send(bytes(node,'utf-8'))
