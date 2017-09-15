import zmq.green as zmq
import json
import gevent
from flask_sockets import Sockets
from flask import Flask, render_template
import logging
from gevent import monkey

monkey.patch_all()

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

sockets = Sockets(app)
context = zmq.Context()

ZMQ_LISTENING_PORT = 12000

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/2')
def index2():
    return render_template('index2.html')

@sockets.route('/test')
def get_data(ws):
    while not ws.closed:
        print("hello")
        print(ws.receive())

@sockets.route('/zeromq')
def send_data(ws):
    while not ws.closed:
        logger.info('Got a websocket connection, sending up data from zmq')
        socket = context.socket(zmq.DEALER)
        print("connecting")
        socket.connect('tcp://localhost:{PORT}'.format(PORT=ZMQ_LISTENING_PORT))
        socket.send_string("greets")
        gevent.sleep()
        while not ws.closed:
            print("rxing")
            data = socket.recv_json()
            logger.info(data)
            ws.send(json.dumps(data))
            gevent.sleep()
        print("DEADED")

if __name__ == '__main__':
    from gevent import pywsgi
    from geventwebsocket.handler import WebSocketHandler
    server = pywsgi.WSGIServer(('', 25000), app, handler_class=WebSocketHandler)
    server.serve_forever()
