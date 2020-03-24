from lxml import etree
import yaml
from error import *

class m_if:
    def __init__(self, default_handler=None):
        self.interface = {}
        self.spec = {}
        self.default_handler = default_handler

    def canonical_name(self, function_name, namespace=None):
        if namespace is None:
            ans = []
            for n in self.interface:
                if function_name in self.interface[n]["inputs"]:
                    return "{}.{}".format(n, function_name)
            return None
        else:
            return "{}.{}".format(namespace, function_name)
        
    def add_interface(self,handlers):
        for fn in handlers:
            self.interface[fn] = handlers[fn]

    def get_spec(self):
        return {name:{c:self.interface[name][c] for c in self.interface[name] if c != 'handlers'} for name in self.interface}
                
    def get_function(self, function_name):
        if function_name in self.interface:
            return self.interface[function_name]
        elif not self.default_handler is None:
            return self.default_handler
        raise m_error(m_error.VALIDATION_ERROR,"Unknown function")
    
    def validate(self, function_name, args):
        self.get_function(function_name)
