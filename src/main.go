package main

import (
	"log"
	"net/http"
)

func main() {
	db := getDBConnection()

	http.HandleFunc("/trigger", makeTriggerHandlers(db))

	http.HandleFunc("/__ping__", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	log.Println("Server listening on port 3003")
	err := http.ListenAndServe(":3003", nil)

	if err != nil {
		log.Println(err)
	}
}
