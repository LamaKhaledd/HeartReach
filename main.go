package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/LamaKhaledd/HeartReach/internal/db"
	"github.com/LamaKhaledd/HeartReach/internal/services/auth"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sqlx.Connect("sqlite3", "./sql/schema.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)
	jwtKey := os.Getenv("JWT_SECRET")

	registerService := services.RegisterService{Queries: queries, JwtKey: jwtKey}
	loginService := services.LoginService{Queries: queries, JwtKey: jwtKey}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email, UserName, Password, PhoneNumber, Role, Location string
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := registerService.Register(context.Background(),
			req.Email, req.UserName, req.Password, req.PhoneNumber, req.Role, req.Location,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email, Password string
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := loginService.Login(context.Background(), req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
