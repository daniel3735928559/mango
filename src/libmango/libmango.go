package libmango

import (
	"os"
	"fmt"
	"sync"
	"time"
	"strings"
	"math/rand"
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
)

type MangoHandler func(map[string]interface{}) (string, map[string]interface{}, error)
type MangoReplyHandler func(string, map[string]interface{})
type MangoDefaultHandler func(string, map[string]interface{}) (string, map[string]interface{}, error)

type NodeStatus int

const (
	NODE_STATUS_STOPPED NodeStatus = iota+1
	NODE_STATUS_RUNNING
)

type MsgHeader struct {
	Source string  `json:"source"`
	MessageId string  `json:"mid"`
	Command string `json:"command"`
	Cookie string `json:"cookie,omitempty"`
	Format string  `json:"format"`
}

type Msg struct {
	MessageId string
	Command string
	Cookie string
	Data map[string]interface{}
}

type Node struct {
	Name string
	Server string
	ctx *zmq.Context
	status NodeStatus
	cookie string
	currentMid int
	handlers map[string]MangoHandler
	DefaultHandler MangoDefaultHandler
	outstanding map[string]MangoReplyHandler
	send_mutex *sync.Mutex
	socket *zmq.Socket
	outgoing chan interface{}
	done_internal chan bool
	done_reactor chan interface{}
	Done chan bool
}

func NewNode(name string, handlers map[string]MangoHandler) (*Node, error) {
	n := &Node{
		status: NODE_STATUS_STOPPED,
		currentMid: 0,
		Name: name,
		handlers: handlers,
		send_mutex: &sync.Mutex{},
		outstanding: make(map[string]MangoReplyHandler),
		outgoing: make(chan interface{}, 1000),
		done_internal: make(chan bool),
		done_reactor: make(chan interface{}),
		Done: make(chan bool)}
	n.DefaultHandler = n.default_default_handler
	if srv, ok := os.LookupEnv("MANGO_SERVER"); ok {
		n.Server = srv
	} else {
		return nil, fmt.Errorf("ERROR: environment var `MANGO_SERVER` not found")
	}
	if c, ok := os.LookupEnv("MANGO_COOKIE"); ok {
		n.cookie = c
	} else {
		return nil, fmt.Errorf("ERROR: environment var `MANGO_COOKIE` not found")
	}
	n.handlers["heartbeat"] = n.heartbeat
	n.handlers["exit"] = n.exit
	
	return n, nil
}

func (n *Node) Send(cmd string, args map[string]interface{}) error {
	n.send_mutex.Lock()
	defer n.send_mutex.Unlock()
	fmt.Println("[LIBMANGO] sending",cmd,args)
	h := MsgHeader{
		MessageId: fmt.Sprintf("%d",n.currentMid),
		Cookie: n.cookie,
		Command: cmd,
		Format: "json"}
	
	n.currentMid++
	header_bytes, _ := json.Marshal(h)
	data_bytes, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("Failed to marshal arguments:%v", args)
	}

	n.outgoing <- fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes))
	// n.socket.Send("",zmq.SNDMORE)
	// n.socket.Send(fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes)), 0)
	fmt.Println("[LIBMANGO] sent")
	return nil
}

func (n *Node) SendForReply(cmd string, args map[string]interface{}, reply_handler MangoReplyHandler) (string, error) {
	n.send_mutex.Lock()
	defer n.send_mutex.Unlock()
	fmt.Println("[LIBMANGO] sending for reply",cmd,args)
	h := MsgHeader{
		MessageId: fmt.Sprintf("%d",n.currentMid),
		Cookie: n.cookie,
		Command: cmd,
		Format: "json"}

	n.outstanding[h.MessageId] = reply_handler
	
	n.currentMid++
	header_bytes, _ := json.Marshal(h)
	data_bytes, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal arguments:%v", args)
	}

	n.outgoing <- fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes))
	// n.socket.Send("",zmq.SNDMORE)
	// n.socket.Send(fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes)), 0)
	return h.MessageId, nil
}

func (n *Node) send_reply(mid, cmd string, args map[string]interface{}) {
	n.send_mutex.Lock()
	defer n.send_mutex.Unlock()
	fmt.Println("[LIBMANGO] sending reply",cmd,args)
	h := MsgHeader{
		MessageId: mid,
		Cookie: n.cookie,
		Command: cmd,
		Format: "json"}
	
	header_bytes, err := json.Marshal(h)
	data_bytes, err := json.Marshal(args)
	if err != nil {
		n.handle_error(fmt.Errorf("Failed to marshal arguments:%v", args))
	} else {
		n.outgoing <- fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes))
		// n.socket.Send("",zmq.SNDMORE)
		// n.socket.Send(fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes)), 0)
	}
}

