package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from Backend 2 (port 8082)\n")
	})

	log.Println("Backend 2 running on :8082")
	err := http.ListenAndServe(":8082", mux)
	if err != nil {
		log.Fatal(err)
	}
}
