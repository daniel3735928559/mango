package serializer


import (
	"fmt"
	"strings"
	"errors"
	"encoding/json"
)

type MCMessage struct {
	Sender string
	MessageId string
	Command string
	Data map[string]interface{}
}

type MCHeader struct {
	Source string  `json:"source"`
	MessageId string  `json:"mid"`
	Command string `json:"command"`
	Format string  `json:"format"`
}

type MCTransport interface {
	RunServer(register func(*MCMessage, MCTransport) bool)
	Tx(string, []byte)
}

func (msg *MCMessage) Serialize() string {
	header, _ := json.Marshal(&MCHeader{Source: msg.Sender, MessageId: msg.MessageId, Command: msg.Command, Format: "json"})
	body, _ := json.Marshal(msg.Data)
	return fmt.Sprintf("%s\n%s", header, body)
}

func (msg *MCMessage) RawHeader() map[string]string {
	return map[string]string{"source": msg.Sender, "mid": msg.MessageId, "command": msg.Command}
}

func MakeMessage(src, mid, command string, args map[string]interface{}) *MCMessage {
	return &MCMessage {
		Sender: src,
		MessageId: mid,
		Command: command,
		Data: args}
}

func ParseMessage(data string) (*MCMessage, error) {
	parts := strings.SplitN(data, "\n", 2)
	if len(parts) < 2 {
		return nil, errors.New(fmt.Sprintf("Invalid data received: %s", data))
	}
	header, body := parts[0], parts[1]
	var header_info MCHeader
	json.Unmarshal([]byte(header), &header_info)
	fmt.Println("AA",header_info)
	if header_info.Format == "json" {
		var body_info map[string]interface{}
		json.Unmarshal([]byte(body), body_info)
		return &MCMessage{
			Sender:header_info.Source,
			MessageId: header_info.MessageId,
			Command: header_info.Command,
			Data: body_info}, nil
	}
	return nil, errors.New(fmt.Sprintf("Failed to parse message: Invalid format: %s", header_info.Format))
}
