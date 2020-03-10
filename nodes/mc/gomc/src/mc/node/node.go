package node

import (
	"fmt"
	"os/exec"
	"time"
	"math/rand"
	"mc/transport"
	"github.com/google/shlex"
)

type NodeStatus int

const (
	NODE_STATUS_READY NodeStatus = iota+1
	NODE_STATUS_RUNNING
	NODE_STATUS_UNRESPONSIVE
	NODE_STATUS_DEAD
)

type Node struct {
	Id string
	Name string
	Group string
	NodeType string
	Transport transport.MangoTransport
	Command string
	LastHeartbeat int64
	Status NodeStatus
	Proc *exec.Cmd
}

func MakeNode(name, group, typename, command string, trans transport.MangoTransport) *Node {
	id := ""
	for i := 0; i < 4; i++ {
		id += fmt.Sprintf("%016x", rand.Uint64())
	}
	return &Node{
		Id: id,
		Name: name,
		Group: group,
		NodeType: typename,
		Command: command,
		Status: NODE_STATUS_READY,
		Transport: trans}
}

// func (n *Node) Handler(src, mid, command string, args map[string]interface{}) {
// 	data := serializer.MakeMessage(src, mid, command, args).Serialize()
// 	n.Transport.Tx(dst, []byte(data))
// }

func (n *Node) RestartWorker() {
	n.Proc.Wait()
	fmt.Println("RESTARTING")
	n.Start()
}

func (n *Node) Start() {
	args, err := shlex.Split(n.Command)
	if err != nil {
		fmt.Println("ERROR starting command",n.Command,err)
	}
	n.Proc = exec.Command(args[0], args[1:]...)
	n.Proc.Start()
	n.Status = NODE_STATUS_RUNNING
}

func (n *Node) Heartbeat() {
	n.LastHeartbeat = time.Now().UnixNano()
}

func (n *Node) HeartbeatFailed() {
	n.Status = NODE_STATUS_UNRESPONSIVE
}

func (n *Node) Kill() {
	n.Status = NODE_STATUS_DEAD
}

func (n *Node) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}

