from libmango import *

class tttt(m_node):
    def __init__(self):
        super().__init__(debug=True)
        self.reset()
        self.interface.add_interface({'move':self.move,'restart':self.restart,'state':self.state,'rules':self.rules})
        self.win_lines = [[[x,y] for x in range(3)] for y in range(3)] + [[[x,y] for y in range(3)] for x in range(3)] + [[[x,x] for x in range(3)],[[x,2-x] for x in range(3)]]
        self.run()

    def board_str(self):
        return "\n"+"\n\n".join(["\n".join(["  ".join(["".join(self.board[by][bx][y]) for bx in range(3)]) for y in range(3)]) for by in range(3)])
        
    def check_win(self, board):
        for l in self.win_lines:
            x,y = l[0]
            winner = board[y][x]
            if winner == "_": continue
            for point in l[1:]:
                x,y = point
                if board[y][x] != winner:
                    winner = "_"
                    break
            if winner != "_":
                return winner
        return "_"

    def rules(self,args):
        rules = """Super tic-tac-toe:
This game consists of a tic-tac-toe board of tic-tac-toe boards.

To win an individual board, you must win the game on that board. Once a board is won, it may no longer be played upon. To win the whole game, you must win a 3-in-a-row line of boards.

First player may play on any board.

Whatever square of an individual board is played on, the next player must play on the corresponding board to that square. (For example, if a player plays on the middle square of any board, the next played must play their next move on the middle board.) 

The exception to this is that if that next board is already won, the next player is free to play on any available board."""
        return "message",{"text":rules}
    
    def state(self,args):
        return "state",{"turn":self.turn,"board":self.board,"boardstr":self.board_str(),"nextboard":"any open" if self.next_bx == -1 else "{},{}".format(self.next_bx,self.next_by)}
    
    def restart(self,args):
        self.restart_requests[args["turn"]] = True
        if self.restart_requests["X"] and self.restart_requests["Y"]:
            self.reset()
            return "message",{"text":"game has been reset"}
    
    def reset(self):
        self.board = [[[['_' for _ in range(3)] for _ in range(3)] for _ in range(3)] for _ in range(3)]
        self.winners = [['_' for _ in range(3)] for _ in range(3)]
        self.winner = "_"
        self.turn = "X"
        self.next_bx,self.next_by = -1,-1
        self.restart_requests = {"X":False,"O":False}
    
    def move(self,args):
        if self.winner != "_":
            return "message",{"text":"{} has already won!".format(self.winner)}
        if args['turn'] != self.turn:
            return "message",{"text":"It's not your turn!"}
        x,y,bx,by = args['x'],args['y'],args['bx'],args['by']
        moves = {0,1,2}
        if not (x in moves and y in moves and bx in moves and by in moves):
            return "message",{"text":"All coordinates must be in range 0-2!"}
        if self.winners[by][bx] != "_":
            return "message",{"text":"That board is already won"}
        if self.next_bx != -1 and self.next_by != -1 and (bx != self.next_bx or by != self.next_by):
            return "message",{"text":"Must move on board {},{}!".format(self.next_bx,self.next_by)}
        if self.board[by][bx][y][x] != "_":
            return "message",{"text":"That spot is taken!"}
        self.board[by][bx][y][x] = self.turn

        self.turn = "O" if self.turn == "X" else "X"
        
        board_winner = self.check_win(self.board[by][bx])
        if board_winner != "_":
            self.winners[by][bx] = board_winner
            self.m_send("message",{"text":"Board {},{} is now closed (won by {})!".format(bx,by,board_winner)})
            if self.check_win(self.winners) != "_":
                self.winner = board_winner
                return "message",{"text":"The winner is {}!".format(board_winner)}

        # Set the board the next player must play if that board isn't already won
        self.m_send("message",{"text":self.board_str()})
        
        if self.winners[y][x] == "_":
            self.next_bx,self.next_by = x,y
            return "message",{"text":"{}, it is your turn. Must play on board {},{}".format(self.turn,x,y)}
        else:
            self.next_bx,self.next_by = -1,-1
            return "message",{"text":"{}, it is your turn. May play on any available board".format(self.turn)}

tttt()
