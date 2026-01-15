package main

import (
	"log"
	"net/http"
)

func main() {
	connectDB()

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/api/auth/register", register)
	mux.HandleFunc("/api/auth/login", login)

	// Protected Routes
	mux.Handle("/api/address", authMiddleware(http.HandlerFunc(addressRootHandler)))
	mux.Handle("/api/address/", authMiddleware(http.HandlerFunc(addressByIDHandler)))

	// Default 404 เป็น JSON
	mux.HandleFunc("/", notFoundHandler)

	// Middleware
	handler := errorMiddleware(mux)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
