package auth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	token, err := getClientCredentialsToken()
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	fmt.Println("Access token:", token)
}

func getClientCredentialsToken() (string, error) {
	authURL := "http://localhost:8080/realms/hypertube/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "hypertube-api")        // 🔁 your client ID here
	data.Set("client_secret", "YOUR_SECRET_HERE") // 🔁 your client secret

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil // contains access_token JSON
}
