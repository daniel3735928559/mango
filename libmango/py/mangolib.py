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
        self.ports = {}
        self.outstanding = {}
        self.invalid_handler = lambda h,c: print("message not understood: "+ str(h))
        self.default_handler = lambda h,c: print("type not understood: "+ h['type'])
        self.default_cmd = lambda h,c: print("\n".join([x+": "+c[x] for x in c.keys() if x != 'command']))
        self.interfaces = {}
        self.interface.add_interface('/home/zoom/suit/mango/libmango/node_if.yaml',{'get_if':self.get_if,'reg':self.reg})
        if not server is None:
            self.local_gateway = m_ZMQ_transport(server,self.context,self.poller)
            s = self.local_gateway.socket
            self.dataflows[s] = m_dataflow(self.interface,self.local_gateway,self.serialiser,self.dispatch,self.handle_reply,self.handle_error)
            self.poller.register(s,zmq.POLLIN)
            self.ports["stdio"] = self.dataflows[s]
            self.ports["mc"] = self.dataflows[s]
            self.m_send('excite',{'str':'foo'},port="mc",reply_callback=print)

    def dispatch(self,header,args,dataflow):
        result = self.interface.interface[header['function']]['handler'](header,args)
        if not result is None:
            self.m_reply(header,result,dataflow)
        
    def get_if(self,header,args):
        print("GET IF")
        i = args['if']
        if i in self.interfaces.keys():
            return {'result':'success','if':self.interfaces[i]}
        else:
            return {'result':'failure'}

    def handle_reply(self,header,args):
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

    def handle_error(self,err):
        print(err)
        return None

    def get_mid(self):
        self.mid = (self.mid + 1)%(2**63)
        return self.mid

    def reg(self,header,args):
        self.key = args["key"]
        self.node_id = args["node_id"]
        self.local_gateway.set_id()
        print('my new node id')
        print(self.node_id)
        print("registered as " + self.node_id)

    def make_header(self,port):
        return {'port':port,
                'src_node':self.node_id,
                'src_port':port,
                'mid':self.get_mid(),
                'command':'call'}
    def m_send(self,command,msg,port="stdio",reply_callback=lambda h,a:None,async=True):
        print('sending',msg)
        mid = self.get_mid()
        header = self.make_header(port)
        header['function'] = command
        self.outstanding[mid] = reply_callback
        self.ports[port].send(header,msg)
        if not async and not reply_callback is None:
            self.ports[port].recv()
        print("outstanding",self.outstanding)

    # def m_error(self,error_msg,src,mid,dataflow):
    #     # print(msg_dict)
    #     msg_dict = {}
    #     msg_dict["command"] = "error"
    #     msg_dict["mid"] = mid
    #     msg_dict["source"] = src
    #     msg_dict["error"] = error_msg
    #     dataflow.send(msg_dict,src,mid)

    def m_recv(self,dport,h,c,raw,dataflow):
        return self.ports[dport].recv()
    
    def m_reply(self,header,msg,dataflow):
        # print(msg_dict)
        print("M REPLY",header,msg)
        h = {'command':'reply',
             'mid':header['mid'],
             'src_node':self.node_id,
             'src_port':header['port']}
        dataflow.send(header,msg)

    def reg(self,header,args):
        self.key = args["key"]
        self.node_id = args["node_id"]
        self.local_gateway.set_id(self.node_id)
        #print('my new node id')
        #print(self.node_id)
        #print("registered as " + self.node_id)

        
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
            #time.sleep(.01) 

# n = m_node("a","localhost:2323")
# n.m_send(0,{"a":"basdasd","c":bytearray([3])})
# n.m_send(0,{"hello":"2","blah":bytearray("abc","ASCII")})

# print(n.m_pack({"a":"basdasd","c":bytearray([3])}))
# print(n.m_unpack("MM 0.1 49 28\nTarget:Hello\nSource:blah\nMID:12129\nConnection:393\n5:5 hello:hello\n3:3 abc:abc\n"))
# print(n.m_prep(393,{"a":"basdasd","c":bytearray([3])}))
