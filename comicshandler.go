package main

import (
	"net/http"
	"log"
	"fmt"
	"regexp"
	"sync"
	"html/template"
	"time"
	"strconv"
)

func comicsHandler(w http.ResponseWriter, r *http.Request) {
	prev := ""
	next := ""
	timedate := time.Now()

	c := Comics {
		Comic{ 
			Title:	"Dilbert",
			URL:	"http://www.dilbert.com",
			Regex: 	`http(s)?://assets\.amuniversal\.com/[\w\d]+`,
			Nav:	func(t time.Time) string { return t.Format("http://www.dilbert.com/strip/2006-01-02") },
		},
		Comic{ 
			Title:	"Pearls before swine",
			URL:	"http://www.gocomics.com/pearlsbeforeswine",
			Regex: 	`http(s)?://assets\.amuniversal\.com/[\w\d]+`,
			Nav:	func(t time.Time) string { return t.Format("http://www.gocomics.com/pearlsbeforeswine/2006/01/02") },
		},
		Comic{ 
			Title:	"Wizard of ID",
			URL: 	"http://www.gocomics.com/wizardofid",
			Regex:	`http(s)?://assets\.amuniversal\.com/[\w\d]+`,
			Nav:	func(t time.Time) string { return t.Format("http://www.gocomics.com/wizardofid/2006/01/02") },
		},
		Comic{ 
			Title:	"Fokke & Sukke",
			URL: 	"http://www.foksuk.nl",
			Regex:	`content/formfield_files/formcartoon_[\w\d_]+\.(jpg|gif)`,
			Nav:	func(t time.Time) string { return fmt.Sprintf("http://www.foksuk.nl/?ctime=%d", t.Unix()) },
		},
	}

	fullpath := r.URL.Path

	rex, err := regexp.Compile(`(\d{4})(\d\d)(\d\d)`)
	if err != nil {
		log.Fatal(err)
	}

	m := rex.FindStringSubmatch(fullpath)
	if m != nil {
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])

		timedate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	}
	if timedate.Add(24*time.Hour).After(time.Now()) {
		next = time.Now().Format("20060102")
	} else {
		next = timedate.Add(24*time.Hour).Format("20060102")
	}

	prev = timedate.Add(-24*time.Hour).Format("20060102")

	tmpl := Template {
		Comics: &c,
		Nav: Navigation{
			Prev: prev,
			Next: next,
		},
	}

	var wg sync.WaitGroup
	// range returns a copy of the elements.
	// We need the actual element here since we want to modify it (set the Image)
	for i, _ := range c {
		wg.Add(1)
		go func(c *Comic) {
			defer wg.Done()
			c.Fetch(&timedate)
		}(&c[i])
	}
	wg.Wait()
	t, _ := template.ParseFiles("comics.html")
	t.Execute(w, tmpl) 
}
