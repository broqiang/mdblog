package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
		return
	})

	server := &http.Server{Addr: ":8080"}

	log.Fatalln(server.ListenAndServe())
}
