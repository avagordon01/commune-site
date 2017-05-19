package main

import (
	"bytes"
	"github.com/blevesearch/bleve"
	"github.com/dustin/go-humanize"
	"html/template"
	"log"
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
	input := struct {
		Posts []Post
		After uint64
        Before uint64
		Start uint64
	}{
		Start: after,
		After: after + 20,
	}
    if int64(after)-20 >= 0 {
        input.Before = after - 20
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
		Title:     template.HTML("commune"),
		Content:   template.HTML(content.String()),
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
		Title:     template.HTML(post.Title),
		Content:   template.HTML(content.String()),
		Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users.page_counter++
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

	search_query := r.FormValue("query")
	query := bleve.NewMatchQuery(search_query)
	search_req := bleve.NewSearchRequest(query)
	search_res, err := text_index.Search(search_req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(search_res)
	after := uint64(0)
	results := struct {
		Query   string
		Results []Post
		After   uint64
        Before uint64
		Start   uint64
	}{
		Query:   search_query,
		Results: []Post{posts[0], posts[1]},
		After:   after + 20,
		Start:   after,
	}
	var content bytes.Buffer
	err = templates.ExecuteTemplate(&content, "search", results)
	if err != nil {
		log.Fatal(err)
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
		Title:     template.HTML("\"" + search_query + "\""),
		Content:   template.HTML(content.String()),
		Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
