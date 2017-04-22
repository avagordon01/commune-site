package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func hsts(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		f(w, r)
	})
}

func user_cookie(f func(w http.ResponseWriter, r *http.Request, user_id uint64)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil {
			rand_id := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
			cookie = &http.Cookie{Name: "user_id", Value: strconv.FormatUint(rand_id, 10), Expires: time.Unix(1<<63-1, 0), Secure: true, HttpOnly: true}
			users.user_counter++
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
		cookie, err := r.Cookie("freshness")
		if err != nil {
			cookie = &http.Cookie{Name: "freshness", Value: strconv.FormatUint(2, 10), Expires: time.Now().Add(time.Hour * 24), Secure: true}
			http.SetCookie(w, cookie)
		}
		freshness, err := strconv.ParseUint(cookie.Value, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		f(w, r, freshness)
	})
}
