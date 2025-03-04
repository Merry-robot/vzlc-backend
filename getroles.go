package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const apiURL = "https://api.vatusa.net/v2/user"

// User represents the structure of the JSON response from the API
type User struct {
	CID       int    `json:"cid"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Rating    int    `json:"rating"`
	Facility  string `json:"facility"`
	Status    string `json:"status"`
}

func getUserData(cid string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("VATUSA_API_KEY")
	if apiKey == "" {
		log.Fatalf("API key not set in environment variables")
	}

	reqURL := fmt.Sprintf("%s/%s?apikey=%s", apiURL, url.PathEscape(cid), url.QueryEscape(apiKey))

	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Non-200 status code received:", resp.Status)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error parsing JSON response")
		return
	}

	userData, exists := data["data"].(map[string]interface{})
	if !exists {
		userData = data
	}

	fmt.Printf("User Info:\nCID: %v\nName: %s %s\nRating: %v\nFacility: %s\nStatus: %s\n",
		userData["cid"], userData["fname"], userData["lname"], userData["rating"], userData["facility"], userData["status"])
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter CID: ")
	cid, _ := reader.ReadString('\n')
	cid = strings.TrimSpace(cid)

	getUserData(cid)
}
