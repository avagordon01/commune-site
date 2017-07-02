package main

import (
	"bytes"
	"encoding/gob"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
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

func next_post_id() uint64 {
	s := uint64(0)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts_comments"))
		s, err = b.NextSequence()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return s
}

func next_comment_id(post_id uint64) uint64 {
	s := uint64(0)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts_comments"))
		post_bucket := b.Bucket(enc_id(post_id))
		s, err = post_bucket.NextSequence()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return s
}

func insert_post(p Post) uint64 {
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts_comments"))
		post_bucket, err := b.CreateBucket(enc_id(p.Id))
		if err != nil {
			return err
		}
		s, err := post_bucket.NextSequence()
		if err != nil {
			return err
		}
		post_bucket.Put(enc_id(s), enc_post(p))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return p.Id
}

func insert_comment(post_id uint64, parent_id uint64, c Comment) uint64 {
	s := uint64(0)
	err = db.Update(func(tx *bolt.Tx) error {
		posts_bucket := tx.Bucket([]byte("posts_comments"))
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

func return_similar(topic_id uint64) []Topic {
	topics := make([]Topic, 10)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("similar_topics"))
		b_topic := b.Bucket(enc_id(topic_id))
		c := b_topic.Cursor()
		i := uint64(0)
		for k, v := c.First(); i < 10 && k != nil; k, v = c.Next() {
			topic := dec_topic(v)
			topics = append(topics, topic)
		}
		return nil
	})
	return topics
}

func return_trending() []Topic {
	topics := make([]Topic, 10)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("trending_topics"))
		c := b.Cursor()
		i := uint64(0)
		for k, v := c.First(); i < 10 && k != nil; k, v = c.Next() {
			topic := dec_topic(v)
			topics = append(topics, topic)
		}
		return nil
	})
	return topics
}

func get_posts(freshness uint64, start uint64, num uint64) ([]Post, bool) {
	var posts []Post
	var more bool
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("freshness0"))
		c := b.Cursor()
		_, v := c.First()
		posts_bucket := tx.Bucket([]byte("posts_comments"))
		i := uint64(0)
		var k []byte
		for ; i < start; k, v = c.Next() {
		}
		for k, v = c.Next(); k != nil && i < num; k, v = c.Next() {
			post := dec_post(posts_bucket.Get(v))
			posts = append(posts, post)
		}
		more = (k != nil)
		return nil
	})
	return posts, more
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
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts_comments"))
		post_bucket := b.Bucket(enc_id(post_id))
		v := post_bucket.Get(enc_id(1))
		post = dec_post(v)
		return nil
	})
	return post, nil
}

/*db.Update(func(tx *bolt.Tx) error {
	b := tx.Bucket([]byte("posts_comments"))
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
	b := tx.Bucket([]byte("posts_comments"))
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("key=%s, value=%s\n", k, v)
	}
	return nil
})*/

func enc_freshness(f float32, id uint64, c rune) []byte {
	return append(append(enc_float(f), enc_id(id)...), []byte(string(c))...)
}
func enc_float(f float32) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&f)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
func dec_float(b []byte) float32 {
	buf := *bytes.NewBuffer(b)
	dec := gob.NewDecoder(&buf)
	var f float32
	err = dec.Decode(&f)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
func enc_id(v uint64) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&v)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
func dec_id(b []byte) uint64 {
	buf := *bytes.NewBuffer(b)
	dec := gob.NewDecoder(&buf)
	var v uint64
	err = dec.Decode(&v)
	if err != nil {
		log.Fatal(err)
	}
	return v
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
func enc_topic(t Topic) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&t)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
func dec_topic(b []byte) Topic {
	buf := *bytes.NewBuffer(b)
	dec := gob.NewDecoder(&buf)
	var t Topic
	err = dec.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
