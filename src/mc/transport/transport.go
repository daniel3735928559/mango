package transport

import (
	"mc/serializer"
)


type MangoTransport interface {
	RunServer()
	Tx(string, serializer.Msg) error
}


type WrappedMessage struct {
	Transport MangoTransport
	Identity string
	Message serializer.Msg
}
