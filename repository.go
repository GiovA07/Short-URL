package main

import "strconv"

// Esta va ser la base de datos por ahora hasta que conecte con mysql
type URLDataBase struct {
	urls   map[string]URL
	nextID int
}

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
