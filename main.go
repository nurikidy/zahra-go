package main
import (
  "net/http"
  "strings"
  "fmt"
)

import "github.com/go-redis/redis"

func shortner(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
  message = strings.TrimPrefix(message, "/")
  shorturl := redisGet(message)
  
  if shorturl !=""{
  	w.WriteHeader(http.StatusFound)
    fmt.Fprint(w, "page found")
  	
  }else{
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "page not found")
   }
  
}

func main() {
  http.HandleFunc("/", shortner)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}


func redisGet(url string) string{

  client := redis.NewClient(&redis.Options{
  Addr :     "localhost:6379",
  Password: "", // no password set
  DB:       0,  // use default DB
  })
  
  val, err := client.Get(url).Result()
  if err != nil {
	//panic(err)
	val= ""
  }
   
  return val
}
