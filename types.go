package main

import (
	"html/template"
	"time"
)

type Page struct {
	Title     string
	Content   interface{}
	Freshness uint64
}

type Post struct {
	ID           uint64
	Rand         uint64
	Title        string
	Snippet      string
	Time         time.Time
	Views        uint64
	Value        float64
	Username     string
	Html         template.HTML
	CommentCount uint64
	Comments     []Comment
}

type Comment struct {
	ID       uint64
	Time     time.Time
	Value    float64
	Username string
	Html     template.HTML
	Comments []Comment
}

type Topic struct {
	ID         uint64
	Value      float64
	Similarity float64
	Content    string
}
