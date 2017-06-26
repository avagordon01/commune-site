package main

import (
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func user_name(now time.Time, user_id uint64, post_id uint64) string {
	rand.Seed(int64(now.UnixNano()) ^ int64(user_id) ^ int64(post_id))
	return adjectives[rand.Intn(len(adjectives))] +
		colours[rand.Intn(len(colours))] +
		animals[rand.Intn(len(animals))]
}

func submit_post(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, title := render_text(r.FormValue("text"))
	snippet := bluemonday.StrictPolicy().Sanitize(html_san)

	snip_length := 200
	if len(snippet) < snip_length {
		snip_length = len(snippet)
	}
	post_id := uint64(len(posts))
	now := time.Now()
	post := Post{
		Id:           post_id,
		Title:        title,
		Snippet:      snippet[:snip_length],
		Time:         now,
		Value:        0,
		Username:     user_name(now, user_id, post_id),
		Html:         template.HTML(html_san),
		CommentCount: 0,
		Comments:     []Comment{},
	}

	http.Redirect(w, r, "/post/"+strconv.FormatUint(post.Id, 10), http.StatusSeeOther)
}

func submit_comment(w http.ResponseWriter, r *http.Request, user_id uint64) {
	html_san, _ := render_text(r.FormValue("text"))

	post_id, err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	comment_id := posts[post_id].CommentCount
	now := posts[post_id].Time
	comment := Comment{
		Id:       comment_id,
		Time:     time.Now(),
		Value:    0,
		Username: user_name(now, user_id, post_id),
		Html:     template.HTML(html_san),
		Comments: []Comment{},
	}
	parent_id, err := strconv.ParseUint(r.FormValue("comment_id"), 10, 64)
	log.Println(parent_id)

	http.Redirect(w, r, "/post/"+strconv.FormatUint(post_id, 10)+"#"+strconv.FormatUint(comment.Id, 10), http.StatusSeeOther)
}
