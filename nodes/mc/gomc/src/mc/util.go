package main

import (
	"fmt"
	"errors"
	"encoding/json"
)
func ParseMessage(header, body string) (*MCMessage, error) {
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
	return nil, errors.New("Failed to parse message: No valid format")
}
