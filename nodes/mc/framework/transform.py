import sys,re,json,copy
            
class Transform:
      def __init__(self, ast):
            self.ast = ast
            #print("AST",ast)
            self.kind = self.ast[0];
            self.env = {}
            self.node_types = ['add','and','div','edit','eq','filter','ge','gt','le','like','list','lt','map','mul','exp','neg','not','or','replace','re_sub','re_subi','script','sub','ternary','test','value','var','var_value','pop','assign']
            self.operand_types = ['list','num','bytes','map','bool']
            self.evals = {}
            self.strs = {}
            self.checks = {}
            self.operation_args_types = {
                  'add':[('list','list'), ('num','num'), ('bytes','bytes')],
                  'sub':[('num','num')],
                  'and':[('bool','bool')]
            }
            for x in self.node_types:
                  self.evals[x] = getattr(self,'eval_'+x)
                  self.strs[x] = getattr(self,'str_'+x)

      def __repr__(self):
            return self.s(self.ast)
                  
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
            if 'name' in e and self.env.get('name','') != e['name']: return ans
            if 'script' in e: ans = self.e(e['script'], d)
            if 'newname' in e: self.env['name'] = e['newname']
            return ans
            
      def eval_filter(self, n, d):
            f = n[1]
            if 'name' in f and self.env.get('name','') != f['name']: return False
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
            if 'name' in r and self.env.get('name','') != r['name']: return d
            ans = self.e(r['map'], d)
            if 'newname' in r: self.env['name'] = r['newname']
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
            #print('DATA: ',d)
            #print('ENV: ',self.env)
            ans = self.evals[n[0]](n, d)
            #print('ANS: ',ans,n)
            return ans

      def evaluate(self, env, d):
            self.env = env
            return self.e(self.ast, d)

      def s(self, n):
            #print("S",n)
            return self.strs[n[0]](n)

      def str_add(self, n):
            return "{} + {}".format(self.s(n[1]), self.s(n[2]))

      def str_and(self, n):
            return "{} && {}".format(self.s(n[1]), self.s(n[2]))

      def str_div(self, n):
            return "{} / {}".format(self.s(n[1]), self.s(n[2]))

      def str_edit(self, n):
            ans,sep = "",""
            if 'newname' in n[1]: ans,sep = "{}".format(n[1]['newname'])," "
            if 'script' in n[1]: ans += "{}{}".format(sep,self.s(n[1]['script']))
            if 'name' in n[1]: ans = "{} {}".format(n[1]['name'],ans)
            ans = "e {}".format(ans)
            return ans
      
      def str_eq(self, n):
            return "{} = {}".format(self.s(n[1]), self.s(n[2]))

      def str_filter(self, n):
            ans,sep = "",""
            if 'name' in n[1]: ans,sep = '{}'.format(n[1]['name']),' and '
            if 'test' in n[1]: ans += "{}{}".format(sep, self.s(n[1]['test']))
            return "f {}".format(ans)

      def str_ge(self, n):
            return "{} >= {}".format(self.s(n[1]), self.s(n[2]))

      def str_gt(self, n):
            return "{} > {}".format(self.s(n[1]), self.s(n[2]))

      def str_le(self, n):
            return "{} <= {}".format(self.s(n[1]), self.s(n[2]))

      def str_like(self, n):
            return "{} ~ {}".format(self.s(n[1]), self.s(n[2]))

      def str_list(self, n):
            ans = ", ".join([self.s(x) for x in n[1]])
            return "[{}]".format(ans)

      def str_lt(self, n):
            return "{} < {}".format(self.s(n[1]), self.s(n[2]))

      def str_map(self, n):
            return "{" + ",".join(["{}:{}".format(x['key'],self.s(x['value'])) for x in n[1]]) + "}"

      def str_mul(self, n):
            return "{} * {}".format(self.s(n[1]), self.s(n[2]))

      def str_exp(self, n):
            return "{}^{}".format(self.s(n[1]), self.s(n[2]))

      def str_neg(self, n):
            return "-{}".format(self.s(n[1]))

      def str_not(self, n):
            return "not({})".format(self.s(n[1]), self.s(n[2]))

      def str_or(self, n):
            return "{} || {}".format(self.s(n[1]), self.s(n[2]))

      def str_replace(self, n):
            ans,sep = "",""
            if 'newname' in n[1]: ans,sep = "{}".format(n[1]['newname']),' '
            if 'map' in n[1]: ans += "{}{}".format(sep, self.s(n[1]['map']))
            if 'name' in n[1]: ans = "{} {}".format(n[1]['name'],ans)
            ans = "r {}".format(ans)
            return ans

      def str_re_sub(self, n):
            return "s/{}/{}/".format(self.s(n[1]), self.s(n[2]))

      def str_re_subi(self, n):
            return "s/{}/{}/i".format(self.s(n[1]), self.s(n[2]))

      def str_script(self, n):
            return ";".join([self.s(x) for x in n[1]])

      def str_sub(self, n):
            return "{} - {}".format(self.s(n[1]), self.s(n[2]))

      def str_ternary(self, n):
            return "{} ? {} : {}".format(self.s(n[1]), self.s(n[2]), self.s(n[3]))

      def str_test(self, n):
            return "test()".format(self.s(n[1]))

      def str_value(self, n):
            return "{}".format(n[1])

      def str_var(self, n):
            return "{}".format(n[1])

      def str_var_value(self, n):
            return "{}".format(n[1])

      def str_pop(self, n):
            return "pop {}".format(self.s(n[1]))

      def str_assign(self, n):
            return "{} = {}".format(self.s(n[1]), self.s(n[2]))


      
if __name__ == "__main__":
      sys.path.append('../parsers')
      import transform_parser
      tp = transform_parser.transform_parser()
      routes = tp.parse(sys.argv[1])
      print("PARSED ROUTES")
      if routes[0] == 'route':
            for x in routes[1]:
                  print("R",x)
      elif routes[0] == 'pipeline':
            print("P",routes[1])
      else:
            print("WAT")
            print(routes)
      ok = True
      routes = routes[1]
      for route in routes:
            for i in range(len(route)):
                  if route[i][0] in ['edit','replace','filter']:
                        route[i] = (route[i][0],Transform(route[i]))

      name = "hello"
      data = {'x':2, 'y':'asda', 'z':[1,2,3,4], 'w':{'a':'blah'}}
      for t in routes[0][1:-1]:
            t = t[1]
            t.env['name'] = name
            print("APPLYING",t)
            #print("STR",t.s(t.ast))
            if t.kind == 'filter':
                  if t.evaluate(t.env, data):
                        continue
                  else:
                        print("FILTERED")
                        break
            else:
                  data = t.evaluate(t.env, data)
                  name = t.env.get('name','')
      print(name,data)
