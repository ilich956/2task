package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RegistrationData struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const port = ":8080"

func main() {
	http.HandleFunc("/register", handleRequest)
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Server listening on port", port)
	http.ListenAndServe(port, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, r)
	case http.MethodPost:
		handlePostRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	response := ResponseData{
		Status:  "success",
		Message: "GET request received",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, "Invalid JSON format")
		return
	}

	var registrationData RegistrationData
	err = json.Unmarshal(body, &registrationData)
	if err != nil {
		handleError(w, "Invalid JSON format")
		return
	}

	if registrationData.Password != registrationData.ConfirmPassword {
		handleError(w, "Password and confirm password do not match")
		return
	}

	fmt.Printf("Received registration data: %+v\n", registrationData)

	response := ResponseData{
		Status:  "success",
		Message: "Registration data successfully received",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		handleError(w, "Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func handleError(w http.ResponseWriter, message string) {
	response := ResponseData{
		Status:  "400",
		Message: message,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(responseJSON)
}
