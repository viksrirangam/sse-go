package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messagec := make(chan int)

	defer func() {
		close(messagec)
		messagec = nil
		log.Printf("client connection is closed")
	}()

	flusher, _ := w.(http.Flusher)

	log.Printf("listening to client")

	go func(channel chan int) {

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
				message := rand.Intn(100)

				channel <- message
				time.Sleep(2 * time.Second)
			}
		}

	}(messagec)

	for {
		select {

		case <-r.Context().Done():
			return

		case rnum := <-messagec:
			if rnum%3 == 0 {
				fmt.Fprintf(w, "event: foo\n")
			}

			if rnum%5 == 0 {
				fmt.Fprintf(w, "event: bar\n")
			}

			fmt.Fprintf(w, "data: { \"msg\": \"%d\" }\n\n", rnum)
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
