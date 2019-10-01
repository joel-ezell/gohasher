package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/joel-ezell/gohasher/passwords"
	"github.com/joel-ezell/gohasher/statistics"
)

// HashHandler Retrieves a hashed password (GET) or computes and stores a hashed password (POST)
func HashHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleHashGet(w, r)
	case http.MethodPost:
		handleHashPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleHashPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pwd := r.Form.Get("password")
	if pwd == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Please populate a x-www-form-urlencoded field with a key of \"password\" and a value of the password to be hashed")
		log.Printf("Password form field not found")
		return
	}

	index, err := passwords.HashAndStore(pwd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error returned while hashing: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, strconv.Itoa(index))
}

func handleHashGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Path is: %s", r.URL.Path)

	// Parse the URL to find the requested index.
	// It seems like there should be an easier way to do this but this was the best I could find.
	re, _ := regexp.Compile("/hash/(.*)")
	values := re.FindStringSubmatch(r.URL.Path)
	if len(values) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "No index could be found in the URL. Please populate a URL of the form /hash/<index> (e.g. /hash/1)")
		log.Printf("No index could be found in the URL: %s", r.URL.Path)
		return
	}

	index, err := strconv.Atoi(values[1])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("The provided index %s doesn't appear to be a valid integer", values[1])
		io.WriteString(w, msg)
		log.Printf(msg)
		return
	}

	hashedPwd, err := passwords.GetHashedPassword(index)

	if hashedPwd == "" {
		w.WriteHeader(http.StatusNotFound)
		msg := fmt.Sprintf("No hashed password found at the provided index %d", index)
		io.WriteString(w, msg)
		log.Printf(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, hashedPwd)
}

// StatsHandler Retrieves the accumulated statistics thus far and returns them in the body of the 200 OK
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		stats, err := statistics.GetStats()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, stats)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
