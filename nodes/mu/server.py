from flask import Flask, send_from_directory, Response, render_template
import os, json, gevent
import zmq.green as zmq
from flask_sockets import Sockets
from gevent import monkey

monkey.patch_all()

app = Flask(__name__)

sockets = Sockets(app)
context = zmq.Context()

m_port = int(os.getenv("MANGO_SIDECHANNEL_PORT"))
root_dir = os.getenv("MU_ROOT_DIR")
lib_dir = os.path.join(os.path.dirname(os.path.realpath(__file__)),"mangojs/")
http_port = int(os.getenv("MU_HTTP_PORT"))

@app.route('/mangojs/<path:path>')
def do_mangojs(path):
    return send_from_directory(lib_dir, path)

@app.route('/<path:path>')
def others(path):
    return send_from_directory(root_dir, path)

@sockets.route('/mangotx')
def get_data(ws):
    socket = context.socket(zmq.DEALER)
    socket.connect('tcp://localhost:{PORT}'.format(PORT=m_port))
    socket.send_string("tx")
    print('helloing')
    while not ws.closed:
        data = ws.receive()
        print("D",data)
        if data:
            socket.send_string(data)

@sockets.route('/mangorx')
def send_data(ws):
    while not ws.closed:
        socket = context.socket(zmq.DEALER)
        print("connecting")
        socket.connect('tcp://localhost:{PORT}'.format(PORT=m_port))
        socket.send_string("rx")
        gevent.sleep()
        while not ws.closed:

            print("rxing")
            data = socket.recv()
            ws.send(data)
            gevent.sleep()
        print("DEADED")

if __name__ == '__main__':
    from gevent import pywsgi
    from geventwebsocket.handler import WebSocketHandler
    server = pywsgi.WSGIServer(('', http_port), app, handler_class=WebSocketHandler)
    server.serve_forever()
