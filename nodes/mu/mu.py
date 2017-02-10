import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, json, traceback, yaml, subprocess, shlex
from libmango import *
from mu_dataflows import *
from mu_transport import *

class mu(m_node):    
    def __init__(self):
        super().__init__(debug=True)
        self.server_sock = mu_server_ws("0.0.0.0",int(os.getenv("MU_WS_PORT")))
        self.server_dataflow = mu_server_dataflow(None,self.server_sock,self.serialiser,self.ui_to_world,self.handle_error,self)
        self.poller.register(self.server_sock.socket.fileno(),zmq.POLLIN)
        self.dataflows[self.server_sock.socket.fileno()] = self.server_dataflow
        self.client_ws_dataflows = {}
        mu_if_file = os.getenv("MU_IF",None)
        if not mu_if_file is None:
            try:
                print("IF_FILE", os.getcwd(), mu_if_file)
                with open(mu_if_file,"r") as f:
                    new_if = yaml.load(f)
                if_dict = {x:self.world_to_ui for x in new_if['inputs']}
                print("INTERFACE",if_dict,new_if)
                self.interface.add_interface(mu_if_file,if_dict)
            except:
                exc_type, exc_value, exc_traceback = sys.exc_info()
                traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
                traceback.print_exception(exc_type, exc_value, exc_traceback,file=sys.stdout)

        subprocess.Popen(["python", "server.py"], cwd=os.path.dirname(os.path.realpath(__file__)), env={"MU_HTTP_PORT":os.getenv("MU_HTTP_PORT"),"MU_WS_PORT":os.getenv("MU_WS_PORT"),"MU_ROOT_DIR":os.getenv("MU_ROOT_DIR")})
        self.ready()
        self.run()#self.main_loop)

    def ui_to_world(self, header, args):
        print("UI2W",header,args)
        self.m_send(header['name'], args)

    def world_to_ui(self, header, args):
        print("W2UI", header, args)
        for x in self.client_ws_dataflows:
            df = self.client_ws_dataflows[x]
            if df.transport.alive:
                df.send(header,args)
        #print(args)
        #for x in args:
        #    if(not isinstance(args[x],str)): args[x] = args[x].decode("UTF-8")
        #self.put_queue.put((args,header['source']))
mu()
