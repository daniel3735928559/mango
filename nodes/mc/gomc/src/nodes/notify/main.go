package main

import (
	"fmt"
	"time"
	"github.com/gen2brain/beeep"
	"libmango"
)

func Notify(args map[string]interface{}) (string, map[string]interface{}, error) {
	if args["beep"].(bool) {
		beeep.Notify(args["title"].(string),args["message"].(string),"")
		beeep.Beep(440.0, 1000)
	} else {
		beeep.Notify(args["title"].(string),args["message"].(string),"")
	}
	return "",nil,nil
}

func main() {
	n, err := libmango.NewNode("notify",map[string]libmango.MangoHandler{"notify":Notify})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
	for {
		time.Sleep(10*time.Second)
	}
}
