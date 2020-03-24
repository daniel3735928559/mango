package main

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"encoding/json"
	"libmango"
	zmq "github.com/pebbe/zmq4"
)

type MxHandler func(map[string]interface{}, chan string)

type MxMsg struct {
	Command string
	Args map[string]interface{}
}

type MxAgent struct {
	node *libmango.Node
	sock *zmq.Socket
	route_ids []string
	target_type string
	MsgQueue []MxMsg
	handlers map[string]MxHandler
}

func NewMx() *MxAgent {
	n, err := libmango.NewNode("mx",map[string]libmango.MangoHandler{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	
	sock, _ := zmq.NewSocket(zmq.ROUTER)
	mx := &MxAgent{
		node:n,
		sock:sock,
		route_ids: make([]string, 0),
		target_type: "",
		MsgQueue: make([]MxMsg, 0),
		handlers: make(map[string]MxHandler)}

	mx.handlers["send"] = mx.Send
	mx.handlers["start"] = mx.StartNode
	mx.handlers["emp"] = mx.Emp
	mx.handlers["connect"] = mx.Connect
	mx.handlers["disconnect"] = mx.Disconnect
	mx.handlers["help"] = mx.Help
	mx.handlers["types"] = mx.Types
	mx.handlers["nodes"] = mx.Nodes
	mx.handlers["routes"] = mx.Routes
	mx.handlers["list"] = mx.List
	mx.handlers["pop"] = mx.Pop
	mx.handlers["get"] = mx.Get

	n.DefaultHandler = mx.HandleFromWorld
	n.Start()
	return mx
}

func (mx *MxAgent) HandleFromWorld(command string, args map[string]interface{}) (string, map[string]interface{}, error) {
	fmt.Println("Got")
	mx.MsgQueue = append(mx.MsgQueue, MxMsg{Command:command, Args:args})
	return "", nil, nil
}

func (mx *MxAgent) Connect(req map[string]interface{}, rep chan string) {
	if len(mx.route_ids) > 0 {
		for _, route_id := range mx.route_ids {
			mx.node.Send("routedel", map[string]interface{}{"control":true,"id":route_id})
		}
		mx.route_ids = []string{}
	}
	ConnectReplyHandler := func(c string, a map[string]interface{}) {
		fmt.Println("[MX AGENT] ConnectReply", c, a)
		routes := a["routes"].([]interface{})
		for _, rt := range routes {
			mx.route_ids = append(mx.route_ids, rt.(map[string]interface{})["id"].(string))
		}
		fmt.Println("connected",mx.route_ids)
		rep <- "Connected"
		close(rep)
		
		// Now get the type of the newly connected target
		FindReplyHandler := func(c string, a map[string]interface{}) {
			fmt.Println("[MX AGENT] FindReply", c, a)
			foundnodes := a["nodes"].([]interface{})
			if len(foundnodes) == 1 {
				f := foundnodes[0].(map[string]interface{})
				mx.target_type = f["type"].(string)
			}
		}
		mx.node.SendForReply("findnodes", map[string]interface{}{"control":true,"name":req["target"].(string)}, FindReplyHandler)
	}
	mx.node.SendForReply("routeadd", map[string]interface{}{"control":true,"group":"system","spec":fmt.Sprintf("system/cli <> %s", req["target"].(string))}, ConnectReplyHandler)
}
func (mx *MxAgent) Disconnect(req map[string]interface{}, rep chan string) {
	if len(mx.route_ids) == 0 {
		fmt.Println("[MX AGENT] Not connected")
		rep <- "Not connected"
		close(rep)
	} else {
		for _, route_id := range mx.route_ids {
			mx.node.Send("routedel", map[string]interface{}{"control":true,"id":route_id})
		}
		mx.route_ids = []string{}
		rep <- "Connecting"
		close(rep)
	}
}
func (mx *MxAgent) Help(req map[string]interface{}, rep chan string) {
	if len(mx.target_type) == 0 {
		rep <- "Not connected"
		close(rep)
	}
	args := map[string]interface{}{"control":true,"nodetype":mx.target_type}
	if c, ok := req["command"]; ok {
		args["command"] = c.(string)
	}
	fmt.Println("[MX AGENT] HELP",args,mx.target_type)
	DocReplyHandler := func(c string, a map[string]interface{}) {
		fmt.Println("[MX AGENT] DocReply", c, a)
		rep <- a["doc"].(string)
		close(rep)
	}
	mx.node.SendForReply("doc", args, DocReplyHandler)
}

func (mx *MxAgent) Send(req map[string]interface{}, rep chan string) {
	mx.node.Send(req["command"].(string), req["args"].(map[string]interface{}))
	fmt.Println("[MX AGENT] SENT")
	rep <- "Sent"
	close(rep)
}

func (mx *MxAgent) Emp(req map[string]interface{}, rep chan string) {
	mx.node.Send("emp", map[string]interface{}{
		"control":true,
		"filename":req["empfile"].(string),
		"group":req["group"].(string)})
	fmt.Println("[MX AGENT] SENT")
	rep <- "Sent"
	close(rep)
}

func (mx *MxAgent) StartNode(req map[string]interface{}, rep chan string) {
	fmt.Println("[MX AGENT] StartNode",req)
	d := map[string]interface{}{
		"control":true,
		"name":req["name"].(string),
		"type":req["type"].(string),
		"group":req["group"].(string),
		"args":req["args"].([]interface{})}
	fmt.Println("D",d)
	mx.node.Send("start", d)
	fmt.Println("[MX AGENT] SENT")
	rep <- "Sent"
	close(rep)
}

func (mx *MxAgent) Nodes(req map[string]interface{}, rep chan string) {
	QueryReplyHandler := func(c string, a map[string]interface{}) {
		foundnodes := a["nodes"].([]interface{})
		fmt.Println("[MX AGENT] QueryReply", c, a)
		ans := make([]string, len(foundnodes))
		for i, fi := range foundnodes {
			f := fi.(map[string]interface{})
			ans[i] = fmt.Sprintf("%s %s/%s", f["type"].(string), f["group"].(string), f["name"].(string))
		}
		fmt.Println("[MX AGENT] NODES",ans)
		rep <- strings.Join(ans,"\n")
		close(rep)
	}
	mx.node.SendForReply("findnodes", map[string]interface{}{"control":true}, QueryReplyHandler)
}


func (mx *MxAgent) Types(req map[string]interface{}, rep chan string) {
	QueryReplyHandler := func(c string, a map[string]interface{}) {
		foundtypes := a["types"].([]interface{})
		fmt.Println("[MX AGENT] QueryReply", c, a)
		ans := make([]string, len(foundtypes))
		for i, fi := range foundtypes {
			f := fi.(map[string]interface{})
			ans[i] = fmt.Sprintf(`%s: 
Usage: %s
Command: %s
Interface: %s
--
`, f["name"].(string), f["usage"].(string), f["command"].(string), f["interface"].(string))
		}
		fmt.Println("[MX AGENT] TYPES",ans)
		rep <- strings.Join(ans,"\n")
		close(rep)
	}
	mx.node.SendForReply("findtypes", map[string]interface{}{"control":true}, QueryReplyHandler)
}

func (mx *MxAgent) Routes(req map[string]interface{}, rep chan string) {
	QueryReplyHandler := func(c string, a map[string]interface{}) {
		foundroutes := a["routes"].([]interface{})
		fmt.Println("[MX AGENT] QueryReply", c, a)
		ans := make([]string, len(foundroutes))
		for i, fr := range foundroutes {
			f := fr.(map[string]interface{})
			ans[i] = fmt.Sprintf("%s: %s", f["id"].(string), f["spec"].(string))
		}
		fmt.Println("[MX AGENT] ROUTES",ans)
		rep <- strings.Join(ans,"\n")
		close(rep)
	}
	mx.node.SendForReply("findroutes", map[string]interface{}{"control":true}, QueryReplyHandler)
}
func (mx *MxAgent) List(req map[string]interface{}, rep chan string) {
	ans := make([]string, len(mx.MsgQueue))
	for i, m := range mx.MsgQueue {
		bs, _ := json.Marshal(m.Args)
		ans[i] = fmt.Sprintf("%d: %s %s", i, m.Command, string(bs))
	}
	rep <- strings.Join(ans,"\n")
	close(rep)
}
func (mx *MxAgent) Pop(req map[string]interface{}, rep chan string) {
	if len(mx.MsgQueue) > 0 {
		m := mx.MsgQueue[0]
		bs, _ := json.Marshal(m.Args)
		ans := fmt.Sprintf("%s %s", m.Command, string(bs))
		mx.MsgQueue = mx.MsgQueue[1:]
		rep <- ans
		close(rep)
	} else {
		rep <- "No more messages"
		close(rep)
	}
}
func (mx *MxAgent) Get(req map[string]interface{}, rep chan string) {
	idx, _ := strconv.Atoi(req["id"].(string))
	if idx < len(mx.MsgQueue) {
		m := mx.MsgQueue[idx]
		bs, _ := json.Marshal(m.Args)
		ans := fmt.Sprintf("%s %s %s", idx, m.Command, string(bs))
		rep <- ans
		close(rep)
	} else {
		rep <- fmt.Sprintf("No such message: %d", idx)
		close(rep)
	}
}

func (mx *MxAgent) HandleFromClient(req map[string]interface{}) string {
	op := req["operation"].(string)
	if handler, ok := mx.handlers[op]; ok {
		rep := make(chan string, 10)
		handler(req, rep)
		select {
		case ans := <- rep:
			fmt.Println("[MX AGENT] REPLIED")
			return ans
		case <- time.After(3*time.Second):
			return "Timeout"
		}
	}
	return fmt.Sprintf("Unknown operation: %s", op)
}

func (mx *MxAgent) RunZmqServer(port int) {
	mx.sock, _ = zmq.NewSocket(zmq.ROUTER)
	fmt.Println("[MX AGENT] BIND", fmt.Sprintf("tcp://*:%d", port))
	mx.sock.Bind(fmt.Sprintf("tcp://*:%d", port))
	for {
		fmt.Println("[MX AGENT] Running")
		client_identity, _ := mx.sock.Recv(0)
		fmt.Println("[MX AGENT] got id",client_identity)
		mx.sock.Recv(0)
		fmt.Println("[MX AGENT] got _")
		data, _ := mx.sock.Recv(0)
		fmt.Println("[MX AGENT] RECV",data)
		req := make(map[string]interface{})
		json.Unmarshal([]byte(data), &req)
		response := mx.HandleFromClient(req)
		mx.sock.Send(client_identity, zmq.SNDMORE)
		mx.sock.Send(response, 0)
	}
}

func main() {
	MX := NewMx()
	MX.RunZmqServer(11313)
}
