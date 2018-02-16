package main

import (
	"encoding/binary"
	"github.com/asdine/storm"
	"log"
)

func loc_to_node(x []uint64) storm.Node {
	b := make([]byte, 8*len(x))
	for i, v := range x {
		binary.BigEndian.PutUint64(b[i*8:(i+1)*8], v)
	}
	var bucket storm.Node
	if len(x) > 0 {
		bucket = database.From(string(b))
	} else {
		bucket = database.Node
	}
	return bucket
}

func set_item(item *Post, loc []uint64) {
	err = loc_to_node(loc).Save(item)
	if err != nil {
		log.Fatal(err)
	}
}

func insert_post(post *Post) {
	err = database.Save(post)
	if err != nil {
		log.Fatal(err)
	}
}

func view_post(post_id uint64) Post {
	var post Post
	err = database.One("ID", post_id, &post)
	if err != nil {
		log.Fatal(err)
	}
	return post
}

func view_post_with_comments(post_id uint64) Post {
	var post Post
	//TODO
	return post
}

func insert_comment(post_id uint64, parent_ids []uint64, comment *Comment) {
	err = loc_to_node(append([]uint64{post_id}, parent_ids...)).Save(comment)
	if err != nil {
		log.Fatal(err)
	}
}

func view_comment(post_id uint64, parent_ids []uint64, comment_id uint64) Comment {
	var comment Comment
	err = loc_to_node(append([]uint64{post_id}, parent_ids...)).One("ID", comment_id, &comment)
	if err != nil {
		log.Fatal(err)
	}
	return comment
}

func test() {
	database, err = storm.Open("database/database.bolt")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	p := Post{Title: "test title"}
	insert_post(&p)
	log.Println(view_post(p.ID))
	err = loc_to_node([]uint64{}).Save(&p)
	if err != nil {
		log.Fatal(err)
	}
	set_item(&Post{}, []uint64{10})
	set_item(&Post{}, []uint64{10})
	set_item(&Post{}, []uint64{10})
	var posts []Post
	err = loc_to_node([]uint64{}).All(&posts)
	if err != nil {
		log.Fatal(err)
	}
	for _, post := range posts {
		log.Println(post)
	}
}
