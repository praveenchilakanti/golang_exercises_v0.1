package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

/*
  use below command for HTTP post request :
  curl --data "post body content" http://localhost:8080/computeMD5
*/

func handler(w http.ResponseWriter, r *http.Request) {
	hash := md5.New()
	defer r.Body.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read the request body", 400)
		return
	}
	if len(reqBody) == 0 {
		http.Error(w, "Please send a request body", 400)
		return
	}
	if _, err := io.Copy(hash, bytes.NewReader(reqBody)); err != nil {
		http.Error(w, "Unable to read the request body", 400)
		return
	}
	fmt.Fprintf(w, "Hi there.\nMD5 Hash for your request content is: %x\n", hash.Sum(nil)[:16])
}

func main() {
	http.HandleFunc("/computeMD5", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
