package main

import (
	"fmt"
	"os/exec"
	"libmango"
)

func Notify(args map[string]interface{}) (string, map[string]interface{}, error) {
	u := args["urgency"].(string)
	if u == "urgent" || u == "high" || u == "important" {
		u = "critical"
	} else if u == "non-urgent" || u == "low" || u == "unimportant" {
		u = "low"
	} else {
		u = "normal"
	}
	c := exec.Command("notify-send", args["title"].(string), args["message"].(string), "-u", u)
	c.Run()
	return "",nil,nil
}

func main() {
	n, err := libmango.NewNode("notify",map[string]libmango.MangoHandler{"notify":Notify})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
}
