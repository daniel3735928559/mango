package transport

import (
	mprotocol "libmango/protocol"
)

type MangoTransport interface {
	RunServer(register func(*mprotocol.MangoMessage, MangoTransport) bool)
	Tx(string, []byte)
}
