class Merge: 
    def __init__(self, node_id, group, route, send_fn, args):
        self.node_type = 'merge'
        self.mergepoints = {}
        self.mergesets = []
        self.node_id = node_id
        self.name = self.node_id
        self.group = group
        self.route = route
        self.send_fn = send_fn

    def get_id(self):
        return str(self)
    
    def __repr__(self):
        return "{}/{}".format(self.group, self.node_id)

    # a message came in for this node.  Validate it and then pass it
    # on to the node through the dataflow

    def add_mergepoint(self, name):
        self.mergepoints[name] = {}
        
    def add_mergeset(self, args):
        self.mergesets.append(args)
    
    def recv(self, header, args, mergepoint):
        print("MERGE handling",header,args)
        # Store the input in the associated mergepoint
        mid = header['mid']
        self.mergepoints[mergepoint][mid] = args
        # Check all mergesets for completeness
        print("checking for completions",self.mergepoints,self.mergesets)
        for s in self.mergesets:
            if not (False in {mid in self.mergepoints[mp] for mp in s}):
                h = {'name':self.name,'mid':mid}
                a = {mp:self.mergepoints[mp][mid] for mp in s}
                for mp in s:
                    del self.mergepoints[mp][mid]
                self.send_fn(h,a,bytes(self.route,'ascii')) 

class Mergepoint:
    def __init__(self, merge_node, name):
        self.node_id = merge_node.node_id
        self.name = self.node_id
        self.merge_node = merge_node
        self.merge_name = name
        self.group = self.merge_node.group
        self.node_type = self.merge_node.node_type
        
    def get_id(self):
        return str(self.merge_node)
    
    def __repr__(self):
        return "{} {}".format(str(self.merge_node), self.node_id)
    
    def handle(self, header, args, route=None):
        print("MP handling",header,args)
        self.merge_node.recv(header,args,self.merge_name)
