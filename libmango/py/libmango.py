import io, re, time, signal, os
from serialiser import *
from dataflow import *
from transport import *
from interface import *
from error import *

class m_node: 
    def __init__(self,debug=False):
        self.version = "0.1"
        self.debug = debug
        self.context = zmq.Context()
        self.poller = zmq.Poller()
        self.serialiser = m_serialiser(self.version)
        self.interface = m_if(default_handler=self.reply)
        self.dataflows = {}
        self.node_id = os.getenv('MANGO_ID')
        self.interface.add_interface(os.path.join(os.getenv('PYTHONPATH'),'../node.yaml'),{
            'reg':self.reg,
            'reply':self.reply,
            'heartbeat':self.heartbeat
        })
        self.ports = []
        self.flags = {}
        self.server = os.getenv('MC_ADDR',None)
        print(self.server,os.environ)
        if not self.server is None:
            self.local_gateway = m_ZMQ_transport(self.server,self.context,self.poller)
            s = self.local_gateway.socket
            self.dataflow = m_dataflow(self.interface,self.local_gateway,self.serialiser,self.dispatch,self.handle_error)
            self.dataflows[s] = self.dataflow
            self.poller.register(s,zmq.POLLIN)

    def ready(self):
        iface = self.interface.get_spec()
        self.debug_print("IF",iface)
        self.m_send('hello',{'id':self.node_id,'if':iface,'ports':self.ports,'flags':self.flags},port="mc")
            
    def dispatch(self,header,args):
        self.debug_print("DISPATCH",header,args)
        try:
           result = self.interface.get_function(header['name'])(header,args)
           if not result is None:
               self.m_send("reply",result,port=header['port'])
        except Exception as exc:
           self.handle_error(header['src_node'],str(exc))

    def heartbeat(self,header,args):
        self.m_send('alive',{},port="mc")
            
    def reply(self,header,args):
        print("REPLY",header,args)

    def handle_error(self,src,err):
        self.debug_print('OOPS',src,err)
        self.m_send('error',{'source':src,'message':err},port="mc")

    def reg(self,header,args):
        if header['src_node'] != 'mc':
            print('only accepts reg from mc',header,args)
            return
        self.node_id = args["id"]
        self.debug_print('my new node id')
        self.debug_print(self.node_id)
        self.debug_print("registered as " + self.node_id)

    def make_header(self,name,src_port='stdio'):
        header = {'src_port':src_port,
                  'name':name}
        self.debug_print("H",header)
        return header
    
    def m_send(self,name,msg,port='stdio',async=True):
        self.debug_print('sending',msg)
        header = self.make_header(name,port)
        self.dataflow.send(header,msg)
        if not async and not callback is None:
            self.dataflow.recv()

    def debug_print(self,*args):
        if self.debug: print("[DEBUG] ",*args)
    
    def run(self,f=None):
        if not self.server is None:
            self.ready()
        while True:
            socks = dict(self.poller.poll(1000*1000))
            #self.local_gateway.socket.send()
            for s in socks:
                #self.debug_print(s,socks[s])
                if socks[s] == zmq.POLLIN:
                    self.debug_print("RX",s,self.dataflows[s])
                    self.dataflows[s].recv()
                elif(socks[s] & zmq.POLLERR != 0):
                    self.debug_print('socket error',s,socks[s])
                    self.dataflows[s].die()
            if not f is None: 
                f()
