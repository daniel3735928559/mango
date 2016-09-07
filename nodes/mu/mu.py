import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, json, traceback, yaml
from libmango import *
from mu_dataflows import *
from mu_transport import *

class mu(m_node):    
    def __init__(self):
        super().__init__(debug=True)

        #self.doc_root = sys.argv[1]
        # self.reply_table = {}
        # self.put_queue = Queue()
        # self.xhrpoller = None
        print("ASDA")
        self.server_sock = mu_server_ws("0.0.0.0",int(os.getenv("MU_PORT",8045)))
        self.server_dataflow = mu_server_dataflow(None,self.server_sock,self.serialiser,self.ui_to_world,self.handle_error,self)
        self.poller.register(self.server_sock.socket.fileno(),zmq.POLLIN)
        self.dataflows[self.server_sock.socket.fileno()] = self.server_dataflow
        #self.default_cmd = self.world_to_ui
        mu_if_file = os.getenv("MU_IF",None)
        if not mu_if_file is None:
            try:
                if_dict = {x:self.world_to_ui for x in yaml.load(mu_if_file)}
                self.interface.add_interface(mu_if_file,if_dict)
            except:
                exc_type, exc_value, exc_traceback = sys.exc_info()
                traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
                traceback.print_exception(exc_type, exc_value, exc_traceback,file=sys.stdout)
        self.ready()
        self.run()#self.main_loop)

    def get_cb(self,path):
        fn = (self.doc_root+"/" + ("index.html" if path == "/" else path)) if path[1:4] != 'lib' else "./" + path
        print(fn)
        try:
            my_file = open(fn,'rb')
            file_data = my_file.read()
            my_file.close()
        except:
            exc_type, exc_value, exc_traceback = sys.exc_info()
            traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
            traceback.print_exception(exc_type, exc_value, exc_traceback,file=sys.stdout)
            return
        return file_data

    def poll_cb(self,dataflow):
        if not self.xhrpoller is None and not self.xhrpoller == dataflow: self.xhrpoller.die()
        self.xhrpoller = dataflow;
        self.main_loop()

    # def main_loop(self):
    #     if(self.xhrpoller == None or self.put_queue.empty()): return
    #     try:
    #         if(not self.put_queue.empty()):
    #             msg,src = self.put_queue.get()
    #             self.xhrpoller.send(msg,src)
    #             self.xhrpoller.die()
    #             self.xhrpoller = None
    #     except:
    #         print("poll response exception")
    #         et,ev,etb = sys.exc_info()
    #         traceback.print_tb(etb, limit=1, file=sys.stdout)
    #         traceback.print_exception(et,ev,etb,file=sys.stdout)
    #         self.xhrpoller.die()
    #         self.xhrpoller = None

    def ui_act(self,header,args):
        # Generate ID based on requested name ID and a key, make a
        # Node object based on this, add it to the list, and return
        # the ID and key
        
        return {"node_id":n.ID,"key":n.key}

    def ui_to_world(self,port,header,args,raw,dataflow):
        print("UI2W")
        print(header,args)
        self.m_send(args,reply_callback=lambda h,a: self.world_to_ui(h,a,args['mid']))

    def world_to_ui(self,header,args,msgid=None):
        print("W2UI",header,args)
        if not msgid is None:
            args['reply'] = bytes(str(msgid),"ASCII")
        #print(args)
        #for x in args:
        #    if(not isinstance(args[x],str)): args[x] = args[x].decode("UTF-8")
        #self.put_queue.put((args,header['source']))
mu()
