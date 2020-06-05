package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Generate hmac_sha256_hex
func HmacSha256hex(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Generate hmac_sha256_base64
func HmacSha256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return string(h.Sum(nil))
}

func main() {

       Timestamp := "2019-12-26T19:29:12+03:00"
       Url := "http://test.com"
       TopicArn := "mcs2883541269|bucketA|s3:ObjectCreated:Put"
       Token := "RPE5UuG94rGgBH6kHXN9FUPugFxj1hs2aUQc99btJp3E49tA"

	// Construct sinature responce
	signature := HmacSha256hex(Url, HmacSha256(TopicArn, HmacSha256(Timestamp, Token)))
	fmt.Printf(signature)
}
