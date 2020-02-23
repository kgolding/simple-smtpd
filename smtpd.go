package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/alash3al/go-smtpsrv"
)

var VERSION string
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

	debug("ssmtpd - a simple SMTP server to JSON output by Kevin Golding")
	debug("Version: %s", VERSION)

	srv := &smtpsrv.Server{
		Name:        serverName,
		Addr:        ":" + strconv.Itoa(port),
		MaxBodySize: int64(maxBody) * 1024,
		Handler:     handleSmtpRequest,
	}

	debug("Starting SMTP server on %s, accepting messages with a maximum body size of %d bytes", srv.Addr, srv.MaxBodySize)

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		printErrf("Unable to open port %d: %s", port, err)
	}

	err = srv.Serve(ln)
	printErrf(err.Error())
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
}

func handleSmtpRequest(req *smtpsrv.Request) error {
	result := Result{
		From: req.From,
	}

	t, err := mail.ParseAddress(req.From)
	if err == nil {
		result.From = t.Address
	}

	result.Subject = req.Message.Header.Get("Subject")

	debug("Received email '%s' from %s to %s", result.Subject, req.From, req.To)

	contentType, params, err := parseContentType(req.Message.Header.Get("Content-Type"))
	if err != nil {
		printErrf("Error getting content type: %s", err)
	}

	switch contentType {
	case contentTypeMultipartMixed:
		// email.TextBody, email.HTMLBody, email.Attachments, email.EmbeddedFiles, err = parseMultipartMixed(msg.Body, params["boundary"])
		result.Body, result.Html, _, _, err = parseMultipartMixed(req.Message.Body, params["boundary"])
	case contentTypeMultipartAlternative:
		result.Body, result.Html, _, err = parseMultipartAlternative(req.Message.Body, params["boundary"])
	case contentTypeTextPlain:
		message, _ := ioutil.ReadAll(req.Message.Body)
		result.Body = strings.TrimSuffix(string(message[:]), "\n")
	case contentTypeTextHtml:
		message, _ := ioutil.ReadAll(req.Message.Body)
		result.Html = strings.TrimSuffix(string(message[:]), "\n")
	default:
		b := []byte{}
		b, err = ioutil.ReadAll(req.Message.Body)
		if err == nil {
			result.Body = string(b)
			if strings.HasSuffix(result.Body, "\r\n.\r\n") {
				result.Body = result.Body[0 : len(result.Body)-5]
			}
		}
	}
	if err != nil {
		printErrf("Error decoding body: %s", err)
	}

	for _, To := range req.To {
		result.To = To
		b, _ := json.Marshal(result)
		fmt.Println(string(b))
	}

	return nil
}
