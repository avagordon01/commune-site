package main

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/gob"
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
	err              error
	template_base    *template.Template
	template_home    *template.Template
	template_post    *template.Template
	template_search  *template.Template
	template_preview *template.Template
	text_index       bleve.Index
	database         *storm.DB
)

func main() {
    /*
	text_index, err = bleve.Open("database/search.bleve")
	if err != nil {
		log.Fatal(err)
	}
	defer text_index.Close()
    */

	database, err = storm.Open("database/database.bolt", storm.Codec(gob.Codec))
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

    test()
    return

	func_map := template.FuncMap{"human_time": humanize.Time}
	template_base = template.Must(template.ParseFiles("templates/base.html")).Funcs(func_map)
	template_home = template.Must(template.Must(template_base.Clone()).ParseFiles("templates/home.html"))
	template_post = template.Must(template.Must(template_base.Clone()).ParseFiles("templates/post.html"))
	template_search = template.Must(template.Must(template_base.Clone()).ParseFiles("templates/search.html"))
	template_preview = template.Must(template.Must(template_base.Clone()).ParseFiles("templates/preview.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", hsts(fresh_cookie(home)))
	mux.HandleFunc("/post/", hsts(fresh_cookie(post)))
	mux.HandleFunc("/search/", hsts(fresh_cookie(search)))
    /*
	mux.HandleFunc("/preview", hsts(user_cookie(preview)))
	mux.HandleFunc("/submit_post", hsts(user_cookie(submit_post)))
	mux.HandleFunc("/submit_comment", hsts(user_cookie(submit_comment)))
    */
	mux.Handle("/static/", http.FileServer(http.Dir("./")))

	log.Println("server ready")
	go http.ListenAndServe(":80", http.HandlerFunc(https_redirect))
	go http.Serve(autocert.NewListener("commune.is"), mux)
	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)
	<-close
}
