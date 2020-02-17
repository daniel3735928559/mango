package main

import (
	"fmt"
	// "strings"
	// "strconv"
	// "github.com/docopt/docopt-go"
	// "github.com/google/shlex"
)

// Each server is responsible for registration and mapping node

type MHeader struct {
	Cmd string
	MessageId string
	Format string
}

type MMsg struct {
	header *MHeader
	args map[string]interface{}
}

type MInterface struct {
	Commands map[string]func(*MHeader,interface{})(string,interface{},error) MMsg
}

type MNode struct {
	Version string
	NodeId string
	server string
	m_interface *MInterface
	
}

func (n *MNode) Run() {
	
}

func (n *MNode) Dispatch(header *MHeader, args interface{})   {
	if handler, ok := n.m_interface.Commands[header.Cmd]; ok {
		result_cmd, result, err := handler()
		if err != nil {
			n.HandleError(err)
		} else if result != nil {
			n.MSend(result_cmd, result, header.MessageId)
		}
	} else {
		fmt.Println("No handler for",header.Cmd)
		n.HandleError(errors.New(fmt.Sprintf("No handler for %s",header.Cmd)))
	}
}

func (n *MNode) HandleError(e error) {
	fmt.Println(e)
	n.MSend("error",map[string]string{"source":n.Name,"message":e.ToString()})
}

func (n *MNode) Heartbeat() {
	n.MSend("alive",map[string]interface{}{},header.MessageId)
}

func (n *MNode) MakeHeader(cmd, mid string) *MHeader {
	return &MHeader{
		Cmd: cmd,
		MessageId: mid}
}

func (n *MNode) Exit() {
	os.Exit()
}

func (n *MNode) MSend(name, msg, mid, cmd string) {
	fmt.Println("sending",name,msg,mid)
	header := n.MakeHeader(name,mid,cmd)
	
}

type MCMessage struct {
	Sender string
	MessageId string
	Command string
	Data interface{}
}

type MCHeader struct {
	Source string  `json:"source"`
	MessageId string  `json:"mid"`
	Command string `json:"command"`
	Format string  `json:"format"`
}

type MCNode struct {
	NodeId string
	GroupId string
	Transport *MCTransport
	CurrentMid int
	State int // alive, stalled, dead
}

type MCTransport interface {
	RunServer(register func(*MCMessage, MCTransport) bool)
	Tx(string, []byte)
}

func (n *MCNode) heartbeat_worker() {
	for {
		time.Sleep(5*time.Second)
		n.Transport.Tx("asda")
	}
}

func (n *MCNode) send_message(command string, args interface{}) {
	
}

func (n *MCNode) send_message(command string, args interface{}) {
	header := n.MakeHeader(command)
	body := JSON.Marsal(args)
	n.Transport.Send(header, body)
}

func (mc *MangoCommander) send_message(dest, command string, args interface{}) {
	mc.Nodes[dest].Send(command, )
}

func (mc *MangoCommander) reap(n *MCNode) {
	
}

func (mc *MangoCommander) reaper_worker() {	
	for {
		time.Sleep(5*time.Second)
		for _, n := range mc.Nodes {
			if n.State == 2 {
				mc.reap(n)
			}
		}
	}
}

type MangoCommander struct {
	zmqTransport *ZMQTransport
	socketTransport *SocketTransport
	MessageInput chan MCMessage
	Nodes map[string]*MCNode
	Routes map[string]map[string]string // node_id -> command_name -> dst node_id
}

func (mc *MangoCommander) Register(msg *MCMessage, t MCTransport) bool {
	fmt.Println("Register")
	// Peel off registration commands directly
	var group string
	var ok bool
	if group, ok = msg.Data.(map[string]interface{})["group"].(string); !ok {
		group = "root"
	}
	mc.Nodes[msg.Sender] = &MCNode{
		NodeId: msg.Sender,
		GroupId: group,
		Transport: &t}
	return true
}

func (mc *MangoCommander) Run() {
	go mc.zmqTransport.RunServer(mc.Register)
	go mc.socketTransport.RunServer(mc.Register)

	for msg := range mc.MessageInput {
		fmt.Println("FROM",msg.Sender)
		fmt.Println("DATA",msg.Data)
		if _, ok := mc.Nodes[msg.Sender]; ok {
			fmt.Println("Processing from",msg.Sender)
		} else {
			fmt.Println("No such node",msg.Sender)
		}
	}
}

func main() {
	fmt.Println("hi")
	go test_zmq_client(1919)
	go test_socket_client(1920)
	message_aggregator := make(chan MCMessage, 100)
	MC := &MangoCommander{
		zmqTransport: MakeZMQTransport(1919, message_aggregator),
		socketTransport: MakeSocketTransport(1920, message_aggregator),
		Nodes: make(map[string]*MCNode),
		MessageInput: message_aggregator}
	MC.Run()
}
