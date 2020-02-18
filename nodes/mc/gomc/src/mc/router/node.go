package router

type Node struct {
	Name string
	Send func(map[string]string, map[string]interface{})
}

