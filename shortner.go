package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//I started out by looking at goshorty, but that was a little much for what I wanted
//I found https://gist.github.com/tuxychandru/d5ac22f5cb4a8ca3eeff which was more up my alley
//but it lacked the kind of routing, stats, and connection pooling I wanted

func newPool() *redis.Pool {
	//Orginally started out with individual connections, but who really wants that?
	return &redis.Pool{
		MaxIdle:   20,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

}

var pool = newPool()

func main() {
	//Seed random for code generation later
	rand.Seed(time.Now().UnixNano())

	//Use a pool of redis connections
	c := pool.Get()
	defer c.Close()

	//Routing
	rt := mux.NewRouter()
	rt.HandleFunc("/add", createShortUrl)
	rt.HandleFunc("/{short:[a-zA-Z0-9]+}", redirectShortURL).Methods("GET")
	rt.HandleFunc("/{short:[a-zA-Z0-9]+}/stats", statsForURL).Methods("GET")
	http.Handle("/", rt)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", rt)
}

func redirectShortURL(w http.ResponseWriter, r *http.Request) {
	//Get redis connection
	h := pool.Get()
	defer h.Close()

	//Pull GET params
	params := mux.Vars(r)
	shortUrl := params["short"]

	log.Println("Redirect called for: " + shortUrl)

	//Lookup the shortUrl in redis, 404 if doesn't exist
	url, err := h.Do("HGET", shortUrl, "longurl")
	if err != nil || url == nil {
		http.NotFound(w, r)
		return
	}

	//Increment hits for shortUrl
	h.Do("HINCRBY", shortUrl, "hits", 1)
	http.Redirect(w, r, string(url.([]byte)), 301)
}

func createShortUrl(w http.ResponseWriter, r *http.Request) {
	//Get redis connection
	h := pool.Get()
	defer h.Close()

	code := randCode(7)

	//I think mux should be able to obtain the url POST value, similar to the GET functions
	//but I couldn't get it working so use the request directly
	url := r.PostFormValue("url")

	log.Println("Creating short url for: " + url)

	//Don't need a variable for the return, use _ for "throwaway"
	_, err := h.Do("HMSET", string(code), "longurl", string(url), "hits", 1)
	if err != nil {
		http.Error(w, "Internal error: "+err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "http://localhost:8080/%s\n", string(code))
}

func statsForURL(w http.ResponseWriter, r *http.Request) {
	//Get redis connection
	h := pool.Get()
	defer h.Close()

	//Again, able to use mux here
	params := mux.Vars(r)
	shortUrl := params["short"]

	log.Println(shortUrl)

	//We'll return this value to the end-user so hold it
	hits, err := h.Do("HGET", shortUrl, "hits")
	if err != nil || hits == nil {
		http.NotFound(w, r)
		return
	}

	//Nice and detailed stats ;-)
	longUrl, err := h.Do("HGET", shortUrl, "longurl")
	if err != nil || hits == nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hits for %s, pointer to %s: %s\n", shortUrl, longUrl, hits)
}

func randCode(length int) []byte {
	//Allows chars for shortUrl
	src := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//Build up to length
	result := make([]byte, length)
	for i := range result {
		result[i] = src[rand.Intn(len(src))]
	}
	return result
}
