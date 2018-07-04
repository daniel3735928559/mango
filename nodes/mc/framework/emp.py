class EMP:
    def __init__(self, name, path):
        self.name = name
        self.path = path
    def get_id(self):
        return self.name
    def __repr__(self):
        return "{}:{}".format(self.name,self.path)
