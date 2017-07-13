import io, re, socket, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from mc_dataflows import *
from mc_transport import *
from mc_interface import *
from mc_workers import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node
from lxml import etree
from obj import *
import pijemont.doc

class mc(m_node):
    
    def __init__(self):
        super().__init__(debug=True)
        self.node_types = {}
        self.collect_nodes()
        self.transform_parser = transform_parser()
        
        # Now read xml file to get list of nodes and shell commands to spin them up.  
        # for child in etree.parse("nodes.xml").getroot():
        #     self.node_types[child.get("name")] = NodeType(child.get("name"),child.get("wd"),child.get("runner"))

        self.interface.add_interface("mc.yaml",
                                     {
                                         "route":self.route_add,
                                         "hello":self.hello,
                                         "alive":self.alive,
                                         "doc":self.doc,
                                         "nodes":self.node_list,
                                         "error":self.mc_error,
                                         "launch":self.launch,
                                         "types":self.list_types,
                                         "delnode":self.delnode,
                                         "delroute":self.delroute,
                                         #"remote":self.remote_connect,
                                         #"delremote":self.remote_disconnect,
                                         "routes":self.route_list
                                     })
        print(self.interface.interface['mc']);
        self.uuid = str(self.gen_key())

        self.mc_addr = "tcp://*:"+sys.argv[1]
        self.mc_target = "tcp://localhost:"+sys.argv[1]
        self.local_gateway = m_ZMQ_transport(self.mc_addr,self.context,self.poller,True)
        s = self.local_gateway.socket
        self.dataflow = mc_router_dataflow(self.local_gateway,self.serialiser,self.mc_recv)
        self.dataflows[s] = self.dataflow
        self.poller.register(s,zmq.POLLIN)
        self.routes_to_nodes = {}
        self.nodes = {"mc":Node("mc",0,mc_loopback_dataflow(self.interface,self.mc_dispatch,self.dataflow),bytes("mc","ASCII"),mc_if(self.interface.interface))}

        #self.ports["stdio"] = self.dataflows[s]

        # Remote listening stuff: 
        if len(sys.argv) > 2:
            print("Adding remote port: ",sys.argv[2])
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
        
        self.run()

    def do_heartbeat(self,node_name):
        if node_name == b"mc":
            print("REAP")
            to_reap = []
            now = time.time()
            for n in self.nodes:
                if self.nodes[n].node_id != "mc" and (now - self.nodes[n].last_alive_time) > self.too_old and (now - self.nodes[n].last_heartbeat_time) < self.too_old:
                    to_reap.append(n)
            for n in to_reap:
                self.delete_node(n)
        else:
            print("HB",node_name)
            header = self.make_header("heartbeat")
            self.nodes[node_name.decode()].last_heartbeat_time = time.time()
            self.dataflow.send(header,{},self.nodes[node_name.decode()].route)

    def alive(self,header,args):
        if header['src_node'] in self.nodes:
            self.nodes[header['src_node']].last_alive_time = time.time()
        print("ALIVE",header,args)


    def delnode(self, header, args):
        node_name = args['node']
        if node_name in self.nodes:
            self.delete_node(node_name)
            return "success",{}
        else:
            return "error",{"message":"no such node"}
        
    def delete_node(self, node_name):
        node = self.nodes[node_name]
        print("REAPING",node.node_id)
        for n in self.nodes:
                self.nodes[n].del_route_to(node_name)
                
        if not node.hb_stopper is None:
            node.hb_stopper.set()
        del self.nodes[node_name]
        for r in self.routes_to_nodes:
            if self.routes_to_nodes[r] == node_name:
                del self.routes_to_nodes[r]
    
    def collect_nodes(self):
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
        
    def hello(self,header,args):
        print("HELLO",header,args)
        return "reg",{'id':args['id']}
        
    def mc_error(self,header,args):
        print("ERR",header,args)

    def mc_dispatch(self,header,args,route):
        print("MC DISPATCH",header,args)
        result = self.interface.get_function(header['name'])(header,args)
        print("RES",result)
        if not result is None:
            name,resp_data = result
            resp_header = self.make_header(name)
            resp_header['src_node'] = 'mc'
            self.dataflow.send(resp_header,resp_data,route)
        
    def remote_recv(self,h,c,raw,route,dataflow):
        pass
    
    def mc_recv(self,h,c,raw,route,dataflow):
        print("MC got",h,c,raw,route)
        print(str(route),h,c)
        if not route in self.routes_to_nodes:
            # If we don't have a Node for this already, we expect the
            # first message to be "hello".  Otherwise we ignore it.
            # If it is, make a Node object for it and send it the

            if h['name'] == '_mc_hello':
                if not c['group'] in self.groups:
                    print("Joining non-existent group",c['group'])
                    return
                # Make the ID for the node object
                new_id = self.gen_id(c['group'], c['id'])
                print("New node: " + c['id'] + " id = " + new_id)
                c['id'] = new_id
                h['src_node'] = new_id
                # Make the Node object
                iface = mc_if(c.get("if",{}))
                flags = c.get("flags", {})
                n = Node(new_id,self.gen_key(), self.dataflow, route, iface, flags=flags)

                # Start heartbeating thread
                n.hb_stopper = threading.Event()
                n.hb_thread = threading.Thread(target=mc_heartbeat_worker, daemon=True, args=(new_id,self.mc_hb_addr,self.hb_time,n.hb_stopper,self.context))
                n.hb_thread.start()
                
                # Send the "reg" message
                # header = self.make_header("reg")
                #dataflow.send(header,{"id":new_id,"key":n.key},route)
                
                # Add the Node object to our registery
                self.nodes[n.node_id] = n
                self.routes_to_nodes[route] = n
                print("passing along init msg finally",h,c,raw,c['id'])
                self.routes_to_nodes[route].send(raw,h,c)
            else:
                print("Non-hello init message: \n\n{}\n{}\n\nIgnoring".format(h,c))
        else:
            src_node = self.routes_to_nodes[route].node_id
            h['src_node'] = src_node
            print("passing along",h,c,raw,src_node)
            self.routes_to_nodes[route].send(raw,h,c)
        
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

    def node_del(self,header,args):
        nid = args["node"]
        if nid in self.nodes and self.nodes[nid].local:
            del self.nodes[nid]
            return "success",{}
        return "error",{'message':'fail'}

    def node_list(self,header,args):
        print("NODES")
        print(",".join(self.nodes.keys()))
        return "nodes",{"list":[x for x in self.nodes.keys() if (x != 'mc' and (not 'pattern' in args or args['pattern'] in x))]}

    def node_flags(self,header,args):
        nid = args["ID"];
        f = int(args["flags"])
        self.nodes[nid].flags = f
        pass

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
        return "error",{"message":"Failed to delete route"}
    
    def route_list(self,header,args):
        everything = r".*"
        sn = args.get('src_node',everything)
        dn = args.get('dest_node',everything)
        sns = [self.nodes[n] for n in self.nodes if re.match(sn,n)]
        rs = []
        for n in sns:
            rs += [n.routes[r] for r in n.routes if re.match(dn,n.routes[r].end.node_id) and re.match(dn,n.routes[r].end.name)]
        return "routes",{'routes':[str(r) for r in rs]}

    def route_add(self,header,args):
        print("BUILDING ROUTE",header,args)
        for r in self.transform_parser.parse(args['rt']):
            src = self.nodes[r[0][1]]
            dest = self.nodes[r[-1][1]]
            new_route = Route(src, dest, r[1:-1])
        self.nodes[sn].add_route(new_route)
        return "success",{}
            
    def rt_del(self,header,args):
        sn = self.find_node(args['src'].decode())
        res = sp.del_route(args['dest'].decode())
        return {'result':'success' if res else 'fail'}

    def launch(self,header,args):
        n = args['node']
        nid = args.get('id',n)
        base_path = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../')
        lib_base_path = os.path.abspath(os.path.dirname(os.path.abspath(__file__))+'/../../libmango')
        if n in self.node_types:
            lang = self.node_types[n]['lang']
            node_base = os.path.join(base_path,n if not 'dir' in self.node_types[n] else self.node_types[n]['dir'])
            lib_path = os.path.join(lib_base_path,lang)
            nenv = {'MC_ADDR':self.mc_target,'MANGO_ID':nid}
            if 'pathvar' in self.langs.get(lang, {}):
                nenv[self.langs[lang]['pathvar']] = lib_path
            if 'env' in args:
                nenv.update(json.loads(args['env']))
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
            else:
                print("E",nenv,"B",node_base)
                subprocess.Popen(shlex.split(self.node_types[n]['run']), cwd=node_base, env=nenv)
            return "success",{}
        return "error",{'message':'failure to launch'}

    def list_types(self,header,args):
        return "types",{'types':[x for x in self.node_types]}
    
    def find_if(self,header,args):
        #idf = m_send_sync("stdio",{'name':'get_if'})
        #return {'if':idf}
        pass

os.environ["MANGO_ID"] = "mc"
m = mc()
