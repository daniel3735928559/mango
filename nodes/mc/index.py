# We want to store objects based on various ways of indexing their
# properties and to have all the stored information remain consistent
# throughout its lifetime

# Current limitations:
# - objects must not have "uuid" property

class multiindex:
      def __init__(self, indices):
            """ indices = {group_name: {index1_name: [index1 property, subindex1 property, ...], index2_name: [index2 property, ...], ...}}"""
            self.indices = indices
            self.multiindex = {x:{n:{} for n in self.indices[x]} for x in self.indices}
            
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

      def exists(self, group_name, obj):
            try:
                  return len(self.idx_find(group_name, list(self.indices[group_name].keys())[0], obj)) > 0
            except Exception as e:
                  return False
            
      def remove(self, group_name, obj):
            for index_name in self.indices[group_name]:
                  index = self.idx_find(group_name, index_name, obj)
                  del index[obj]

      def add(self, group_name, obj):
            for index_name in self.indices[group_name]:
                  self.idx_add(group_name, index_name, obj)

      def query(self, group_name, index_name, prop_vals):
            """ return a list of all objects in the subtree given by prop_vals """
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
      print(m.exists("all",{}))
      print(m.query("all","by_prop1_2",["qwe","asd"]))
      print(m.query("all","by_prop2",["sdf"]))
      print(m.query("all","by_prop1_2",["ert"]))
