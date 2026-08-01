package main

import (
	"database/sql"
)

// Esta va ser nuestra conexion a la base de datos
type URLMySQLStore struct {
	db *sql.DB
}

// NewURLMySQLStore inicializa la conexión a MySQL
func NewURLMySQLStore(db *sql.DB) *URLMySQLStore {
	return &URLMySQLStore{
		db: db,
	}
}
