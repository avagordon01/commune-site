package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Page struct {
	Title     template.HTML
	Content   interface{}
	Freshness uint64
}

func home(w http.ResponseWriter, r *http.Request, freshness uint64) {
	start, err := strconv.ParseUint(r.FormValue("start"), 10, 64)
	if err != nil {
		start = uint64(0)
	}
	input := struct {
		Posts []Post
		Start uint64
		Next  uint64
		Prev  uint64
	}{
		Start: start,
	}
	if start > page_length {
		input.Prev = start - page_length
	} else if start <= page_length && start > 0 {
		input.Prev = 0
	} else {
		input.Prev = start
	}
	if start+page_length < uint64(len(index[freshness])) {
		input.Next = start + page_length
	} else {
		input.Next = start
	}
	for i := start; i < uint64(len(index[freshness])) && i < start+page_length; i++ {
		input.Posts = append(input.Posts, posts[index[freshness][i]])
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
}

func post(w http.ResponseWriter, r *http.Request, freshness uint64) {
	post_id, err := strconv.ParseUint(r.URL.Path[len("/post/"):], 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	post := view_post(post_id)
	err = templates["post.html"].Execute(w, Page{
		Title:     template.HTML(post.Title),
		Content:   post,
		Freshness: freshness,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func search(w http.ResponseWriter, r *http.Request, freshness uint64) {
	start, err := strconv.ParseUint(r.FormValue("start"), 10, 64)
	if err != nil {
		start = uint64(0)
	}
	search_query := r.FormValue("query")
	if search_query == "" {
		home(w, r, freshness)
		return
	}
	text_search(search_query)
	results := struct {
		Query   string
		Results []Post
		Start   uint64
		Next    uint64
		Prev    uint64
	}{
		Query:   search_query,
		Results: []Post{posts[0], posts[1]},
		Start:   start,
	}
	if start > page_length {
		results.Prev = start - page_length
	} else if start <= page_length && start > 0 {
		results.Prev = 0
	} else {
		results.Prev = start
	}
	if start+page_length < uint64(len(results.Results)) {
		results.Next = start + page_length
	} else {
		results.Next = start
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
}
