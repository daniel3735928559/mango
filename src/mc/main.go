package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"github.com/docopt/docopt-go"
	"mc/route"
	"mc/emp"
	"mc/node"
	"mc/nodetype"
	"mc/value"
	"mc/serializer"
	"mc/transport"
	"mc/transport/mzmq"
	"mc/registry"
	"io/ioutil"
	"libmango"
	// "libmango/transport/msocket"
)

type MangoCommander struct {
	zmqTransport *mzmq.ZMQTransport
	//socketTransport *msocket.SocketTransport
	MessageInput chan transport.WrappedMessage
	Commands map[string]libmango.MangoHandler
	Registry *registry.Registry
	Self node.Node
}

func (mc *MangoCommander) Send(mid, command string, args map[string]interface{}) {
	if len(mid) == 0 {
		mid = fmt.Sprintf("%d",int(rand.Int()))
	}
	mc.MessageInput <- transport.WrappedMessage {
		Transport: nil,
		Identity: "",
		Message: serializer.Msg{
			Sender: "system/mc",
			MessageId: mid,
			Command: command,
			Cookie: mc.Self.GetId(),
			Data: args}}
}

func (mc *MangoCommander) reaper() {
	for {
		time.Sleep(60*time.Second)
		for _, n := range mc.Registry.GetNodes() {
			ago := n.SecsAgo()
			if ago > 300 || ago < 0 {
				fmt.Println("[MC] REAP",n.ToString())
				mc.Registry.DelNode(n.GetId())
			}
		}
	}
}

func (mc *MangoCommander) Run() {
	go mc.zmqTransport.RunServer()
	//go mc.socketTransport.RunServer()
	go mc.reaper()

	// Start system EMP
	time.Sleep(1*time.Second)
	err := mc.EMP("system", "emp/system.emp", map[string]string{"AGENT_PORT":"11313"})
	if err != nil {
		fmt.Println("Problem running system.emp:", err)
	}

	for wrapped_msg := range mc.MessageInput {
		msg := wrapped_msg.Message
		// fmt.Println("[MC] FROM",msg.Sender)
		// fmt.Println("[MC] DATA",msg.Data)
		src := mc.Registry.FindNodeById(msg.Cookie)
		if src == nil {
			fmt.Println("[MC] ERROR: node not found:", msg.Cookie)
			continue
		}
		// Since we got here, the node we're talking to is
		// using this transport/id combo, so set those
		// properties here (in case it has reconnected since
		// last time or something):
		if wrapped_msg.Transport != nil {
			src.GotAlive(wrapped_msg.Identity, wrapped_msg.Transport)
		}
		
		// validate message is of acceptable format for node
		// output:
		cmd := msg.Command
		if cmd == "alive" {
			// Alive messages are already handled
			continue
		}
		incoming_val, err := value.FromObject(msg.Data)
		if err != nil {
			fmt.Println("[MC] ERROR: failed to convert incoming value",err)
			continue
		}
		
		src_type := mc.Registry.FindNodeType(src.GetType())
		if src_type == nil {
			fmt.Println("[MC] ERROR: Could not find type: ",src.GetType())
			continue
		}

		validated_val := incoming_val
		if src_type.Validate {
			v, err := src_type.ValidateOutput(cmd, incoming_val)
			if err != nil {
				fmt.Println("[MC] ERROR: output message failed to validate from:",src_type.Name)
				if incoming_val != nil {
					fmt.Println("[MC] ERROR: Invalid message:",incoming_val.ToString())
				}
				continue
			}
			fmt.Println("[MC] VALIDATION SUCCESSFUL")
			validated_val = v
		}
		
		// If we made it here, the value is legitimate. Send
		// on all routes originating at src:
		routes := mc.Registry.FindRoutesBySrc(src.ToString())
		fmt.Printf("[MC] SENDING ON %d ROUTES\n", len(routes))
		for _, rt := range routes {
			fmt.Println("[MC] ROUTE",rt.ToString())
			dst := mc.Registry.FindNodeByName(rt.GetDest())
			if dst == nil {
				fmt.Println("[MC] ERROR: destination node not found:", rt.GetDest())
				continue
			}
			fmt.Println("[MC] DST",dst.ToString())
			result_cmd, result_val, result_err := rt.Run(cmd, validated_val)
			if result_err != nil {
				fmt.Println("[MC] ERROR: failed running route", rt.ToString(), result_err)
				continue
			}
				
			outval := result_val
			if outval == nil {
				fmt.Println("[MC] Filtered out...moving along")
				continue
			}
			
			// Validate, if needed
			fmt.Println("[MC] getting type",dst.GetType())
			dst_type := mc.Registry.FindNodeType(dst.GetType())
			if dst_type.Validate {
				v, err := dst_type.ValidateInput(result_cmd, result_val)
				if err != nil {
					fmt.Println("[MC] ERROR: route output failed validation", result_cmd, result_val, err)
					continue
				}
				outval = v
			}
			outmsg := serializer.Msg{
				Sender: src.ToString(),
				MessageId: msg.MessageId,
				Command: result_cmd,
				Data: outval.ToObject().(map[string]interface{})}
			fmt.Println("[MC] SEND TO",dst.ToString())
			dst.SendToNode(outmsg)
		}
	}
}

