package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var countingServiceURL string
var port string

func main() {
	port = getEnvOrDefault("PORT", "80")
	portWithColon := fmt.Sprintf(":%s", port)

	countingServiceURL = getEnvOrDefault("COUNTING_SERVICE_URL", "http://localhost:9001")

	fmt.Printf("Starting server on http://0.0.0.0:%s\n", port)
	fmt.Println("(Pass as PORT environment variable)")
	fmt.Printf("Using counting service at %s\n", countingServiceURL)
	fmt.Println("(Pass as COUNTING_SERVICE_URL environment variable)")

	router := mux.NewRouter()
	router.PathPrefix("/socket.io/").Handler(startWebsocket())
	router.HandleFunc("/health", HealthHandler)
	router.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("assets").HTTPBox()))

	log.Fatal(http.ListenAndServe(portWithColon, router))
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// HealthHandler returns a succesful status and a message.
// For use by Consul or other processes that need to verify service health.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, you've hit %s\n", r.URL.Path)
}

func startWebsocket() *gosocketio.Server {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	fmt.Println("Starting websocket server...")
	server.On(gosocketio.OnConnection, handleConnection)
	server.On("send", handleSend)

	return server
}

func handleConnection(c *gosocketio.Channel) {
	fmt.Println("New client connected")
	c.Join("visits")
	handleSend(c, Count{})
}

func handleSend(c *gosocketio.Channel, msg Count) string {
	count, err := getAndParseCount()
	if err != nil {
		count = Count{Count: -1, Message: err.Error(), Hostname: "[Unreachable]"}
	}
	fmt.Println("Fetched count", count.Count)
	c.Ack("message", count, time.Second*10)
	return "OK"
}

// Count stores a number that is being counted and other data to send to
// websocket clients.
type Count struct {
	Count    int    `json:"count"`
	Message  string `json:"message"`
	Hostname string `json:"hostname"`
}

func getAndParseCount() (Count, error) {
	url := countingServiceURL

	// NOTE: We use short timeouts so round robin load balancing
	// of services can be experienced during lab demos.
	tr := &http.Transport{
		IdleConnTimeout: time.Second * 1,
	}

	httpClient := http.Client{
		Timeout:   time.Second * 2,
		Transport: tr,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Count{}, err
	}

	req.Header.Set("User-Agent", "HashiCorp Training Lab")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return Count{}, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return Count{}, readErr
	}

	defer res.Body.Close()
	return parseCount(body)
}

func parseCount(body []byte) (Count, error) {
	textBytes := []byte(body)

	count := Count{}
	err := json.Unmarshal(textBytes, &count)
	if err != nil {
		fmt.Println(err)
		return count, err
	}
	return count, err
}
