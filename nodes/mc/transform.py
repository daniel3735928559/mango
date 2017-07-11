import ply.yacc as yacc
import ply.lex as lex
import sys,re,json,copy

class transform_lexer:
      def __init__(self):
            self.tokens = ['NAME','FLOAT','INT','F','E','R','FE','FR','TRUE','FALSE','POP','STRING','REGEX','GE','EQ','LE','EXP','AND','OR','PE','ME','TE','DE','RE','BI','ID']
            self.literals = ['<','>',';','+','-','{','}','[',']','(',')','=','!','&','|','~',',',':','*','/','?']
            self.t_RE = '~='
            self.t_PE = '\+='
            self.t_ME = '-='
            self.t_TE = '\*='
            self.t_DE = '/='
            self.t_GE = '>='
            self.t_LE = '<='
            self.t_EQ = '=='
            self.t_BI = '<>'
            self.t_EXP = '\*\*'
            self.t_AND = '&&'
            self.t_OR = '\|\|'
            #self.t_ID = r'[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)+'
            self.lexer = lex.lex(module=self)
            
      def t_FLOAT(self,t):
            r'-?[0-9]*\.[0-9]+'
            t.value = float(t.value)
            return t
      
      def t_INT(self,t):
            r'-?[0-9]+'
            t.value = int(t.value)
            return t
      
      def t_ID(self, t):
            r'[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)*'
            r = re.compile(r'[a-zA-Z_][a-zA-Z0-9_]*')
            reserved = {
                  'f' : 'F',
                  'e' : 'E',
                  'r' : 'R',
                  'fe' : 'FE',
                  'fr' : 'FR',
                  'true' : 'TRUE',
                  'false' : 'FALSE',
                  'pop' : 'POP'
            }
            if t.value in reserved:
                  t.type = reserved[t.value]
            elif r.match(t.value):
                  t.type = 'NAME'
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

