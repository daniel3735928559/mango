package node

import (
	"fmt"
	"math/rand"
	"mc/serializer"
	"mc/transport"
)

type CallbackNodeHandler func(serializer.Msg) error

type CallbackNode struct {
	Id string
	Name string
	Group string
	NodeType string
	Handler CallbackNodeHandler
}

func MakeCallbackNode(group, name, typename string, handler CallbackNodeHandler) *CallbackNode {
	id := ""
	for i := 0; i < 4; i++ {
		id += fmt.Sprintf("%016x", rand.Uint64())
	}
	return &CallbackNode{
		Id: id,
		Name: name,
		Group: group,
		NodeType: typename,
		Handler: handler}
}

func (n *CallbackNode) GetName() string {
	return n.Name
}
func (n *CallbackNode) GetId() string {
	return n.Id
}
func (n *CallbackNode) GetGroup() string {
	return n.Group
}
func (n *CallbackNode) GetType() string {
	return n.NodeType
}

func (n *CallbackNode) GotAlive(tid string, t transport.MangoTransport) {
}

func (n *CallbackNode) SecsAgo() int {
	return 0
}

func (n *CallbackNode) LastSeen() string {
	return "now"
}

func (n *CallbackNode) SendToNode(m serializer.Msg) error {
	return n.Handler(m)
}

func (n *CallbackNode) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}
