package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

func insert_post(p Post) uint64 { return uint64(0) }

func insert_comment(post_id uint64, parent_id uint64, c Comment) uint64 { return uint64(0) }

//func update_post(post_id uint64, update_values) {}
//func update_comment(post_id uint64, comment_id uint64, update_values) {}
func update_value(post_id uint64, comment_id uint64, value float32) {}

//use the posts current freshness to index into the freshness indices, to update them
//(delete and reinsert with new freshness)
func return_slice(freshness uint64, start uint64, end uint64) []Post { return nil }

func text_search(query string) []Post { return nil }

func view_post(post_id uint64) Post {
	var post_buf []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		v := b.Get(enc_id(post_id))
		copy(post_buf, v)
		return nil
	})
	return dec_post(post_buf)
}

func test() {
	gob.Register(Post{})
	gob.Register(Comment{})
	db, err := bolt.Open("database/posts.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		//update a post/comment
		//add to their value
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		//read only
		//maybe for getting a single post/comment
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("posts"))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		err := b.Put([]byte("answer"), []byte("42"))
		return err
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		id, _ := b.NextSequence()
		post := Post{
			Id: id,
		}

		gob.Register(Post{})
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(&post)
		if err != nil {
			log.Fatal(err)
		}
		return b.Put(enc_id(post.Id), buf.Bytes())
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}

func enc_id(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
func dec_id(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
func enc_post(p Post) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
func dec_post(b []byte) Post {
	buf := *bytes.NewBuffer(b)
	dec := gob.NewDecoder(&buf)
	var p Post
	err = dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
func enc_comment(c Comment) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&c)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
func dec_comment(b []byte) Comment {
	buf := *bytes.NewBuffer(b)
	dec := gob.NewDecoder(&buf)
	var c Comment
	err = dec.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
