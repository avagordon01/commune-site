package main

import (
	"encoding/gob"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/dustin/go-humanize"
	"golang.org/x/crypto/acme/autocert"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Post struct {
	Id           uint64
	Title        string
	Snippet      string
	Time         time.Time
	Views        uint64
	Value        float64
	Username     string
	Html         template.HTML
	CommentCount uint64
	Comments     []Comment
}

type Comment struct {
	Id       uint64
	Time     time.Time
	Value    float64
	Username string
	Html     template.HTML
	Comments []Comment
}

type Topic struct {
	Id         uint64
	Value      float64
	Similarity float64
	Content    string
}

var (
	err          error
	posts        []Post
	index        [5][]uint64
	user_counter uint64
	templates    map[string]*template.Template
	text_index   bleve.Index
	db           *bolt.DB
)

const page_length uint64 = 50

func main() {
	gob.Register(Post{})
	gob.Register(Comment{})
	gob.Register(Topic{})
	db, err = bolt.Open("database/database.bolt", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		v, err := tx.CreateBucketIfNotExists([]byte("vars"))
		if err != nil {
			log.Println(err)
		}
		if v.Get([]byte("user_counter")) == nil {
			v.Put([]byte("user_counter"), enc_id(0))
		}
		if v.Get([]byte("page_counter")) == nil {
			v.Put([]byte("page_counter"), enc_id(0))
		}
		_, err = tx.CreateBucketIfNotExists([]byte("trending_topics"))
		if err != nil {
			log.Println(err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("similar_topics"))
		if err != nil {
			log.Println(err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("posts_comments"))
		if err != nil {
			log.Println(err)
		}
		f, err := tx.CreateBucketIfNotExists([]byte("freshness_index"))
		if err != nil {
			log.Println(err)
		}
		for i := uint64(0); i < 5; i++ {
			_, err := f.CreateBucketIfNotExists(enc_id(i))
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	text_index, err = bleve.Open("database/search.bleve")
	if err != nil {
		log.Fatal(err)
	}
	defer text_index.Close()

	templates = make(map[string]*template.Template)
	templates["base.html"] = template.Must(template.ParseFiles("templates/base.html")).Funcs(template.FuncMap{"human_time": humanize.Time})
	templates["home.html"] = template.Must(template.Must(templates["base.html"].Clone()).ParseFiles("templates/home.html"))
	templates["post.html"] = template.Must(template.Must(templates["base.html"].Clone()).ParseFiles("templates/post.html"))
	templates["search.html"] = template.Must(template.Must(templates["base.html"].Clone()).ParseFiles("templates/search.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", hsts(fresh_cookie(home)))
	mux.HandleFunc("/post/", hsts(fresh_cookie(post)))
	mux.HandleFunc("/search/", hsts(fresh_cookie(search)))
	mux.HandleFunc("/submit_post", hsts(user_cookie(submit_post)))
	mux.HandleFunc("/submit_comment", hsts(user_cookie(submit_comment)))
	mux.Handle("/static/", http.FileServer(http.Dir("./")))

	go http.ListenAndServe(":80", http.HandlerFunc(https_redirect))
	go http.Serve(autocert.NewListener("commune.is"), mux)
	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)
	<-close
}
