package main

import (
	"fmt"
	"libmango"
)

func Excite(args map[string]interface{}) (string, map[string]interface{}, error) {
	fmt.Println("EXCITING",args)
	return "excited",map[string]interface{}{"message":args["message"].(string)+"!"},nil
}

func main() {
	n, err := libmango.NewNode("excite",map[string]libmango.MangoHandler{"excite":Excite})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
}
