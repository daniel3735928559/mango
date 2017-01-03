from lxml import etree
import yaml
from error import *

class m_if:
    def __init__(self, default_handler=None):
        self.interface = {}
        self.default_handler = default_handler
        self.loaders = {"yaml":m_YAML_if()}

    def canonical_name(self, function_name, namespace=None):
        print(self.interface)
        if namespace is None:
            ans = []
            for n in self.interface:
                if function_name in self.interface[n]["inputs"]:
                    return "{}.{}".format(n, function_name)
            return None
        else:
            return "{}.{}".format(namespace, function_name)
        
    def add_interface(self,if_file,handlers,if_type=None):
        if if_type is None:
            if_type = if_file.rsplit(".",1)[1]
        iface = self.loaders[if_type].load(if_file)
        name = iface["name"]
        if name in self.interface:
            raise m_error(m_error.INVALID_INTERFACE, "Namespace already in use: " + iface["name"])
        missing,extra = self.compare_dict_keys(iface["inputs"],handlers)
        if len(missing) > 0:
            raise m_error(m_error.INVALID_INTERFACE, "Functions not implemented: " + ", ".join(missing))
                    
        self.interface[name] = {"inputs":{},"outputs":{},"handlers":{}}
        for f in iface.get("inputs",[]):
            self.interface[name]["handlers"][f] = handlers[f]
            if not iface["inputs"][f] is None:
                self.interface[name]["inputs"][f] = iface["inputs"][f]
        for f in iface.get("outputs",[]):
            if not iface["outputs"][f] is None:
                self.interface[name]["outputs"][f] = iface["outputs"][f]

    def get_spec(self):
        return {name:{c:self.interface[name][c] for c in self.interface[name] if c != 'handlers'} for name in self.interface}
                
    def get_function(self, function_name):
        if not "." in function_name:
            function_name = self.canonical_name(function_name)
            if function_name is None:
                return self.default_handler
            
        name,fn = function_name.rsplit(".",1)
        print("NF",name,fn)
        if name in self.interface and fn in self.interface[name]["handlers"]:
            return self.interface[name]["handlers"][fn]
            
        raise m_error(m_error.VALIDATION_ERROR,"Unknown function")
    
    def validate(self, function_name, args):
        self.get_function(function_name)

    def compare_dict_keys(self, d1, d2):
        """Returns [things in d1 not in d2, things in d2 not in d1]"""
        return [k for k in d1 if not k in d2], [k for k in d2 if not k in d1]
            
class m_YAML_if():
    def load(self,if_file):
        with open(if_file,'r') as f:
            ans = yaml.load(f)
        return ans
