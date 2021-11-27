package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("listening to client")

	messagec := make(chan int)

	go func(channel chan int) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered", r)
			}
		}()

		for {
			if channel != nil {
				message := rand.Intn(100)

				channel <- message
				time.Sleep(2 * time.Second)
			}
		}
	}(messagec)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		close(messagec)
		messagec = nil
		log.Printf("client connection is closed")
	}()

	flusher, _ := w.(http.Flusher)

	for {

		select {

		case rnum := <-messagec:
			if rnum%3 == 0 {
				fmt.Fprintf(w, "event: foo\n")
			}

			if rnum%5 == 0 {
				fmt.Fprintf(w, "event: bar\n")
			}

			fmt.Fprintf(w, "data: { \"msg\": \"%s\" }\n\n", strconv.Itoa(rnum))
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile("index.htm")
	if err != nil {
		fmt.Println("Err")
	}

	fmt.Fprintf(w, "%v", string(content))
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/hello", hello)

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", nil))
}
