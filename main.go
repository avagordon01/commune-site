package main

import (
	"github.com/asdine/storm"
	"github.com/blevesearch/bleve"
	"github.com/dustin/go-humanize"
	"golang.org/x/crypto/acme/autocert"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	err        error
	templates  map[string]*template.Template
	text_index bleve.Index
	database   *storm.DB
)

func main() {
	text_index, err = bleve.Open("database/search.bleve")
	if err != nil {
		log.Fatal(err)
	}
	defer text_index.Close()

	database, err = storm.Open("database/database.bolt")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	templates = make(map[string]*template.Template)
	func_map := template.FuncMap{"human_time": humanize.Time}
	templates["base"] = template.Must(template.ParseFiles("templates/base.html")).Funcs(func_map)
	templates["home"] = template.Must(template.Must(templates["base"].Clone()).ParseFiles("templates/home.html"))
	templates["post"] = template.Must(template.Must(templates["base"].Clone()).ParseFiles("templates/post.html"))
	templates["search"] = template.Must(template.Must(templates["base"].Clone()).ParseFiles("templates/search.html"))
	templates["preview"] = template.Must(template.Must(templates["base"].Clone()).ParseFiles("templates/preview.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", hsts(fresh_cookie(home)))
	mux.HandleFunc("/post/", hsts(fresh_cookie(post)))
	mux.HandleFunc("/search/", hsts(fresh_cookie(search)))
	mux.HandleFunc("/preview", hsts(user_cookie(preview)))
	mux.HandleFunc("/submit_post", hsts(user_cookie(submit_post)))
	mux.HandleFunc("/submit_comment", hsts(user_cookie(submit_comment)))
	mux.Handle("/static/", http.FileServer(http.Dir("./")))

	log.Println("server ready")
	go http.ListenAndServe(":80", http.HandlerFunc(https_redirect))
	go http.Serve(autocert.NewListener("commune.is"), mux)
	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)
	<-close
}
