package node

import (
	"mc/serializer"
	"mc/transport"
)

type NodeStatus int

const (
	NODE_STATUS_READY NodeStatus = iota+1
	NODE_STATUS_RUNNING
	NODE_STATUS_UNRESPONSIVE
	NODE_STATUS_DEAD
)

type Node interface {
	GotAlive(identity string, transport transport.MangoTransport)
	LastSeen() string
	GetId() string
	GetGroup() string
	GetName() string
	GetType() string
	ToString() string
	SendToNode(serializer.Msg) error
}
