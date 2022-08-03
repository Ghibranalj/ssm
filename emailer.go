package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Email struct {
	To      []string
	From    string
	Subject string
	Body    string
	CC      []string
	BCC     []string
}

func (e *Email) toString() string {

	to := strings.Join(e.To, ",")
	cc := strings.Join(e.CC, ",")
	bcc := strings.Join(e.BCC, ",")
	return fmt.Sprintf(
		`From: %s
To: %s
Subject: %s
CC: %s
BCC: %s

%s`, e.From, to, e.Subject, cc, bcc, e.Body)
}

func (e *Email) toBytes() []byte {
	return []byte(e.toString())
}

func (e *Email) Send(pass, server, port string) error {
	s := fmt.Sprintf("%s:%s", server, port)
	auth := smtp.PlainAuth("", e.From, pass, server)
	err := smtp.SendMail(s, auth, e.From, e.To, e.toBytes())
	if err != nil {
		return err
	}
	return nil
}
