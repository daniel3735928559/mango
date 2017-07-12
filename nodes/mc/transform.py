import sys,re,json,copy
            
class Transform:
      def __init__(self, ast):
            self.ast = ast
            self.env = {}
            self.node_types = ['add','and','div','edit','eq','filter','ge','gt','le','like','list','lt','map','mul','exp','neg','not','or','replace','re_sub','re_subi','script','sub','ternary','test','value','var','var_value','pop','assign']
            self.operand_types = ['list','num','bytes','map','bool']
            self.evals = {}
            self.checks = {}
            self.operation_args_types = {
                  'add':[('list','list'), ('num','num'), ('bytes','bytes')],
                  'sub':[('num','num')],
                  'and':[('bool','bool')]
            }
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
            ans = self.evals[n[0]](n, d)
            #print('ANS: ',ans,n)
            return ans
      
      def evaluate(self, env, d):
            self.env = env
            return self.e(self.ast, d)
      
if __name__ == "__main__":
      import transform_parser
      tp = transform_parser.transform_parser()
      routes = tp.parse(sys.argv[1])
      ok = True
      for route in routes:
            for i in range(len(route)):
                  if route[i][0] in ['edit','replace','filter']:
                        route[i] = (route[i][0],Transform(route[i]))

      name = "hello"
      data = {'x':2, 'y':'asda', 'z':[1,2,3,4], 'w':{'a':'blah'}}
      for t in routes[0][1:-1]:
            if t[0] == 'filter':
                  if t[1].evaluate(name, data):
                        continue
                  else:
                        print("FILTERED")
            else:
                  data = t[1].evaluate(name, data)
                  name = t[1].env.get('name','')
      print(name,data)
