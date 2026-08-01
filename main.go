package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Instanciamos la base de datos de forma global
var store = NewURLData()

// CreateURLRequest representa el JSON entrante
type CreateURLRequest struct {
	URL string `json:"url"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func handleCreateShortURL(w http.ResponseWriter, r *http.Request) {
	var URLreq CreateURLRequest

	// 1. Decodificamos y validamos el JSON entrante
	err := json.NewDecoder(r.Body).Decode(&URLreq)
	if err != nil {
		http.Error(w, `{"error": "URL es requerida o JSON inválido"}`, http.StatusBadRequest)
		return
	}

	//2. Genero los datos
	now := time.Now().UTC()
	shortCode := generateShortCode(6)

	//Creamos el objeto usando el struct de URL
	newURL := URL{
		OriginalUrl: URLreq.URL,
		ShortCode:   shortCode,
		CreateAt:    now,
		UpdateAt:    now,
	}

	store.Save(&newURL)

	// 3. Respondo al cliente
	w.Header().Set("Content-Type", "application/json") //aviso que mando JSON
	w.WriteHeader(http.StatusCreated)                  //code 201 created
	json.NewEncoder(w).Encode(newURL)

}

func handleGetURL(w http.ResponseWriter, r *http.Request) {
	//  GET /shorten/{shortCode}
	shortCode = r.PathValue("shortCode")
	
	URL := repository.

	fmt.Fprintf(w, "Endpoint para obtener la URL original")
}

func handleUpdateURL(w http.ResponseWriter, r *http.Request) {
	// PUT /shorten/{shortCode}
	fmt.Fprintf(w, "Endpoint para actualizar URL")
}

func handleDeleteURL(w http.ResponseWriter, r *http.Request) {
	// DELETE /shorten/{shortCode}
	fmt.Fprintf(w, "Endpoint para eliminar URL")
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	// GET /shorten/{shortCode}/stats
	fmt.Fprintf(w, "Endpoint para estadisticas")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handleCreateShortURL)
	mux.HandleFunc("GET /shorten", handleGetURL)
	mux.HandleFunc("PUT /shorten/{shortCode}", handleUpdateURL)
	mux.HandleFunc("DELETE /shorten/{shortCode}", handleDeleteURL)
	mux.HandleFunc("GET /shorten/{shortCode}/status", handleGetStats)

	fmt.Println("listening on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
