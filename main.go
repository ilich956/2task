package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	port        = ":8080"
	connStr     = "postgres://postgres:bayipket@localhost/adv_database?sslmode=disable"
	driverName  = "postgres"
	tableName   = "user_table"
	createTable = `
		CREATE TABLE IF NOT EXISTS user_table (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			email VARCHAR(255),
			username VARCHAR(255),
			password VARCHAR(255)
		);
	`
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

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open(driverName, connStr)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Error creating user_table:", err)
		return
	}

	http.HandleFunc("/register", handleRequest)
	http.HandleFunc("/userList", handleGetRequest)
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Server listening on port", port)
	http.ListenAndServe(port, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/userList" {
		handleGetRequest(w, r)
		return
	}
	if r.Method == http.MethodPost {
		handlePostRequest(w, r)
		return
	}
	http.Error(w, "Method not all", http.StatusMethodNotAllowed)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	users, err := getFromDB()
	if err != nil {
		http.Error(w, "Status", http.StatusInternalServerError)
		return
	}
	ts, err := template.ParseFiles("userList.html")
	if err != nil {
		http.Error(w, "Loh", http.StatusInternalServerError)
		return
	}
	ts.Execute(w, users)

}

func getFromDB() ([]RegistrationData, error) {
	rows, err := db.Query(`SELECT "username","password" FROM user_table`)
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %s", err)
	}
	defer rows.Close()

	var users []RegistrationData

	for rows.Next() {
		var u RegistrationData
		err := rows.Scan(&u.Username, &u.Password)
		if err != nil {
			return nil, fmt.Errorf("Error scanning row: %s", err)
		}
		users = append(users, u)
	}
	fmt.Print(users)
	return users, nil
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

	// Insert user registration data into the database
	err = insertUser(registrationData)
	if err != nil {
		handleError(w, "Error inserting user data into the database")
		return
	}

	fmt.Printf("Received registration data: %+v\n", registrationData)

	response := ResponseData{
		Status:  "success",
		Message: "Registration data successfully received and inserted into the database",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		handleError(w, "Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func insertUser(data RegistrationData) error {
	_, err := db.Exec("INSERT INTO "+tableName+" (name, email, username, password) VALUES ($1, $2, $3, $4)",
		data.Name, data.Email, data.Username, data.Password)
	return err
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
