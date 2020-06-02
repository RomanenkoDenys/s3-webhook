package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type subscriptionrequest struct {
	Timestamp        *string `json:"Timestamp"`
	Type             *string `json:"Type"`
	Message          *string `json:"Message"`
	TopicArn         *string `json:"TopicArn"`
	SignatureVersion *int    `json:"SignatureVersion"`
	Token            *string `json:"Token"`
}

// Generate hmac_sha256
func HmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ping(w http.ResponseWriter, req *http.Request) {
	// log request
	log.Printf("[%s] incoming HTTP Ping request from %s\n", req.Method, req.RemoteAddr)
	fmt.Fprintf(w, "Pong\n")
}

func webhook(w http.ResponseWriter, req *http.Request) {
	var subscr subscriptionrequest

	// log request
	log.Printf("[%s] incoming HTTP request from %s\n", req.Method, req.RemoteAddr)
	// Decode json fields
	d := json.NewDecoder(req.Body)
	//	d.DisallowUnknownFields() // catch unwanted fields
	err := d.Decode(&subscr)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log token and uri
	fullURI := "http://" + req.Host + req.URL.Path
	log.Printf("Got token: %s for URI: %s\n", *subscr.Token, fullURI)

	// Construct sinature responce
	signature := HmacSha256(fullURI, HmacSha256(*subscr.TopicArn, HmacSha256(*subscr.Timestamp, *subscr.Token)))
	log.Printf("Generate responce signature: %s \n", signature)

	// Send responce
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"signature\":\"%s\"}", signature)
}

func main() {

	// get command line args
	bindPort := flag.Int("port", 80, "number between 1-65535")
	bindAddr := flag.String("address", "", "ip address in dot format")
	flag.Parse()

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/webhook", webhook)

	log.Fatal(http.ListenAndServe(*bindAddr+":"+strconv.Itoa(*bindPort), nil))
}
