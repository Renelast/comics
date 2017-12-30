package main

import (
	"net/http"
	"log"
	"fmt"
	"io"
	"golang.org/x/net/html"
	"regexp"
	"errors"
	"net/url"
	"time"
)

type Comic struct {
	Title	string
	URL	string
	Regex	string
	Image   string
	Nav	func(time.Time) string
}

type Comics []Comic

func (c *Comic) GetUrl(u string) (io.ReadCloser, error) {
	uri, err := c.UrlParse(u)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	resp, err := http.Get(uri.String())
        if err != nil {
		return nil, err
        }
	if (resp.StatusCode == http.StatusOK) {
		return resp.Body, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%i", resp.StatusCode))
	}
	return nil, errors.New("Could not fetch url")
}

func (c *Comic) Fetch(t *time.Time) {
	uri := c.Nav(*t)
	r, err := c.GetUrl(uri)
	if err != nil {
		log.Print(uri)
		log.Fatal(err)
	}
	defer r.Close()
	// Parse the retunrd html doc into html.Node
	n, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	// Check if the html.Node elements contain the image we are looking for
	ref := c.Parse(n)
	// Get the image if found.
	if (ref != "") {
		refuri, err := c.UrlParse(ref)
		if err != nil {
			log.Fatal(err)
		}
		c.Image = refuri.String()
	}
}

func (c *Comic) UrlParse(u string) (*url.URL, error) {
	uri, err := url.Parse(u)
	if err != nil {
		log.Print("Parse error")
		return nil, err
	}

	if (uri.IsAbs() == false) {
		base, _ := uri.Parse(c.URL)
		uri = base.ResolveReference(uri)
	}
	return uri, nil
}

func (c *Comic) Parse(n *html.Node) string {
	rex, err := regexp.Compile(c.Regex)
	if err != nil {
		log.Fatal(err)
	}

	var ret string

	var f func(*html.Node)

	f = func(n *html.Node){
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" && rex.Match([]byte(a.Val)) {
					ret = a.Val
					break
				}
			}
		}
		for ch := n.FirstChild; ch != nil && ret == ""; ch = ch.NextSibling {
			f(ch)
		}
	}
	f(n)

	return ret
}
