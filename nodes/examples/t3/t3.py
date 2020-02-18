from libmango import *

class t3(m_node):
    def __init__(self):
        super().__init__(debug=True)
        self.board = [['_','_','_'],['_','_','_'],['_','_','_']]
        self.turn = "X"
        self.interface.add_interface('t3.yaml',{'move':self.move,'new':self.new})
        self.run()
    
    def move(self,header,args):
        x,y = args['x'],args['y']
        if board[x][y] == '_':
            board[x][y] = self.turn
            self.turn = "X" if self.turn == "O" else "O"
        self.check_win()
        return "board",{'turn':self.turn,'board':"\n".join("".join(self.board))}
    
t3()
