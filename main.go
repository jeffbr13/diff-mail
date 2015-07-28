package main

import (
	"container/ring"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/aryann/difflib"
)

// Global state, because hey, they seem to like that in Go!
var webAddress string
var emailAddress string
var scrapeCacheRing *ring.Ring

func init() {
	flag.StringVar(&webAddress, "url", "http://www.example.com", "URL to watch for changes")
	flag.StringVar(&emailAddress, "email", "mail@benjeffrey.com", "email address to send diffs")
	flag.Parse()
}

func main() {
	log.Printf("Send %s hourly diffs of %s", emailAddress, webAddress)
	store := newScrapeStore()

	c := time.Tick(1 * time.Hour)
	for now := range c {
		fmt.Printf("Scraping at %v\n", now)
		err := scrape(store)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Emailing %s with diff results.\n", emailAddress)
		err = emailDiff(store)
		if err != nil {
			log.Println(err)
		}
	}
}

// Scrapes webpage and puts the response body in the cache.
func scrape(store *scrapeStore) error {
	// get webpage
	resp, err := http.Get(webAddress)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%d when accessing URL %s", resp.StatusCode, webAddress)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	store.add(body)
	return nil
}

// Email the diff of the latest two scrapes in the cache to the notification address.
func emailDiff(store *scrapeStore) error {
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	to := []string{emailAddress}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: diff-mail for " + webAddress + "\n"

	diff, err := store.htmlDiffPrev()
	if err != nil {
		return err
	}
	msg := mime + subject + "<html><body><h1>Diff between last two scrapes:</h1>" + diff + "</body></html>"
	err = smtp.SendMail(os.Getenv("SMTP_HOST")+":25", auth, os.Getenv("SMTP_USERNAME"), to, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func bytesToStringsOnNewline(data []byte) []string {
	return strings.Split(html.EscapeString(string(data)), "\n")
}

type scrapeStore struct {
	*ring.Ring
}

// NewScrapeStore constructs a scrapeStore with 24 slots.
func newScrapeStore() *scrapeStore {
	store := new(scrapeStore)
	store.Ring = ring.New(24)
	return store
}

func (store *scrapeStore) add(data []byte) {
	if store.Ring.Value != nil {
		store.Ring = store.Ring.Next()
	}

	store.Ring.Value = data
}

func (store *scrapeStore) current() []byte {
	if store.Ring.Value != nil {
		return store.Ring.Value.([]byte)
	}
	return nil
}

func (store *scrapeStore) prev() []byte {
	if store.Ring.Prev().Value != nil {
		return store.Ring.Prev().Value.([]byte)
	}
	return nil
}

func (store *scrapeStore) htmlDiffPrev() (string, error) {
	if store.prev() != nil {
		return "<table>" + difflib.HTMLDiff(bytesToStringsOnNewline(store.prev()), bytesToStringsOnNewline(store.current())) + "</html>", nil
	}
	return "", fmt.Errorf("Can't generate diff with only one scrape.")
}
