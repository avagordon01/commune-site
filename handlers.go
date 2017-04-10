package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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
			rand_id := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
			cookie = &http.Cookie{Name: "user_id", Value: strconv.FormatUint(rand_id, 10), Expires: time.Unix(1<<63-1, 0), Secure: true, HttpOnly: true}
			users.user_counter++
			http.SetCookie(w, cookie)
		}
		user_id, err := strconv.ParseUint(cookie.Value, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f(w, r, freshness)
	})
}

func home(w http.ResponseWriter, r *http.Request, freshness uint64) {
	templates, err := template.New("").
		Funcs(template.FuncMap{
			"human_time": humanize.Time,
		}).
		ParseGlob("templates.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	after, err := strconv.ParseUint(r.FormValue("after"), 10, 64)
	if err != nil {
		after = uint64(0)
	}
	type PageAfter struct {
		Posts []Post
		After uint64
	}
	input := PageAfter{
		After: after + 20,
	}
    if after + 20 >= uint64(len(index[freshness])) {
        input.After = 0
    }
	for i := uint64(0); i+after < uint64(len(index[freshness])) && i < 20; i++ {
		input.Posts = append(input.Posts, posts[index[freshness][i+after]])
	}
	var content bytes.Buffer
	err = templates.ExecuteTemplate(&content, "posts", input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := Page{
		Title:   template.HTML("commune"),
		Content: template.HTML(content.String()),
	}
    pusher, ok := w.(http.Pusher)
    if ok {
        pusher.Push("static/style.css", nil)
        pusher.Push("static/script.js", nil)
        pusher.Push("static/icon.png", nil)
    }
	err = templates.ExecuteTemplate(w, "main", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func search(w http.ResponseWriter, r *http.Request, freshness uint64) {
	templates, err := template.New("").
		Funcs(template.FuncMap{
			"human_time": humanize.Time,
		}).
		ParseGlob("templates.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var p []Post
	after := uint64(0)
	for i := after; i < uint64(len(index[freshness])) && i < after+20; i++ {
		p = append(p, posts[index[freshness][i]])
	}
	var content bytes.Buffer
	err = templates.ExecuteTemplate(&content, "posts", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := Page{
		Title:   template.HTML("commune"),
		Content: template.HTML(content.String()),
	}
    pusher, ok := w.(http.Pusher)
    if ok {
        pusher.Push("static/style.css", nil)
        pusher.Push("static/script.js", nil)
        pusher.Push("static/icon.png", nil)
    }
	err = templates.ExecuteTemplate(w, "main", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func post(w http.ResponseWriter, r *http.Request, freshness uint64) {
	templates, err := template.New("").
		Funcs(template.FuncMap{
			"human_time": humanize.Time,
		}).
		ParseGlob("templates.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := Page{
		Title:   template.HTML(post.Title),
		Content: template.HTML(content.String()),
	}
    pusher, ok := w.(http.Pusher)
    if ok {
        pusher.Push("static/style.css", nil)
        pusher.Push("static/script.js", nil)
        pusher.Push("static/icon.png", nil)
    }
	err = templates.ExecuteTemplate(w, "main", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func user_name(user_id uint64, post_id uint64) string {
	value := user_id + uint64(len(posts)) + names.Salt
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	v := sha1.Sum(buf)
	return names.Colours[binary.LittleEndian.Uint16(v[0:2])%uint16(len(names.Colours))] +
		names.Adjectives[binary.LittleEndian.Uint16(v[2:4])%uint16(len(names.Adjectives))] +
		names.Animals[binary.LittleEndian.Uint16(v[4:6])%uint16(len(names.Animals))]
}

func submit_post(w http.ResponseWriter, r *http.Request, user_id uint64) {
	markdown_raw := r.FormValue("text")
	markdown_san := html.EscapeString(markdown_raw)
	html_raw := string(blackfriday.MarkdownCommon([]byte(markdown_san)))
	html_san := bluemonday.UGCPolicy().Sanitize(html_raw)
	snippet := bluemonday.StrictPolicy().Sanitize(html_raw)

	post := Post{
		Id:           uint64(len(posts)),
		Title:        html.EscapeString(r.FormValue("title")),
		Snippet:      snippet,
		Time:         time.Now(),
		Votes:        0,
		Username:     user_name(user_id, uint64(len(posts))),
		Html:         template.HTML(html_san),
		CommentCount: 0,
		Comments:     []Comment{},
	}
	posts = append(posts, post)
	index[0] = append(index[0], post.Id)
	index[1] = append(index[1], post.Id)
	index[2] = append(index[2], post.Id)
	index[3] = append(index[3], post.Id)
	index[4] = append(index[4], post.Id)
	update_indices()

	http.Redirect(w, r, "/post/"+strconv.FormatUint(post.Id, 10), http.StatusSeeOther)
}

func submit_comment(w http.ResponseWriter, r *http.Request, user_id uint64) {
	markdown_raw := r.FormValue("text")
	markdown_san := html.EscapeString(markdown_raw)
	html_raw := string(blackfriday.MarkdownCommon([]byte(markdown_san)))
	html_san := bluemonday.UGCPolicy().Sanitize(html_raw)

	comment := Comment{
		Votes:    0,
		Html:     template.HTML(html_san),
		Comments: []Comment{},
	}
	post_id, err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	comment.Username = user_name(user_id, post_id)
	comment.Id = posts[post_id].CommentCount
	posts[post_id].CommentCount++
	comment_id, err := strconv.ParseUint(r.FormValue("comment_id"), 10, 64)
	if err == strconv.ErrSyntax {
		posts[post_id].Comments = append(posts[post_id].Comments, comment)
	} else if err != nil {
		posts[post_id].Comments = append(posts[post_id].Comments, comment)
	} else {
		log.Println(comment_id)
	}

	http.Redirect(w, r, "/post/"+strconv.FormatUint(post_id, 10)+"#"+strconv.FormatUint(comment.Id, 10), http.StatusSeeOther)
}
