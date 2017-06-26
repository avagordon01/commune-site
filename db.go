package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"log"
	"math"
	"time"
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
	s := uint64(0)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		s, err = b.NextSequence()
		if err != nil {
			log.Fatal(err)
		}
		post_bucket, err := b.CreateBucket(enc_id(s))
		if err != nil {
			log.Fatal(err)
		}
		s, err = post_bucket.NextSequence()
		if err != nil {
			log.Fatal(err)
		}
		post_bucket.Put(enc_id(s), enc_post(p))
		return nil
	})
	return s
}

func insert_comment(post_id uint64, parent_id uint64, c Comment) uint64 {
	s := uint64(0)
	err = db.Update(func(tx *bolt.Tx) error {
		posts_bucket := tx.Bucket([]byte("posts"))
		post_bucket := posts_bucket.Bucket(enc_id(post_id))
		s, err = post_bucket.NextSequence()
		post_bucket.Put(enc_id(s), enc_comment(c))
		return nil
	})
	return s
}

//func update_post(post_id uint64, update_values) {}
//func update_comment(post_id uint64, comment_id uint64, update_values) {}
//use the posts current freshness to index into the freshness indices, to update them
//(delete and reinsert with new freshness)
func update_value(post_id uint64, comment_id uint64, value float32) {}

func return_slice(freshness uint64, start uint64, end uint64) []Post {
	var posts []Post
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("freshness0"))
		c := b.Cursor()
		_, v := c.First()
		posts_bucket := tx.Bucket([]byte("posts"))
		i := uint64(0)
		for ; i < start; _, v = c.Next() {
		}
		for _, v = c.Next(); i < end; _, v = c.Next() {
			post := dec_post(posts_bucket.Get(v))
			posts = append(posts, post)
		}
		return nil
	})
	return posts
}

func text_search(query string) []Post {
	search_query := bleve.NewMatchQuery(query)
	search_req := bleve.NewSearchRequest(search_query)
	search_res, err := text_index.Search(search_req)
	log.Println(search_res)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

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
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("posts_comments"))
		if err != nil {
			log.Fatal(err)
		}
		f, err := tx.CreateBucketIfNotExists([]byte("freshness_index"))
		if err != nil {
			log.Fatal(err)
		}
		for i := uint64(0); i < 5; i++ {
			_, err := f.CreateBucketIfNotExists(enc_id(i))
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*db.Update(func(tx *bolt.Tx) error {
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
	})*/
}

func enc_freshness(f float32, id uint64, c rune) []byte {
	return append(append(enc_float(f), enc_id(id)...), []byte(string(c))...)
}
func enc_float(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)
	return bytes
}
func dec_float(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	f := math.Float32frombits(bits)
	return f
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
