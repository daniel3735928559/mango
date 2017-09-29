class Group:
    def __init__(self, name):
        self.name = name
        self.route_id = 0
    def rt_id(self):
        self.route_id += 1
        return "rt{}".format(self.route_id)
    def get_id(self):
        return self.name
    def __repr__(self):
        return self.name
