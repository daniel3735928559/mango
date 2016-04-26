import io, re, socket, time, signal, os, sys, random, zmq, subprocess, shlex, json,traceback
from route_parser import route_parser
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from mangolib import m_node
from lxml import etree
from obj import *

class mc(m_node):
    
    def __init__(self,ID):
        super().__init__(ID,server=None)
        self.nodes = {"mc":Node("mc",0,loopback_dataflow(self))}
        self.node_types = {}
        self.route_parser = route_parser()
        
        # Now read xml file to get list of nodes and shell commands to spin them up.  
        # for child in etree.parse("nodes.xml").getroot():
        #     self.node_types[child.get("name")] = NodeType(child.get("name"),child.get("wd"),child.get("runner"))

        self.interface.add_interface("mc_if.yaml",
                                     {
                                         "excite":self.excite
                                     })
                                 # {
                                 #     "rt_list":self.rt_list,
                                 #     "rt_add":self.rt_add,
                                 #     "rt_del":self.rt_del,
                                 #     "port_add":self.port_add,
                                 #     "port_del":self.port_del,
                                 #     "node_del":self.node_del,
                                 #     "node_list":self.node_list,
                                 #     "node_flags":self.node_flags,
                                 #     "start":self.go,
                                 #     "list":self.type_list,
                                 #     "find_if":self.find_if,
                                 #     "remote_connect":self.remote_connect,
                                 #     "remote_disconnect":self.remote_disconnect
                                 # })
        print(self.interface.interface);
        self.uuid = str(self.gen_key())

        self.local_gateway = m_ZMQ_transport("tcp://*:"+sys.argv[1],self.context,self.poller,True)
        s = self.local_gateway.socket
        self.dataflows[s] = mc_router_dataflow(self,self.interface,self.local_gateway,self.serialiser,self.mc_recv,self.handle_reply,self.handle_error)
        self.poller.register(s,zmq.POLLIN)
        self.routes = {}

        #self.ports["stdio"] = self.dataflows[s]

        # Remote listening stuff: 

        # self.remote_gateway = m_srv_sock(int(sys.argv[2]))
        # f = self.remote_gateway.socket.fileno()
        # self.dataflows[f] = mc_remote_srv_dataflow(self,self.remote_recv,self.remote_gateway)
        # self.poller.register(f,zmq.POLLIN)

        self.run()

    def excite(self,header,args):
        print(args['str'])
        return args['str']+'!'
        
    def mc_recv(self,h,c,raw,dataflow):
        print("MC got")#: ", dataflow.route ,"on",dataflow)
        src_node = h['src_node']
        src_port = h['src_port']
        print(h)
        print(c)
        if not src_node in self.nodes:
            route = dataflow.route
            new_id = self.gen_id(src_node)
            print("New node: " + src_node + " id = " + new_id)
            trans = mc_ZMQ_transport(self.local_gateway.socket, bytearray(new_id,"ASCII"))
            df = m_dataflow(self.interface,trans,self.serialiser,self.mc_recv,self.handle_reply,self.handle_error)
            n = Node(new_id,self.gen_key(),df,master=self.nodes["mc"].ports["stdio"])
            header = self.make_header("stdio")
            header['function'] = 'reg'
            dataflow.send(header,{"node_id":new_id,"key":n.key})
            self.routes[route] = df
            self.nodes[n.node_id] = n
            #time.sleep(1) # Nasty hack to give the node time to change its ID
        if not src_port in self.nodes[src_node].ports:
            print("Invalid port: " + src_node+"/"+src_port)
            return
        print("sending anyway",h,c,raw,src_node,src_port)
        self.nodes[src_node].ports[src_port].send(raw,h,c)
        
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
            return {'result':'success'}
        return {'result':'fail'}

    def node_list(self,header,args):
        print("NODES")
        print(",".join(self.nodes.keys()))
        return {"list":", ".join([x for x in self.nodes.keys() if args['pattern'].decode() in x])}

    def node_flags(self,header,args):
        nid = args["ID"];
        f = int(args["flags"])
        self.nodes[nid].flags = f
        pass

    def rt_list(self,header,args):
        sn,sp = parse_port(args['src'])
        dn,dp = parse_port(args['dest'])
        
        sns = [self.nodes[n] for n in self.nodes if re.match(sn,n)]
        sps = [p for p in node.ports if re.match(sp,p) for node in sns]
        rs = [r for r in p.routes if re.match(dn,r.endpoint[0]) and re.match(dp,r.endpoint[1]) for p in sps]
        return {'routes':"\n".join([r.to_string() for r in rs])}

    def rt_add(self,header,args):
        print("BUILDING ROUTE",args)
        rt = self.route_parser.parse(args['chain'])
        if(rt is None):
            return {'result':'error parsing route'}
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
                        if(i != start and p in self.nodes[n].ports):
                            chains += [[r[1] for r in chain[start:i+1]]]
                            start = i
                        elif not self.nodes[n].local:
                            if not '.' in p:
                                p = p+'.stdio'
                                chain[i][1][1]+='.stdio'
                            print("Making a remote port",p,"in",n)
                            self.nodes[n].ports[p] = Port(p,self.nodes[n])
                            if(i != start):
                                chains += [[r[1] for r in chain[start:i+1]]]
                                start = i
                        else:
                            return {'result':'bad port: '+str(n)+'/'+str(p)}
                    else:
                        return {'result':'bad node: '+str(n)}
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
        return {'result':'success'}
            
    def rt_del(self,header,args):
        sn,sp = self.find_port(args['src'].decode())
        res = sp.del_route(args['dest'].decode())
        return {'result':'success' if res else 'fail'}

    def port_add(self,header,args):
        node,sp = self.parse_port(header['source'])
        port = args['name']
        print("PORT ADD: ",node, "PORT", port)
        if(not node in self.nodes):
            return {'result':'no such node'}
        if(port in self.nodes[node].ports):
            return {'result':'port exists'}
        self.nodes[node].ports[port] = Port(port,self.nodes[node])
        return {'result':'success'}

    def port_del(self,header,args):
        s = self.find_port(args['port'].decode())
        if s == None:
            return {'result':'fail'}
        node,port = s
        del self.nodes[node].ports[port]
        return {'result':'success'}

    def go(self,header,args):
        n = args['node'].decode()
        nid = args['id'].decode()
        if n in self.avail_nodes:
            subprocess.Popen(shlex.split(self.node_types[n].runner.replace("$ID",nid)))
            return {'result':'success'}
        return {'result':'fail'}

    def find_if(self,header,args):
        #idf = m_send_sync("stdio",{'command':'get_if'})
        #return {'if':idf}
        pass

    def type_list(self,header,args):
        return {'types':",".join(self.node_types.keys())}

m = mc("mc")
