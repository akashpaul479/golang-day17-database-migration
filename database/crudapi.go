package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at"`
}

var db *sql.DB

func Crudapi() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file ")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/users", Createusers).Methods("POST")
	r.HandleFunc("/users", Getusers).Methods("GET")
	r.HandleFunc("/users/{id}", updateusers).Methods("PUT")
	r.HandleFunc("/users/{id}", Deleteusers).Methods("DELETE")

	log.Println("Server running on port :8080")
	http.ListenAndServe(":8080", r)
}

func Createusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	_, err := db.Exec("INSERT INTO users(name , email, age) VALUES (? , ? , ?)", u.Name, u.Email, u.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func Getusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	rows, err := db.Query("SELECT id , name , email, age, created_at FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.CreatedAt)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func updateusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	json.NewDecoder(r.Body).Decode(&u)
	_, err := db.Exec("UPDATE users SET name=? , email=?, age=? WHERE id =?", u.Name, u.Email, u.Age, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Deleteusers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
