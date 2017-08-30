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
from route import *
from transform_parser import *
import pijemont.doc
from index import multiindex

class mc(m_node):
    
    def __init__(self):
        super().__init__(debug=True)
        self.transform_parser = transform_parser()
        
        # Now read xml file to get list of nodes and shell commands to spin them up.  
        # for child in etree.parse("nodes.xml").getroot():
        #     self.node_types[child.get("name")] = NodeType(child.get("name"),child.get("wd"),child.get("runner"))

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
        self.local_gateway = m_ZMQ_transport(self.mc_addr,self.context,self.poller,True)
        s = self.local_gateway.socket
        self.dataflow = mc_router_dataflow(self.local_gateway,self.serialiser,self.mc_recv)
        self.dataflows[s] = self.dataflow
        self.poller.register(s,zmq.POLLIN)
        self.index = multiindex({
            "routes":{
                "src":["src","dst"]
                "dst":["dst","src"]
                "group_src":["group","src"]
                "group_dst":["group","dst"]
            },
            "nodes":{
                "type":["node_type"]
                "name":["name","type"]
                "group":["group"]
                "group_name":["group","name"]
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
        index.add("nodes",Node("mc","system",0,mc_loopback_dataflow(self.interface,self.mc_dispatch,self.dataflow),bytes("mc","ASCII"),mc_if(self.interface.interface),"root"))
        index.add("groups",Group("system"))

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
        self.hb_gateway = m_ZMQ_transport(self.mc_hb_addr,self.context,self.poller,True)
        h = self.hb_gateway.socket
        self.dataflows[h] = mc_heartbeat_dataflow(self.hb_gateway,self.do_heartbeat)
        self.poller.register(h,zmq.POLLIN)

        # Start reaper thread
        self.reap_time = 150
        self.reaper_thread = threading.Thread(target=mc_heartbeat_worker, daemon=True, args=("mc",self.mc_hb_addr,self.reap_time,threading.Event(),self.context))
        self.reaper_thread.start()
        self.too_old = 200
        with open("./init.emp","r") as f:
            self.mp(f.read(), "system")
        self.run()

    ## Internal functions
        
    def do_heartbeat(self,node_name):
        if node_name == b"mc":
            self.debug_print("REAP",node_name)
            to_reap = []
            now = time.time()
            for n in self.nodes:
                if self.nodes[n].node_id != "mc" and (now - self.nodes[n].last_alive_time) > self.too_old and (now - self.nodes[n].last_heartbeat_time) < self.too_old:
                    to_reap.append(n)
            for n in to_reap:
                self.delete_node(n)
        else:
            self.debug_print("HB",node_name)
            header = self.make_header("heartbeat")
            self.nodes[node_name.decode()].last_heartbeat_time = time.time()
            self.dataflow.send(header,{},self.nodes[node_name.decode()].route)
            
    def mc_error(self,header,args):
        self.deubg_print("ERR",header,args)

    def mc_dispatch(self,header,args,route):
        self.debug_print("MC DISPATCH",header,args)
        result = self.interface.get_function(header['name'])(header,args)
        self.debug_print("RES",result)
        if not result is None:
            name,resp_data = result
            resp_header = self.make_header(name)
            resp_header['src_node'] = 'mc'
            self.dataflow.send(resp_header,resp_data,route)
        
    def remote_recv(self,h,c,raw,route,dataflow):
        pass
    
    def mc_recv(self,h,c,raw,route,dataflow):
        self.debug_print("MC got",h,c,raw,route)
        nodes = self.index.query("nodes","routes",[route])
        if len(nodes) == 1 
            src_node = nodes[0]
            if h['type'] == 'system':
                h['src_node'] = src_node.node_id
                self.nodes["mc"].handle(h, c, route)
            else:
                h['src_node'] = src_node.node_id
                print("passing along",h,c,raw,src_node)
                routes = self.index.query("routes","src",[src_node.node_id])
                for r in routes:
                    print("SENDING ON",str(r))
                    self.routes[r].send(raw,h,c)
                #src_node.emit(raw,h,c)
        else:
            print("Received message from non-existent or ambiguous node",h,c,raw,route)

    ## Helper functions

    def initialise_types(self):
        manifest_dir = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../')
        manifest_path = manifest_dir+'/manifest.json'
        with open(manifest_path) as f:
            manifest = json.loads(f.read())
        self.node_types = manifest['nodes']
        self.langs = manifest['langs']
        for ext in manifest['extensions']:
            manifest_path = manifest_dir+'/'+ext+'/manifest.json'
            with open(manifest_path) as f:
                manifest = json.loads(f.read())
            for n in manifest.get('nodes',{}):
                self.node_types[n] = manifest['nodes'][n]
                self.node_types[n]['dir'] = ext + '/' + self.node_types[n].get('dir',n)
            self.langs.update(manifest.get('langs',{}))

    ### Remove objects
    
    def delete_node(self, node):
        # Remove the node
        self.index.remove("nodes",node)
        
        # Stop heartbeat thread
        if not node.hb_stopper is None:
            node.hb_stopper.set()

        # Remove all routes from or to this node
        routes = self.index.query("routes","dst",[node]) + self.index.query("routes","src",[node])
        for r in routes:
            self.index.remove("routes",r)
    
    def delete_route(self, route):
        self.index.remove("routes",r)

    def delete_group(self, group):
        # Delete all nodes
        nodes = self.index.query("nodes","group",[group_name])
        for n in nodes:
            self.delete_node(n)

        # Delete all routes
        routes = self.index.query("routes","group",[group_name])
        for r in routes:
            self.delete_route(r)

        # Delete the group
        self.index.remove("groups",group)

    def gen_id(self, group, ID_req):
        if(not (ID_req in self.nodes.keys())):
            return ID_req
        for i in range(0,255):
            if(not (ID_req+"."+str(i) in self.nodes.keys())):
                return ID_req+"."+str(i)
        return -1
        
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
    def get_node(self, name, group="system"):
            if "/" in name: group,name = name.split("/")
            ans = self.groups.get(group,{}).get(name,None)
            if ans is None:
               raise Exception("No such node: {}/{}".format(group, name))
            return ans
    
    def add_node(self, name, group, route, dataflow, iface):
        n = Node(name, group, self.gen_key(), self.dataflow, route, iface)
        self.index.add("nodes",n)
        # Start heartbeating thread
        n.hb_stopper = threading.Event()
        n.hb_thread = threading.Thread(target=mc_heartbeat_worker, daemon=True, args=(new_id,self.mc_hb_addr,self.hb_time,n.hb_stopper,self.context))
        n.hb_thread.start()
    
    def add_group(self, name):
        self.index.add("groups",Group(name))

    def create_routes(self, route_spec, group):
        new_routes = []
        for r in self.transform_parser.parse(route_spec):
            src = self.get_node(r[0][1])
            dest = self.get_node(r[-1][1])
            new_route = {"src":src, "dest":dest, "route":Route(src, dest, r[1:-1])}
            
        return new_routes

    def emp(self, header, args):
        return self.mp(args['program'], args['group'])
    
    def mp(self, program, group):
        prog_lines = program.split("\n")
        prog = {}
        var = {}
        nodes = []
        routes = []
        section = None

        for l in prog_lines:
            l = l.strip()
            m = re.match(r'^\[[a-zA-Z_][a-zA-Z0-9_]*\]$', l)
            if m:
                section = l[1:-1]
                prog[section] = []
            elif section is None:
                print("Line without section header!")
            else:
                prog[section] += [l]

        # Config lines are of the form:
        # name = value

        # Current config settings are:
        # group: The name of the group under which these nodes will be added (will be created if doesn't exist)
        
        for l in prog.get('config',[]):
            l = l.strip()
            if l == "":
                continue
            ll = l.split('=')
            var[ll[0].strip()] = ll[1].strip()
            
        # Nodes are of the form:
        # [group/]node --arg1=value1 --arg2=value2 ...
        for l in prog.get('nodes',[]):
            l = l.strip()
            if l == "": continue
            ll = shlex.split(l)
            print(ll)
            new_node = {"group": group, "node": ll[1], "args":[], "type": ll[0]}
            for a in ll[2:]:
                m = re.match(a, '--([a-zA-Z_][a-zA-Z_0-9]*)=(.*)')
                if m:
                    new_nodes['args'][m.group(1)] = m.group(2)
                else:
                    print("Malformed argument: {}".format(a))
            if not group in self.groups:
                return {"success":False, "message":"No such group: {}".format(new_node['group'])}

        # Routes are lines following the mc routing spec
        new_routes = []
        for l in prog.get('routes',[]):
            new_routes += self.create_routes(l.strip())

        # Having got to here with no exceptions raised, everything is
        # validated and we can start actually creating nodes and
        # routes:

        self.add_group(var['group'])
            
        for n in new_nodes:
            self.launch_node(n['type'], n['name'], n['group'], n['args'])
            
        for r in new_route:
            r['src'].add_route(r['route'])

    
    def launch_node(self, node_type, node_id, node_group, env):
        n = node_type
        nid = node_id
        base_path = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../')
        lib_base_path = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../../libmango')
        if n in self.node_types:
            lang = self.node_types[n]['lang']
            node_base = os.path.join(base_path,n if not 'dir' in self.node_types[n] else self.node_types[n]['dir'])
            lib_path = os.path.join(lib_base_path,lang)
            nenv = {'MC_ADDR':self.mc_target,'MANGO_ID':nid}
            if 'pathvar' in self.langs.get(lang, {}):
                nenv[self.langs[lang]['pathvar']] = lib_path
            nenv.update(env)
            self.add_node(node_id, node_group, self.gen_route_key(), self.router_dataflow, self.node_types[n['name']]['if'])
            if lang == "mu":
                mu_path = os.path.join(base_path, 'mu/mu.py')
                print(mu_path)
                lib_path = os.path.join(lib_base_path,"py")
                nenv.update({
                    "MU_WS_PORT":self.node_types[n]['ws_port'],
                    "MU_HTTP_PORT":self.node_types[n]['http_port'],
                    "MU_IF":os.path.join(node_base, self.node_types[n]['if']),
                    "MU_ROOT_DIR":node_base
                })
                nenv[self.langs["py"]['pathvar']] = lib_path
                print("nb",node_base,nenv,mu_path,shlex.split("python " + mu_path))
                subprocess.Popen(shlex.split("python " + mu_path), cwd=node_base, env=nenv)
            
                print("E",nenv,"B",node_base)
                subprocess.Popen(shlex.split(self.node_types[n]['run']), cwd=node_base, env=nenv)
            return "success",{}
        return "error",{'message':'no such node type: {}'.format(n)}

    def success(self, successful=True, message=""):
        ans = {"success":successful}
        if message != "":
            ans['message'] = message
        return "success",ans

    def find_nodes(self, name, group, node_type):
        return [
            {"name":x.name, "group": x.group, "type":x.node_type}
            for x in self.nodes if
            re.match(x.name,name) and
            re.match(x.group,group) and
            re.match(x.node_type,node_type)
        ]
    
    def find_routes(self, src, dst, group):
        return [
            {"src":x.src, "dst": x.dst, "group": x.group, "str":x.str}
            for x in self.routes if
            re.match(x.src,src) and
            re.match(x.dst,dst) and
            re.match(x.group,group)
        ]
    
    ## API functions

    def alive(self,header,args):
        if header['src_node'] in self.nodes:
            self.nodes[header['src_node']].last_alive_time = time.time()
        self.debug_print("ALIVE",header,args)

    def addgroup(self,header,args):
        if self.add_group(args['name']):
            return self.success()
        return self.success(False, "group already exists")

    def delgroup(self,header,args):
        pass

    def addnode(self,header,args):
        return self.launch_node(args['node'], args.get('name', args['node']), args['group'], json.loads(args.get('env',"{}")))

    def query(self,header,args):
        ans = {}
        if 'type' in args:
            ans["types"] = [x for x in self.node_types if re.match(x,args['type']['name'])]
        if 'group' in args:
            ans["groups"] = [x for x in self.groups if re.match(x,args['type']['name'])]
        if 'node' in args:
            ans["nodes"] = self.find_nodes(args['node']['name'],args['node']['group'],args['node']['type'])
        if 'route' in args:
            ans["routes"] = self.find_routes(args['route']['src'],args['route']['dst'],args['route']['group'])
        return "info",ans

    def delnode(self,header,args):
        # Delete the node and all routes relating to it
        node_name = args['node']
        if node_name in self.nodes:
            self.delete_node(node_name)
            return self.success()
        else:
            return self.success(False,"no such node")

    def addroute(self,header,args):
        print("BUILDING ROUTE",header,args)
        for r in self.transform_parser.parse(args['rt']):
            src = self.nodes[r[0][1]]
            dest = self.nodes[r[-1][1]]
            new_route = Route(src, dest, r[1:-1])
        self.nodes[sn].add_route(new_route)
        return self.success()
    
    def delroute(self,header,args):
        if not args['src_node'] in self.nodes:
            print(1)
            return "error",{"message":"Source node not found"}
        sn = self.nodes[args['src_node']]
        if not args['dest_node'] in self.nodes:
            print(3)
            return "error",{"message":"Target node not found"}
        dn = self.nodes[args['dest_node']]
        if sn.del_route_to(dn):
            return "success",{}
        return self.success(False,"Failed to delete route")
    
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
