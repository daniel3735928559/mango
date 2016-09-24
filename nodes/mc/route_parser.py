import ply.yacc as yacc
import ply.lex as lex
import sys,re,json

class route_lexer:
      def __init__(self):
            self.lexer = lex.lex(module=self)

      tokens = ('FROM',
                'TO',
                'CHAIN_DELIM',
                'ADD',
                'ADDN',
                'DEL',
                'COMM',
                'FILTER',
                'SH',
                'PORT',
                'JSON_DICT',
                'JSON_LIST',
                'JSON_STRING',
                'RAW')

      t_FROM = r"<"
      t_TO = r">"
      t_CHAIN_DELIM = r";"
      t_ADD = r"\+|add"
      t_ADDN = r"\*|addn"
      t_DEL = r"-|del"
      t_COMM = r"@|comm"
      t_FILTER = r"\?|filter"
      t_SH = r'sh "(?:[^"\\]|\\.)*"'
      t_RAW = r"%[a-z]+"

      def t_PORT(self,t):
            r"[_A-Za-z0-9.:/]+"
            if '/' in t.value:
                  t.value = list(t.value.rsplit('/',1))
            else:
                  t.value = (t.value,'stdio')
            return t

      def t_JSON_DICT(self,t):
            r'{("(?:[^"\\]|\\.)*":"(?:[^"\\]|\\.)*")(,"(?:[^"\\]|\\.)*":"(?:[^"\\]|\\.)*")*}'
            t.value = json.loads(t.value)
            return t

      def t_JSON_LIST(self,t):
            r'\["(?:[^"\\]|\\.)*"(,"(?:[^"\\]|\\.)*")*\]'
            t.value = json.loads(t.value)
            return t

      def t_JSON_STRING(self,t):
            r'"([a-zA-Z_][a-zA-Z0-9_]*)"'
            t.value = t.value[1:-1]
            return t

      t_ignore = " \t\n"


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
            

class route_parser:
      def __init__(self):
            self.lexer = route_lexer()
            self.tokens = self.lexer.tokens
            self.parser = yacc.yacc(module=self,write_tables=0,debug=False)
            self.transmogrifiers = ["add","addn","del","filter","sh"]

      def parse(self,data):
            if data:
                  return self.parser.parse(data,self.lexer.lexer,0,0,None)
            else:
                  return []

      def p_error(self,p):
            print('Error!',p)

      def p_chains(self,p):
            ''' chains : chain CHAIN_DELIM chains'''
            p[0] = [p[1]] + p[3]
      def p_chains_chain(self,p):
            ''' chains : chain'''
            p[0] = [p[1]]
      def p_chains_twoway_plus(self,p):
            ''' chains : twoway CHAIN_DELIM chains'''
            p[0] = p[1] + p[3]
      def p_chains__twoway(self,p):
            ''' chains : twoway'''
            p[0] = p[1]
      def p_twoway(self,p):
            ''' twoway : PORT FROM TO PORT'''
            p[0] = [[['port',p[1]],['port',p[4]]],[['port',p[4]],['port',p[1]]]]
      def p_chain_target(self,p):
            ''' chain : target'''
            p[0] = p[1]
      def p_chain_to(self,p):
            ''' chain : target TO chain'''
            p[0] = p[1]+p[3]
      def p_target_transmogrifier(self,p):
            ''' target : transmogrifier '''
            p[0] = p[1]
      def p_target_port(self,p):
            ''' target : PORT '''
            p[0] = [['port',p[1]]]
      def p_target_transmogrifiers(self,p):
            ''' target : PORT transmogrifiers'''
            p[0] = p[2]+[['port',p[1]]]
      def p_transmogrifiers_transmogrifier(self,p):
            ''' transmogrifiers : transmogrifier '''
            p[0] = p[1]
      def p_transmogrifiers(self,p):
            ''' transmogrifiers : transmogrifier transmogrifiers '''
            p[0] = p[1] + p[2]
      def p_transmogrifier_sub(self,p):
            ''' transmogrifier : JSON_DICT '''
            p[0] = [['edit',('sub',p[1])]]
      def p_transmogrifier_add(self,p):
            ''' transmogrifier : ADD JSON_DICT '''
            p[0] = [['edit',('add',p[2])]]
      def p_transmogrifier_addn(self,p):
            ''' transmogrifier : ADDN JSON_DICT '''
            p[0] = [['edit',('addn',p[2])]]
      def p_transmogrifier_del(self,p):
            ''' transmogrifier : DEL JSON_LIST '''
            p[0] = [['edit',('del',p[2])]]
      def p_transmogrifier_comm(self,p):
            ''' transmogrifier : COMM JSON_STRING '''
            p[0] = [['edit',('comm',p[2])]]
      def p_transmogrifier_filter(self,p):
            ''' transmogrifier : FILTER JSON_STRING '''
            p[0] = [['edit',('filter',p[2])]]
      def p_transmogrifier_raw(self,p):
            ''' transmogrifier : RAW '''
            p[0] = [['edit',('raw',p[1][1:])]]
      def p_transmogrifier_sh(self,p):
            ''' transmogrifier : SH '''
            p[0] = [['edit',('sh',p[1][4:-1])]]
