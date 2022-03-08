package main

import (
	"fmt"
	"net/http"
)

func GetRouter() http.Handler {
	sm := http.NewServeMux()
	sm.HandleFunc("/line1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Knock knock\n"))
	})
	sm.HandleFunc("/line2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Who's there\n"))
	})
	sm.HandleFunc("/line3", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Race condtion\n"))
	})
	return sm
}

func main() {
	fmt.Println("Starting server")
	http.ListenAndServe(":8080", GetRouter())
}
