package main

import (
	"strconv"
	"sync"
	"time"
)

type URLRepository interface {
	Save(u *URL) error
	Get_URL(shortCode string) (*URL, bool)
	UpdateURL(shortCode string, newOriginalURL string) (*URL, bool)
	DeleteURL(shortCode string) bool
	GetAndIncrement(shortCode string) (*URL, bool)
}

// Esta va ser la base de datos por ahora hasta que conecte con mysql
type URLDataBase struct {
	mu     sync.RWMutex // Para concurrencia
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
func (s *URLDataBase) Save(u *URL) error {
	s.mu.Lock()         // Bloqueamos la escritura para otras peticiones
	defer s.mu.Unlock() // Nos aseguramos de desbloquear al salir

	u.ID = strconv.Itoa(s.nextID)
	s.nextID++
	s.urls[u.ShortCode] = *u

	return nil
}

// Obtener una URL
func (s *URLDataBase) Get_URL(shortCode string) (*URL, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, exists := s.urls[shortCode]
	if !exists {
		return nil, false
	}
	return &url, true
}

// ACTUALIZAR una URL
func (s *URLDataBase) UpdateURL(shortCode string, newOriginalURL string) (*URL, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.urls[shortCode]
	if !exists {
		return false
	}

	delete(s.urls, shortCode)

	return true
}

// busca la URL, incrementa en 1 el contador y la devuelve
func (s *URLDataBase) GetAndIncrement(shortCode string) (*URL, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, exists := s.urls[shortCode]
	if !exists {
		return nil, false
	}

	url.AccessCount++
	s.urls[shortCode] = url

	return &url, true
}