func (n *Node) Start() {
	n.ctx, _ = zmq.NewContext()
	fmt.Println("[LIBMANGO] STARTING", n.Server)
	n.socket, _ = n.ctx.NewSocket(zmq.DEALER)
	rand.Seed(time.Now().UnixNano())
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	fmt.Println("[LIBMANGO]", identity)
	n.socket.SetIdentity(identity)
	n.socket.Connect(n.Server)
	n.status = NODE_STATUS_RUNNING
	n.Send("alive", map[string]interface{}{})
	go n.heartbeat_worker()
	n.run()
}

func (n *Node) Stop() {
	n.status = NODE_STATUS_STOPPED
}

func (n *Node) handle_error(err error) {
	fmt.Printf("[LIBMANGO] ERROR: %v", err)
}

func (n *Node) serialize(mid, command string, args map[string]interface{}) ([]byte, error) {
	header, _ := json.Marshal(&MsgHeader{
		MessageId: mid,
		Command: command,
		Format: "json"})
	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s\n%s", header, body)), nil
}

func (n *Node) deserialize(data string) (*Msg, error) {
	parts := strings.SplitN(data, "\n", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid data received: %s", data)
	}
	header, body := parts[0], parts[1]
	fmt.Println("BODY",body)
	var header_info MsgHeader
	json.Unmarshal([]byte(header), &header_info)
	fmt.Println("HEADER",header_info)
	if header_info.Format == "json" {
		var body_info map[string]interface{}
		json.Unmarshal([]byte(body), &body_info)
		return &Msg{
			MessageId: header_info.MessageId,
			Command: header_info.Command,
			Data: body_info}, nil
	}
	return nil, fmt.Errorf("Failed to parse message: Invalid format: %s", header_info.Format)
}

func (n *Node) run() {
	incoming_handler := func(zmq.State) error {
		data, err := n.socket.Recv(0)
		if err != nil {
			n.handle_error(fmt.Errorf("Problem receiving on socket: %v",err))
		}
		fmt.Println("[LIBMANGO] RECVED",data)

		msg, err := n.deserialize(data)
		if err != nil {
			n.handle_error(fmt.Errorf("Failed to deserialize: %s",data))
		}

		if handler, ok := n.outstanding[msg.MessageId]; ok {
			fmt.Println("[LIBMANGO] GOT REPLY",msg.MessageId,handler)
			handler(msg.Command, msg.Data)
			delete(n.outstanding, msg.MessageId)
		} else if handler, ok := n.handlers[msg.Command]; ok {
			// Handle other registered commands
			reply_cmd, reply_args, reply_err := handler(msg.Data)
			if reply_err != nil {
				n.handle_error(reply_err)
			}
			if len(reply_cmd) > 0 {
				n.send_reply(msg.MessageId, reply_cmd, reply_args)
			}
			
		} else if n.DefaultHandler != nil {
			reply_cmd, reply_args, reply_err := n.DefaultHandler(msg.Command, msg.Data)
			if reply_err != nil {
				n.handle_error(reply_err)
			}
			if len(reply_cmd) > 0 {
				n.send_reply(msg.MessageId, reply_cmd, reply_args)
			}
			
		}
		return nil
	}
	outgoing_handler := func(data interface{}) error {
		fmt.Println("[LIBMANGO] SENDING: ",data.(string))
		n.socket.Send(data.(string), 0)
		return nil
	}
	reactor_done_handler := func(data interface{}) error {
		return fmt.Errorf("done")
	}

	reactor := zmq.NewReactor()
	reactor.AddChannel(n.done_reactor, 1000, reactor_done_handler)
	reactor.AddChannel(n.outgoing, 1000, outgoing_handler)
	reactor.AddSocket(n.socket, zmq.POLLIN, incoming_handler)
	reactor.Run(time.Second)
	fmt.Println("EXITING")
}

func (n *Node) heartbeat_worker() {
	ticker := time.NewTicker(10*time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("[LIBMANGO] heartbeat")
			n.Send("alive", map[string]interface{}{})
		case <- n.done_internal:
			// Stop the reactor
			n.done_reactor <- true
			
			// Signal anyone else who cares that we are done
			n.Done <- true
			break
		}
	}
}

func (n *Node) default_default_handler(cmd string, args map[string]interface{}) (string, map[string]interface{}, error) {
	n.handle_error(fmt.Errorf("No handler for command: %s", cmd))
	return "", nil, nil
}

func (n *Node) heartbeat(args map[string]interface{}) (string, map[string]interface{}, error) {
	return "alive", map[string]interface{}{}, nil
}

func (n *Node) exit(args map[string]interface{}) (string, map[string]interface{}, error) {
	n.done_internal <- true
	return "", nil, nil
}
