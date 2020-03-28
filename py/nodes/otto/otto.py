import re, os, sys, random, time
from threading import Timer
from libmango import *

class otto(m_node):
    def __init__(self):
        self.root = os.getenv("OTTOROOT")        
        if self.root is None or len(root) == 0:
            self.root = os.path.join(os.getenv("HOME"),".otto/")

        self.data = {
            "joke": [],
            "insult": [],
            "quote": []}
        
        self.logfile = os.path.join(self.root, "otto.log")
        self.blacklist = []
        for ty in self.data:
            for l in self.read_whole_file(ty).split("\n"):
                if len(l) == 0:
                    continue
                self.data[ty] += [l.split("...")]
            
        super().__init__(debug=True)
        self.interface.add_interface({
            'joke':lambda args: self.getdata('joke',args),
            'insult':lambda args: self.getdata('insult',args),
            'quote':lambda args: self.getdata('quote',args)})
        
        self.debug_print("running",self.data)
        self.run()

    def write_to_logfile(self,s):
        with open(self.logfile,'a') as f:
            f.write(time.strftime("%Y%m%d:%H%M%S %Z",time.localtime()) + ": " + s + "\n")

    def read_whole_file(self,fn):
        with open(os.path.join(self.root, fn)) as f: return f.read()

    def getdata(self,ty,args):
        j = random.choice(self.data[ty])
        
        if "addr" in args:
            self.write_to_logfile(ty + " ping: " + str(args["addr"]))
            for line in j:
                self.m_send(ty,{'addr':str(args["addr"]),'text':line})
                time.sleep(2)
        else:
            self.write_to_logfile(ty + " ping: <anon>")
            j = " ...".join(j)
            return ty,{'text':j}

otto()
