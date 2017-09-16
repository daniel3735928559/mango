import io, re, socket, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from mc_dataflows import *
from mc_transport import *
from mc_interface import *
from mc_workers import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node
from lxml import etree
from node import *
from merge import *
from split import *
from node_type import *
from route import *
from group import *
from transform_parser import *
import pijemont.doc
from index import multiindex
import yaml

class mc(m_node):
    
    def __init__(self):
        super().__init__(debug=True)
        self.transform_parser = transform_parser()
        
        # Now read xml file to get list of nodes and shell commands to spin them up.  
        # for child in etree.parse("nodes.xml").getroot():
        #     self.node_types[child.get("name")] = NodeType(child.get("name"),child.get("wd"),child.get("runner"))

        self.route = self.gen_route_key()
        self.interface.add_interface("mc.yaml",
                                     {
                                         "addgroup":self.addgroup,
                                         "delgroup":self.delgroup,
                                         "addroute":self.addroute,
                                         "delroute":self.delroute,
                                         "query":self.query,
                                         "addnode":self.addnode,
                                         "delnode":self.delnode,
                                         "alive":self.alive,
                                         "doc":self.doc,
                                         "emp":self.emp,
                                         "error":self.mc_error,
                                     })
        self.uuid = str(self.gen_key())

        self.mc_addr = "tcp://*:"+sys.argv[1]
        self.mc_target = "tcp://localhost:"+sys.argv[1]
        self.local_gateway = m_ZMQ_transport(self.mc_addr,self.context,self.poller,None,True)
        s = self.local_gateway.socket
        self.dataflow = mc_router_dataflow(self.local_gateway,self.serialiser,self.mc_recv)
        self.dataflows[s] = self.dataflow
        self.poller.register(s,zmq.POLLIN)
        self.index = multiindex({
            "routes":{
                "src":["src_name","dst_name"],
                "dst":["dst_name","src_name"],
                "group":["group"],
                "group_src":["group","src_name","dst_name"],
                "group_id":["group","route_id"],
                "group_dst":["group","dst_name","src_name"]
            },
            "nodes":{
                "type":["node_type"],
                "name":["node_id","node_type"],
                "group":["group"],
                "group_name":["group","node_id"],
                "route":["route"]
            },
            "groups":{
                "name":["name"]
            },
            "types":{
                "name":["name"]
            }
        })
        self.initialise_types()
        self.index.add("nodes",Node("mc", "system", "mc", 0, mc_loopback_dataflow(self.interface,self.serialiser,self.mc_dispatch,self.mc_recv), self.route, mc_if(self.interface.interface)))
        #self.index.add("groups",Group("system"))

        # Remote listening stuff: 
        if len(sys.argv) > 2:
            self.debug_print("Adding remote port: ",sys.argv[2])
            self.mc_remote_addr = "tcp://*:"+sys.argv[2]
            self.remote_gateway = m_ZMQ_transport(self.mc_remote_addr,self.context,self.poller,True)
            f = self.remote_gateway.socket
            self.dataflows[f] = mc_remote_dataflow(self.remote_gateway,self.serialiser,self.remote_recv)
            self.poller.register(f,zmq.POLLIN)

        # Heartbeat listening stuff:
        self.mc_hb_addr = "inproc://hb"
        self.hb_time = 100
        self.hb_gateway = m_ZMQ_transport(self.mc_hb_addr,self.context,self.poller,None,True)
        h = self.hb_gateway.socket
        self.dataflows[h] = mc_heartbeat_dataflow(self.hb_gateway,self.do_heartbeat)
        self.poller.register(h,zmq.POLLIN)

        # Start reaper thread
        self.reap_time = 150
        self.reaper_thread = threading.Thread(target=mc_heartbeat_worker, daemon=True, args=("mc",self.mc_hb_addr,self.reap_time,threading.Event(),self.context))
        self.reaper_thread.start()
        self.too_old = 200
        with open("./init.emp","r") as f:
            print(self.mp(f.read(), "system"))
        self.run()

    ## Internal functions
        
    def do_heartbeat(self,node):
        node = node.decode()
        self.debug_print("HB",node)
        if node == "mc":
            self.debug_print("REAPING")
            to_reap = []
            now = time.time()
            node_list = self.index.query("nodes")
            for n in node_list:
                if n.node_id != "mc" and (not n.node_type in ['merge','split']) and (now - n.last_alive_time) > self.too_old and (now - n.last_heartbeat_time) < self.too_old:
                    to_reap.append(n)
            for n in to_reap:
                self.delete_node(n)
        else:
            header = self.make_header("heartbeat")
            n = self.index.query("nodes","route",[node])[0]
            n.last_heartbeat_time = time.time()
            self.dataflow.send(header,{},n.route)
            
    def mc_error(self,header,args):
        self.deubg_print("ERR",header,args)

    def mc_system_dispatch(self,header,args,route):
        self.debug_print("MC SYSTEM DISPATCH",header,args, route)
        src_node = self.index.query("nodes","route",[route])[0]
        header['src_node'] = src_node
        self.interface.get_function(header['name'])(header,args)
        
    def mc_dispatch(self,header,args,route):
        self.debug_print("MC DISPATCH",header,args, route)
        src_node = self.index.query("nodes","route",[route])[0]
        header['src_node'] = src_node
        result = self.interface.get_function(header['name'])(header,args)
        self.debug_print("RES",result)
        if not result is None:
            name,resp_data = result
            resp_header = self.make_header(name)
            resp_header['src_node'] = 'mc'
            self.mc_recv(resp_header, resp_data, bytes(self.route,'ascii'))
            #self.dataflow.send(resp_header,resp_data,route)
        
    def remote_recv(self,h,c,raw,route,dataflow):
        pass
    
    def mc_recv(self,h,c,route):
        if not 'mid' in h:
            print("MC adding MID")
            h['mid'] = self.gen_route_key()+"_"+str(time.time())
        self.debug_print("MC got",h,c,route)
        nodes = self.index.query("nodes","route",[route.decode('ascii')])
        print("nodes",route.decode('ascii'),nodes)
        #print(self.index.multiindex)
        if len(nodes) == 1:
            if h.get('type','') == 'system':
                self.mc_system_dispatch(h,c,route.decode('ascii'))
            else:
                src_node = nodes[0]
                self.node_send(src_node, h, c)
        else:
            print("Received message from non-existent or ambiguous node",h,c,route)

    def node_send(self, n, h, c):
        h['src_node'] = n.node_id
        print("passing along",h,c,n)
        routes = self.index.query("routes","src",[str(n)])
        print("RS",routes)
        for r in routes:
            print("SENDING ON",str(r))
            raw = self.serialiser.serialise(h,c)
            r.send(raw,h,c)
        
                
    ## Helper functions

    def load_types(self, base, spec):
        for t in spec:
            iface_path = spec[t]['if']
            with open(os.path.join(base, iface_path),'r') as f:
                spec[t]['if'] = mc_if(yaml.load(f.read()))
            spec[t]['dir'] = base
            self.index.add("types",NodeType(t,base,spec[t].get('run',''),spec[t]['if'],spec[t]['lang'],spec[t]))
            
    def initialise_types(self):
        manifest_dir = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)),'..'))
        manifest_path = manifest_dir+'/config.yaml'
        with open(manifest_path) as f:
            manifest = yaml.load(f.read())
        print("M",manifest)
        dirs = manifest.get('nodes',[])
        self.langs = manifest.get('langs',{})
        for d in dirs:
            print("D",d)
            node_path = os.path.join(manifest_dir,d)
            with open(os.path.join(node_path,'mango.yaml')) as f:
                node_spec = yaml.load(f.read())
                print("S",node_spec)
                self.load_types(node_path, node_spec)
    ### Remove objects
    
    def delete_node(self, node):
        # Send the node an exit message
        node.handle(self.make_header("exit"),{})
        
        # Remove the node
        self.index.remove("nodes",node)
        
        # Stop heartbeat thread
        if not node.hb_stopper is None:
            node.hb_stopper.set()

        # Remove all routes from or to this node
        routes = self.index.query("routes","dst",[str(node)]) + self.index.query("routes","src",[str(node)])
        for r in routes:
            self.index.remove("routes",r)
    
    def delete_route(self, route):
        self.index.remove("routes",route)

    def find_group(self, group_name):
        ans = self.index.query("groups","name",[group_name])
        return ans[0] if len(ans) > 0 else None
        
    def find_node(self, node_spec):
        if "/" in node_spec: 
            ans = self.index.query("nodes","group_name",node_spec.split("/"))
            if len(ans) == 1:
                return ans[0]
            else:
                return None
        else:
            return None
        
    def delete_group(self, group):
        # Delete all nodes
        nodes = self.index.query("nodes","group",[group])
        print("DELETING NODES: ",nodes)
        for n in nodes:
            self.delete_node(n)

        # Delete all routes
        routes = self.index.query("routes","group",[group])
        print("DELETING ROUTES: ",nodes)
        for r in routes:
            self.delete_route(r)

        # Delete the group
        self.index.remove("groups",group)

    def gen_key(self):
        return random.randint(0,2**64-1)

    # def node_add(self,header,args):
    #     # Generate ID based on requested name ID and a key, make a
    #     # Node object based on this, add it to the list, and return
    #     # the ID and key
    #     print("HELLO from")
    #     print(header["source"])
    #     print(args)
    #     n = Node(self.gen_id(args["id_request"]),self.gen_key(),self.ports["stdio"],self.nodes["mc"]) # replace local_gateway with the actual socket that the command came in on
    #     print("New node: "+n.node_id)
    #     self.nodes[n.node_id] = n
    #     return {"node_id":n.node_id,"key":n.key}

    # Finds a node with a given specified name (e.g. group/name)
    # relative to a given group
    def has_node(self, name, group="system"):
        return len(self.index.query("nodes","group_name", [group, name])) > 0
    
    def get_node(self, name, group="system"):
        print('GN',name, group)
        nodes = self.index.query("nodes","group_name", [group, name])
        if len(nodes) == 0:
            raise Exception("No such node: {}/{}".format(group, name))
        return nodes[0]

    def add_merge(self, name, group, args):
        n = Merge(name, group, self.gen_route_key(), self.mc_recv, args)
        self.index.add("nodes",n)
        
    def add_split(self, name, group, args):
        n = Split(name, group, self.gen_route_key(), self.mc_recv, args)
        self.index.add("nodes",n)
    
    def add_node(self, name, group, node_type, route, dataflow, iface):
        n = Node(name, group, node_type, self.gen_key(), dataflow, route, iface)
        self.index.add("nodes",n)
        # Start heartbeating thread
        n.hb_stopper = threading.Event()
        n.hb_thread = threading.Thread(target=mc_heartbeat_worker, daemon=True, args=(route, self.mc_hb_addr,self.hb_time,n.hb_stopper,self.context))
        n.hb_thread.start()
    
    def add_group(self, name):
        self.index.add("groups",Group(name))

    def create_routes(self, route_spec, group):
        print("CR",route_spec)
        grp = self.find_group(group)
        for r in self.transform_parser.parse(route_spec):
            src_name,src_group = r[0][1]['name'],r[0][1].get('group',group)
            print(src_name, src_group)
            if "/" in src_name: src_group,src_name = src_name.split("/")
            src = self.get_node(src_name, src_group)
            
            dst_name,dst_group = r[-1][1]['name'],r[-1][1].get('group',group)
            if "/" in dst_name: dst_group,dst_name = dst_name.split("/")
            dest = self.get_node(dst_name, dst_group)

            print("R",r)
            if src.node_type == 'merge':
                src.add_mergeset(r[0][1]['args'])
            if dest.node_type == 'merge':
                dest.add_mergepoint(r[-1][1]['args'][0])
                dest = Mergepoint(dest, r[-1][1]['args'][0])
            
            print("new route",src,dest)
            
            new_route = {"src":src, "dest":dest, "route":Route(grp.rt_id(), src, [Transform(t) for t in r[1:-1]], dest, group, route_spec)}
            self.index.add("routes",new_route['route'])

    def emp(self, header, args):
        return self.mp(args['mp'], args['group'])
    
    def mp(self, program, group):
        prog_lines = program.split("\n")
        prog = {}
        var = {}
        new_nodes = []
        new_routes = []
        section = None

        if len(self.index.query("groups","name",[group])) > 0:
            return {"success":False, "message":"Group already exists: {}".format(group)}

        # The group did not exist.  Add it: 
        
        try:
            self.add_group(group)
            
            for l in prog_lines:
                l = l.strip()
                if len(l) == 0: continue
                m = re.match(r'^\[[a-zA-Z_][a-zA-Z0-9_]*\]$', l)
                if m:
                    section = l[1:-1]
                    prog[section] = []
                elif section is None:
                    print("Line without section header!")
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
                print(ll)
                new_node = {"group": group, "node": ll[1], "args":{}, "type": ll[0]}
                for a in ll[2:]:
                    m = re.match('--([a-zA-Z_][a-zA-Z_0-9]*)=(.*)',a)
                    if m:
                        new_node['args'][m.group(1)] = m.group(2)
                    else:
                        print("Malformed argument: {}".format(a))

                self.launch_node(new_node['type'], new_node['node'], new_node['group'], new_node['args'])
                #new_nodes.append(new_node)
            
            # Routes are lines following the mc routing spec
            for l in prog.get('routes',[]):
                self.create_routes(l.strip(),group)
            
            print("Starting", group)
            
        except Exception as e:
            # If there was an exception, tear down the whole group we just created:
            traceback.print_exc()
            self.delete_group(group)


    def gen_route_key(self):
        return str(random.randint(0,2**64))
    
    def launch_node(self, node_type, node_id, node_group, env):
        n = node_type
        if n == 'merge':
            self.add_merge(node_id, node_group, env)
            return
        if n == 'split':
            self.add_split(node_id, node_group, env)
            return
        nid = node_id
        base_path = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../')
        lib_base_path = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../../libmango')
        t = self.index.query("types","name",[n])[0]
        print("T",t)
        lang = t.lang
        node_base = t.wd
        lib_path = os.path.join(lib_base_path,lang)
        new_route = self.gen_route_key()
        nenv = {'MC_ADDR':self.mc_target,'MANGO_ID':nid,'MANGO_ROUTE':new_route}
        if 'pathvar' in self.langs.get(lang, {}):
            nenv[self.langs[lang]['pathvar']] = lib_path
        nenv.update(env)
        self.add_node(node_id, node_group, n, new_route, self.dataflow, t.iface)
        if lang == "mu":
            mu_path = os.path.join(base_path, 'mu/mu.py')
            print(mu_path)
            lib_path = os.path.join(lib_base_path,"py")
            nenv[self.langs["py"]['pathvar']] = lib_path
            print("nb",node_base,nenv,mu_path,shlex.split("python " + mu_path))
            subprocess.Popen(shlex.split("python " + mu_path), cwd=node_base, env=nenv)
        else:
            print("E",nenv,"B",node_base,"R",t.run)
            subprocess.Popen(shlex.split(t.run), cwd=node_base, env=nenv)
        return "success",{}
        return "error",{'message':'no such node type: {}'.format(n)}

    def success(self, successful=True, message=""):
        ans = {"success":successful}
        if message != "":
            ans['message'] = message
        return "success",ans

    ## API functions

    def alive(self,header,args):
        header['src_node'].last_alive_time = time.time()
        self.debug_print("ALIVE",header,args)

    def addgroup(self,header,args):
        if self.find_group(args['name']):
            return self.success(False, "group already exists")
        self.add_group(args['name'])
        return self.success()

    def delgroup(self,header,args):
        grp = self.find_group(args['name'])
        if grp:
            self.delete_group(grp)
            return self.success()
        return self.success(False, "No such group")
            

    def addnode(self,header,args):
        return self.launch_node(args['node'], args.get('name', args['node']), args['group'], json.loads(args.get('env',"{}")))

    def query(self,header,args):
        ans = {s:self.index.search(s+'s',args[s]) for s in ['type','group','node','route'] if s+'s' in args}
            
        if len(list(ans.keys())) == 0:
            ns = self.index.search("nodes")
            rs = self.index.search("routes")
            gs = self.index.search("groups")
            ts = self.index.search("types")
            ans = {"nodes":[str(x) for x in ns],
                   "routes":[str(x) for x in rs],
                   "groups":[str(x) for x in gs],
                   "types":[str(x) for x in ts]}
        return "info",ans

    def delnode(self,header,args):
        # Delete the node and all routes relating to it
        node_name = args['node']
        n = self.find_node(node_name)
        if n:
            self.delete_node(n)
            return self.success()
        else:
            return self.success(False,"no such node")

    def addroute(self,header,args):
        print("BUILDING ROUTE",header,args)
        try:
            self.create_routes(args['spec'],args['group'])
            return self.success()
        except Exception:
            return self.success(False,"Failed to create route")
    
    def delroute(self,header,args):
        rs = self.index.search("routes",args)
        if len(rs) == 0:
            return self.success(False, "No such route")
        elif len(rs) > 1:
            return self.success(False, "Multiple routes match: " + str(rs))
        else:
            self.delete_route(rs[0])
            return self.success()
    
    def doc(self,header,args):
        n = args['node']
        
        if n in self.nodes:
            to_doc = self.nodes[n].interface.interface
        else:
            raise m_error(m_error.BAD_ARGUMENT,"Node not found: {}".format(n))

        if 'function' in args:
            f = args['function']
            if f in to_doc:
                to_doc = {f:to_doc[f]}
            else:
                raise m_error(m_error.BAD_ARGUMENT,"Function not found: {}".format(f))
            
        # if 'element' in args:
        #     e = args['element']
        #     for ei in e.split("."):
        #         if ei in to_doc:
        #             to_doc = to_doc[ei]
        #         else:
        #             raise m_error(m_error.BAD_ARGUMENT,"Function not found: {}".format(ei))
        return "doc",{"doc":pijemont.doc.doc_gen(to_doc)}

os.environ["MANGO_ID"] = "mc"
m = mc()
