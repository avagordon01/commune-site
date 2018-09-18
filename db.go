package main

import (
	"encoding/binary"
	"github.com/asdine/storm"
	"log"
)

func loc_to_node(x []uint64) storm.Node {
	if len(x) > 0 {
        b := make([]byte, 8)
        binary.BigEndian.PutUint64(b, x[0])
        return database.From(string(b))
	} else {
		return database.Node
	}
}

func loc_to_key(x []uint64) string {
	b := make([]byte, 8*len(x))
	for i, v := range x {
		binary.BigEndian.PutUint64(b[i*8:(i+1)*8], v)
	}
	return string(b)
}

func set_item(item *Post, loc []uint64) {
	err = loc_to_node(loc).Save(item)
	if err != nil {
		log.Fatal(err)
	}
}

func get_item(loc []uint64) *Post {
    var items []Post
    err = loc_to_node(loc).All(&items)
    if err != nil {
        log.Fatal(err)
    }
    return &items[0]
}

func get_items(loc []uint64) []Post {
    var items []Post
    err = loc_to_node(loc).From().All(&items)
    if err != nil  {
        log.Fatal(err)
    }
    return items
}

func insert_post(post *Post) {
	err = database.Save(post)
	if err != nil {
		log.Fatal(err)
	}
}

func view_post(post_id uint64) (Post, error) {
	var post Post
	err = database.One("ID", post_id, &post)
    return post, err
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

func view_comment(id []uint64) (Comment, error) {
	var comment Comment
	err = loc_to_node(id).One("ID", uint64(0), &comment)
    return comment, err
}

func get_posts(freshness uint64, start uint64, page_length uint64) ([]Post, bool) {
    //TODO
    return []Post{}, false
}

func text_search(search_query string, freshness uint64, start uint64, page_length uint64) ([]Post, bool) {
    //TODO
    return []Post{}, false
}

func test() {

    //create a "root" bucket
    root := database.From("root")
    //create a bucket with key "key"
    bucket := root.From("key")
    err = bucket.Set("key", "test-sub-key", "value")
    if err != nil {
        log.Fatal(err)
    }
    //test if the key "key" exists in bucket "root"
    /*
    is_found, err := database.KeyExists("root", "key")
    if err != nil {
        log.Fatal(err)
    }
    log.Println(is_found)
    */
    //set the key "key" in bucket "root" to value "test"
    err = database.Set("root", "key", "test")
    if err != nil {
        log.Fatal(err)
    }

	p := Post{Title: "test title 0"}
    _ = loc_to_node([]uint64{0}).Save(&p)
    p.Title = "test title 0 0"
    _ = loc_to_node([]uint64{0, 0}).Save(&p)
    p.Title = "test title 0 1"
    _ = loc_to_node([]uint64{0, 1}).Save(&p)
    p.Title = "test title 0 2"
    _ = loc_to_node([]uint64{0, 2}).Save(&p)

    var posts []Post
    _ = loc_to_node([]uint64{0}).All(&posts)
    log.Println(posts)




    /*
	insert_post(&p)
    log.Println(p.ID)
	log.Println(view_post(p.ID))
    log.Println(get_item([]uint64{}))
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
    */
}
