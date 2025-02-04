package main

import (
	"html/template"
	"net/http"
	"sync"
)

type Page struct {
	Title  string
	Body   template.HTML
	Params map[string]interface{}
	Site   Site
}

func NewPage() *Page {
	var p Page
	p.Params = make(map[string]interface{})
	p.Site = site
	return &p
}

var templates = template.Must(template.ParseFiles("templates/_layout.html", "templates/index.html"))

var (
	subscriptions      = make(map[string]Subscription)
	subscriptionsMutex sync.Mutex
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	p := NewPage()
	p.Title = "Home"
	if err := templates.ExecuteTemplate(w, "_layout.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
