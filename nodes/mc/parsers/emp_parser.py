import re, shlex

class emp_parser:
    def __init__(self):
        pass

    # returns: 
    def parse(self, program, args=None):
        ans = {'nodes':[],'routes':[]}
        prog = {}
        var = {}
        new_nodes = []
        new_routes = []
        section = None
        prog_lines = program.split("\n")
        for l in prog_lines:
            l = l.strip()
            if len(l) == 0 or l[0] == "#": continue
            m = re.match(r'^\[[a-zA-Z_][a-zA-Z0-9_]*\]$', l)
            if m:
                section = l[1:-1]
                prog[section] = []
            elif section is None:
                raise Exception("Line without section header!")
            elif l != "":
                prog[section] += [l]
        
        print(prog)
        # Config lines are of the form:
        # name = value
        
        # Current config settings are:
        # group: The name of the group under which these nodes will be added (will be created if doesn't exist)
        
        for l in prog.get('config',[]):
            ll = l.split('=')
            var[ll[0].strip()] = ll[1].strip()
            
        # Nodes are of the form:
        # [group/]node --arg1=value1 --arg2=value2 ...
        for l in prog.get('nodes',[]):
            ll = shlex.split(l)
            new_node = {"node": ll[1], "args":{}, "type": ll[0]}
            for a in ll[2:]:
                m = re.match('--([a-zA-Z_][a-zA-Z_0-9]*)=(.*)',a)
                if m:
                    new_node['args'][m.group(1)] = m.group(2)
                else:
                    self.debug_print("Malformed argument: {}".format(a))

            ans['nodes'].append(new_node)
        
        # Routes are lines following the mc routing spec
        for l in prog.get('routes',[]):
            ans['routes'].append(l.strip())
        
        return ans
