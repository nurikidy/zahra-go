package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

import "github.com/go-redis/redis"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func customerHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	shorturl := redisGet(message)

	if shorturl != "" {
		http.Redirect(w, r, shorturl, 302)

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "URL does not exist")
	}

}

type request struct {
	LongUrl  string `json: "longUrl"`
	ShortUrl string `json: "shortUrl"`
}

type response struct {
	StatusCode string `json: "statusCode"`
	ShortUrl   string `json: "shortUrl"`
	Timestamp  int32  `json: "unixtimestamp"`
}

func marketingOfficerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var u request
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if u.ShortUrl == "" {
			u.ShortUrl = RandomString(6)
		}
		redisSet(u.ShortUrl, u.LongUrl)
		s := response{StatusCode: "OK", ShortUrl: u.ShortUrl, Timestamp: int32(time.Now().Unix())}
		json.NewEncoder(w).Encode(s)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can not do that.")
	}
}
func main() {

	finish := make(chan bool)

	serverMarketingOfficer := http.NewServeMux()
	serverMarketingOfficer.HandleFunc("/", marketingOfficerHandler)
	serverCustomer := http.NewServeMux()
	serverCustomer.HandleFunc("/", customerHandler)

	go func() {
		http.ListenAndServe(":8080", serverMarketingOfficer)
	}()
	go func() {
		http.ListenAndServe(":80", serverCustomer)
	}()
	<-finish

}

func redisSet(shorturl string, longurl string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Set(shorturl, longurl, 0).Err()
	if err != nil {
		panic(err)
	}
}

func redisGet(url string) string {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := client.Get(url).Result()
	if err != nil {
		//panic(err)
		val = ""
	}

	return val
}
