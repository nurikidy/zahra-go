package main

import (

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

//creates randomstring based on length put in field

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func CustomerHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	shorturl := RedisGet(message)

	//if shorturl keyword exists in Redis server than it redirects to 302

	if shorturl != "" {
		http.Redirect(w, r, shorturl, 302)
		fmt.Println("keyword:", message, "does exist and redirects to", shorturl)

		//if shorturl keyword does not exist in redis server it redirects to 404

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "URL does not exist")
		fmt.Println("keyword:", message, "does not exist")

	}

}

//creating request struct that takes LongUrl and ShortUrl as fields. 'json:: "longUrl" tells JSON encoder and decode to use these names instead of capatlised names

type request struct {
	LongUrl  string `json: "longUrl"`
	ShortUrl string `json: "shortUrl"`
}

//creating response struct that takes StatusCode,Timestamp, and ShortUrl as fields. 'json:: "statusCode" tells JSON encoder and decode to use these names instead of capatlised names

type response struct {
	StatusCode string `json: "statusCode"`
	ShortUrl   string `json: "shortUrl"`
	Timestamp  int32  `json: "unixtimestamp"`
}

func main() {
	  http.HandleFunc("/", CustomerHandler)
	  if err := http.ListenAndServe(":80", nil); err != nil {
	    panic(err)
	  }
	}
	

func RedisSet(shorturl string, longurl string) {

	//creates connection with redis server

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

//gets the longurl that is matched with the keyword

func RedisGet(url string) string {

	//creates connection with redis server

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
