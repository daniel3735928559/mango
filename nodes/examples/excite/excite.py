from libmango import *

class excite(m_node):
    def __init__(self):
        super().__init__(debug=True)
        self.interface.add_interface('excite.yaml',{'excite':self.excite})
        self.run()
    def excite(self,header,args):
        return {'excited':args['str']+'!'}
    
excite()
