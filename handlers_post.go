package main

import (
	"crypto/sha1"
	"encoding/binary"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func user_name(user_id uint64, post_id uint64) string {
	value := user_id + uint64(len(posts))
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	v := sha1.Sum(buf)
	return names.Adjectives[binary.LittleEndian.Uint16(v[0:2])%uint16(len(names.Adjectives))] +
		names.Colours[binary.LittleEndian.Uint16(v[2:4])%uint16(len(names.Colours))] +
		names.Animals[binary.LittleEndian.Uint16(v[4:6])%uint16(len(names.Animals))]
}

func submit_post(w http.ResponseWriter, r *http.Request, user_id uint64) {
    html_san := render_markdown(r.FormValue("text"))
	snippet := bluemonday.StrictPolicy().Sanitize(html_san)

	post := Post{
		Id:           uint64(len(posts)),
		Title:        html.EscapeString(r.FormValue("title")),
		Snippet:      snippet,
		Time:         time.Now(),
		Value:        0,
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
    html_san := render_markdown(r.FormValue("text"))

	post_id, err := strconv.ParseUint(r.FormValue("post_id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	comment := Comment{
		Id:       posts[post_id].CommentCount,
		Time:     time.Now(),
		Value:    0,
		Username: user_name(user_id, post_id),
		Html:     template.HTML(html_san),
		Comments: []Comment{},
	}
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
