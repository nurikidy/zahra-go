package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	//"strings"
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

//marketingOfficerHandler handles the post

func MarketingOfficerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	//only accepts POST if not goes to default

	case "POST":
		var u request

		//if there is nothing in the request body sends message to customer to send request body. It will return 400.

		if r.Body == nil {
			fmt.Println("There is nothing in request body")
			http.Error(w, "Please send a request body", 400)
			return
		}

		//should have same fields as Request body, if does not match Request payload throws error and returns 400.

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			fmt.Println("Does not have same fields as request body")
			http.Error(w, err.Error(), 400)
			return
		}

		//if there is no shorturl given in request payload, generates a random 6 string keyword for shorturl

		if u.ShortUrl == "" {
			u.ShortUrl = RandomString(6)
			fmt.Println("Random short url:", u.ShortUrl)
			if RedisGet(u.ShortUrl) != "" {
				fmt.Println("Duplicate keyword")
				a := false
				for a == false {
					u.ShortUrl = RandomString(6)
					if RedisGet(u.ShortUrl) != "" {
						a = false
					} else {
						a = true
					}
				}
			}
		}

		//sets u.shorturl

		RedisSet(u.ShortUrl, u.LongUrl)
		s := response{StatusCode: "OK", ShortUrl: u.ShortUrl, Timestamp: int32(time.Now().Unix())}
		json.NewEncoder(w).Encode(s)
		fmt.Println(s)
	//if other than POST returns message that "I can not do that."	And redirects to 405.

	default:
		fmt.Println("Did not submit request in Post")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can not do that.")
	}
}
func main() {
	  http.HandleFunc("/", MarketingOfficerHandler)
	  if err := http.ListenAndServe(":8080", nil); err != nil {
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
