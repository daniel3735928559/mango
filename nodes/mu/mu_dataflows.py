import io, re, socketserver, socket, zmq, subprocess
from dataflow import *
from libmango import *
from serialiser import *
from transport import *
from mu_transport import *
#from http_parser.parser import HttpParser

# class mu_srv_dataflow(m_dataflow):
#     def __init__(self,owner,get_cb,poll_cb,recv_cb,transport,serialiser="json"):
#         super().__init__(owner,recv_cb,transport,serialiser)
#         self.get_cb = get_cb
#         self.poll_cb = poll_cb

#     def recv(self):
#         c,a = self.transport.rx()
#         df = mu_dataflow(self.owner,self.get_cb,self.poll_cb,self.recv_cb,m_client_sock(c),self.serialiser_method,a)
#         ip = str(a[0])
#         addr = ip+":"+str(a[1])
#         self.owner.dataflows[c.fileno()] = df
#         self.owner.ports[addr] = df
#         self.owner.poller.register(c.fileno(),zmq.POLLIN)
#         print("A",a)

#     def send(self,msg_dict,sport,dport=None,reply_callback=None,serialiser="json",route=None):
#         pass

# class mu_dataflow(m_dataflow):
#     def __init__(self,owner,get_cb,poll_cb,recv_cb,transport,serialiser,addr):
#         super().__init__(owner,recv_cb,transport,serialiser)
#         self.get_cb = get_cb
#         self.poll_cb = poll_cb
#         self.addr = addr
#         self.straddr = str(self.addr[0])+":"+str(self.addr[1])
#         self.fd = self.transport.socket.fileno()

#     def die(self):
#         try: del self.owner.dataflows[self.fd]
#         except: pass
#         try: self.owner.poller.unregister(self.fd)
#         except: pass
#         self.transport.die()

#     def recv(self):
#         data = self.transport.rx()
#         print("MU data",data)
#         if(data is None or len(data) == 0): 
#             self.die()
#             return
#         self.hp = HttpParser()
#         self.hp.execute(data,len(data))
#         method = self.hp.get_method()
#         print(self.hp.get_url())
#         if method == 'GET':
#             path = self.hp.get_url()
#             file_data = self.get_cb(path)
#             if file_data is None: return
#             try:
                
#                 hdr = bytes("HTTP/1.1 200 OK\nContent-Length: "+str(len(file_data)) + "\n\n","ASCII")
#                 self.transport.tx(hdr+file_data)
#             except:
#                 err = bytes("Not Found","ASCII")
#                 hdr = bytes("HTTP/1.1 404 Not Found\nContent-Length: "+str(len(err)) + "\n\n","ASCII")
#                 self.transport.tx(hdr+err)
#         elif method == 'POST':
#             body = self.hp.recv_body()
#             #print(body)
#             if('poller' in self.hp.get_url()):
#                 self.poll_cb(self)
#             else:
#                 try:
#                     data = body.decode()
#                     print("decoded data",data)
#                     ver,dport,h,a = self.owner.serialiser.unpack(body)
#                     if(ver != self.owner.version): self.err("VERSION MISMATCH",True)
#                     resp = bytes("Thank you","ASCII")
#                     hdr = bytes("HTTP/1.1 200 OK\nContent-Length: "+str(len(resp)) + "\n\n","ASCII")
#                     self.transport.tx(hdr+resp)
#                     self.recv_cb(dport,h,a,body,self)
#                 except:
#                     exc_type, exc_value, exc_traceback = sys.exc_info()
#                     traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
#                     traceback.print_exception(exc_type, exc_value, exc_traceback,file=sys.stdout)
#                     resp = bytes("Bad request","ASCII")
#                     hdr = bytes("HTTP/1.1 401 Bad\nContent-Length: "+str(len(resp)) + "\n\n","ASCII")
#                     self.transport.tx(hdr+resp)
#                     print("bad")    


#     def send(self,msg_dict,src,mid=None,dport=None,reply_callback=None):
#         print("remote sending",msg_dict,"from",src)
#         mid=self.owner.get_mid() if mid is None else mid
#         n,p = self.owner.parse_port(src)
#         data = self.owner.serialiser.pack(self.serialiser_method,msg_dict,src,mid,dport)
#         print("data")
#         print(data)
#         hdr = bytes("HTTP/1.1 200 OK\nContent-Type: text/plain\nContent-Length: "+str(len(data)) + "\n\n","ASCII")
#         self.transport.tx(hdr+data)
#         if not reply_callback is None:
#             self.owner.outstanding[mid] = reply_callback

#     def send_raw(self,data):
#         print("remote sending raw",data)
#         print("data")
#         print(data)
#         hdr = bytes("HTTP/1.1 200 OK\nContent-Type: text/plain\nContent-Length: "+str(len(data)) + "\n\n","ASCII")
#         self.transport.tx(hdr+data)
#         if not reply_callback is None:
#             self.owner.outstanding[mid] = reply_callback


class mu_server_dataflow(m_dataflow):
    def __init__(self,interface,transport,serialiser,dispatch_cb,error_cb,owner):
        super().__init__(interface,transport,serialiser,dispatch_cb,error_cb)
        self.owner = owner
    
    def recv(self):
        c,a = self.transport.rx()
        df = mu_client_dataflow(self.interface,mu_client_ws(c),self.serialiser,self.dispatch_cb,self.error_cb,a)
        self.owner.dataflows[c.fileno()] = df
        self.owner.client_ws_dataflows[c.fileno()] = df
        self.owner.poller.register(c.fileno(),zmq.POLLIN)
        print("A",a)

    def send(self,header,msg):
        pass

class mu_client_dataflow(m_dataflow):
    def __init__(self,interface,transport,serialiser,dispatch_cb,error_cb,addr):
        super().__init__(interface,transport,serialiser,dispatch_cb,error_cb)
        self.addr = addr
        self.straddr = str(self.addr[0])+":"+str(self.addr[1])

    def recv(self):
        data = self.transport.rx()
        print("MU DATA",data);
        if data is None:
            return None
        try:
            header,args = self.serialiser.deserialise(data)
            result = self.dispatch_cb(header,args)
        except m_error as exc:
            # self.error_cb(header['src_node'],str(exc))
            print(header['src_node'],str(exc))
            return None
        
    # def send(self,header,msg):
    #     try:
    #         print("remote sending",msg,"from",header['src'])
    #         data = self.serialiser.serialize(header,msg)
    #         print("data")
    #         print(data)
    #         self.transport.tx(data)
    #     except m_error as exc:
    #         self.error_cb(header['src_node'],str(exc))
    #         return None



