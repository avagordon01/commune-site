package main

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func user_name(post_time time.Time, user_id uint64, post_id uint64) string {
	return "test_name"
	/*rand.Seed(post_time.UnixNano() ^ int64(names_seed) ^ int64(user_id) ^ int64(post_id))
	return adjectives[rand.Intn(len(adjectives))] +
		colours[rand.Intn(len(colours))] +
		animals[rand.Intn(len(animals))]*/
}

func preview(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, _ := render_text(r.FormValue("text"), false)
	var input struct {
		Parent    Post
		Child     Comment
		Markdown  string
		Parent_id uint64
		Post_id   uint64
	}
	post_id, post_err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
	parent_id, parent_err := strconv.ParseUint(r.FormValue("parent_id"), 10, 64)
	if post_err == nil && parent_err == nil {
		parent_comment, _ := view_comment(post_id, parent_id)
		input.Parent = Post{
			ID:       parent_comment.ID,
			Time:     parent_comment.Time,
			Value:    parent_comment.Value,
			Username: parent_comment.Username,
			Html:     parent_comment.Html,
		}
	} else if post_err == nil {
		input.Parent, _ = view_post(post_id)
	} else {
		input.Parent = Post{}
	}
	input.Child = Comment{
		Html:     html_san,
		Username: user_name(input.Parent.Time, user_id, input.Parent.Rand),
		Time:     time.Now(),
	}
	input.Markdown = r.FormValue("text")
	input.Parent_id = parent_id
	input.Post_id = post_id
	err = templates["preview"].Execute(w, Page{
		Content: input,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func submit_post(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, title := render_text(r.FormValue("text"), true)
	snippet := bluemonday.StrictPolicy().Sanitize(string(html_san))
	snip_length := 200
	if len(snippet) > snip_length {
		snippet = snippet[:snip_length]
	}
	post := Post{
		Title:   title,
		Snippet: snippet,
		Time:    time.Now(),
		Value:   0,
		Html:    html_san,
		Rand:    rand.Uint64(),
	}
	post.Username = user_name(post.Time, user_id, post.Rand)
	post_id := insert_post(post)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", post_id), http.StatusSeeOther)
}

func submit_comment(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, _ := render_text(r.FormValue("text"), false)
	post_id, err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	comment := Comment{
		Time:  time.Now(),
		Value: 0,
		Html:  html_san,
	}
	post, _ := view_post(post_id)
	comment.Username = user_name(post.Time, user_id, post.Rand)
	parent_id, err := strconv.ParseUint(r.FormValue("parent_id"), 10, 64)
	insert_comment(post_id, parent_id, comment)
	http.Redirect(w, r, fmt.Sprintf("/post/%d#%d", post_id, comment.ID), http.StatusSeeOther)
}
