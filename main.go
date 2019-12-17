package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

import "github.com/go-redis/redis/v7"

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI)

	if r.URL.Path == "/__health" {
		fmt.Print(w, "hello")
		return
	}

	key := r.URL.Path

	switch r.Method {
	case "GET", "HEAD":
		log.Println("GET", key)
		handleGet(w, key)
	case "POST", "PUT", "PATCH":
		log.Println("SET", key)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Print(w, "Error reading body", err)
		} else {
			handleSet(w, key, body)
		}
	case "DELETE":
		log.Println("DELETE", key)
		handleDelete(w, key)
	default:
		fmt.Print(w, "didn't do anything")
	}
}

func handleGet(w http.ResponseWriter, key string) {
	val, err := rdb.Get(key).Result()
	if err == nil {
		w.Write([]byte(val))
	} else if err == redis.Nil {
		w.Write([]byte("<nil>"))
	} else {
		fmt.Print(w, "Error", err)
	}
}

func handleSet(w http.ResponseWriter, key string, value []byte) {
	if err := rdb.Set(key, value, time.Hour*24).Err(); err != nil {
		fmt.Print(w, "Error", err)
		return
	}
	w.Write([]byte("ok"))
}

func handleDelete(w http.ResponseWriter, key string) {
	if err := rdb.Del(key).Err(); err != nil {
		fmt.Print(w, "Error", err)
		return
	}
	w.Write([]byte("ok"))
}

var rdb *redis.Client

func init() {
	fmt.Println("env", os.Environ())

	opts, err := redis.ParseURL(os.Getenv("FLY_REDIS_CACHE_URL"))
	if err != nil {
		log.Println("error parsing FLY_REDIS_CACHE_URL", err)
		os.Exit(1)
	}

	rdb = redis.NewClient(opts)
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
