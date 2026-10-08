package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	//1. Usamos la base de datos en memoria como un mock
	store = NewURLData()

	//2. Preparamos el cuerpo de la peticion para simular a un usuario
	jsonBody := []byte(`{"url": "https://test.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	//3. Creamos un recorder para agarrar la respuesta del handle
	rr := httptest.NewRecorder()
	//4. Llamamos al handler directamente
	handleCreateShortURL(rr, req)

	//5. Verificamos que el codigo HTTP sae 201 Created
	if rr.Code != http.StatusCreated {
		t.Errorf("El codigo HTTP es incorrecto obtuve %d, esperaba %d", rr.Code, http.StatusCreated)
	}

	//6. Revisamos el json que nos devolvio para ver si los datos son correctos
	var responseURL URL
	err := json.NewDecoder(rr.Body).Decode(&responseURL)
	if err != nil {
		t.Fatalf(" Error al leer la respuesta: %v", err)
	}

	if responseURL.ShortCode == "" {
		t.Errorf("El codigo generado esta vacio")
	}

	if responseURL.OriginalUrl != "https://test.com" {
		t.Errorf("La URL original no coincide. Obtuve: %s", responseURL.OriginalUrl)
	}

	// Si se llega aca es porque el test paso
	t.Logf("Test exitoso. ShortCode generado: %s", responseURL.ShortCode)
}
