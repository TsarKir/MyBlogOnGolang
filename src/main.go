package main

import (
	"ex02/server"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

const (
	rateLimit  = 100
	burstLimit = 10
)

func rateLimiter(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rateLimit), burstLimit)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	start()
	r := mux.NewRouter()
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images/")))) // handler for logo
	r.Use(rateLimiter)
	r.HandleFunc("/login", server.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/admin", server.AdminHandler).Methods("GET", "POST")
	r.HandleFunc("/posts/", server.PostsHandler)
	r.HandleFunc("/posts/post/", server.PostHandler)

	log.Println("Server is working on port 8888")
	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Printf("Server launch error: %s\n", err)
	}
}

func start() {
	cmd := exec.Command("unzip", "additional_files.zip")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}
