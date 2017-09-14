class NodeType: 
    def __init__(self,name,wd,run,iface,lang,props):
        self.name = name
        self.wd = wd
        self.run = run
        self.lang = lang
        self.iface = iface
        self.props = props
    def __repr__(self):
        return self.name
