package main

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/SherClockHolmes/webpush-go"
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

var templates = template.Must(template.ParseGlob("templates/*.html"))

var (
	subscriptions      = make(map[string]webpush.Subscription)
	subscriptionsMutex sync.Mutex
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	p := NewPage()
	p.Title = "Home"
	p.Body = template.HTML("<h1>Welcome to the home page</h1>")
	if err := templates.ExecuteTemplate(w, "_layout.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
