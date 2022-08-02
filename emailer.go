package main

import (
	"fmt"
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
	return fmt.Sprintf(`
From: %s
To: %s
Subject: %s
CC: %s
BCC: %s
%s
`, e.From, to, e.Subject, cc, bcc, e.Body)
}

func (e *Email) toBytes() []byte {
	return []byte(e.toString())
}

func (e *Email) Send(pass, server string) error {
	fmt.Println(e.toString())

	return nil
}