class transform_parser:
      def __init__(self):
            self.precedence = (
                  ('nonassoc', '<', '>', 'LE', 'GE', 'EQ'),
                  ('left', 'AND', 'OR'),
                  ('left', '+', '-'),
                  ('left', '*', '/'),
                  ('right', 'EXP'),
                  ('right', 'UMINUS'),
            )
            self.lexer = transform_lexer()
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

      # A route is a list of actual routes, e.g.:
      # "node1 > node2 > node3"
      # parses to: 
      # [[node1, node2], [node2, node3]]
            
      def p_route_node_node(self,p):
            ''' route : node '>' node'''
            p[0] = [[p[1], p[3]]]
            
      def p_route_trans_node(self,p):
            ''' route : transform '>' node'''
            p[0] = [p[1] + [p[3]]]
            
      def p_route_node_route(self,p):
            ''' route : node '>' route'''
            if p[3][0][0][0] == 'node':
                  p[0] = [[p[1],p[3][0][0]]] + p[3]
            else:
                  p[0] = [[p[1]]+p[3][0]]+p[3][1:]
            
      def p_route_trans_route(self,p):
            ''' route : transform '>' route'''
            if p[3][0][0][0] == 'node':
                  p[0] = [p[1]+[p[3][0][0]]] + p[3]
            else:
                  p[0] = [p[1]+p[3][0]]+p[3][1:]

      def p_route_bidirectional(self,p):
            ''' route : node BI node'''
            p[0] = [[p[1], p[3]], [p[3], p[1]]]
            
      def p_route_bidirectional_mode(self,p):
            ''' route : node BI route'''
            if p[3][0][0][0] == 'node':
                  p[0] = [[p[1],p[3][0][0]], [p[3][0][0]],p[1]] + p[3]
            else:
                  raise ParseError

      def p_node(self,p):
            ''' node : NAME'''
            p[0] = ('node',{'name':p[1]})

      def p_transform_filter(self,p):
            ''' transform : F filter'''
            p[0] = [p[2]]

      def p_transform_edit(self,p):
            ''' transform : E edit '''
            p[0] = [p[2]]
            
      def p_transform_replace(self,p):
            ''' transform : R replace '''
            p[0] = [p[2]]
            
      def p_transform_filter_edit(self,p):
            ''' transform : FE filter edit '''
            p[0] = [p[2],p[3]]
            
      def p_transform_filter_replace(self,p):
            ''' transform : FR filter replace '''
            p[0] = [p[2],p[3]]
            
      def p_filter_name(self,p):
            ''' filter : NAME '''
            p[0] = ('filter',{'name':p[1]})
            
      def p_filter_script(self,p):
            ''' filter : '{' test '}' '''
            p[0] = ('filter',{'test':p[2]})
            
      def p_filter_name_script(self,p):
            ''' filter : NAME '{' test '}' '''
            p[0] = ('filter',{'name':p[1],'test':p[3]})
            
      def p_edit_name(self,p):
            ''' edit : NAME '''
            p[0] = ('edit',{'newname':p[1]})

      def p_edit_name_name (self,p):
            ''' edit : NAME NAME '''
            p[0] = ('edit',{'name':p[1],'newname':p[2]})
            
      def p_edit_script(self,p):
            ''' edit : script '''
            p[0] = ('edit',{'script':p[1]})
            
      def p_edit_name_script(self,p):
            ''' edit : NAME script '''
            p[0] = ('edit',{'name':p[1],'script':p[2]})
            
      def p_edit_name_name_script(self,p):
            ''' edit : NAME NAME script '''
            p[0] = ('edit',{'name':p[1],'newname':p[2],'script':p[3]})

      def p_replace_script(self,p):
            ''' replace : map '''
            p[0] = ('replace',{'newname':p[1]})
            
      def p_replace_name_script(self,p):
            ''' replace : NAME map '''
            p[0] = ('replace',{'newname':p[1],'map':p[2]})

      def p_replace_name_name_script(self,p):
            ''' replace : NAME NAME map '''
            p[0] = ('replace',{'name':p[1],'newname':p[2],'map':p[3]})

      def p_script(self,p):
            ''' script : '{' statements '}' '''
            p[0] = ('script',p[2])
            
      def p_statements_statement(self,p):
            ''' statements : statement '''
            p[0] = [p[1]]
            
      def p_statements(self,p):
            ''' statements : statement ';' statements '''
            p[0] = [p[1]] + p[3]
            
      def p_statement_del(self,p):
            ''' statement : POP var '''
            p[0] = ('pop',p[2])
            
      def p_statement_eq(self,p):
            ''' statement : var '=' expr '''
            p[0] = ('assign',p[1],p[3])
            
      def p_statement_peq(self,p):
            ''' statement : var PE expr '''
            p[0] = ('assign',p[1],('add',('var_value',p[1][1]),p[3]))
                    
      def p_statement_meq(self,p):
            ''' statement : var ME expr '''
            p[0] = ('assign',p[1],('sub',('var_value',p[1][1]),p[3]))
                    
      def p_statement_teq(self,p):
            ''' statement : var TE expr '''
            p[0] = ('assign',p[1],('mul',('var_value',p[1][1]),p[3]))
                    
      def p_statement_deq(self,p):
            ''' statement : var DE expr '''
            p[0] = ('assign',p[1],('div',('var_value',p[1][1]),p[3]))
            
      def p_statement_req(self,p):
            ''' statement : var RE REGEX ',' REGEX '''
            p[0] = ('assign',p[1],('re_sub',('var_value',p[1][1]),p[3],p[5]))
            
      def p_statement_reqf(self,p):
            ''' statement : var RE REGEX ',' REGEX ',' NAME'''
            flags = []
            if p[7] == 'i' or p[7] == 'I': 
                  p[0] = ('assign',p[1],('re_subi',('var_value',p[1][1]),p[3],p[5]))
            else:
                  raise ParseError                  
                    
      def p_expr_test(self,p):
            ''' expr : test '''
            p[0] = ('test',p[1])

      def p_expr_ternary(self,p):
            ''' expr : test '?' expr ':' expr'''
            p[0] = ('ternary',p[1],p[3],p[5])
            
      def p_expr_add(self,p):
            ''' expr : expr '+' expr '''
            p[0] = ('add',p[1],p[3])
            
      def p_expr_sub(self,p):
            ''' expr : expr '-' expr '''
            p[0] = ('sub',p[1],p[3])
            
      def p_expr_mul(self,p):
            ''' expr : expr '*' expr '''
            p[0] = ('mul',p[1],p[3])
            
      def p_expr_div(self,p):
            ''' expr : expr '/' expr '''
            p[0] = ('div',p[1],p[3])
            
      def p_expr_exp(self,p):
            ''' expr : expr EXP expr '''
            p[0] = ('exp',p[1],p[3])
            
      def p_expr_neg(self,p):
            ''' expr : '-' expr %prec UMINUS'''
            p[0] = ('neg',p[2])
            
      def p_expr_list(self,p):
            ''' expr : list '''
            p[0] = p[1]

      def p_list(self,p):
            ''' list : '[' expr_list ']' '''
            p[0] = ('list',p[2])

      def p_expr_map(self,p):
            ''' expr : map '''
            p[0] = p[1]
            
      def p_expr_paren(self,p):
            ''' expr : '(' expr ')' '''
            p[0] = p[2]
            
      def p_expr_value(self,p):
            ''' expr : value '''
            p[0] = p[1]
            
      def p_expr_name(self,p):
            ''' expr : NAME '''
            p[0] = ('var_value',p[1])
            
      def p_expr_ID(self,p):
            ''' expr : ID '''
            p[0] = ('var_value',p[1])
            
      def p_test_eq(self,p):
            ''' test : expr EQ expr '''
            p[0] = ('eq',p[1],p[3])
            
      def p_test_like(self,p):
            ''' test : expr '~' REGEX '''
            p[0] = ('like',p[1],re.compile(p[3]))
            
      def p_test_ge(self,p):
            ''' test : expr GE expr '''
            p[0] = ('ge',p[1],p[3])
            
      def p_test_gt(self,p):
            ''' test : expr '>' expr '''
            p[0] = ('gt',p[1],p[3])
            
      def p_test_lt(self,p):
            ''' test : expr '<' expr '''
            p[0] = ('lt',p[1],p[3])
            
      def p_test_le(self,p):
            ''' test : expr LE expr '''
            p[0] = ('le',p[1],p[3])
            
      def p_test_paren(self,p):
            ''' test : '(' test ')' '''
            p[0] = p[2]
            
      def p_test_and(self,p):
            ''' test : test AND test '''
            p[0] = ('and',p[1],p[3])
            
      def p_test_not(self,p):
            ''' test : '!' test '''
            p[0] = ('not',p[2])
            
      def p_test_or(self,p):
            ''' test : test OR test '''
            p[0] = ('or',p[1],p[3])
            
      def p_expr_list_list(self,p):
            ''' expr_list : expr ',' expr_list '''
            p[0] = [p[1]]+p[3]
            
      def p_expr_list_expr(self,p):
            ''' expr_list : expr '''
            p[0] = [p[1]]

      def p_map(self,p):
            ''' map : '{' mappings '}' '''
            p[0] = ('map',p[2])
      
      def p_mappings_more(self,p):
            ''' mappings : NAME ':' expr ',' mappings '''
            p[0] = [{'key':p[1],'value':p[3]}] + p[5]
            
      def p_mappings_end(self,p):
            ''' mappings : NAME ':' expr '''
            p[0] = [{'key':p[1],'value':p[3]}]
            
      def p_var(self,p):
            ''' var : NAME '''
            p[0] = ('var',p[1])
      
      def p_value_float(self,p):
            ''' value : FLOAT '''
            p[0] = ('value',p[1])
            
      def p_value_int(self,p):
            ''' value : INT '''
            p[0] = ('value',p[1])
            
      def p_value_string(self,p):
            ''' value : STRING '''
            p[0] = ('value',p[1])
            
      def p_value_true(self,p):
            ''' value : TRUE'''
            p[0] = ('value',True)
            
      def p_value_false(self,p):
            ''' value : FALSE'''
            p[0] = ('value',False)
            
