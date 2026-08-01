package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

/*
	ME BASE EN ESTE PROYECTO
	https://roadmap.sh/projects/url-shortening-service
*/

// Instanciamos la base de datos de forma global
var store URLRepository = NewURLData()

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

	err = store.Save(&newURL)
	if err != nil {
		http.Error(w, `{"error": "Error interno al guardar la URL"}`, http.StatusInternalServerError)
		return
	}

	// 3. Respondo al cliente
	w.Header().Set("Content-Type", "application/json") //aviso que mando JSON
	w.WriteHeader(http.StatusCreated)                  //code 201 created
	json.NewEncoder(w).Encode(newURL)

}

func handleGetURL(w http.ResponseWriter, r *http.Request) {
	//  GET /shorten/{shortCode}
	shortCode := r.PathValue("shortCode")

	url, exist := store.Get_URL(shortCode)
	if !exist {
		http.Error(w, `{"error": "URL no encontrada"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //code 200 OK
	json.NewEncoder(w).Encode(url)
}

func handleUpdateURL(w http.ResponseWriter, r *http.Request) {
	// PUT /shorten/{shortCode}
	shortCode := r.PathValue("shortCode")

	var newURL CreateURLRequest

	// 1. Decodificamos y validamos el JSON entrante
	err := json.NewDecoder(r.Body).Decode(&newURL)
	if err != nil {
		http.Error(w, `{"error": "URL es requerida o JSON inválido"}`, http.StatusBadRequest)
		return
	}

	updateURL, exist := store.UpdateURL(shortCode, newURL.URL)
	if !exist {
		http.Error(w, `{"error": "URL no encontrada"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateURL)

}

func handleDeleteURL(w http.ResponseWriter, r *http.Request) {
	// DELETE /shorten/{shortCode}
	shortCode := r.PathValue("shortCode")
	deleted := store.DeleteURL(shortCode)
	if !deleted {
		http.Error(w, `{"error": "URL no encontrada"}`, http.StatusNotFound)
		return
	}
	// Si fue exitoso, devolvemos 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	// GET /shorten/{shortCode}/stats
	shortCode := r.PathValue("shortCode")
	url, exist := store.Get_URL(shortCode)
	if !exist {
		http.Error(w, `{"error": "URL no encontrada"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(url)

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")
	// 1. Buscamos y sumamos 1 al mismo tiempo
	url, exists := store.GetAndIncrement(shortCode)
	if !exists {
		http.Error(w, "Enlace no encontrado", http.StatusNotFound)
		return
	}

	// 2. Redireccionamos al usuario a la URL original
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handleCreateShortURL)
	mux.HandleFunc("GET /shorten/{shortCode}", handleGetURL)
	mux.HandleFunc("PUT /shorten/{shortCode}", handleUpdateURL)
	mux.HandleFunc("DELETE /shorten/{shortCode}", handleDeleteURL)
	mux.HandleFunc("GET /shorten/{shortCode}/stats", handleGetStats)
	mux.HandleFunc("GET /{shortCode}", handleRedirect)

	fmt.Println("listening on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
