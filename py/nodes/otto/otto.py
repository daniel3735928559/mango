import re, os, sys, random, time
from threading import Timer
from libmango import *

class otto(m_node):
    def __init__(self):
        self.root = os.getenv("OTTOROOT")        
        if self.root is None or len(root) == 0:
            self.root = os.path.join(os.getenv("HOME"),".otto/")

        self.jokes = []
        self.logfile = os.path.join(self.root, "otto.log")
        self.blacklist = []
        for l in self.read_whole_file("jokes").split("\n"):
            if len(l) == 0:
                continue
            self.jokes += [l.split("...")]
            
        super().__init__(debug=True)
        self.interface.add_interface({'joke':self.joke})
        self.debug_print("running")
        self.run()

    def write_to_logfile(self,s):
        with open(self.logfile,'a') as f:
            f.write(time.strftime("%Y%m%d:%H%M%S %Z",time.localtime()) + ": " + s + "\n")

    def read_whole_file(self,fn):
        with open(os.path.join(self.root, fn)) as f: return f.read()

    def joke(self,args):
        j = random.choice(self.jokes)
        
        if "addr" in args:
            self.write_to_logfile("Joke ping: " + str(args["addr"]))
            for line in j:
                self.m_send("joke",{'addr':str(args["addr"]),'joke':line})
                time.sleep(2)
        else:
            j = " ...".join(j)
            return "joke",{'joke':j}
    
otto()
