package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request to Backend 1")
		time.Sleep(10 * time.Second)
		fmt.Fprintf(w, "Response from Backend 1 (port 8081, after delay)\n")
	})

	log.Println("Backend 1 running on :8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal(err)
	}
}
