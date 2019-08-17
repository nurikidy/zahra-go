/***
 * This package is used to handle request from common users to map given shortened URL to its original one
 */
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

import "github.com/go-redis/redis"

var r *rand.Rand

// init
func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	log.SetPrefix("ZahraLog: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.Println("init started")
}

//RandomString creates randomstring based on length put in field
func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

//CustomerHandler hanling request for customer
func CustomerHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	shorturl := RedisGet(message)

	//if shorturl keyword exists in Redis server than it redirects to 302

	if shorturl != "" {
		http.Redirect(w, r, shorturl, 302)
		log.Println("Path:", message, " Status: 302 (Redirect) Target:", shorturl)

		//if shorturl keyword does not exist in redis server it redirects to 404

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "URL does not exist")
		log.Println("Path:", message, " Status: 404 (does not exist)")

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

	log.Println("Starting server at port 8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Panicln("ERROR STARTING SERVER !!")
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
		log.Fatalln("ERROR! Unable to access Redis Server")
		//panic(err)
		val = ""
	}

	return val
}
