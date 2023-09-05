// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

var consulPersistPath string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	portWithColon := fmt.Sprintf(":%s", port)

	consulPersistPath = os.Getenv("COUNTING_PATH")
	fmt.Printf("Saving count at \"%s\"\n(Pass as COUNTING_PATH environment variable)\n", consulPersistPath)

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler)

	var index uint64

	getStartingIndex(&index)

	router.PathPrefix("/").Handler(CountHandler{index: &index})

	// Serve!
	fmt.Printf("Serving at http://localhost:%s\n(Pass as PORT environment variable)\n", port)
	go http.ListenAndServe(portWithColon, router)

	waitForSignal(&index)
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

func getStartingIndex(index *uint64) {
	if consulPersistPath == "" {
		*index = uint64(0)
		return
	}
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	pair, _, err := kv.Get(consulPersistPath, nil)
	if err != nil {
		panic(err)
	}
	*index = uint64(0)
	if pair != nil {
		_, err = fmt.Sscan(string(pair.Value), index)
	}
	fmt.Printf("starting value is :%d\n", *index)
}

func putFinishingIndex(index *uint64) {
	if consulPersistPath == "" {
		return
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	s := fmt.Sprintf("%d", *index)
	p := &api.KVPair{Key: consulPersistPath, Value: []byte(s)}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("saved value was :%d\n", *index)
}

func waitForSignal(index *uint64) {
	xsig := make(chan os.Signal)
	signal.Notify(xsig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	hsig := make(chan os.Signal)
	signal.Notify(hsig, syscall.SIGHUP)
	for {
		select {
		case s := <-xsig:
			putFinishingIndex(index)
			log.Fatalf("Got signal: %v, exiting.", s)
		case s := <-hsig:
			putFinishingIndex(index)
			log.Printf("Got signal: %v, continue.", s)
		}
	}
}
