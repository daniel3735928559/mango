import sys, os, re, signal, time, zmq, traceback, json
from transport import *
from dataflow import *
from libmango import *
from lxml import etree
import jinja2

# Usage: 
#   ./mx -h blah [print [arg1 ...]|file filename [arg1 ...]|forward conn]
#
#   To set the handler for the blah command to print (or just print
#   one argument, if you like) or forward to the connection with
#   number conn
#
#   ./mx [-t type] conn [-c] cmd -argname1 argval1 ...
#
#   To send the command cmd to connection conn with the specified
#   arguments

class mx_dataflow(m_dataflow):
      def __init__(self,recv_cb,transport):
            super().__init__(None,transport,None,recv_cb,None)
      def recv(self):
            data = self.transport.rx()
            #print("D",data)
            self.dispatch_cb(data)

      def send(self,msg):
            self.transport.tx(msg)

class zmq_rep_transport():
      def __init__(self,owner,context,poller):
            self.owner = owner
            self.socket = context.socket(zmq.REP)
            self.port = self.socket.bind_to_random_port("tcp://*",min_port=10000,max_port=30000)
            #print("P",self.port)
            poller.register(self.socket)

      def tx(self,payload):
            self.socket.send_unicode(payload)
            return "sent"

      def rx(self):
            return self.socket.recv()


