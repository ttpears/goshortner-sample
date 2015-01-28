package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln("Error connecting to redis:", err)
	}
	http.ListenAndServe(":8080", &handler{db})
}

type handler struct{ redis.Conn }

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.redirectShortURL(w, r)
	case "POST":
		h.createShortUrl(w, r)
	}
}

func (h *handler) redirectShortURL(w http.ResponseWriter, r *http.Request) {
	url, err := h.Do("HGET", r.URL.Path[1:], "longurl")
	if err != nil || url == nil {
		http.NotFound(w, r)
		return
	}
	h.Do("HINCRBY", r.URL.Path[1:], "hits", 1)
	http.Redirect(w, r, string(url.([]byte)), 301)
}

func (h *handler) createShortUrl(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL must be provided", 400)
		return
	}
	code := randCode(7)
	_, err := h.Do("HMSET", string(code), "longurl", string(url), "hits", 1)
	if err != nil {
		http.Error(w, "Something failed internally: "+err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "http://phoebe.webmodule.org:8080/%s\n", string(code))
}

func randCode(length int) []byte {
	src := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = src[rand.Intn(len(src))]
	}
	return result
}
