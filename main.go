package main

import (
	"fmt"
	"net/http"
)

var (
	port = ":8080"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Connection success!")
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Printf("Listening on port: %s\n", port)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/createRoom", createRoomHandler)
	http.HandleFunc("/joinRoom", joinRoomHandler)
	http.HandleFunc("/vote", voteHandler)
	http.ListenAndServe(port, nil)
}
