package main

import (
	"fmt"
	"time"
	"libmango"
)

func Excite(args map[string]interface{}) (string, map[string]interface{}, error) {
	return "excited",map[string]interface{}{"message":args["msg"].(string)+"!"},nil
}

func main() {
	n, err := libmango.NewNode("excite",map[string]libmango.MangoHandler{"excite":Excite})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
	for {
		time.Sleep(10*time.Second)
		fmt.Println("working...")
	}
}
