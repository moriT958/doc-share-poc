package main

import (
	"doc-share-poc/internal"
	"log"
	"net/http"
)

func main() {
	hub := internal.NewHub()
	go hub.Start()

	http.HandleFunc("/", internal.ServeHome)
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		internal.ServeWS(hub, w, r)
	})

	log.Println("Server starting on port :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
