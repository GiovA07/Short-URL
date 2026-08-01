package main

import (
	"strconv"
	"time"
)

// Esta va ser la base de datos por ahora hasta que conecte con mysql
type URLDataBase struct {
	urls   map[string]URL
	nextID int
}

/* ESTO SEGURAMENTE SI ES CONCURRENTE PODRIAN OCURRIR ERRORES*/
func NewURLData() *URLDataBase {
	return &URLDataBase{
		urls:   make(map[string]URL),
		nextID: 1,
	}
}

// Guarda la URL generada en el MAP
func (s *URLDataBase) Save(u *URL) {
	u.ID = strconv.Itoa(s.nextID)
	s.nextID++
	s.urls[u.ShortCode] = *u
}

// Obtener una URL
func (s *URLDataBase) Get_URL(shortCode string) (*URL, bool) {
	url, exists := s.urls[shortCode]
	if !exists {
		return nil, false
	}
	return &url, true
}

// ACTUALIZAR una URL
func (s *URLDataBase) UpdateURL(shortCode string, newOriginalURL string) (*URL, bool) {
	url, exists := s.urls[shortCode]
	if !exists {
		return nil, false
	}

	url.OriginalUrl = newOriginalURL
	url.UpdateAt = time.Now().UTC()
	s.urls[shortCode] = url

	return &url, true
}

// ELIMINAR una URL
func (s *URLDataBase) DeleteURL(shortCode string) bool {
	_, exists := s.urls[shortCode]
	if !exists {
		return false
	}

	delete(s.urls, shortCode)

	return true
}

// busca la URL, incrementa en 1 el contador y la devuelve
func (s *URLDataBase) GetAndIncrement(shortCode string) (*URL, bool) {
	url, exists := s.urls[shortCode]
	if !exists {
		return nil, false
	}

	url.AccessCount++
	s.urls[shortCode] = url

	return &url, true
}
