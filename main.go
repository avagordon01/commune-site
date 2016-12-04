package main

import (
    "os"
    "log"
    "math"
    "time"
    "sort"
    "net/http"
    "html/template"
    "encoding/json"
)

type Post struct {
    Id uint64
    Title string
    Time time.Time
    Votes uint64
    Username string
    Html template.HTML
    Comments []Comment
}

type Comment struct {
    Id uint64
    Time time.Time
    Votes uint64
    Username string
    Html template.HTML
    Comments []Comment
}

type Page struct {
    After uint64
    Title template.HTML
    Content template.HTML
}

type (
    Freshness0 []uint64
    Freshness1 []uint64
    Freshness2 []uint64
    Freshness3 []uint64
    Freshness4 []uint64
    Freshness5 []uint64
    Freshness6 []uint64
)

var (
    templates *template.Template
    err error
    posts []Post
    index [7][]uint64
    user_counter uint64
    posts_encoder json.Encoder
)

func value(freshness float64, post Post) float64 {
    return float64(post.Votes) * math.Pow(0.75, freshness * float64(10 - post.Time.Unix()))
}

func init() {
    f, err := os.Open("posts.json")
    if err != nil {
        log.Fatal(err)
    }
    err = json.NewDecoder(f).Decode(&posts)
    if err != nil {
        log.Fatal(err)
    }
    posts_encoder = *json.NewEncoder(f)

    for i := 0; i < 7; i++ {
        index[i] = make([]uint64, len(posts))
        for j := 0; j < len(posts); j++ {
            index[i][uint64(j)] = uint64(j)
        }
    }
    sort.Stable(Freshness0(index[0]))
    sort.Stable(Freshness1(index[1]))
    sort.Stable(Freshness2(index[2]))
    sort.Stable(Freshness3(index[3]))
    sort.Stable(Freshness4(index[4]))
    sort.Stable(Freshness5(index[5]))
    sort.Stable(Freshness6(index[6]))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/",                 hsts(fresh_cookie(home)))
    mux.HandleFunc("/search",           hsts(fresh_cookie(search)))
    mux.HandleFunc("/post/",            hsts(fresh_cookie(post)))
    mux.HandleFunc("/submit_post",      hsts(user_cookie(submit_post)))
    mux.HandleFunc("/submit_comment",   hsts(user_cookie(submit_comment)))
    mux.HandleFunc("/submit_upvote",    hsts(user_cookie(submit_upvote)))
    mux.Handle("/static/",              http.FileServer(http.Dir("./")))

    err = http.ListenAndServeTLS(":443", "cert.pem", "privkey.pem", mux)
    if err != nil {
        log.Fatal(err)
    }
}
