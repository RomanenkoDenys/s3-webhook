package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
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
	//signature := HmacSha256("http://test.com", HmacSha256("mcs2883541269|bucketA|s3:ObjectCreated:Put", HmacSha256("2019-12-26T19:29:12+03:00", "RPE5UuG94rGgBH6kHXN9FUPugFxj1hs2aUQc99btJp3E49tA")))
	log.Printf("Generate responce signature: %s \n", signature)

	// Send responce
	type Signature struct {
		Signature string `json:"signature"`
	}
	var s Signature
	s.Signature = signature
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
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
