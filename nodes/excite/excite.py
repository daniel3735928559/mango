from libmango import *
import time

class excite(m_node):
    def __init__(self):
        super().__init__("excite","tcp://localhost:61453",debug=True)
        self.interface.add_interface('/home/zoom/suit/mango/nodes/excite/excite.yaml',{
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
