package main

import (
	"net/http"
)

type Navigation struct {
	Prev	string
	Next	string
}

type Template struct {
	*Comics
	Nav Navigation
}

func main() {
	http.HandleFunc("/", comicsHandler)
	http.ListenAndServe(":8080", nil)
}
