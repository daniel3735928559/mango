from libmango import *
import time

class test(m_node):
    def __init__(self):
        super().__init__("test","tcp://localhost:61453")
        self.interface.add_interface('/home/zoom/suit/mango/nodes/test/test.yaml',{
            'print':self.output
        })
        self.m_send('route',{'spec':'test <> excite'},callback="print",port="mc")
        self.m_send('excite',{'str':'foo'},callback="print")
        self.run()
    def output(self,header,args):
        print("GOTIT",header,args)
t = test()
