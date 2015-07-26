package main

import (
	"flag"
	"log"
)

var webAddress string
var emailAddress string

func init() {
	flag.StringVar(&webAddress, "url", "http://www.example.com", "URL to watch")
	flag.StringVar(&emailAddress, "email", "me@example.com", "email address to send diffs")

	flag.Parse()
}

func main() {
	log.Printf("webAddress: %s", webAddress)
	log.Printf("emailAddress: %s", emailAddress)
}
