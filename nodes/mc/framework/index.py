# We want to store objects based on various ways of indexing their
# properties and to have all the stored information remain consistent
# throughout its lifetime

# Current limitations:
# - objects must not have "uuid" property

class multiindex:
      def __init__(self, indices):
            """ indices = {group_name: {index1_name: [index1 property, subindex1 property, ...], index2_name: [index2 property, ...], ...}}"""
            self.indices = indices
            self.properties = {g:set() for g in self.indices}
            for g in self.indices:
                  for i in self.indices[g]:
                        self.properties[g] = self.properties[g].union(set(self.indices[g][i]))
            self.multiindex = {x:{n:{} for n in self.indices[x]} for x in self.indices}
            self.flat = {g:set() for g in self.indices}
            
      # internal methods
      
      def idx_find(self, group_name, index_name, obj):
            props = self.indices[group_name][index_name]
            index = self.multiindex[group_name][index_name]
            for p in props:
                  pval = getattr(obj, p)
                  if not pval in index: index[pval] = {}
                  index = index[pval]
            return index
      
      def idx_add(self, group_name, index_name, obj):
            self.idx_find(group_name, index_name, obj)[obj] = obj

      # external methods

      # props = {"prop1":str, "prop2":str, "prop3":{"subprop1":str, "subprop2":str}}
      def summary(self, group_name, props):
            ans = {}
            def prop_dict(x, props):
                  d = {}
                  for p in props:
                       if props[p] in [str,float,int]: d[p] = getattr(x,p)
                       elif type(props[p]) == dict: d[p] = prop_dict(getattr(x,p), props[p])
                       else: d[p] = props[p](getattr(x,p))
                  # d = {p:getattr(x,p) for p in props if props[p] == str or props[p] == float or props[p] == int}
                  # d.update({p:prop_dict(getattr(x,p), props[p]) for p in props if type(props[p]) == dict})
                  return d
            
            for x in self.flat[group_name]:
                  i = x.get_id()
                  ans[i] = prop_dict(x, props)
                  
            return ans
      
      def exists(self, group_name, obj):
            return obj in self.flat[group_name]
            
      def remove(self, group_name, obj):
            if obj in self.flat[group_name]:
                  self.flat[group_name].remove(obj)
                  for index_name in self.indices[group_name]:
                        index = self.idx_find(group_name, index_name, obj)
                        del index[obj]
                        return True
            return False
      
      def add(self, group_name, obj):
            if obj in self.flat[group_name]:
                  return False
            self.flat[group_name].add(obj)
            for index_name in self.indices[group_name]:
                  self.idx_add(group_name, index_name, obj)
            return True

      def get_all(self, group_name):
            nodes = self.query(group_name)
            return [{p:getattr(x,p) for p in self.properties[group_name]} for x in nodes]
      
      def search(self, group_name, props=[]):
            ans = self.query(group_name)
            for p in props:
                  ans = [x for x in ans if getattr(x,p) == props[p]]
            return ans
                  
      def query(self, group_name, index_name=None, prop_vals=[]):
            """ return a list of all objects in the subtree given by prop_vals """
            if index_name is None:
                  index_name = list(self.indices[group_name].keys())[0]
            index = self.multiindex[group_name][index_name]
            for p in prop_vals:
                  index = index.get(p,{})
            ans = list(index.values())
            for i in range(len(self.indices[group_name][index_name]) - len(prop_vals)):
                  new_ans = []
                  for x in ans: new_ans += list(x.values())
                  ans = new_ans
            return ans

if __name__ == "__main__":
      print("hi")
      class testobj:
            def __init__(self, p1, p2):
                  self.prop1 = p1
                  self.prop2 = p2
            def __repr__(self):
                  return "{}.{}".format(self.prop1, self.prop2)
      t1 = testobj("qwe", "asd")
      t2 = testobj("qwe", "sdf")
      t3 = testobj("ert", "sdf")
      m = multiindex({"all":{"by_prop1_2":["prop1", "prop2"], "by_prop2_1":["prop2","prop1"],"by_prop2":["prop2"]}})
      m.add("all",t1)
      m.add("all",t2)
      m.add("all",t3)
      print(m.multiindex)
      print(m.query("all","by_prop1_2",["qwe","asd"]))
      print(m.exists("all",t2))
      print(m.query("all","by_prop1_2",["qwe"]))
      print(m.query("all","by_prop2",["sdf"]))
      print(m.query("all","by_prop1_2",["ert"]))
      m.remove("all", t2)
      print(m.exists("all",t2))
      print(m.exists("all","qwe.asd"))
      print(m.query("all","by_prop1_2",["qwe","asd"]))
      print(m.query("all","by_prop2",["sdf"]))
      print(m.query("all","by_prop1_2",["ert"]))
