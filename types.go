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
	ID           uint64 `storm:"id,increment"`
	Rand         uint64
	Title        string
	Snippet      string
	Time         time.Time
	Views        uint64
	Value        float64 `storm:"index"`
	Username     string
	Html         template.HTML
	CommentCount uint64
	Comments     []Comment `storm:"inline"`
}

type Comment struct {
	ID       uint64 `storm:"id,increment"`
	Time     time.Time
	Value    float64 `storm:"index"`
	Username string
	Html     template.HTML
	Comments []Comment `storm:"inline"`
}

type Topic struct {
	ID         uint64  `storm:"id,increment"`
	Value      float64 `storm:"index"`
	Similarity float64
	Content    string
}
