package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/mail"
	"os"

	"github.com/DusanKasan/parsemail"
	"github.com/kgolding/smtpd"
)

var VERBOSE bool

func debug(format string, args ...interface{}) {
	if VERBOSE {
		fmt.Printf(format+"\n", args...)
	}
}

func main() {
	var port int
	var maxBody int
	var serverName string

	flag.BoolVar(&VERBOSE, "v", false, "verbose logging")
	flag.IntVar(&port, "port", 25, "listening port")
	flag.IntVar(&maxBody, "max", 256, "max body size KB")
	flag.StringVar(&serverName, "name", "mail-server", "server name")

	flag.Parse()

	debug("smtpd - a simple SMTP server to JSON output by Kevin Golding")
	debug("Version: 1.1.0")

	debug("Starting SMTP server on port %d", port)
	err := smtpd.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mailHandler, serverName, "")
	if err != nil {
		printErrf("Unable to open port %d: %s", port, err)
	}
}

func printErrf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
}

type Result struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
	Html    string `json:"html,omitempty"`
	Error   string `json:"errors,omitempty"`
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	debug("Received mail from %s for %s with subject %s", from, to[0], subject)

	result := Result{
		Subject: subject,
	}

	t, err := mail.ParseAddress(from)
	if err == nil {
		result.From = t.Address
	}

	mail, err := parsemail.Parse(bytes.NewReader(data))
	if err != nil {
		result.Error = err.Error()
	} else {
		result.Body = mail.TextBody
		result.Html = mail.HTMLBody
	}

	for _, t := range to {
		result.To = t
		b, _ := json.Marshal(result)
		fmt.Println(string(b))
	}
}