func (mc *MangoCommander) MsgHandler(m serializer.Msg) error {
	fmt.Println("[MC] SELF RX",m)
	if handler, ok := mc.Commands[m.Command]; ok {
		fmt.Println("[MC] found")
		retcmd, retval, err := handler(m.Data)
		if err != nil {
			mc.Send("", "error", map[string]interface{}{"message":fmt.Sprintf("%v", err)})
		} else {
			mc.Send(m.MessageId, retcmd, retval)
		}
	} else {
		return fmt.Errorf("[MC] ERROR: Invalid command %s\n",m.Command)
	}
	return nil
}

func (mc *MangoCommander) Doc(args map[string]interface{}) (string, map[string]interface{}, error) {
	node_type := mc.Registry.FindNodeType(args["nodetype"].(string))
	if node_type == nil {
		return "", nil, fmt.Errorf("Type not found: %s", args["nodetype"].(string))
	}
	if command_if, ok := args["command"]; ok {
		command := command_if.(string)
		if inp, ok := node_type.Interface.Inputs[command]; ok {
			return "doc", map[string]interface{}{"doc":inp.ToString()},nil
		}
		return "", nil, fmt.Errorf("Input not found: %s", command)
	}
	return "doc", map[string]interface{}{"doc":node_type.Interface.ToString()}, nil
}

func (mc *MangoCommander) RouteAdd(args map[string]interface{}) (string, map[string]interface{}, error) {
	spec := args["spec"].(string)
	group := args["group"].(string)
	id := ""
	for i := 0; i < 4; i++ {
		id += fmt.Sprintf("%016x", rand.Uint64())
	}
	routes, err := route.Parse(spec, group)
	if err != nil {
		return "", nil, fmt.Errorf("[MC] ERROR: Routes failed to parse: %v", err)
	}
	ans := make([]interface{}, 0)
	for i, rt := range routes {
		rt.Id = fmt.Sprintf("%s_%d", id, i)
		rt.Group = group
		mc.Registry.AddRoute(rt)
		ans = append(ans, map[string]interface{}{"id":rt.Id,"src":rt.GetSource(),"dst":rt.GetDest(),"spec":rt.ToString()})
	}
	return "routeinfo",map[string]interface{}{"routes":ans},nil
}

func (mc *MangoCommander) StartNode(nodetype, group, name string, args []interface{}) error {
	if new_type := mc.Registry.FindNodeType(nodetype); new_type != nil {
		nodeargs := make([]string, len(args))
		for i, a := range args {
			nodeargs[i] = a.(string)
		}
		new_node := node.MakeExecNode(group, name, nodetype, fmt.Sprintf("%s %s", new_type.Command, strings.Join(nodeargs, " ")), new_type.Environment, mc.zmqTransport)
		if new_node == nil {
			return fmt.Errorf("[MC] ERROR Failed to make node: `%s`", name)
		}
		
		mc.Registry.AddNode(new_node)
		new_node.Start(mc.zmqTransport.GetServerAddr())
	} else {
		return fmt.Errorf("[MC] ERROR Failed to find type: `%s`", nodetype)
	}
	return nil
}

