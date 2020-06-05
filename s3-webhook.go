package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"net/http"
	"strconv"
)


// Generate hmac_sha256_hex
func HmacSha256hex(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Generate hmac_sha256
func HmacSha256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return string(h.Sum(nil))
}

// Send subscription confirmation
func SubscriptionConfirmation(w http.ResponseWriter, req *http.Request , body string) {
	type S3Request struct {
		Timestamp        string `json:"Timestamp"`
		Type             string `json:"Type"`
		Message          string `json:"Message"`
		TopicArn         string `json:"TopicArn"`
		SignatureVersion int    `json:"SignatureVersion"`
		Token            string `json:"Token"`
	}

	var s3req S3Request

	// Decode json fields
	d := json.NewDecoder(req.Body)
	//	d.DisallowUnknownFields() // catch unwanted fields
	err := d.Decode(&s3req)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log token and uri
	fullURI := "http://" + req.Host + req.URL.Path
	log.Printf("Got timestamp: %s TopicArn: %s Token: %s URL: %s\n", s3req.Timestamp, s3req.TopicArn, s3req.Token, fullURI)

	// Construct sinature responce
	signature := HmacSha256hex(fullURI, HmacSha256(s3req.TopicArn, HmacSha256(s3req.Timestamp, s3req.Token)))
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

// Send subscription confirmation
//func GotRecords(w http.ResponseWriter, req *http.Request , body string) {

//}

// Liveness probe
func Ping(w http.ResponseWriter, req *http.Request) {
	// log request
	log.Printf("[%s] incoming HTTP Ping request from %s\n", req.Method, req.RemoteAddr)
	fmt.Fprintf(w, "Pong\n")
}

//Webhook
func Webhook(w http.ResponseWriter, req *http.Request) {

	// Read body
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
    
	// log request
	log.Printf("[%s] incoming HTTP request from %s\n", req.Method, req.RemoteAddr)
	log.Printf("%s\n",string(body))
	// check if we got subscription confirmation request
	if strings.Contains(string(body),"\"Type\":\"SubscriptionConfirmation\"") {
		SubscriptionConfirmation(w,req,string(body))	
	} else {
		log.Printf("Unknown request\n")
	}

}

func main() {

	// get command line args
	bindPort := flag.Int("port", 80, "number between 1-65535")
	bindAddr := flag.String("address", "", "ip address in dot format")
	flag.Parse()

	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/webhook", Webhook)

	log.Fatal(http.ListenAndServe(*bindAddr+":"+strconv.Itoa(*bindPort), nil))
}
