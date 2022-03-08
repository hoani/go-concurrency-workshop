package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
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

func StartServer(addr string) *http.Server {
	srv := &http.Server{Addr: addr}
	go http.ListenAndServe(addr, GetRouter())
	return srv
}

func main() {
	fmt.Println("Starting server")
	srv := StartServer(":8080")
	defer srv.Close()

	fmt.Print("press enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