func (mc *MangoCommander) EMP(group, emp_file string, args map[string]string) error {
	data, err := ioutil.ReadFile(emp_file)
	if err != nil {
		return fmt.Errorf("[MC] ERROR reading EMP file: %s, %v",emp_file,err)
	}
	emp_data := string(data)
	e, err := emp.Parse(emp_data)
	if err != nil {
		return fmt.Errorf("[MC] ERROR parsing EMP: `%s` %v",emp_data, err)
	}

	// Validate that all arguments were passed:
	fmt.Println("EMP ARGS", args)
	arg_replacements := make([]string, 0)
	for _, param := range e.ParamNames {
		if _, ok := args[param]; !ok {
			return fmt.Errorf("[MC] ERROR parsing EMP: %s -- Argument `%s` not found", emp_data, param)
		}
		arg_replacements = append(arg_replacements, fmt.Sprintf("{%s}",param), args[param])
	}
	arg_replacer := strings.NewReplacer(arg_replacements...)
	fmt.Println("made replacer")
	
	// Parse all the nodes
	new_nodes := make([]node.Node, 0)
	for _, n := range e.Nodes {
		if n.TypeName == "dummy" {
			new_dummy_node := node.MakeDummyNode(fmt.Sprintf("%s_%s",group,n.Name), group, n.Name, mc.MessageInput)
			new_nodes = append(new_nodes, new_dummy_node)
		} else if n.TypeName == "merge" {
			new_merge_nodes := node.MakeMergeNode(group, n.Name, n.Args, mc.MessageInput)
			new_nodes = append(new_nodes, new_merge_nodes...)
		} else if new_type := mc.Registry.FindNodeType(n.TypeName); new_type != nil {
			fmt.Printf("[MC] EXEC %s %s",new_type.Command, strings.Join(n.Args, " "))
			new_args := make([]string, 0)
			for _, a := range n.Args {
				new_args = append(new_args, arg_replacer.Replace(a))
			}
			
			new_node := node.MakeExecNode(group, n.Name, n.TypeName, fmt.Sprintf("%s %s", new_type.Command, strings.Join(new_args, " ")), new_type.Environment, mc.zmqTransport)
			if new_node == nil {
				return fmt.Errorf("[MC] ERROR Failed to make node: `%s`", n.Name)
			}
			new_nodes = append(new_nodes, new_node)
		} else {
			return fmt.Errorf("[MC] ERROR Failed to find type: `%s`", n.TypeName)
		}
	}

	// Parse all the given routes
	new_routes := make([]*route.Route, 0)
	route_idx := 0
	for _, spec := range e.Routes {
		// Replace arguments as they occur in route spec:
		argspec := arg_replacer.Replace(spec)

		// Parse route
		routes, err := route.Parse(argspec, group)
		if err != nil {
			return fmt.Errorf("[MC] ERROR: Routes failed to parse: %v: %s -> %s", err, spec, argspec)
		}
		for _, rt := range routes {
			rt.Id = fmt.Sprintf("%s_r%d", group, route_idx)
			route_idx++
			rt.Group = group
			new_routes = append(new_routes, rt)
		}
	}
	for _, n := range new_nodes {
		mc.Registry.AddNode(n)
	}
	for _, rt := range new_routes {
		mc.Registry.AddRoute(rt)
	}
	for _, n := range new_nodes {
		if en, ok := n.(*node.ExecNode); ok {
			en.Start(mc.zmqTransport.GetServerAddr())
		}
	}
	
	return nil
}

func (mc *MangoCommander) RunEMP(args map[string]interface{}) (string, map[string]interface{}, error) {
	group := args["group"].(string)
	if _, ok := mc.Registry.Groups[group]; ok {
		return "", nil, fmt.Errorf("[MC] ERROR: group already exists: `%s`" ,group)
	}
	filename := args["filename"].(string)
	emp_args := make(map[string]string)
	for _, a := range args["args"].([]interface{}) {
		arg := a.(map[string]interface{})
		emp_args[arg["name"].(string)] = arg["value"].(string)
	}
	err := mc.EMP(group, filename, emp_args)
	return "", nil, err
}

