package main

import (
    "fmt"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"os"
	"time"
    "encoding/json"
	"bytes"
	"encoding/binary"
	"math"
)

type Post struct {
	Id           uint64
	Title        string
	Snippet      string
	Time         time.Time
	Views        uint64
	Value        float64
	Username     string
	Html         template.HTML
	Comments     []Comment
}

type Comment struct {
	Id       uint64
	Time     time.Time
	Value    float64
	Username string
	Html     template.HTML
	Comments []Comment
}

type Topic struct {
    Id uint64
    Value float64
    Similarity float64
    Content string
}

var (
	err          error
	posts        []Post
	db           bolt.DB
)

const page_length uint64 = 50

func main() {
	f, err := os.Open("database/posts.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(f).Decode(&posts)
	if err != nil {
		log.Fatal(err)
	}
    f.Close()



	gob.Register(Post{})
	gob.Register(Comment{})
    gob.Register(Topic{})
	db, err := bolt.Open("database/posts.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	err = db.Update(func(tx *bolt.Tx) error {
		v, err := tx.CreateBucketIfNotExists([]byte("vars"))
		if err != nil {
			log.Fatal(err)
		}
		if v.Get([]byte("user_counter")) == nil {
			v.Put([]byte("user_counter"), enc_id(0))
		}
        _, err = tx.CreateBucketIfNotExists([]byte("trending_topics"))
        if err != nil {
            log.Fatal(err)
        }
        _, err = tx.CreateBucketIfNotExists([]byte("similar_topics"))
        if err != nil {
            log.Fatal(err)
        }
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



    var traverse_comment func(post_id uint64, parent_id uint64, comment Comment)
    var traverse_post func(post_id uint64, post Post)
    traverse_comment = func(post_id uint64, parent_id uint64, comment Comment) {
        //insert_comment(post_id, parent_id, comment)
        fmt.Printf("\t%v\n", comment.Id)
        for i := uint64(0); i < uint64(len(comment.Comments)); i++ {
            traverse_comment(post_id, comment.Id, comment.Comments[i])
        }
    }
    traverse_post = func(post_id uint64, post Post) {
        //insert_post(post)
        fmt.Printf("%v\n", post.Id)
        for i := uint64(0); i < uint64(len(post.Comments)); i++ {
            traverse_comment(post_id, 0, post.Comments[i])
        }
    }
    for i := uint64(0); i < uint64(len(posts)); i++ {
        traverse_post(i, posts[i])
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
        p.Id = s
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
        if err != nil {
            log.Fatal(err)
        }
        c.Id = s
		post_bucket.Put(enc_id(s), enc_comment(c))
		return nil
	})
	return s
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
