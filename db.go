package main

import (
	"github.com/blevesearch/bleve"
	"log"
	"math"
)

func value(freshness float64, post Post) float64 {
	return float64(post.Value) * math.Pow(1.25, freshness*float64(post.Time.Unix()))
}

func insert_post(p Post) uint64 {
	database.Save(&p)
	return p.ID
}

func insert_comment(post_id uint64, parent_id uint64, c Comment) uint64 {
	return 0
}

func get_posts(freshness uint64, start uint64, num uint64) ([]Post, bool) {
	var posts []Post
	query := database.Select().Limit(int(num)).Skip(int(start)).OrderBy("Value").Reverse()
	err = query.Find(&posts)
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
	var post Post
	err = database.One("ID", post_id, &post)
	return post, err
}

func view_comment(post_id uint64, comment_id uint64) (Comment, error) {
	comment := Comment{Html: "test"}
	return comment, nil
}
