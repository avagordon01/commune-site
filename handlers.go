package commune

import (
    "log"
    "bytes"
    "time"
    "strconv"
    "strings"
    "net/http"
    "html/template"
    "golang.org/x/net/html"
    "github.com/microcosm-cc/bluemonday"
    "github.com/russross/blackfriday"
)

func hsts(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Strict-Transport-Security", "max-age=31536000")
        f(w, r)
    })
}

func user_cookie(f func(w http.ResponseWriter, r *http.Request, user_id uint64)) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("user_id")
        if err != nil {
            cookie = &http.Cookie{Name: "user_id", Value: strconv.FormatUint(user_counter, 10), Expires: time.Unix(1<<63-1, 0), Secure: true, HttpOnly: true}
            http.SetCookie(w, cookie)
        }
        user_id, err := strconv.ParseUint(cookie.Value, 10, 64)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError);
            return
        }
        f(w, r, user_id)
    })
}

func fresh_cookie(f func(w http.ResponseWriter, r *http.Request, freshness uint64)) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("freshness")
        if err != nil {
            cookie = &http.Cookie{Name: "freshness", Value: strconv.FormatUint(2, 10), Expires: time.Unix(1<<63-1, 0), Secure: true}
            http.SetCookie(w, cookie)
        }
        freshness, err := strconv.ParseUint(cookie.Value, 10, 64)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError);
            return
        }
        f(w, r, freshness)
    })
}

func home(w http.ResponseWriter, r *http.Request, freshness uint64) {
    templates, err = template.ParseGlob("templates.html");
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }

    var p []Post
    for i := 0; i < len(index[0]) && i < 20; i++ {
        p = append(p, posts[index[0][i]])
    }
    var content bytes.Buffer
    err = templates.ExecuteTemplate(&content, "posts", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
    page := Page{
        Title: template.HTML("commune"),
        Content: template.HTML(content.String()),
    }
    err = templates.ExecuteTemplate(w, "main", page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
}

func search(w http.ResponseWriter, r *http.Request, freshness uint64) {
    templates, err = template.ParseGlob("templates.html");
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }

    var p []Post
    after := uint64(0)
    for i := after; i < uint64(len(index[freshness])) && i < after + 20; i++ {
        p = append(p, posts[index[freshness][i]])
    }
    var content bytes.Buffer
    err = templates.ExecuteTemplate(&content, "posts", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
    page := Page{
        Title: template.HTML("commune"),
        Content: template.HTML(content.String()),
    }
    err = templates.ExecuteTemplate(w, "main", page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
}

func post(w http.ResponseWriter, r *http.Request, freshness uint64) {
    templates, err = template.ParseGlob("templates.html");
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }

    path := r.URL.Path[len("/path/"):]
    //post_id, err := base64.URLEncoding.DecodeString(path)
    post_id, err := strconv.ParseUint(path, 10, 32)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    post, ok := posts[post_id]
    if !ok {
        http.NotFound(w, r)
        return
    }
    var content bytes.Buffer
    err = templates.ExecuteTemplate(&content, "post", post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
    page := Page{
        Title: template.HTML("commune"),
        Content: template.HTML(content.String()),
    }
    err = templates.ExecuteTemplate(w, "main", page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
}

func submit_post(w http.ResponseWriter, r *http.Request, user_id uint64) {
    r.ParseMultipartForm(1<<20)
    post_text_s, ok := r.Form["textarea"]
    if !ok {
        log.Fatal("no text")
    }
    if len(post_text_s) == 0 {
        log.Fatal("empty text")
    }
    user_name := "user"
    markdown_raw := post_text_s[0]
    markdown_san := []byte(html.EscapeString(markdown_raw))
    html_raw := blackfriday.MarkdownCommon(markdown_san)
    html_san := string(bluemonday.UGCPolicy().SanitizeBytes(html_raw))

    //get title
    z := html.NewTokenizer(strings.NewReader(html_san))
    title := ""
    for {
        tt := z.Next()
        switch {
            case tt == html.ErrorToken:
                return
            case tt == html.StartTagToken:
                if z.Token().Data == "h1" {
                    z.Next()
                    title = z.Token().Data
                    return
                }
        }
    }
    if title == "" {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest);
        return
    }

    //create post
    post := Post{
        Title: title,
        Votes: 0,
        Username: user_name,
        Html: template.HTML(html_san),
        Comments: []Comment{},
    }
    log.Println(post)
    posts_encoder.Encode(posts)

    post_id := uint64(0)
    http.Redirect(w, r, "/post/" + strconv.FormatUint(post_id, 10), http.StatusSeeOther)
}

func submit_comment(w http.ResponseWriter, r *http.Request, user_id uint64) {
    r.ParseMultipartForm(1<<20)
    post_id_s, ok := r.Form["post_id"]
    if !ok {
        log.Fatal("no post_id")
    }
    if len(post_id_s) == 0 {
        log.Fatal("empty post_id")
    }
    post_id, err := strconv.ParseUint(post_id_s[0], 10, 64)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
    comment_id_s, ok := r.Form["comment_id"]
    if !ok {
        log.Fatal("no comment_id")
    }
    if len(comment_id_s) == 0 {
        log.Fatal("empty comment_id")
    }
    var comment_id uint64
    if comment_id_s[0] != "" {
        comment_id, err = strconv.ParseUint(comment_id_s[0], 10, 64)
    } else {
        comment_id = 0
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }
    log.Println(post_id)
    log.Println(comment_id)

    if strings.ToLower(r.Method) != "post" {
        http.NotFound(w, r)
        return
    }
    user_name := "user"
    markdown_raw := "empty"
    markdown_san := []byte(html.EscapeString(markdown_raw))
    html_raw := blackfriday.MarkdownCommon(markdown_san)
    html_san := string(bluemonday.UGCPolicy().SanitizeBytes(html_raw))

    //create comment
    comment := Comment{
        Votes: 0,
        Username: user_name,
        Html: template.HTML(html_san),
        Comments: []Comment{},
    }
    log.Println(comment)
    posts_encoder.Encode(posts)

    http.Redirect(w, r, "/post/" + strconv.FormatUint(post_id, 10) + "#" + strconv.FormatUint(comment_id, 10), http.StatusSeeOther)
}
