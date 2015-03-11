package main

import (
	"log"
	"net/http"
)

func main() {
	db := getDBConnection()

	http.HandleFunc("/trigger", makeTriggerHandlers(db))

	http.HandleFunc("/__ding__", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("dong\n"))
	})

	log.Println("Server listening on port 3003")
	err := http.ListenAndServe("localhost:3003", nil)

	if err != nil {
		log.Println(err)
	}
}
