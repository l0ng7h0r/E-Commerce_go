//go:build ignore

package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	secretKey := "$2a$10$9SJjEs7KYVq6tjQ970bdN.GhOB3UNdiqJR8ZuHFuB7qesiDEbjnuW"
	txID := "849aeadf-4a24-49f9-a0c5-8d50d8295c31"

	url := fmt.Sprintf("https://payment-gateway.phajay.co/v1/api/payment/check-transaction/status/%s", txID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Request error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	encodedKey := base64.StdEncoding.EncodeToString([]byte(secretKey))
	req.Header.Set("Authorization", "Basic "+encodedKey)
	req.Header.Set("secretKey", secretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Do error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll error: %v", err)
	}

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Raw Response Body:\n%s\n", string(body))
}
