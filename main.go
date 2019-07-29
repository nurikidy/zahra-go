package main
import (
  "net/http"
  "strings"
)

import "github.com/go-redis/redis"

func sayHello(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
  message = strings.TrimPrefix(message, "/")
  shorturl := redisGet(message)
  if shorturl !=""{
  	message = "Hello " + message + " --> " + shorturl
  	w.Write([]byte(message))
  	
  }else{
     w.Write([]byte("Record not found"))
   }
  
}

func main() {
  http.HandleFunc("/", sayHello)
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
