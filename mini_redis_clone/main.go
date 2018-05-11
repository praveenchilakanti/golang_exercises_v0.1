package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	redis "github.com/go-redis/redis"
)

var client *redis.Client

func main() {
	http.HandleFunc("/redis", handler)
	http.HandleFunc("/redisCounts", countsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		if len(r.URL.Query()) > 1 {
			http.Error(w, "Bad Request, PUT accepts only one value pair.", 400)
			log.Println("Found more then one value pair in PUT request.")
			break
		}
		for key, value := range r.URL.Query() {
			err := PUTReq(key, value[0])
			if err != nil {
				http.Error(w, "Unable to Add the value into Redis DB.", 400)
				log.Println(r.Method, ":", err)
				break
			}
			fmt.Fprintf(w, "%s: Added into DB.", value[0])
		}
	case "GET":
		key := r.URL.RawQuery
		value, err := GETReq(key)
		if err != nil {
			http.Error(w, "Value Not Found.", 400)
			log.Println(r.Method, ":", err)
			break
		}
		fmt.Fprintf(w, "::: Results ::::\nKey: %s\nValue: %s", key, value)

	case "DELETE":
		key := r.URL.RawQuery
		DelOK, err := DELReq(key)
		if err != nil {
			http.Error(w, "Unable to delete value.", 400)
			log.Println(r.Method, ":", err)
			break
		}
		if DelOK == 0 {
			fmt.Fprintf(w, "::: Results ::::\nNo value found for key: %s", key)
			break
		}
		fmt.Fprintf(w, "::: Results ::::\nDeleted value for Key: %s", key)

	default:
		http.Error(w, "Request method "+r.Method+" Not Supported", 400)
		log.Println("Request method ", r.Method, " Not Supported")
	}
}

func countsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad request, "+r.Method+" Not Supported", 400)
		return
	}
	countKey := strings.TrimSpace(r.URL.RawQuery)
	if len(countKey) == 0 {
		countKey = "*"
	} else {
		countKey = "*" + countKey + "*"
	}
	result, err := client.Keys(countKey).Result()
	if err != nil {
		http.Error(w, "Bad Request, Unable to get #records", 400)
		log.Println(r.Method, ":", err)
	}
	fmt.Fprintf(w, "::: Results ::::\n#Records : %d", len(result))
}

func PUTReq(key string, value string) error {
	return client.Set(key, value, 0).Err()
}

func GETReq(key string) (string, error) {
	return client.Get(key).Result()
}

func DELReq(key string) (int64, error) {
	return client.Del(key).Result()
}

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Redis DB started...")
}
