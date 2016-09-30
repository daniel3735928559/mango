from libmango import *

class excite(m_node):
    def __init__(self):
        super().__init__(debug=True)
        self.interface.add_interface('excite.yaml',{
            'excite':self.excite,
            'print':self.output
        })
        self.run()
    def output(self,header,args):
        print("PRINTING",header,args)
    def excite(self,header,args):
        print("EXCITING",header,args)
        return {'excited':args['str']+'!'}
t = excite()
