package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func worker(channel chan string, eventT string, pauseT int, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered", r)
		}
	}()

	for {
		select {

		case <-r.Context().Done():
			return

		default:
			rnum := rand.Intn(100)

			message := ""

			if rnum%3 == 0 {
				message = fmt.Sprintf("event: %s\ndata: { \"msg\": \"%d\" }\n\n", eventT, rnum)
			} else {
				message = fmt.Sprintf("data: { \"msg\": \"%d\" }\n\n", rnum)
			}

			channel <- message
			time.Sleep(time.Duration(pauseT) * time.Second)
		}
	}

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messagec := make(chan string)

	defer func() {
		close(messagec)
		messagec = nil
		log.Printf("client connection is closed")
	}()

	flusher, _ := w.(http.Flusher)

	log.Printf("listening to client")

	go worker(messagec, "foo", 2, r)
	go worker(messagec, "bar", 4, r)

	for {
		select {

		case <-r.Context().Done():
			return

		case message := <-messagec:

			fmt.Fprintf(w, message)

			flusher.Flush()
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.htm")
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/hello", hello)

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", nil))
}
