import io, re, socketserver, socket, zmq, subprocess
from dataflow import m_dataflow
from mangolib import m_node
from serialiser import *
from transport import *
from obj import *

class mc_loopback_dataflow(m_dataflow):
    def __init__(self,interface,dispatch_cb,reply_df):
        self.interface = interface
        self.dispatch_cb = dispatch_cb
        self.reply_df = reply_df

    def send(self,header,msg,route):
        print("SS",header,msg)
        self.interface.validate(header['function'],msg)
        result = self.dispatch_cb(header,msg)
        self.reply_df.send(header,result,route)
        
    def send_raw(msg,route):
        pass#self.owner.ports["stdio"].send_raw(msg)

class mc_router_dataflow():
    def __init__(self,transport,serialiser,dispatch_cb):
        self.transport = transport
        self.serialiser = serialiser
        self.dispatch_cb = dispatch_cb
    
    def recv(self):
        route = self.transport.rx()
        data = self.transport.rx()
        header,args = self.serialiser.deserialise(data)
        self.dispatch_cb(header,args,data,route,self)

    def send(self,header,msg,route):
        print("MC ROUTER SEND",header,msg,"route =",route)
        data = self.serialiser.serialise(header,msg)
        self.transport.socket.send(route,zmq.SNDMORE)
        self.transport.tx(data)
        
    def send_raw(self,data,route):
        print("MC ROUTER SEND",header,msg,"route =",route)
        self.transport.socket.send(route,zmq.SNDMORE)
        self.transport.tx(data)
    

# class mc_dataflow(m_dataflow):
#     def __init__(self,owner,recv_cb,transport,serialiser,route):
#         super().__init__(owner,recv_cb,transport,serialiser)
#         self.route=route

#     def recv(self):
#         data = self.transport.rx()
#         print("MC DATAFLOWDATA",data)
#         ver,dport,h,a = self.owner.serialiser.unpack(data)
#         if(ver != self.owner.version): print("VERSION MISMATCH")
#         self.recv_cb(dport,h,a,data,self)

#     def send(self,msg_dict,src,mid=None,dport=None,reply_callback=None):
#         print("MC SEND",msg_dict,src,dport,"route",self.route)
#         mid=self.owner.get_mid() if mid is None else mid
#         data = self.owner.serialiser.pack(self.serialiser_method,msg_dict,src,mid,dport)
#         self.transport.socket.send(self.route,zmq.SNDMORE)
#         self.transport.tx(data)
#         if not reply_callback is None:
#             self.owner.outstanding[mid] = reply_callback

#     def send_raw(self,data,route):
#         self.transport.socket.send(bytearray(route),zmq.SNDMORE)
#         self.transport.tx(bytearray(data,"ASCII"))

# class mc_remote_srv_dataflow(m_dataflow):
#     def recv(self):
#         c,a = self.transport.rx()
#         df = mc_remote_dataflow(self.owner,self.recv_cb,m_client_sock(c),self.serialiser_method,a)
#         ip = str(a[0])
#         addr = ip+":"+str(a[1])
#         self.owner.routes[addr] = df
#         self.owner.dataflows[c.fileno()] = df
#         self.owner.ports[addr] = df
#         self.owner.nodes[addr] = Node(addr,0,df,self.owner.nodes["mc"],False)
#         self.owner.poller.register(c.fileno(),zmq.POLLIN)
#         print("A",a)

#     def send(self,msg_dict,sport,dport=None,reply_callback=None,serialiser="json",route=None):
#         pass

# class mc_remote_dataflow(m_dataflow):
#     def __init__(self,owner,recv_cb,transport,serialiser,addr):
#         super().__init__(owner,recv_cb,transport,serialiser)
#         self.addr = addr
#         self.straddr = str(self.addr[0])+":"+str(self.addr[1])

#     def send(self,msg_dict,src,mid=None,dport=None,reply_callback=None):
#         print("remote sending",msg_dict,"from",src)
#         mid=self.owner.get_mid() if mid is None else mid
#         n,p = self.owner.parse_port(src)
#         src=self.straddr + '/' + n + "." + p
#         data = self.owner.serialiser.pack(self.serialiser_method,msg_dict,src,mid,dport)
#         print("data")
#         print(data)
#         self.transport.tx(data)
#         if not reply_callback is None:
#             self.owner.outstanding[mid] = reply_callback

