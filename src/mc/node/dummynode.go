package node

import (
	"fmt"
	"mc/serializer"
	"mc/transport"
)

type DummyNode struct {
	Identity string
	Group string
	Name string
	OutputChannel chan transport.WrappedMessage
}

func MakeDummyNode(identity, group, name string, ch chan transport.WrappedMessage) *DummyNode {
	return &DummyNode {
		Identity: identity,
		Group: group,
		Name: name,
		OutputChannel: ch}
}

func (n *DummyNode) GotAlive(identity string, transport transport.MangoTransport) {
	n.Identity = identity
}

func (n *DummyNode) LastSeen() string {
	return "now"
}

func (n *DummyNode) GetId() string {
	return n.Identity
}

func (n *DummyNode) GetGroup() string {
	return n.Group
}
func (n *DummyNode) GetName() string {
	return n.Name
}
func (n *DummyNode) GetType() string {
	return "dummy"
}
func (n *DummyNode) ToString() string {
	return fmt.Sprintf("%s/%s", n.Group, n.Name)
}
func (n *DummyNode) SendToNode(m serializer.Msg) error {
	m.Cookie = n.Identity
	wmsg := transport.WrappedMessage {
		Identity: n.Identity,
		Transport: nil,
		Message: m}
	n.OutputChannel <- wmsg
	return nil
}
