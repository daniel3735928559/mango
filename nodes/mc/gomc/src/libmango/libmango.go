package libmango

import (
	"os"
	"fmt"
	"sync"
	"strings"
	"math/rand"
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
)

type MangoHandler func(map[string]interface{}) (string, map[string]interface{}, error)
type MangoReplyHandler func(string, map[string]interface{})

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
	status NodeStatus
	cookie string
	currentMid int
	handlers map[string]MangoHandler
	outstanding map[string]MangoReplyHandler
	send_mutex sync.Mutex
	socket *zmq.Socket
}

func NewNode(name string, handlers map[string]MangoHandler) (*Node, error) {
	n := &Node{
		status: NODE_STATUS_STOPPED,
		currentMid: 0,
		Name: name,
		handlers: handlers,
		send_mutex: sync.Mutex{},
		outstanding: make(map[string]MangoReplyHandler)}
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
	
	n.socket.Send("",zmq.SNDMORE)
	n.socket.Send(fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes)), 0)
	fmt.Println("[LIBMANGO] sent")
	return nil
}

func (n *Node) SendForReply(cmd string, args map[string]interface{}, reply_handler MangoReplyHandler) (string, error) {
	n.send_mutex.Lock()
	defer n.send_mutex.Unlock()
	fmt.Println("sending",cmd,args)
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
	
	n.socket.Send("",zmq.SNDMORE)
	n.socket.Send(fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes)), 0)
	return h.MessageId, nil
}

func (n *Node) send_reply(mid, cmd string, args map[string]interface{}) {
	n.send_mutex.Lock()
	defer n.send_mutex.Unlock()
	fmt.Println("sending reply",cmd,args)
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
		n.socket.Send("",zmq.SNDMORE)
		n.socket.Send(fmt.Sprintf("%s\n%s",string(header_bytes), string(data_bytes)), 0)
	}
}

func (n *Node) Start() {
	fmt.Println("[LIBMANGO] STARTING", n.Server)
	n.socket, _ = zmq.NewSocket(zmq.DEALER)
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	fmt.Println(identity)
	//n.socket.SetIdentity(identity)
	n.socket.Connect(n.Server)
	n.status = NODE_STATUS_RUNNING
	go n.run()
}

func (n *Node) Stop() {
	n.status = NODE_STATUS_STOPPED
}

func (n *Node) handle_error(err error) {
	fmt.Printf("ERROR: %v", err)
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
	for n.status == NODE_STATUS_RUNNING {
		data, err := n.socket.Recv(0)
		if err != nil {
			n.handle_error(fmt.Errorf("Problem receiving on socket: %v",err))
		}
		fmt.Println("RECVED",data)

		msg, err := n.deserialize(data)
		if err != nil {
			n.handle_error(fmt.Errorf("Failed to deserialize: %s",data))
		}

		if msg.Command == "reply" {
			// Handle replies
			if handler, ok := n.outstanding[msg.MessageId]; ok {
				handler(msg.MessageId, msg.Data)
			} else {
				n.handle_error(fmt.Errorf("Unexpected reply for message id %s", msg.MessageId))
			}
		} else if handler, ok := n.handlers[msg.Command]; ok {
			// Handle other registered commands
			reply_cmd, reply_args, reply_err := handler(msg.Data)
			if reply_err != nil {
				n.handle_error(reply_err)
			}
			if len(reply_cmd) > 0 {
				n.send_reply(msg.MessageId, reply_cmd, reply_args)
			}
			
		} else {
			n.handle_error(fmt.Errorf("No handler for command: %s", msg.Command))
		}
		
	}
}

func (n *Node) heartbeat(args map[string]interface{}) (string, map[string]interface{}, error) {
	return "alive", map[string]interface{}{}, nil
}

func (n *Node) exit(args map[string]interface{}) (string, map[string]interface{}, error) {
	os.Exit(688)
	return "", nil, nil
}
