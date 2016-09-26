import time, os
print(os.environ,os.getcwd())
from libmango import *

class excite(m_node):
    def __init__(self):
        super().__init__(debug=True)
        self.interface.add_interface(os.path.join(os.path.abspath(os.path.dirname(os.path.abspath(__file__))),'excite.yaml'),{
            'excite':self.excite,
            'print':self.output
        })
        self.ready()
        self.run()
    def output(self,header,args):
        print("PRINTING",header,args)
    def excite(self,header,args):
        print("EXCITING",header,args)
        return {'excited':args['str']+'!'}
t = excite()