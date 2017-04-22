package main

import (
	"bytes"
	"github.com/dustin/go-humanize"
	"html/template"
	"net/http"
	"strconv"
)

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
        Start uint64
	}
	input := PageAfter{
        Start: after,
		After: after + 20,
	}
	if after+20 >= uint64(len(index[freshness])) {
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
	pusher, ok := w.(http.Pusher)
	if ok {
		pusher.Push("static/style.css", nil)
		pusher.Push("static/script.js", nil)
		pusher.Push("static/icon.png", nil)
	}
	err = templates.ExecuteTemplate(w, "main", Page{
		Title:   template.HTML("commune"),
		Content: template.HTML(content.String()),
        Freshness: freshness,
	})
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

	post_id, err := strconv.ParseUint(r.URL.Path[len("/post/"):], 10, 64)
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
	pusher, ok := w.(http.Pusher)
	if ok {
		pusher.Push("static/style.css", nil)
		pusher.Push("static/script.js", nil)
		pusher.Push("static/icon.png", nil)
	}
	err = templates.ExecuteTemplate(w, "main", Page{
		Title:   template.HTML(post.Title),
		Content: template.HTML(content.String()),
        Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users.page_counter++
}

func search(w http.ResponseWriter, r *http.Request, freshness uint64) {
}
