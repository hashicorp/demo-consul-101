package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	portWithColon := fmt.Sprintf(":%s", port)

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler)

	var index uint64
	router.PathPrefix("/").Handler(CountHandler{index: &index})

	// Serve!
	fmt.Printf("Serving at http://localhost:%s\n(Pass as PORT environment variable)\n", port)
	log.Fatal(http.ListenAndServe(portWithColon, router))
}

// HealthHandler returns a succesful status and a message.
// For use by Consul or other processes that need to verify service health.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, you've hit %s\n", r.URL.Path)
}

// Count stores a number that is being counted and other data to
// return as JSON in the API.
type Count struct {
	Count    uint64 `json:"count"`
	Hostname string `json:"hostname"`
}

// CountHandler serves a JSON feed that contains a number that increments each time
// the API is called.
type CountHandler struct {
	index *uint64
}

func (h CountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(h.index, 1)
	hostname, _ := os.Hostname()
	index := atomic.LoadUint64(h.index)

	count := Count{Count: index, Hostname: hostname}

	responseJSON, _ := json.Marshal(count)
	fmt.Fprintf(w, string(responseJSON))
}
