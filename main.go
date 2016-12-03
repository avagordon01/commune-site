package commune

import (
    "os"
    "log"
    "math"
    "sort"
    "net/http"
    "html/template"
    "encoding/json"
)

type Post struct {
    Id string
    Title string
    Time int64
    Votes uint64
    Username string
    Html template.HTML
    Comments []Comment
}

type Comment struct {
    Id string
    Time int64
    Votes uint64
    Username string
    Html template.HTML
    Comments []Comment
}

type Page struct {
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
    posts map[uint64]Post
    index [7][]uint64
    user_counter uint64
    posts_encoder json.Encoder
)

func value(freshness float64, post Post) float64 {
    return float64(post.Votes) * math.Pow(0.75, freshness * float64(10 - post.Time))
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
    h := http.NewServeMux()
    h.HandleFunc("/",               hsts(fresh_cookie(home)))
    h.HandleFunc("/search",         hsts(fresh_cookie(search)))
    h.HandleFunc("/post/",          hsts(fresh_cookie(post)))
    h.HandleFunc("/submit_post",    hsts(user_cookie(submit_post)))
    h.HandleFunc("/submit_comment", hsts(user_cookie(submit_comment)))
    h.Handle("/static/",            http.FileServer(http.Dir("./")))

    err = http.ListenAndServeTLS(":443", "cert.pem", "privkey.pem", h)
    if err != nil {
        log.Fatal(err)
    }
}
