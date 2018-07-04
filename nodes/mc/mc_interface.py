from lxml import etree
from error import *
import json
import pijemont.verifier

class mc_if:
    def __init__(self,iface):
        self.interface = iface

    def __repr__(self):
        print("repr",self.interface)
        return json.dumps(self.interface)
        
    def validate(self, function_name, args):
        """
        Returns: modified_input, list_of_errors
        where:
        - modified_input is the input populated with default values where applicable
        - success is a boolean true if there were no problems and false otherwise
        - list_of_errors is as in verify_helper
        """
        return args
        if(function_name in self.interface):
            print("MC VALIDATING",function_name,"AGAINST",self.interface[function_name])
            if not 'args' in self.interface[function_name]:
                return args
            args, messages = pijemont.verifier.verify_helper("", args, {'type':'dict','values':self.interface[function_name]['args']})
        
            if len(messages)>0:
                print("MC NOPE",args,messages)
                raise m_error(m_error.VALIDATION_ERROR,"\n".join(['{}: {}'.format(m['name'], m['message']) for m in messages]))
            else:
                print("MC YEP",args)
                return args
        else:
            raise m_error(m_error.VALIDATION_ERROR,"Unknown function")

    def compare_dict_keys(self, d1, d2):
        """
        Returns [things in d1 not in d2, things in d2 not in d1]
        """
        return [k for k in d1 if not k in d2], [k for k in d2 if not k in d1]
