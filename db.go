package main

import (
    //"github.com/asdine/storm"
	"github.com/blevesearch/bleve"
	"log"
	"math"
)

func value(freshness float64, post Post) float64 {
	return float64(post.Value) * math.Pow(1.25, freshness*float64(post.Time.Unix()))
}

func compare(freshness float64) func(i, j int) bool {
	return func(i, j int) bool {
		return value(freshness, posts[i]) < value(freshness, posts[j])
	}
}

func insert_post(p Post) uint64 {
	return 0
}

func insert_comment(post_id uint64, parent_id uint64, c Comment) uint64 {
	return 0
}

func get_posts(freshness uint64, start uint64, num uint64) ([]Post, bool) {
    posts := []Post{Post{Title: "test"}}
	return posts, false
}

func text_search(query string, freshness uint64, start uint64, num uint64) ([]Post, bool) {
	search_query := bleve.NewMatchQuery(query)
	search_req := bleve.NewSearchRequest(search_query)
	search_res, err := text_index.Search(search_req)
	log.Println(search_res)
	if err != nil {
		log.Fatal(err)
	}
	return nil, false
}

func view_post(post_id uint64) (Post, error) {
    post := Post{Title: "test"}
	return post, nil
}
