from lxml import etree
import yaml
from error import *

class m_if:
    def __init__(self):
        self.interface = {}
        self.loaders = {"yaml":m_YAML_if()}
        
    def add_interface(self,if_file,handlers,if_type=None,namespace=None):
        if if_type is None:
            if_type = if_file.rsplit(".",1)[1]
        iface = self.loaders[if_type].load(if_file)
        missing,extra = self.compare_dict_keys(iface,handlers)
        if len(missing) > 0:
            raise m_error(m_error.INVALID_INTERFACE, "Functions not implemented: " + ", ".join(missing))
        if len(extra) > 0:
            raise m_error(m_error.INVALID_INTERFACE, "Functions not in interface: " + ", ".join(extra))
        
        existing = [f for f in iface if f in self.interface]
        if len(existing) > 0:
            raise m_error(m_error.INVALID_INTERFACE, "Functions already defined: " + ", ".join(existing))
                    
        for f in iface:
            if not namespace is None:
                f = namespace + "." + f
            self.interface[f] = {'handler':handlers[f]}
            if not iface[f] is None:
                if 'args' in iface[f]:
                    self.interface[f]['args'] = iface[f]['args']
                if 'rets' in iface[f]:
                    self.interface[f]['rets'] = iface[f]['rets']
        
    def validate(self, function_name, args):
        if(not function_name in self.interface):
            raise m_error(m_error.VALIDATION_ERROR,"Unknown function")

    def compare_dict_keys(self, d1, d2):
        """Returns [things in d1 not in d2, things in d2 not in d1]"""
        return [k for k in d1 if not k in d2], [k for k in d2 if not k in d1]
            
class m_YAML_if():
    def load(self,if_file):
        with open(if_file,'r') as f:
            ans = yaml.load(f)
        return ans