class transform:
      def __init__(self, ast):
            self.ast = ast
            self.node_types = ['add','and','div','edit','eq','filter','ge','gt','le','like','list','lt','map','mul','exp','neg','not','or','replace','re_sub','re_subi','script','sub','test','value','var','var_value','pop','assign']
            self.evals = {}
            for x in self.node_types:
                  self.evals[x] = getattr(self,'eval_'+x)

      def eval_add(self, n, d):
            return self.e(n[1], d) + self.e(n[2], d)

      def eval_and(self, n, d):
            return self.e(n[1], d) and self.e(n[2], d)

      def eval_div(self, n, d):
            return self.e(n[1], d) / self.e(n[2], d)

      def eval_eq(self, n, d):
            return self.e(n[1], d) == self.e(n[2], d)

      def eval_ge(self, n, d):
            return self.e(n[1], d) >= self.e(n[2], d)

      def eval_gt(self, n, d):
            return self.e(n[1], d) > self.e(n[2], d)

      def eval_le(self, n, d):
            return self.e(n[1], d) <= self.e(n[2], d)

      def eval_like(self, n, d):
            return n[2].match(self.e(n[1], d)) != None

      def eval_list(self, n, d):
            return [self.e(x, d) for x in n[1]]

      def eval_lt(self, n, d):
            return self.e(n[1], d) < self.e(n[2], d)

      def eval_map(self, n, d):
            return {x['key']:self.e(x['value'], d) for x in n[1]}

      def eval_mul(self, n, d):
            return self.e(n[1], d) * self.e(n[2], d)

      def eval_neg(self, n, d):
            return -self.e(n[1], d)
      
      def eval_not(self, n, d):
            return not self.e(n[1], d)

      def eval_or(self, n, d):
            return self.e(n[1], d) or self.e(n[2], d)

      def eval_sub(self, n, d):
            return self.e(n[1], d) - self.e(n[2], d)
      
      def eval_exp(self, n, d):
            return self.e(n[1], d) ** self.e(n[2], d)

      def eval_value(self, n, d):
            return n[1]
      
      def eval_re_sub(self, n, d):
            return re.sub(n[2], n[3], self.e(n[1], d))
            
      def eval_re_subi(self, n, d):
            return re.sub(n[2], n[3], self.e(n[1], d),flags=re.I)
      
      def eval_ternary(self, n, d):
            return self.e(n[2], d) if self.e(n[1], d) else self.e(n[3], d)
      
      def eval_edit(self, n, d):
            e = n[1]
            ans = copy.deepcopy(d)
            if 'name' in e and d['_name'] != e['name']: return ans
            if 'script' in e: ans = self.e(e['script'], d)
            if 'newname' in e: ans['_name'] = e['newname']
            return ans
            
      def eval_filter(self, n, d):
            f = n[1]
            if 'name' in f and d['_name'] != f['name']: return False
            if 'test' in f: return self.e(f['test'], d)
            return True

      def eval_script(self, n, d):
            self.editing = copy.deepcopy(d)
            for x in n[1]:
                  self.e(x, d)
            return self.editing
                  
      def eval_assign(self, n, d):
            var = self.e(n[1], d)
            val = self.e(n[2], d)
            tmp = self.editing
            for i in range(len(var)):
                  if i == len(var) - 1:
                        tmp[var[i]] = val
                  else:
                        tmp = tmp[var[i]]
                        
      def eval_pop(self, n, d):
            var = self.e(n[1], d)
            tmp = self.editing
            for i in range(len(var)):
                  if i == len(var) - 1:
                        del tmp[var[i]]
                  else:
                        tmp = tmp[var[i]]

      def eval_replace(self, n, d):
            r = n[1]
            if 'name' in r and d['_name'] != r['name']: return d
            ans = self.e(r['map'], d)
            if 'newname' in r: ans['_name'] = r['newname']
            else: ans = {'_name':d['_name']}
            return ans

      def eval_test(self, n, d):
            return self.e(n[1], d)

      def eval_var(self, n, d):
            return n[1].split(".")
      
      def eval_var_value(self, n, d):
            var = n[1].split(".")
            tmp = d
            for i in range(len(var)):
                  if i == len(var) - 1:
                        return tmp[var[i]]
                  else:
                        tmp = tmp[var[i]]
            return None
      
      def e(self, n, d):
            #print('NODE: ',n)
            ans = self.evals[n[0]](n, d)
            #print('ANS: ',ans,n)
            return ans
      
      def evaluate(self, d):
            return self.e(self.ast, d)

tp = transform_parser()
routes = tp.parse(sys.argv[1])
ok = True
for route in routes:
      if ok:
            for x in route:
                  if x[0] in ['edit','replace','filter']:
                        t = transform(x)
                        ok = False
                        break
      else: break


ans = t.evaluate({'_name': 'hello', 'x':2, 'y':'asda', 'z':[1,2,3,4], 'w':{'a':'blah'}})
print(ans)
