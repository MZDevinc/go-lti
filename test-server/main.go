package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/MZDevinc/go-lti/ltiservice"
)

var ltis ltiservice.LTIService

func init() {
}

func main() {
	fmt.Println("Testing go-lti on port 8345")

	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/launch", launch)

	log.Fatal(http.ListenAndServe(":8345", nil))
}

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "root: %s\n", req.URL.Path[1:])
}

func login(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "login: %s\n", req.URL.Path[1:])
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "req: %s\n", body)
}

func launch(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "launch: %s\n", req.URL.Path[1:])
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "req: %s\n", body)
}