class mx_agent(m_node):
      def __init__(self):  
            super().__init__(debug=True)
            self.flags = {"strict":False}
            self.fc = mx_dataflow(self.handler,zmq_rep_transport(self,self.context,self.poller))
            self.dataflows[self.fc.transport.socket] = self.fc
            time.sleep(1)
            #self.load_if("mc")
            self.pid = os.getpid()
            self.command_cbs = {}
            self.interface.add_interface('./mx.yaml',{
                  'reply':self.mx_handle_reply
            })

            self.handlers = {'print':self.print_handler,'file':self.file_handler,'forward':self.forward_handler,'nop':self.nop_handler}
            self.handler_arg_defaults = {'file':{'filename':'/dev/stdout'},'print':{'format':None},'nop':{}}
            #os.mkfifo('.run/'+str(self.pid)+'_out')
            self.sh_init = open('.run/'+str(self.pid),'w')
            with open('.init','r') as f:
                  sh_init_base = f.read()
            self.sh_init.write(sh_init_base.format(pid=os.getpid(),port=self.fc.transport.port))
            if len(sys.argv) >= 4:
                  f = open(sys.argv[3],"r")
                  self.sh_init.write("\n"+f.read())
                  f.close()
            self.sh_init.flush()
            #print("asd")
            self.output = open('.run/'+str(self.pid)+'_out','w')#os.open('.run/'+str(self.pid)+'_out',os.O_WRONLY)
            self.output.write("asd2\n")
            self.output.flush()
            #print("asd3")
            self.run()


      def mx_handle_reply(self,header,args):
            #print("AAAAAAA",header,args)
            self.output.write("H "+str(header))
            self.output.write("A "+str(args))
            self.output.flush()
            
      def setup_if(self,h,a):
            src_node,src_port = self.parse_port(h['source'])
            self.output.write("HEAD",h,"SOURCE",src_node,"PORT",src_port)
            self.output.write(a)
            if(a['result'] == 'success'):
                  self.output.write("IF_NAME",if_name)
                  self.sh_init.write(str(etree.XSLT(etree.parse('mx.xsl'))(etree.XML(a['if']),node_name="'"+src_node+"'",port=("'mc'" if if_name=="mc" else "'stdio'"))))
                  self.sh_init.flush()
                  self.output.write("Got it")
            else:
                  self.output.write("Fail")
            self.output.flush()
            

      def load_if(self,if_name):
            if if_name == "mc":
                  self.m_send({"command":"get_if","if":if_name}, callback="setup_if", port="mc")
            else:
                  self.m_send({"command":"get_if","if":if_name}, callback="setup_if")
            
            # def c_1(h,a,t):
            #       return ({"command":"get_if","if":if_name},"mc","mc"),()
            # def c_2(h,a,t):
            #       print(a)
            #       if(a['result'] == 'success'):
            #             self.sh_init.write(str(etree.XSLT(etree.parse('mx.xsl'))(etree.XML(a['if']),conn='mc')))
            #             self.sh_init.flush()
            # return self.chain_cbs({},{},(),[c_1,c_2])
            #print("registered")

      # def reg_cb(self,header,reply):
      #       self.m_send({"command":"mc/get_if","if":"mc"},reply_callback=self.if_next,False)

      # def if_next(self,header,args,conn='0'):
      #       print(args)
      #       self.sh_init.write(str(etree.XSLT(etree.parse('mx.xsl'))(etree.XML(args['if'].decode()),conn=conn)))
      #       self.sh_init.flush()
      #       #print('source .run/'+str(self.pid))

      def file_handler(self,header,args,handler_args):
            with open(filename) as f: 
                  f.write("\n".join([x+':'+args[x].decode('utf-8') for x in args.keys()]))
            self.output.write("Done")
            self.output.flush()

      def nop_handler(self,header,args,handler_args):
            self.output.write("")

      def print_handler(self,header,args,handler_args):
            print(header,args,handler_args)
            if(handler_args is None):
                  self.output.write("\n".join([x+':'+args[x].decode('utf-8') for x in args.keys()]))
            else:
                  self.output.write(re.sub(r'([^\\]*)\{([a-zA-Z0-9_]+)\}',lambda m: m.group(1) + (args[m.group(2)].decode() if m.group(2) in args.keys() else ""),handler_args['format']))
            self.output.flush()
          
      def forward_handler(self,header,message,conns):
            for c in conns:
                  m_send(0,c,message,lambda h,c:print(c),False,mtype=header['type'])
  
      def parse_callback(self,cb):
            rcb = self.reply_cb
            ha = {}
            h = cb.split(" ",maxsplit=1)
            if h[0] in self.handlers.keys():
                  rcb = self.handlers[h[0]]
                  try:
                        for x in h[1].split(','):
                              y = x.split('=')
                              ha[y[0].strip()] = y[1].strip()
                  except:
                        exc_type, exc_value, exc_traceback = sys.exc_info()
                        traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
                        traceback.print_exception(exc_type, exc_value, exc_traceback,limit=10, file=sys.stdout)
                        self.output.write('Broken reply callback specifier')
                        return None
            else:
                  self.output.write('Bad reply callback specifier')
                  return None
            return ha,rcb

      def handler(self,msg):
            self.fc.send("ack")
            #print(self.outstanding)
            #print(msg)
            args = json.loads(msg.decode())
            #print(args)
            cmd = args["mx_cmd"]
            del args["mx_cmd"]
            if(cmd == 'send'):
                  rcb = self.reply_cb
                  ha = {}
                  if 'p' in args:
                        port = args['p']
                        del args['p']
                  else:
                        port = "stdio"
                  if "r" in args:
                        new_cb = self.parse_callback(args["r"])
                        if(not new_cb is None):
                              ha,rcb = new_cb
                        del args["r"]
                  c = args.pop('command')
                  self.m_send(c,args,port=port)
            if(cmd == 'mc'):
                  rcb = self.reply_cb
                  ha = {}
                  port = "mc"
                  if "r" in args:
                        new_cb = self.parse_callback(args["r"])
                        if(not new_cb is None):
                              ha,rcb = new_cb
                        del args["r"]
                  c = args.pop('command')
                  self.m_send(c,args,port=port)
            elif cmd == 'handle':
                  try:
                        ha,rcb = {},self.reply_cb
                        new_cb = self.parse_callback(args["handler"])
                        if(not new_cb is None):
                              ha,rcb = new_cb
                        self.command_cbs[args["command"]] = lambda header,reply: rcb(header,reply,ha)
                  except:
                        self.output.write('Broken reply callback specifier')
            elif cmd == 'load_if':
                  self.load_if(args["if"])
                  # Load an interface file so that tab-completion for its functions works
                  pass
            elif cmd == 'unload_if':
                  pass

      #       elif cmd == 'connect':
      #             print(arr)
                       
      #             try:
      #                   self.mconnect(arr[1])
      #                   #self.m_send(0,0,{'command':'mc/cx_list'},lambda h,a: self.cx_list_cb(h,a,arr[1]))
      #             except:
      #                   exc_type, exc_value, exc_traceback = sys.exc_info()
      #                   traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
      #                   traceback.print_exception(exc_type, exc_value, exc_traceback,limit=2, file=sys.stdout)
      #                   traceback.print_exc()
      #                   self.fc.send('Bad')

      # def mconnect(self,target): 
      #       print(target)
      #       def c_1(h,a,t):
      #             try:
      #                   return ({'command':'mc/cx_list'}),t
      #             except:
      #                   exc_type, exc_value, exc_traceback = sys.exc_info()
      #                   traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
      #                   traceback.print_exception(exc_type, exc_value, exc_traceback,limit=2, file=sys.stdout)
      #                   traceback.print_exc()
      #                   self.fc.send('Bad')
      #       def c_2(h,a,t):
      #             print(h,a,t)
      #             l = a['list'].split(';')
      #             cs = [x.split(',') for x in l if t[0] in x]
      #             if(len(cs) == 0):
      #                   self.fc.send('No connection')
      #             elif(len(cs) == 1):
      #                   print("GOT IT")
      #                   return (0,0,{'command':'mc/cx_add','conn_id':cs[0][1]}),(cs[0][0],int(cs[0][1]))
      #             else:
      #                   self.fc.send('Possibilities: \n' + "\n".join(l))

      #       def c_3(h,a,t):
      #             print("CONN_CB",t,str(t[1]))
      #             return (0,t[1],{'command':t[0]+'/get_if','if':t[0]}),(str(t[1]))
            
      #       def c_4(h,a,t):
      #             print(a,t)
      #             self.sh_init.write(str(etree.XSLT(etree.parse('mx.xsl'))(etree.XML(a['if']),conn=t[0])))
      #             self.sh_init.flush()
      #             self.fc.send('Success')
                  
      #       self.chain_cbs({},{},[target],[c_1,c_2,c_3,c_4])

      def chain_cbs(self,h,a,t,cb_list):
            print('args',h,a,t,len(cb_list))
            if(len(cb_list) == 0):
                  print("RETING",h,a,t)
                  return h,a,t
            retval = cb_list[0](h,a,t)
            print("RET",retval,retval[0])
            self.m_send(retval[0][0], port=retval[0][1], dport=retval[0][2], reply_callback=lambda h,a: self.chain_cbs(h,a,retval[1],cb_list[1:]))

      def reply_cb(self,header,reply,handler_args):
            #print(reply)
            #print(header)
            
            self.output.write("HEADER:\n"+"\n".join([k+":"+str(header[k]) for k in header.keys()])+"\nBODY:\n"+"\n".join([k+":"+reply[k] for k in reply.keys()]))
            

mx_agent()
