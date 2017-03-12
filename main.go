package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"time"
    "syscall"
    "os/signal"
)

type Post struct {
	Id           uint64
	Title        string
    Snippet      string
	Time         time.Time
	Votes        uint64
	Username     string
	Html         template.HTML
	CommentCount uint64
	Comments     []Comment
}

type Comment struct {
	Id       uint64
	Time     time.Time
	Votes    uint64
	Username string
	Html     template.HTML
	Comments []Comment
}

type Page struct {
	After   uint64
	Title   template.HTML
	Content template.HTML
}

type Users struct {
	user_counter uint64
	page_counter uint64
}

type Names struct {
	Salt       uint64
	Animals    []string
	Colours    []string
	Adjectives []string
}

var (
	err   error
	posts []Post
	index [5][]uint64
	users Users
	names Names
)

func value(freshness float64, post Post) float64 {
	return float64(post.Votes) * math.Pow(0.75, freshness * float64(10 - post.Time.Unix()))
}

func compare(freshness float64) func(i, j int) bool {
    return func(i, j int) bool {
        return value(freshness, posts[i]) < value(freshness, posts[j])
    }
}

func main() {
	f, err := os.Open("posts.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(f).Decode(&posts)
	if err != nil {
		log.Fatal(err)
	}
    f.Close()

    f, err = os.Open("names.json")
    if err != nil {
        log.Fatal(err)
    }
    err = json.NewDecoder(f).Decode(&names)
    if err != nil {
        log.Fatal(err)
    }
    f.Close()

	for i := 0; i < len(index); i++ {
		index[i] = make([]uint64, len(posts))
		for j := 0; j < len(posts); j++ {
			index[i][uint64(j)] = uint64(j)
		}
	}
	sort.SliceStable(index[0], compare(0.0))
	sort.SliceStable(index[1], compare(0.1))
	sort.SliceStable(index[2], compare(0.2))
	sort.SliceStable(index[3], compare(0.5))
	sort.SliceStable(index[4], compare(1.0))

	mux := http.NewServeMux()
	mux.HandleFunc("/",                 hsts(fresh_cookie(home)))
	mux.HandleFunc("/search",           hsts(fresh_cookie(search)))
	mux.HandleFunc("/post/",            hsts(fresh_cookie(post)))
	mux.HandleFunc("/submit_post",      hsts(user_cookie(submit_post)))
	mux.HandleFunc("/submit_comment",   hsts(user_cookie(submit_comment)))
	//mux.HandleFunc("/submit_upvote",    hsts(user_cookie(submit_upvote)))
	mux.Handle("/static/",              http.FileServer(http.Dir("./")))

    srv := &http.Server{Addr: ":443", Handler: mux}
    go func() {
        err = srv.ListenAndServeTLS("cert.pem", "privkey.pem")
        if err != nil {
            log.Fatal(err)
        }
    }()

    close := make(chan os.Signal)
    signal.Notify(close, os.Interrupt, syscall.SIGTERM)
    <-close
    srv.Close()

	f, err = os.Open("posts.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewEncoder(f).Encode(&posts)
	if err != nil {
		log.Fatal(err)
	}
    log.Println(posts)
    log.Println("graceful")
}
