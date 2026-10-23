package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/GiovA07/SHORT-URL/storage"
)

/*
	ME BASE EN ESTE PROYECTO
	https://roadmap.sh/projects/url-shortening-service
*/

// Instanciamos la base de datos de forma global
var store storage.URLRepository

// CreateURLRequest representa el JSON entrante
type CreateURLRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
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

	//2 Valido que la URL sea del formato correcto (http://)
	//ParseREquestURI comprueba que el formato sea correcto
	parsedURL, err := url.ParseRequestURI(URLreq.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, `{"error": "URL inválida. Asegúrate de incluir http:// o https://"}`, http.StatusBadRequest)
		return
	}

	//3 Determinamos el short code
	var shortCode string

	if URLreq.CustomCode != "" {
		shortCode = URLreq.CustomCode
		_, exists := store.Get_URL(shortCode)
		if exists {
			http.Error(w, `{"error": "Este codigo personalizado ya esta en uso"}`, http.StatusConflict)
			return
		}
	} else {
		// si no mando nada, lo generamos aleatoriamente
		shortCode = generateShortCode(6)
	}

	//3. Genero los datos
	now := time.Now().UTC()

	//Creamos el objeto usando el struct de URL
	newURL := storage.URL{
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

	// 4. Respondo al cliente
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
	var err error
	// 1. Inicializo nuestro nuevo motor de MySQL
	store, err = storage.NewURLMySQLStore()
	if err != nil {
		log.Fatalf("Error al iniciar MySQL: %v", err)
	}
	log.Println("Conectado a MySQL!")

	mux := http.NewServeMux()
	// 2. Registro las rutas
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 2. Sirve el HTML principal exactamente en la raiz
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	//3. Mis rutas
	mux.HandleFunc("POST /r/shorten", handleCreateShortURL)
	mux.HandleFunc("GET /r/shorten/{shortCode}", handleGetURL)
	mux.HandleFunc("PUT /r/shorten/{shortCode}", handleUpdateURL)
	mux.HandleFunc("DELETE /r/shorten/{shortCode}", handleDeleteURL)
	mux.HandleFunc("GET /r/shorten/{shortCode}/stats", handleGetStats)
	mux.HandleFunc("GET /r/{shortCode}", handleRedirect)

	log.Println("Servidor corriendo en el puerto 8080")

	// 4. Prendo el server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
