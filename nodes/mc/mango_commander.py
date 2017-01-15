import io, re, socket, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from route_parser import route_parser
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
        self.route_parser = route_parser()
        
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
                                         #"delnode":self.delnode,
                                         "ports":self.port_list,
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
        self.routes = {}
        self.nodes = {"mc":Node("mc",0,mc_loopback_dataflow(self.interface,self.mc_dispatch,self.dataflow),bytes("mc","ASCII"),mc_if(self.interface.interface))}

        #self.ports["stdio"] = self.dataflows[s]

        # Remote listening stuff: 

        self.mc_remote_addr = "tcp://*:"+sys.argv[2]
        self.remote_gateway = m_ZMQ_transport(self.mc_remote_addr,self.context,self.poller,True)
        f = self.remote_gateway.socket
        self.dataflows[f] = mc_remote_dataflow(self.remote_gateway,self.serialiser,self.remote_recv)
        self.poller.register(f,zmq.POLLIN)

        # Heartbeat listening stuff:
        self.mc_hb_addr = "inproc://hb"
        self.hb_time = 10
        self.hb_gateway = m_ZMQ_transport(self.mc_hb_addr,self.context,self.poller,True)
        h = self.hb_gateway.socket
        self.dataflows[h] = mc_heartbeat_dataflow(self.hb_gateway,self.do_heartbeat)
        self.poller.register(h,zmq.POLLIN)

        # Start reaper thread
        self.reap_time = 15
        self.reaper_thread = threading.Thread(target=mc_heartbeat_worker, args=("mc",self.mc_hb_addr,self.reap_time,self.context))
        self.reaper_thread.start()
        self.too_old = 25
        
        self.run()

    def do_heartbeat(self,node_name):
        if node_name == b"mc":
            print("REAP")
            now = time.time()
            for n in self.nodes:
                if self.nodes[n].node_id != "mc" and (now - self.nodes[n].alive_time) > self.too_old:
                    print("REAPING",self.nodes[n].node_id)
        else:
            print("HB",node_name)
            header = self.make_header("heartbeat",src_port="mc")
            self.dataflow.send(header,{},self.nodes[node_name.decode()].route)

    def alive(self,header,args):
        if header['src_node'] in self.nodes:
            self.nodes[header['src_node']].alive_time = time.time()
        print("ALIVE",header,args)
        
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
        return {"doc":pijemont.doc.doc_gen(to_doc)}
        
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
            resp_header = self.make_header(name, src_port=header['port'])
            resp_header['src_node'] = 'mc'
            resp_header['port'] = 'mc'
            self.dataflow.send(resp_header,resp_data,route)
        
    def remote_recv(self,h,c,raw,route,dataflow):
        pass
    
    def mc_recv(self,h,c,raw,route,dataflow):
        print("MC got",h,c,raw,route)
        src_port = h['src_port']
        print(str(route),h,c)
        if not route in self.routes:
            # If we don't have a Node for this already, we expect the
            # first message to be "hello".  Otherwise we ignore it.
            # If it is, make a Node object for it and send it the

            if h['name'] == 'hello' and h['src_port'] == 'mc' and 'id' in c:
                # Make the ID for the node object
                new_id = self.gen_id(c['id'])
                print("New node: " + c['id'] + " id = " + new_id)
                c['id'] = new_id
                h['src_node'] = new_id
                # Make the Node object
                iface = mc_if(c.get("if",{}))
                ports = c.get("ports",[])
                flags = c.get("flags", {})
                n = Node(new_id,self.gen_key(), self.dataflow, route, iface, master=self.nodes["mc"].ports["stdio"], ports=ports, flags=flags)

                # Start heartbeating thread
                n.hb_thread = threading.Thread(target=mc_heartbeat_worker, args=(new_id,self.mc_hb_addr,self.hb_time,self.context))
                n.hb_thread.start()
                
                # Send the "reg" message
                # header = self.make_header("reg")
                #dataflow.send(header,{"id":new_id,"key":n.key},route)
                
                # Add the Node object to our registery
                self.nodes[n.node_id] = n
                self.routes[route] = n
                print("passing along init msg finally",h,c,raw,c['id'],src_port)
                self.routes[route].ports[src_port].send(raw,h,c)
            else:
                print("Non-hello init message: \n\n{}\n{}\n\nIgnoring".format(h,c))
        else:
            src_node = self.routes[route].node_id
            h['src_node'] = src_node
            if not src_port in self.routes[route].ports:
                # If there is no matching Port object in the specified Node, fail
                print("Invalid port: " + src_node+"/"+src_port)
                return
            # Now that we are guaranteed either way a matching Node/Port
            # combo for the incoming message, send it!,
            print("passing along",h,c,raw,src_node,src_port)
            self.routes[route].ports[src_port].send(raw,h,c)
        
    def gen_id(self,ID_req):
        if(not (ID_req in self.nodes.keys())):
            return ID_req
        for i in range(0,255):
            if(not (ID_req+"."+str(i) in self.nodes.keys())):
                return ID_req+"."+str(i)
        return -1
        
    def gen_key(self):
        return random.randint(0,2**64-1)

    def find_port(self,port_id):
        if not "/" in port_id:
            node = port_id
            port = "stdio"
        else:
            s = port_id.split("/")
            node = s[0]
            port=  s[1]
        if(not node in self.nodes):
            return None
        if not port in self.nodes[node].ports:
            return None
        return self.nodes[node].ports[port]

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

    def port_list(self,header,args):
        if not args['node'] in self.nodes:
            return "ports",{}
        n = self.nodes[args['node']]
        return "ports",{'ports':[p for p in n.ports]}

    def delroute(self,header,args):
        if not args['src_node'] in self.nodes:
            print(1)
            return "error",{"message":"Source node not found"}
        sn = self.nodes[args['src_node']]
        if not args['src_port'] in sn.ports:
            print(2)
            return "error",{"message":"Source port not found"}
        sp = sn.ports[args['src_port']]
        if not args['dest_node'] in self.nodes:
            print(3)
            return "error",{"message":"Target node not found"}
        dn = self.nodes[args['dest_node']]
        if not args['dest_port'] in dn.ports:
            print(4)
            return "error",{"message":"Target port not found"}
        dp = dn.ports[args['dest_port']]
        print("A",sn.ports,sp,dn.ports,dp,sp.routes)
        if sp.del_route_to(dp):
            return "success",{}
        return "error",{"message":"Failed to delete route"}
    
    def route_list(self,header,args):
        everything = r".*"
        sn = args.get('src_node',everything)
        sp = args.get('src_port',everything)
        dn = args.get('dest_node',everything)
        dp = args.get('dest_port',everything)
        sns = [self.nodes[n] for n in self.nodes if re.match(sn,n)]
        sps = []
        for n in sns:
            sps += [n.ports[p] for p in n.ports if re.match(sp,p)]
        rs = []
        for p in sps:
            rs += [p.routes[r] for r in p.routes if re.match(dn,p.routes[r].endpoint.owner.node_id) and re.match(dp,p.routes[r].endpoint.name)]
        return "routes",{'routes':[r.to_string() for r in rs]}

    def route_add(self,header,args):
        print("BUILDING ROUTE",header,args)
        rt = self.route_parser.parse(args['map'])
        if(rt is None):
            return "success",{}
        chains = []
        print("RT",rt)
        for chain in rt:
            start = 0
            for i in range(len(chain)):
                print('chaini',chain[i])
                if(chain[i][0] == 'port'):
                    n,p = chain[i][1]
                    print('port',n,p)
                    if(n in self.nodes):
                        print("NNN",i,start,p,p in self.nodes[n].ports,self.nodes[n].ports)
                        if not self.nodes[n].local:
                            if not '.' in p:
                                p = p+'.stdio'
                                chain[i][1][1]+='.stdio'
                            print("Making a remote port",p,"in",n)
                            self.nodes[n].ports[p] = Port(p,self.nodes[n])
                            if(i != start):
                                chains += [[r[1] for r in chain[start:i+1]]]
                                start = i
                        elif p in self.nodes[n].ports:
                            if i != start:
                                chains += [[r[1] for r in chain[start:i+1]]]
                                start = i
                        else:
                            return "error",{'message':'bad port: '+str(n)+'/'+str(p)}
                    else:
                        return "error",{'message':'bad node: '+str(n)}
        print("chains",chains)
        # Now check that all routes are valid
        for c in chains:
            print("C",c)
            sn,sp = c[0]
            dn,dp = c[-1]
            print("src:",sn,sp,"dest:",dn,dp)
            new_route = Route(self.nodes[sn].ports[sp],self.nodes[dn].ports[dp])
            for t in c[1:-1]:
                new_route.transmogrifiers += [t]
            self.nodes[sn].ports[sp].add_route(new_route)
        return "success",{}
            
    def rt_del(self,header,args):
        sn,sp = self.find_port(args['src'].decode())
        res = sp.del_route(args['dest'].decode())
        return {'result':'success' if res else 'fail'}

    def port_add(self,header,args):
        node,sp = self.parse_port(header['source'])
        port = args['name']
        print("PORT ADD: ",node, "PORT", port)
        if(not node in self.nodes):
            return "error",{'message':'no such node'}
        if(port in self.nodes[node].ports):
            return "error",{'message':'port exists'}
        self.nodes[node].ports[port] = Port(port,self.nodes[node])
        return "success",{}

    def port_del(self,header,args):
        s = self.find_port(args['port'].decode())
        if s == None:
            return "error",{'message':'port not found'}
        node,port = s
        del self.nodes[node].ports[port]
        return "success",{}

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
            if 'pathvar' in self.langs[lang]:
                nenv[self.langs[lang]['pathvar']] = lib_path
            if 'env' in args:
                nenv.update(json.loads(args['env']))
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