func (mc *MangoCommander) RunNode(args map[string]interface{}) (string, map[string]interface{}, error) {
	group := args["group"].(string)
	nodetype := args["type"].(string)
	name := args["name"].(string)
	nodeargs := args["args"].([]interface{})
	err := mc.StartNode(nodetype, group, name, nodeargs)
	return "", nil, err
}

func (mc *MangoCommander) FindTypes(args map[string]interface{}) (string, map[string]interface{}, error)  {
	var types []*nodetype.NodeType
	name, by_name := args["name"]
	if by_name {
		types = []*nodetype.NodeType{mc.Registry.FindNodeType(name.(string))}
	} else {
		for _, ty := range mc.Registry.NodeTypes {
			types = append(types, ty)
		}
	}
	
	ans := map[string]interface{}{"types":make([]interface{}, len(types))}
	for i, ty := range types {
		if ty != nil {
			ans["types"].([]interface{})[i] = map[string]interface{}{
				"name":ty.Name,
				"usage":ty.Usage,
				"command":ty.Command,
				"interface":ty.Interface.ToString()}
		}
	}
	return "typeinfo",ans,nil
}

func (mc *MangoCommander) FindGroups(args map[string]interface{}) (string, map[string]interface{}, error)  {
	groups := mc.Registry.GetGroups()
	return "groupinfo",map[string]interface{}{"groups":groups},nil
}

func (mc *MangoCommander) FindRoutes(args map[string]interface{}) (string, map[string]interface{}, error)  {
	// group, by_group := args["group"]
	// name, by_name := args["name"]
	// var nodes []node.Node
	// if by_group && !by_name {
	// 	nodes = mc.Registry.FindNodesByGroup(group.(string))
	// } else if !by_group && !by_name {
	// 	nodes = mc.Registry.GetNodes()
	// } else if by_name {
	// 	nodes = []node.Node{mc.Registry.FindNodeByName(name.(string))}
	// }
	var routes []*route.Route
	routes = mc.Registry.GetRoutes()
	
	ans := map[string]interface{}{"routes":make([]interface{}, len(routes))}
	for i, rt := range routes {
		ans["routes"].([]interface{})[i] = map[string]interface{}{"src":rt.GetSource(),"dst":rt.GetDest(),"id":rt.Id,"spec":rt.ToString()}
	}
	return "routeinfo",ans,nil
}

func (mc *MangoCommander) FindNodes(args map[string]interface{}) (string, map[string]interface{}, error)  {
	group, by_group := args["group"]
	name, by_name := args["name"]
	var nodes []node.Node
	if by_group && !by_name {
		nodes = mc.Registry.FindNodesByGroup(group.(string))
	} else if !by_group && !by_name {
		nodes = mc.Registry.GetNodes()
	} else if by_name {
		nodes = []node.Node{mc.Registry.FindNodeByName(name.(string))}
	}
	
	ans := map[string]interface{}{"nodes":make([]interface{}, len(nodes))}
	for i, no := range nodes {
		if no != nil {
			ans["nodes"].([]interface{})[i] = map[string]interface{}{
				"type":no.GetType(),
				"name":no.GetName(),
				"group":no.GetGroup(),
				"lastseen":no.LastSeen()}
		}
	}
	return "nodeinfo",ans,nil
}

func (mc *MangoCommander) GetGroup(args map[string]interface{}) (string, map[string]interface{}, error)  {
	nodes := mc.Registry.FindNodesByGroup(args["name"].(string))
	ans := map[string]interface{}{"nodes":make([]map[string]interface{}, len(nodes))}
	for i, no := range nodes {
		ans["nodes"].([]map[string]interface{})[i] = map[string]interface{}{"type":no.GetType(),"name":no.GetName(),"group":no.GetGroup()}
	}
	return "nodeinfo",ans,nil
}

