from libmango import *
import time, sys

class test(m_node):
    def __init__(self):
        name = sys.argv[1]
        super().__init__(name,"tcp://localhost:61453")
        self.interface.add_interface('/home/zoom/suit/mango/nodes/test/test.yaml',{
            'print':self.output
        })
        self.m_send('route',{'spec':'{} > excite > +{{"a":"{}"}} > {}'.format(name, name, name)},callback="print",port="mc")
        self.m_send('excite',{'thestring':'foo'},callback="print")
        self.m_send('excite',{'str':'foo'},callback="print")
        self.run()
    def output(self,header,args):
        print("GOTIT",header,args)
t = test()
