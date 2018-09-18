package main

import (
	c_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	m_rand "math/rand"
	"net/http"
	"strconv"
	"time"
)

func https_redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://commune.is", http.StatusMovedPermanently)
}

func hsts(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		f(w, r)
	})
}

func user_name(post_time time.Time, user_id uint64, post_id uint64) string {
    names_salt := uint64(0)
    adjectives := []string{"Testing"}
    colours := []string{"Green"}
    plants := []string{"Tree"}
	m_rand.Seed(int64(uint64(post_time.UnixNano()) ^ names_salt ^ user_id ^ post_id))
	return adjectives[m_rand.Intn(len(adjectives))] +
		colours[m_rand.Intn(len(colours))] +
		plants[m_rand.Intn(len(plants))]
}

func user_cookie(f func(w http.ResponseWriter, r *http.Request, user_id uint64)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil {
			b := make([]byte, 8)
			_, err := c_rand.Read(b)
			if err != nil {
				log.Fatal(err)
			}
			rand_id := binary.BigEndian.Uint64(b)
			cookie = &http.Cookie{Name: "user_id", Value: fmt.Sprintf("%d", rand_id), Expires: time.Unix(1<<63-1, 0), Secure: true, HttpOnly: true}
			http.SetCookie(w, cookie)
		}
		user_id, err := strconv.ParseUint(cookie.Value, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f(w, r, user_id)
	})
}

func fresh_cookie(f func(w http.ResponseWriter, r *http.Request, freshness uint64)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		freshness := uint64(2)
		if r.FormValue("freshness") != "" {
			freshness, err = strconv.ParseUint(r.FormValue("freshness"), 10, 64)
		} else if cookie, err := r.Cookie("freshness"); err == nil {
			freshness, err = strconv.ParseUint(cookie.Value, 10, 64)
		}
		cookie := &http.Cookie{Value: fmt.Sprintf("%d", freshness), Expires: time.Now().Add(time.Hour * 24), Secure: true}
		http.SetCookie(w, cookie)
		f(w, r, freshness)
	})
}
