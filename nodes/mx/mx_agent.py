import sys, os, re, time, zmq, traceback, json
from libmango import *
from collections import OrderedDict
from datetime import datetime

COLORS = {'RED':'\033[01;31m','GREEN':'\033[01;32m','YELLOW':'\033[01;33m','NORMAL':'\033[0m'}

class mx_agent(m_node):
      def __init__(self):  
            super().__init__(debug=True)
            self.pid = os.getpid()
            self.ipc_filename = "ipc://.run/{}_ipc.tmp".format(str(self.pid))
            self.init_filename = '.run/{}_init.tmp'.format(str(self.pid))
            self.output_filename = '.run/{}_out.tmp'.format(str(self.pid))
            
            self.interface.default_handler = self.recv

            self.messages = [] # message:{time:num, name:str, args:str, type:["<",">"]}

            self.shell_socket = self.context.socket(zmq.ROUTER)
            self.shell_socket.bind(self.ipc_filename)
            self.add_socket(self.shell_socket, self.shell_input, self.shell_error)
            
            self.sh_init = open(self.init_filename,'w')
            with open('.init','r') as f:
                  sh_init_base = f.read()
            self.sh_init.write(sh_init_base.format(pid=os.getpid(), ipc=self.ipc_filename))
            self.sh_init.close()
            self.output = open(self.output_filename,'w')
            self.shell_output()
            self.run()

      def recv(self,header,args):
            self.messages.append({"time":time.time(), "name":header['name'], "args":args, "type":"<"})
            self.shell_output()
            
      def shell_error(self,msg):
            self.debug_print("SHELLSHOCK DX")
            
      def shell_input(self):
            msg = self.shell_socket.recv_multipart()
            self.debug_print("MSG",msg)
            msg = msg[1]
            args = json.loads(msg.decode())
            cmd = args.pop("_mx_cmd")
            if(cmd == 'send'):
                  c = args.pop('_name')
                  for a in args:
                        try: args[a] = json.loads(args[a])
                        except: pass
                  self.messages.append({"time":time.time(), "name":c, "args":args, "type":">"})
                  self.shell_output()
                  self.m_send(c,args)

      def shell_output(self):
            self.output.write('\033[2J\033[0;0H')
            for m in self.messages:
                  color = COLORS['GREEN'] if m['type'] == ">" else COLORS['YELLOW']
                  timestr = datetime.fromtimestamp(m['time']).strftime("%Y-%m-%d %H:%M:%S")
                  margs = " ".join(["--{}={}".format(x, m['args'][x]) for x in m['args']])
                  self.output.write('{color}{type} [{timestr}] {name} {margs}{NORMAL}\n'.format(color=color, timestr=timestr, margs=margs, **m, **COLORS))
            self.output.flush()

mx_agent()
