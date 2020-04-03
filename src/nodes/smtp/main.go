package main

import (
	"fmt"
	"log"
	"strconv"
	"libmango"
	"net/smtp"
	"github.com/docopt/docopt-go"
)

type SMTPAccount struct {
	username string
	password string
	port int
	server string
	default_from string
}

func (a *SMTPAccount) Send(args map[string]interface{}) (string, map[string]interface{}, error) {
	from := args["from"].(string)
	if len(from) == 0 {
		from = a.default_from
	}
	subject := args["subject"].(string)
	body := args["body"].(string)
	to := args["to"].(string)
	subject_body := fmt.Sprintf("Subject: %s\n\n%s", subject, body)
	status := smtp.SendMail(
		fmt.Sprintf("%s:%d",a.server,a.port),
		smtp.PlainAuth("", a.username, a.password, a.server),
		from,
		[]string{to},
		[]byte(subject_body))
	if status != nil {
		return "error",map[string]interface{}{"message":fmt.Sprintf("Error from SMTP Server: %s", status)}, nil
	}
	log.Print("Email Sent Successfully")
	return "", nil, nil
}

func main() {
	usage := `Usage: smtp <server> <username> <password> <default_from> [<port>]`
	args, err := docopt.ParseDoc(usage)
	srv := args["<server>"].(string)
	un := args["<username>"].(string)
	pw := args["<password>"].(string)
	df := args["<default_from>"].(string)
	port := 587
	if v, ok := args["<port>"]; ok && v != nil {
		port, _ = strconv.Atoi(v.(string))
	}

	a := &SMTPAccount{
		username: un,
		password: pw,
		port: port,
		server: srv,
		default_from: df}
	
	n, err := libmango.NewNode("sendemail",map[string]libmango.MangoHandler{"send":a.Send})
	if err != nil {
		fmt.Println(err)
		return
	}
	n.Start()
}
