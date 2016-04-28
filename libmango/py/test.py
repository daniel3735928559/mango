from libmango import *

class test(m_node):
    def __init__(self):
        super().__init__("test","tcp://localhost:61453")
        self.run()

t = test()
