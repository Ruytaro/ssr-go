package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecat/go-libs/log"
)

var templates map[string]*template.Template

type Data struct {
	Title   string
	Message string
	Users   []User
}

type User struct {
	Name  string
	Email string
}

func init() {
	//if os.Getenv("APP_ENV") == "production" {
	loadTemplates()
	//}
}

func loadTemplates() {
	templates = make(map[string]*template.Template)

	base := template.Must(template.ParseFiles("views/base"))

	// Add other templates
	templateFiles, _ := filepath.Glob("views/*")
	for _, file := range templateFiles {
		if file == "views/base" {
			continue
		}
		name := strings.Split(filepath.Base(file), ".")[0]
		t := template.Must(template.Must(base.Clone()).ParseFiles(file))
		templates[name] = t
	}
	log.Debug(fmt.Sprintf("loaded %d templates", len(templates)))
}
func main() {
	// Parse templates at startup

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split((r.URL.Path[1:]), ".")[0]
		log.Debug(p)
		t, ok := templates[p]
		if !ok {
			t = templates["index"]
			p = "index"
		}
		err := t.ExecuteTemplate(w, p, Data{
			Title:   "Home Page",
			Message: "Welcome to net/http SSR!",
			//Users:   []User{{Name: "Admin", Email: "admin@domain.com"}, {Name: "Guest", Email: ""}},
		})
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Debug(fmt.Sprintf("rendered %s", p))
	})

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error(err.Error())
	}
}
