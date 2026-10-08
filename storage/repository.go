package storage

import (
	"time"
)

// La estructura principal de la URL
type URL struct {
	ID          string    `json:"id"`
	OriginalUrl string    `json:"url"`
	ShortCode   string    `json:"shortCode"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
	AccessCount int       `json:"accessCount"`
}

/*
URLRepository permite desacoplar la logica (handlers) del almacenamiento de datos.
Asi los handlers solo saben que existen los metodos (ej save) pero no donde se estan guardando.
*/
type URLRepository interface {
	Save(u *URL) error
	Get_URL(shortCode string) (*URL, bool)
	UpdateURL(shortCode string, newOriginalURL string) (*URL, bool)
	DeleteURL(shortCode string) bool
	GetAndIncrement(shortCode string) (*URL, bool)
}
