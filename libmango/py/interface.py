from lxml import etree
import yaml

class m_if:
    def __init__(self):
        self.interface = {}
        self.loaders = {"xml":m_XML_if(),
                        "yaml":m_YAML_if()}
        
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
            self.interface[f] = {'handler':handlers[f],'args':iface[f]}
        
    def validate(self, function_name, args):
        """
        Returns: modified_input, list_of_errors
        where:
        - modified_input is the input populated with default values where applicable
        - success is a boolean true if there were no problems and false otherwise
        - list_of_errors is as in verify_helper
        """

        if(function_name in self.interface):
            args, messages = self.verify_helper("", args, {'type':'dict','values':self.interface[function_name]})
        
            if len(messages)>0:
                raise m_error(m_error.VALIDATION_ERROR,"\n".join([m['name']+': ' +m['message'] for m in messages]))
            else:
                return args
        else:
            raise m_error(m_error.VALIDATION_ERROR,"Unknown function")

    def verify_helper(self, name, input_element, reference_dict):
        """
        Returns: modified_input,list_of_errors
        where:
        - modified_input is the input populated with default values
        - list_of_errors is: [{name: name, message: ...}, ...]
        """
        ans = []
        if reference_dict['type'] == 'dict':
            if not isinstance(input_element, (dict)):
                ans += [{"name":name, "message":"invalid dict"}]
            else:
                l1,l2 = self.compare_dict_keys(input_element, reference_dict['values'])
                if len(l1) > 0:
                    ans += [{"name":name, "message":"extra keys in input: " + ",".join(l1)}]
                else:
                    ok = True
                    for k in l2:
                        if 'default' in reference_dict['values'][k]:
                            input_element[k] = reference_dict['values'][k]['default']
                            if reference_dict['values'][k]['type'] == 'num':
                                input_element[k] = float(input_element[k])
                            elif (not 'optional' in reference_dict['values'][k]) or reference_dict['values'][k]['optional'] == False:
                                ans += [{"name":name+'/'+k, "message":"required key is absent"}]
                                ok = False
                    if(ok):
                        for k in input_element:
                            input_element[k], temp_ans = self.verify_helper(name + '/' + k, input_element[k], reference_dict['values'][str(k)])
                            ans += temp_ans

        elif reference_dict['type'] == 'list':
            if not isinstance(input_element, (list)):
                ans += [{"name":name, "message":"invalid list"}]
            else:
                for i in range(len(input_element)):
                    input_element[i],temp_ans = self.verify_helper(name+'/'+str(i), input_element[i], reference_dict['values'])
                    ans += temp_ans

        elif reference_dict['type'] == 'boolean':
            if not isinstance(input_element, (bool)):
                ans += [{"name":name, "message":"invalid boolean"}]

        elif reference_dict['type'] == 'num':
            if not isinstance(input_element, (int, long, float)):
                ans += [{"name":name, "message":"invalid number"}]
                
        elif reference_dict['type'] == 'str' or reference_dict['type'] == 'multiline':
            if not isinstance(input_element, (str, unicode)):
                ans += [{"name":name, "message":"expected a string, got {}".format(type(input_element))}]
            elif 'values' in reference_dict and not input_element in reference_dict['values']:
                ans += [{"name":name, "message":"argument must be one of the specified strings: "+", ".join(reference_dict['values'])}]

        elif reference_dict['type'] == 'oneof':
            count = 0
            for k in reference_dict['values']:
                if k in input_element:
                    count += 1
                    if count > 1:
                        ans += [{"name":name+"/"+k,"message":"More than one argument specified for 'oneof arg: " + name}]
            if count == 0:
                if 'default' in reference_dict:
                    input_element = reference_dict['default']
                else:
                    ans += [{"name":name, "message":"no argument provided for 'oneof' arg"}]

        return input_element,ans

    def compare_dict_keys(self, d1, d2):
        """
        Returns [things in d1 not in d2, things in d2 not in d1]
        """
        return [k for k in d1 if not k in d2], [k for k in d2 if not k in d1]


class m_XML_if():
    def load(self,if_file):
        with open(if_file,'r') as f:
            s = f.read()
        root = etree.fromstring(s);
        doc = s;
        iface = {}
        for c in [child for child in root if child.tag == "cmd"]:
            cname = c.get("name")

            args = {}
            rets = {}

            for a in [child for child in c if child.tag == "arg"]:
                args[a.get("name")] = load_helper(a)
            
            for r in [child for child in c if child.tag == "return"]:
                rets[r.get("name")] = load_helper(r)
            
            iface[cname] = {"args":args,"returns":rets}
        return iface

    def load_helper(n):
        ans = {"type":a.get("type"),
               "doc":a.get("desc","")}
        s = a.get("set")
        if not s is None:
            ans['set'] = s
        ans['val'] = [load_helper(v) for v in n if v.tag == "val"]
        return ans
            
class m_YAML_if():
    def load(self,if_file):
        with open(if_file,'r') as f:
            ans = yaml.load(f)
        return ans
