package main

import (
    "log"
    "bytes"
    "time"
    "strconv"
    "net/http"
    "html/template"
    "golang.org/x/net/html"
    "github.com/microcosm-cc/bluemonday"
    "github.com/russross/blackfriday"
)

func hsts(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
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
            cookie = &http.Cookie{Name: "freshness", Value: strconv.FormatUint(2, 10), Expires: time.Now().Add(time.Hour * 12), Secure: true}
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

    after, err := strconv.ParseUint(r.FormValue("after"), 10, 64)
    if err != nil {
        after = uint64(0)
    }
    var p []Post
    for i := uint64(0); i + after < uint64(len(index[0])) && i < 20; i++ {
        p = append(p, posts[index[0][i + after]])
    }
    type PageAfter struct {
        Posts []Post
        After uint64
    }
    input := PageAfter{
        Posts: p,
        After: after + 20,
    }
    var content bytes.Buffer
    err = templates.ExecuteTemplate(&content, "posts", input)
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
    post_id, err := strconv.ParseUint(path, 10, 64)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    post := posts[post_id]
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
    user_name := "user"
    markdown_raw := r.FormValue("text")
    markdown_san := []byte(html.EscapeString(markdown_raw))
    html_raw := blackfriday.MarkdownCommon(markdown_san)
    html_san := string(bluemonday.UGCPolicy().SanitizeBytes(html_raw))

    post := Post{
        Title: r.FormValue("title"),
        Votes: 0,
        Username: user_name,
        Html: template.HTML(html_san),
        Comments: []Comment{},
    }
    post.Id = uint64(len(posts))
    posts[post.Id] = post;
    posts_encoder.Encode(posts)

    http.Redirect(w, r, "/post/" + strconv.FormatUint(post.Id, 10), http.StatusSeeOther)
}

func submit_comment(w http.ResponseWriter, r *http.Request, user_id uint64) {
    user_name := "user"
    markdown_raw := r.FormValue("text")
    markdown_san := []byte(html.EscapeString(markdown_raw))
    html_raw := blackfriday.MarkdownCommon(markdown_san)
    html_san := string(bluemonday.UGCPolicy().SanitizeBytes(html_raw))

    comment := Comment{
        Votes: 0,
        Username: user_name,
        Html: template.HTML(html_san),
        Comments: []Comment{},
    }
    post_id, err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    comment_id, err := strconv.ParseUint(r.FormValue("comment_id"), 10, 64)
    if err == strconv.ErrSyntax {
        posts[post_id].Comments = append(posts[post_id].Comments, comment)
    } else if err != nil {
        http.NotFound(w, r)
        return
    } else {
        log.Println(comment_id)
        //find comment and add reply
    }
    posts_encoder.Encode(posts)

    http.Redirect(w, r, "/post/" + strconv.FormatUint(post_id, 10) + "#" + strconv.FormatUint(comment.Id, 10), http.StatusSeeOther)
}

func submit_upvote(w http.ResponseWriter, r *http.Request, user_id uint64) {
}
