package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joel-ezell/gohasher/handlers"
	"github.com/joel-ezell/gohasher/passwords"
)

var srv http.Server

func main() {
	const defaultPort = "8080"
	const portEnv = "HASHER_PORT"

	port := defaultPort

	if p := os.Getenv("HASHER_PORT"); p != "" {
		port = p
	}

	log.Printf("Hash server is starting to listen on port %s", port)

	/*
	* Basic design:
	*
	* The statistics package has a singleton that tracks the total number of invocations and the average execution time.
	* This data is protected by a mutex.
	*
	* The passwords package features a singleton that maps the index number to the corresponding hashed password. There is also a singleton
	* counter to track the current index. Both of these are protected via mutex.
	* Another singleton waitgroup exists. When a request to perform a new has arrives, the password package spins up a
	* goroutine and adds the worker to the waitgroup.
	*
	* A call to a shutdown immediately stops the web server to prevent new requests from arriving. It then spins up a goroutine
	* which waits for the group to complete. After it completes, the process exits
	*
	 */

	srv := &http.Server{Addr: ":" + port}
	http.HandleFunc("/hash", handlers.HashHandler)
	http.HandleFunc("/hash/", handlers.HashHandler)
	http.HandleFunc("/stats", handlers.StatsHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server listen error: %+v", err)
	}
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	srv.Shutdown(context.Background())
	w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	// io.WriteString(w, `{"alive": true}`)
	go func() {
		passwords.WaitToComplete()
		os.Exit(1)
	}()
}