func (mc *MangoCommander) GroupDel(args map[string]interface{}) (string, map[string]interface{}, error)  {
	mc.Registry.DelGroup(args["name"].(string))
	return "",nil,nil
}

func (mc *MangoCommander) NodeDel(args map[string]interface{}) (string, map[string]interface{}, error)  {
	mc.Registry.DelNode(args["id"].(string))
	return "",nil,nil
}

func (mc *MangoCommander) RouteDel(args map[string]interface{}) (string, map[string]interface{}, error)  {
	mc.Registry.DelRoute(args["id"].(string))
	return "",nil,nil
}

func (mc *MangoCommander) HandleError(args map[string]interface{}) (string, map[string]interface{}, error)  {
	return "",nil,nil
}

func (mc *MangoCommander) Echo(args map[string]interface{}) (string, map[string]interface{}, error)  {
	return "echo",args,nil
}

func main() {
	usage := `Usage: mc [-t] [--manifest=<manifest>]

Options:
-m --manifest manifest_file  A file from which to read the available node types`
	args, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Println("[MC] ERROR parsing args: ",err)
	}
	
	fmt.Println("hi")

	// Create MC
	
	//go test_socket_client(1920)
	message_aggregator := make(chan transport.WrappedMessage, 100)
	
	MC := &MangoCommander{
		zmqTransport: mzmq.MakeZMQTransport(1919, message_aggregator),
		//socketTransport: MakeSocketTransport(1920, message_aggregator),
		MessageInput: message_aggregator,
		Registry: registry.MakeRegistry()}	
	
	MC.Commands = map[string]libmango.MangoHandler{
		"start":MC.RunNode,
		"routeadd":MC.RouteAdd,
		"findgroups":MC.FindGroups,
		"findroutes":MC.FindRoutes,
		"findtypes":MC.FindTypes,
		"findnodes":MC.FindNodes,
		"getgroup":MC.GetGroup,
		"groupdel":MC.GroupDel,
		"routedel":MC.RouteDel,
		"nodedel":MC.NodeDel,
		"error":MC.HandleError,
		"echo":MC.Echo,
		"doc":MC.Doc,
		"emp":MC.RunEMP}
	
	// Load node types from manifest:
	fmt.Println(args)
	manifest_filename := "types.manifest"
	if mf, ok := args["--manifest"]; ok && mf != nil {
		manifest_filename = mf.(string)
	}
		manifest_data, err := ioutil.ReadFile(manifest_filename)
	if err != nil {
		fmt.Println("[MC] ERROR reading manifest", manifest_filename, err)
		return
	}
	for _, line := range strings.Split(string(manifest_data), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		nodetype_data, err := ioutil.ReadFile(line)
		if err != nil {
			fmt.Println("[MC] ERROR reading node type file", manifest_filename, err)
			return
		}
		nt, err := nodetype.Parse(string(nodetype_data))
		if err != nil {
			fmt.Println("[MC] ERROR Failed parsing node type: ", string(nodetype_data), err)
			return
		}
		MC.Registry.AddNodeType(nt)
	}
	
	// Add self as a node
	MC.Self = node.MakeCallbackNode("system", "mc", "mc", MC.MsgHandler)
	MC.Registry.AddNode(MC.Self)
	
	// if args["-t"].(bool) {
	// 	TestNode := node.MakeNode("test", "system", "test_node", "", MC.zmqTransport)
	// 	go mzmq.TestZmqClient(1919, TestNode.Id)
		
	// 	MC.Registry.AddNode(TestNode)

	// 	testroutes, rerr := route.Parse(`system/test > ? excite % echo {message += "!";} > system/mc`)
	// 	if rerr != nil {
	// 		fmt.Println("Failed parsing test route: ",rerr)
	// 		return
	// 	}

	// 	for i, rt := range testroutes {
	// 		rt.Group = "system"
	// 		rt.Id = fmt.Sprintf("test_%d", i)
	// 		MC.Registry.AddRoute(rt)
	// 	}
		
	// }
	MC.Run()
}
