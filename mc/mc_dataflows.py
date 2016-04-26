import io, re, socketserver, socket, zmq, subprocess
from dataflow import m_dataflow
from mangolib import m_node
from serialiser import *
from transport import *
from obj import *

class loopback_dataflow(m_dataflow):
    def __init__(self,owner):
        self.owner = owner

    def send(self,msg_dict,src,mid=None,dport=None,reply_callback=None):
        self.owner.m_recv("stdio",{"source":src,"mid":mid},msg_dict,None,self.owner.dataflows[self.owner.local_gateway.socket])
        
    def send_raw(msg):
        pass#self.owner.ports["stdio"].send_raw(msg)

class mc_router_dataflow(m_dataflow):
    def __init__(self,owner,interface,transport,serialiser,dispatch_cb,reply_cb,error_cb):
        super().__init__(interface,transport,serialiser,dispatch_cb,reply_cb,error_cb)
        self.route = bytes()
        self.owner = owner

    def recv(self):
        self.route = self.transport.rx()
        if not self.route in self.owner.routes:
            data = self.transport.rx()
            header,args = self.serialiser.deserialise(data)
            self.dispatch_cb(header,args,data,self)
        else:
            self.owner.routes[self.route].recv()

    def send(self,header,msg,reply_callback=None):
        print("MC ROUTER SEND",header,msg,"route =",self.route)
        data = self.serialiser.serialise(header,msg)
        self.transport.socket.send(self.route,zmq.SNDMORE)
        self.transport.tx(data)
        if not reply_callback is None:
            self.owner.outstanding[mid] = reply_callback
    

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

