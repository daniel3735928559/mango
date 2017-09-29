import sys,re,json,copy
            
class Query:
      def __init__(self, ast):
            self.ast = ast
            self.node_types = ['test','eq','ne','var','and','or','not','xor','like']
            self.evals = {}
            self.strs = {}
            for x in self.node_types:
                  self.evals[x] = getattr(self,'eval_'+x)
                  self.strs[x] = getattr(self,'str_'+x)

      def __repr__(self):
            return self.s(self.ast)
                  
      def eval_test(self, n, d):
            return {x:d[x] for x in d if self.e(n[1], d[x])}

      def eval_and(self, n, d):
            return {x:d[x] for x in d if x in self.e(n[1], d) and x in self.e(n[2], d)}
      
      def eval_not(self, n, d):
            return {x:d[x] for x in d if not x in self.e(n[1], d)}

      def eval_or(self, n, d):
            return {x:d[x] for x in d if x in self.e(n[1], d) or x in self.e(n[2], d)}

      def eval_xor(self, n, d):            
            return {x:d[x] for x in d if x in self.e(n[1], d) or x in self.e(n[2], d) and not (x in self.e(n[1], d) and x in self.e(n[2], d))}
      
      def eval_eq(self, n, d):
            return self.e(n[1], d) == n[2]

      def eval_ne(self, n, d):
            return self.e(n[1], d) != n[2]

      def eval_like(self, n, d):
            return re.match(n[2], self.e(n[1], d)) != None
      
      def eval_var(self, n, d):
            obj = d
            for x in n[1]:
                  obj = obj[x]
            return obj
            return 
            
      def e(self, n, d):
            #print('AST NODE: ',n)
            #print('CANDIDATE NODES: ',d)
            ans = self.evals[n[0]](n, d)
            #print('ANS: ',ans,n)
            return ans

      def evaluate(self, d):
            return self.e(self.ast, d)

      def s(self, n):
            #print("S",n)
            return self.strs[n[0]](n)

      def str_test(self, n):
            return self.s(n[1])

      def str_and(self, n):
            return "{} and {}".format(self.s(n[1]), self.s(n[2]))
      
      def str_not(self, n):
            return "not {}".format(self.s(n[1]))

      def str_or(self, n):
            return "{} or {}".format(self.s(n[1]), self.s(n[2]))

      def str_xor(self, n):            
            return "{} xor {}".format(self.s(n[1]), self.s(n[2]))
      
      def str_eq(self, n):
            return "{} == {}".format(self.s(n[1]), n[2])

      def str_ne(self, n):
            return "{} != {}".format(self.s(n[1]), n[2])

      def str_like(self, n):
            return "{} ~ /{}/".format(self.s(n[1]), n[2])
      
      def str_var(self, n):
            return ".".join(n[1])

if __name__ == "__main__":
      import query_parser
      fake_nodes = {
            "n1":{"A":"1","B":{"C":"4"}},
            "n2":{"A":"1","B":{"C":"5"}},
            "n3":{"A":"2","B":{"C":"6"}},
            "n4":{"A":"3","B":{"C":"4"}}
      }
      qp = query_parser.query_parser()
      query = qp.parse(sys.argv[1])
      q = Query(query)
      print(q)
      print(q.evaluate(fake_nodes))
