package main

import (
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func user_name(post_time time.Time, user_id uint64, post_id uint64) string {
	rand.Seed(post_time.UnixNano() ^ int64(names_seed) ^ int64(user_id) ^ int64(post_id))
	return adjectives[rand.Intn(len(adjectives))] +
		colours[rand.Intn(len(colours))] +
		animals[rand.Intn(len(animals))]
}

func submit_post(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, title := render_text(r.FormValue("text"))
	snippet := bluemonday.StrictPolicy().Sanitize(html_san)
	snip_length := 200
	if len(snippet) > snip_length {
		snippet = snippet[:snip_length]
	}
	post := Post{
		Title:   title,
		Snippet: snippet,
		Time:    time.Now(),
		Value:   0,
		Html:    template.HTML(html_san),
	}
	post.Id = next_post_id()
	post.Username = user_name(post.Time, user_id, post.Id)
	insert_post(post)
	http.Redirect(w, r, "/post/"+strconv.FormatUint(post.Id, 10), http.StatusSeeOther)
}

func submit_comment(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, _ := render_text(r.FormValue("text"))
	post_id, err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	comment := Comment{
		Time:  time.Now(),
		Value: 0,
		Html:  template.HTML(html_san),
	}
	post := view_post(post_id)
	comment.Id = next_comment_id(post_id)
	comment.Username = user_name(post.Time, user_id, post_id)
	parent_id, err := strconv.ParseUint(r.FormValue("parent_id"), 10, 64)
	insert_comment(post_id, parent_id, comment)
	http.Redirect(w, r, "/post/"+strconv.FormatUint(post_id, 10)+"#"+strconv.FormatUint(comment.Id, 10), http.StatusSeeOther)
}
