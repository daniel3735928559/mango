package serializer


import (
	"fmt"
	"strings"
	"errors"
	"encoding/json"
)

type Msg struct {
	Sender string
	MessageId string
	Command string
	Cookie string
	Data map[string]interface{}
}

type MsgHeader struct {
	Source string  `json:"source"`
	MessageId string  `json:"mid"`
	Command string `json:"command"`
	Cookie string `json:"cookie,omitempty"`
	Format string  `json:"format"`
}

func Serialize(sender, mid, command string, payload interface{}) ([]byte, error) {
	header, _ := json.Marshal(&MsgHeader{
		Source: sender,
		MessageId: mid,
		Command: command,
		Format: "json"})
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s\n%s", header, body)), nil
}

func Deserialize(data string) (*Msg, error) {
	parts := strings.SplitN(data, "\n", 2)
	if len(parts) < 2 {
		return nil, errors.New(fmt.Sprintf("Invalid data received: %s", data))
	}
	header, body := parts[0], parts[1]
	fmt.Println("BODY",body)
	var header_info MsgHeader
	json.Unmarshal([]byte(header), &header_info)
	fmt.Println("AA",header_info)
	if header_info.Format == "json" {
		var body_info map[string]interface{}
		json.Unmarshal([]byte(body), &body_info)
		return &Msg{
			Sender:header_info.Source,
			MessageId: header_info.MessageId,
			Command: header_info.Command,
			Cookie: header_info.Cookie,
			Data: body_info}, nil
	}
	return nil, errors.New(fmt.Sprintf("Failed to parse message: Invalid format: %s", header_info.Format))
}
