package main

import (
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"net/http"
	"fmt"
    "strconv"
	"time"
    "math/rand"
)

func user_name(post_time time.Time, user_id uint64, post_id uint64) string {
    return "test_name"
	/*rand.Seed(post_time.UnixNano() ^ int64(names_seed) ^ int64(user_id) ^ int64(post_id))
	return adjectives[rand.Intn(len(adjectives))] +
		colours[rand.Intn(len(colours))] +
		animals[rand.Intn(len(animals))]*/
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
        Rand:   rand.Uint64(),
	}
	post.Username = user_name(post.Time, user_id, post.Rand)
	post_id := insert_post(post)
    http.Redirect(w, r, fmt.Sprintf("/post/%d", post_id), http.StatusSeeOther)
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
	post, _ := view_post(post_id)
	comment.Username = user_name(post.Time, user_id, post.Rand)
	parent_id, err := strconv.ParseUint(r.FormValue("parent_id"), 10, 64)
	insert_comment(post_id, parent_id, comment)
	http.Redirect(w, r, fmt.Sprintf("/post/%d#%d", post_id, comment.ID), http.StatusSeeOther)
}
