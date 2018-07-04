class Split: 
    def __init__(self, node_id, group, route, send_fn, args):
        self.node_type = 'split'
        self.node_id = node_id
        self.name = self.node_id
        self.group = group
        self.route = route
        self.send_fn = send_fn
        self.hb_stopper = None

    def get_id(self):
        return str(self)
        
    def __repr__(self):
        return "{}/{}".format(self.group, self.node_id)

    def handle(self, header, args):
        if header['name'] != 'exit':
            self.send_fn(header,args,bytes(self.route,'ascii'))
        else:
            print("Split goodbye")
