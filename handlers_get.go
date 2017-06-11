package main

import (
	"github.com/blevesearch/bleve"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request, freshness uint64) {
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
		After: after + page_length,
	}
    if int64(after)-int64(page_length) >= 0 {
        input.Before = after - page_length
    }
	if after+page_length >= uint64(len(index[freshness])) {
		input.After = 0
	}
	for i := uint64(0); i+after < uint64(len(index[freshness])) && i < page_length; i++ {
		input.Posts = append(input.Posts, posts[index[freshness][i+after]])
	}
	err = templates["home.html"].Execute(w, Page{
		Title:     template.HTML("commune"),
		Content:   input,
		Freshness: freshness,
	})
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
}

func post(w http.ResponseWriter, r *http.Request, freshness uint64) {
	post_id, err := strconv.ParseUint(r.URL.Path[len("/post/"):], 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
    post := posts[post_id]
	err = templates["post.html"].Execute(w, Page{
		Title:     template.HTML(post.Title),
		Content:   post,
		Freshness: freshness,
	})
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
	users.page_counter++
}

func search(w http.ResponseWriter, r *http.Request, freshness uint64) {
	search_query := r.FormValue("query")
    if search_query == "" {
        home(w, r, freshness)
        return
    }
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
        Before  uint64
		Start   uint64
	}{
		Query:   search_query,
		Results: []Post{posts[0], posts[1]},
		After:   after + page_length,
		Start:   after,
	}
	err = templates["search.html"].Execute(w, Page{
		Title:     template.HTML("\"" + search_query + "\""),
		Content:   results,
		Freshness: freshness,
	})
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
}
