import ply.yacc as yacc
import ply.lex as lex
import sys,re,json,copy

class query_lexer:
      def __init__(self):
            # nodes have: str:name, str:group, str:type, str:lang
            # routes have: str:name, str:group, node:src, node:dst
            # nodesets: * = all niodes currently displayed
            # examples:
            # route lang=="py" and 
            self.tokens = ['NAME','STRING','REGEX','AND','OR','NOT','XOR','EQ','NE']
            self.literals = ['(',')','.','~']
            self.lexer = lex.lex(module=self)
            
      def t_ID(self, t):
            r'[a-zA-Z_][a-zA-Z0-9_]*'#(\.[a-zA-Z_][a-zA-Z0-9_]*)*'
            r = re.compile(r'[a-zA-Z_][a-zA-Z0-9_]*')
            reserved = {
                  'and' : 'AND',
                  'or' : 'OR',
                  'not' : 'NOT',
                  'xor' : 'XOR',
            }
            if t.value in reserved:
                  t.type = reserved[t.value]
            elif r.match(t.value):
                  t.type = 'NAME'
            return t
            
      def t_NE(self,t):
            r'!='
            return t
      
      def t_EQ(self,t):
            r'=='
            return t
      
      def t_STRING(self,t):
            r'"(?:[^"\\]|\\.)*"'
            t.value = t.value[1:-1]
            return t

      def t_REGEX(self,t):
            r'/(?:[^/\\]|\\.)*/(i|g|ig|gi)?[0-9]*'
            t.value = t.value[1:-1]
            return t
      
      t_ignore = ' \t\n'

      def t_error(self,t):
            print('Lex Error!',t,t.value)

      def tokenize(self,data):
            self.lexer.input(data)
            while True:
                  tok = self.lexer.token()
                  if tok:
                        yield tok
                  else:
                        break

class query_parser:
      def __init__(self):
            # self.precedence = (
            #       ('nonassoc', '<', '>', 'LE', 'GE', 'EQ'),
            #       ('left', 'AND', 'OR'),
            #       ('left', '+', '-'),
            #       ('left', '*', '/'),
            #       ('right', 'EXP'),
            #       ('right', 'UMINUS'),
            # )
            self.lexer = query_lexer()
            self.tokens = self.lexer.tokens
            self.parser = yacc.yacc(module=self,write_tables=0,debug=False)

      def parse(self,data):
            if data:
                  ans = self.parser.parse(data,self.lexer.lexer,0,0,None)
                  return ans
            else:
                  return []

      def p_error(self,p):
            print('Error!',p)

      def p_query_test(self,p):
            ''' query : test'''
            p[0] = ('test',p[1])
            
      def p_query_paren(self,p):
            ''' query : '(' query ')' '''
            p[0] = p[2]
            
      def p_query_and(self,p):
            ''' query : query AND query '''
            p[0] = ('and', p[1], p[3])
            
      def p_query_or(self,p):
            ''' query : query OR query '''
            p[0] = ('or', p[1], p[3])
            
      def p_query_xor(self,p):
            ''' query : query XOR query '''
            p[0] = ('xor', p[1], p[3])
            
      def p_query_not(self,p):
            ''' query : NOT query '''
            p[0] = ('not', p[2])
            
      def p_test_eq(self,p):
            ''' test : var EQ STRING '''
            p[0] = ('eq', p[1], p[3])
            
      def p_test_ne(self,p):
            ''' test : var NE STRING '''
            p[0] = ('ne', p[1], p[3])
            
      def p_test_like(self,p):
            ''' test : var '~' REGEX '''
            p[0] = ('like', p[1], p[3])
            
      def p_var_name(self,p):
            ''' var : name '''
            p[0] = ('var', p[1])

      def p_name_name(self,p):
            ''' name : NAME '''
            p[0] = [p[1]]
            
      def p_name_dot(self,p):
            ''' name : NAME '.' name '''
            p[0] = [p[1]] + p[3]
            
if __name__ == "__main__":
      qp = query_parser()
      query = qp.parse(sys.argv[1])
      print(query)
