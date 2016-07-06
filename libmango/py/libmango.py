import io, re, time, signal, os
from serialiser import *
from dataflow import *
from transport import *
from interface import *
from error import *

class m_node: 
    def __init__(self,node_id,server=None):
        self.version = "0.1"
        self.context = zmq.Context()
        self.poller = zmq.Poller()
        self.serialiser = m_serialiser(self.version)
        self.interface = m_if()
        self.dataflows = {}
        self.node_id = node_id
        self.mid = 0
        self.outstanding = {}
        self.invalid_handler = lambda h,c: print("message not understood: "+ str(h))
        self.default_handler = lambda h,c: print("type not understood: "+ h['type'])
        self.default_cmd = lambda h,c: print("\n".join([x+": "+c[x] for x in c.keys() if x != 'command']))
        self.interfaces = {}
        self.interface.add_interface('/home/zoom/suit/mango/libmango/node_if.yaml',{
            'get_if':self.get_if,
            'reg':self.reg,
            'reply':self.reply
        })
        if not server is None:
            self.local_gateway = m_ZMQ_transport(server,self.context,self.poller)
            s = self.local_gateway.socket
            self.dataflow = m_dataflow(self.interface,self.local_gateway,self.serialiser,self.dispatch,self.handle_error)
            self.dataflows[s] = self.dataflow
            self.poller.register(s,zmq.POLLIN)

    def ready():
        self.m_send('hello',{},callback="print",port="mc")
            
    def dispatch(self,header,args):
        print("DISPATCH",header,args)
        result = self.interface.interface[header['command']]['handler'](header,args)
        if not result is None:
            self.m_send(header['callback'],result,port=header['port'],mid=header['mid'])
        
    def get_if(self,header,args):
        print("GET IF")
        i = args['if']
        if i in self.interfaces.keys():
            return {'result':'success','if':self.interfaces[i]}
        else:
            return {'result':'failure'}

    def reply(self,header,args):
        mid = header["mid"]
        print("REPLY HANDELR")
        print(mid,self.outstanding)
        if not int(mid) in self.outstanding.keys():
            print("Fake reply",mid,self.outstanding)
            return None
    
        del args["source"]
        del args["mid"]

        self.outstanding[mid](header,args) # Maybe build in here some checks that the reply contains what we expected?
        del self.outstanding[mid]
        return None

    def handle_error(self,src,err):
        print(err)
        self.m_send('error',{'source':src,'message':err},port="mc")
        return None

    def get_mid(self):
        self.mid = (self.mid + 1)%(2**63)
        return self.mid

    def reg(self,header,args):
        self.key = args["key"]
        self.node_id = args["node_id"]
        #self.local_gateway.set_id(args["node_id"])
        print('my new node id')
        print(self.node_id)
        print("registered as " + self.node_id)

    def make_header(self,command,callback=None,mid=None,src_port='stdio'):
        header = {'src_node':self.node_id,
                  'src_port':src_port,
                  'mid':self.get_mid() if mid is None else mid,
                  'command':command}
        if callback is None and command != 'reply':
            header['callback'] = 'reply'
        elif not callback is None:
            header['callback'] = callback
        print("H",header)
        return header
    
    def m_send(self,command,msg,callback=None,mid=None,port='stdio',reply_callback=None,async=True):
        print('sending',msg)
        header = self.make_header(command,callback,mid,port)
        self.dataflow.send(header,msg)
        if not async and not callback is None:
            self.dataflow.recv()
        print("outstanding",self.outstanding)
        return header['mid']
        
    def run(self,f=None):
        while True:
            socks = dict(self.poller.poll(1000*1000))
            #self.local_gateway.socket.send()
            for s in socks:
                #print(s,socks[s])
                if socks[s] == zmq.POLLIN:
                    self.dataflows[s].recv()

                elif(socks[s] & zmq.POLLERR != 0):
                    print('socket error',s,socks[s])
                    self.dataflows[s].die()
            if not f is None: 
                f()
