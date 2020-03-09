package transport

import (

)


type MangoTransport interface {
	RunServer()
	Tx(string, []byte)
}


type WrappedMessage struct {
	Transport MangoTransport
	Identity string
	Data []byte
}
