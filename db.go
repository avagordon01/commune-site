package main

import (
    "github.com/boltdb/bolt"
    "encoding/gob"
    "bytes"
    "encoding/binary"
    "html/template"
    "time"
    "fmt"
	"log"
	"math"
)

type Post struct {
	Id           uint64
	Title        string
	Snippet      string
	Time         time.Time
	Value        float64
	Username     string
	Html         template.HTML
	CommentCount uint64
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

var (
	err       error
	posts     []Post
	index     [5][]uint64
)

func value(freshness float64, post Post) float64 {
	return float64(post.Value) * math.Pow(0.75, -freshness*float64(post.Time.Unix()))
}

func insert_post(Post p) uint64 {}
func insert_comment(uint64 post_id, uint64 parent_id, Comment c) uint64 {}
func update_post(uint64 post_id, update_values...) {}
func update_comment(uint64 post_id, uint64 comment_id, update_values...) {}
func update_value(uint64 post_id, uint64 comment_id, float value) {}
//use the posts current freshness to index into the freshness indices, to update them
//(delete and reinsert with new freshness)
func return_slice(uint64 freshness, uint64 start, uint64 end) []Post {}
func search(string query) []
func view_post(uint64 post_id) Post {}

func main() {
    log.Println("starting...")
    db, err := bolt.Open("db/posts.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
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
        post := Post {
            Id: id,
        }

        gob.Register(Post{})
        var buf bytes.Buffer
        enc := gob.NewEncoder(&buf)
        err = enc.Encode(&post)
        if err != nil {
            log.Fatal(err)
        }
        return b.Put(itob(post.Id), buf.Bytes())
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

func itob(v uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, v)
    return b
}
