package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const apiURL = "https://api.vatusa.net/v2/user"

type User struct {
	CID       int    `json:"cid"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Rating    int    `json:"rating"`
	Facility  string `json:"facility"`
	Status    string `json:"status"`
}

func getUserData(cid string) (*User, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	apiKey := os.Getenv("VATUSA_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not set in environment variables")
	}

	reqURL := fmt.Sprintf("%s/%s?apikey=%s", apiURL, url.PathEscape(cid), url.QueryEscape(apiKey))

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: non-200 status code received: %s", resp.Status)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error parsing JSON response")
	}

	userData, exists := data["data"].(map[string]interface{})
	if !exists {
		userData = data
	}

	user := &User{
		CID:       int(userData["cid"].(float64)),
		FirstName: userData["fname"].(string),
		LastName:  userData["lname"].(string),
		Rating:    int(userData["rating"].(float64)),
		Facility:  userData["facility"].(string),
		Status:    userData["status"].(string),
	}

	return user, nil
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		http.Error(w, "Missing CID parameter", http.StatusBadRequest)
		return
	}

	user, err := getUserData(cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the VATUSA User Lookup API! Use /user?cid={CID} to fetch user details.\n")
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/user", userHandler)

	port := ":8080"
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
