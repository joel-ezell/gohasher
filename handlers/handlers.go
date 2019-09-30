package handlers

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/joel-ezell/gohasher/passwords"
	"github.com/joel-ezell/gohasher/statistics"
)

// e.g. http.HandleFunc("/health-check", HealthCheckHandler)
func HashHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleHashGet(w, r)
	case http.MethodPost:
		handleHashPost(w, r)
	default:
		// Return error
	}
}

func handleHashPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pwd := r.Form.Get("password")
	if pwd == "" {
		// Return 404
		log.Printf("Password form field not found")
	}

	index, err := passwords.HashAndStore(pwd)

	if err != nil {
		// return 500
		log.Printf("Error returned while hashing: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, strconv.Itoa(index))
}

func handleHashGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Path is: %s", r.URL.Path)

	re, _ := regexp.Compile("/hash/(.*)")
	values := re.FindStringSubmatch(r.URL.Path)
	if len(values) == 0 {
		//TODO: return 404
	}

	index, err := strconv.Atoi(values[1])

	if err == nil {
		// TODO: return error
	}

	hashedPwd, err := passwords.GetHashedPassword(index)

	log.Printf("hashedPwd is %s", hashedPwd)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, hashedPwd)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		stats := statistics.GetStats()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		// In the future we could report back on the status of our DB, or our cache
		// (e.g. Redis) by performing a simple PING, and include them in the response.
		io.WriteString(w, stats)
	default:
		// return error
	}
}
