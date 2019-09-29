package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joel-ezell/gohasher/handlers"
)

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
	* The statistics package has a singleton that tracks the current instance and the average execution time. Both of these must be protected
	* via a mutex.
	* The passwords package features a singleton that maps the increment number (learned from the statistics package) to the hashed password
	* Another singleton waitgroup exists. When a POST comes to the hash, the password package spins off a goroutine and adds the worker to
	* the waitgroup. A call to a shutdown immediately stops the web server to prevent new requests from arriving. It then spins off a goroutine
	* which waits for the group to complete. After it completes, the process exits.
	*
	 */

	http.HandleFunc("/hash", handlers.HashHandler)
	http.HandleFunc("/stats", handlers.StatsHandler)
	http.HandleFunc("/shutdown", handlers.ShutdownHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("server listen error: %+v", err)
	}
}
