import io, re, time, signal, os
from serialiser import *
from dataflow import *
from transport import *
from interface import *
from error import *
import traceback
from inspect import getframeinfo, stack

class m_node: 
    def __init__(self,debug=False):
        self.version = "0.1"
        self.debug = debug
        self.context = zmq.Context()
        self.poller = zmq.Poller()
        self.serialiser = m_serialiser(self.version)
        self.interface = m_if(default_handler=self.default)
        self.dataflows = {}
        self.route = os.getenv('MANGO_ROUTE')
        self.node_id = os.getenv('MANGO_ID')
        self.group_id = os.getenv('MANGO_GROUP','root')
        self.interface.add_interface(os.path.join(os.getenv('PYTHONPATH'),'../node.yaml'),{
            'heartbeat':self.heartbeat,
            'exit':self.end
        })
        self.flags = {}
        self.server = os.getenv('MC_ADDR',None)
        if not self.server is None:
            self.local_gateway = m_ZMQ_transport(self.server,self.context,self.poller,self.route)
            s = self.local_gateway.socket
            self.dataflow = m_dataflow(self.interface,self.local_gateway,self.serialiser,self.dispatch,self.handle_error)
            self.dataflows[s] = self.dataflow
            self.poller.register(s,zmq.POLLIN)

    def add_socket(self, sock, recv_cb, err_cb):
        dataflow = lambda: None
        dataflow.recv = recv_cb
        dataflow.die = err_cb
        self.poller.register(sock, zmq.POLLIN)
        self.dataflows[sock] = dataflow
            
    def ready(self):
        iface = self.interface.get_spec()
        self.debug_print("IF",iface)
            
    def dispatch(self,header,args):
        self.debug_print("DISPATCH",header,args)
        try:
           result = self.interface.get_function(header['name'])(header,args)
           if not result is None:
               self.m_send(result[0],result[1],mid=header.get('mid',None))
        except Exception as exc:
           self.handle_error("unknown",traceback.format_exc())

    def heartbeat(self,header,args):
        self.mc_send('alive',{})
        
    def end(self,header,args):
        exit(0)
            
    def default(self,header,args):
        self.debug_print("UNHANDLED",header,args)

    def handle_error(self,src,err):
        self.debug_print('OOPS',src,err)
        self.mc_send('error',{'source':src,'message':err})

    def make_header(self,name,msg_type=None,mid=None):
        header = {'name':name}
        if mid: header['mid'] = mid
        if msg_type: header['type'] = msg_type
        self.debug_print("H",header)
        return header

    def mc_send(self,name,msg):
        self.dataflow.send(self.make_header(name,'system'),msg)
        
    def m_send(self,name,msg,mid=None):
        self.debug_print('sending',msg)
        self.dataflow.send(self.make_header(name,mid=mid),msg)

    def debug_print(self,*args):
        if self.debug:
            caller = getframeinfo(stack()[1][0])
            print("[{}:{} {} DEBUG] ".format(caller.filename, caller.lineno, self.node_id),*args)
    
    def run(self,f=None):
        if not self.server is None:
            self.ready()
        while True:
            socks = dict(self.poller.poll(1000*1000))
            for s in socks:
                if socks[s] == zmq.POLLIN:
                    self.debug_print("RX",s,self.dataflows[s])
                    self.dataflows[s].recv()
                elif(socks[s] & zmq.POLLERR != 0):
                    self.debug_print('socket error',s,socks[s])
                    self.dataflows[s].die()
            if not f is None: 
                f()
