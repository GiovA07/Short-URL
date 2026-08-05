package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Esta va ser nuestra conexion a la base de datos
type URLMySQLStore struct {
	db *sql.DB
}

// NewURLMySQLStore recibe las credenciales y devuelve la conexion
func NewURLMySQLStore() (*URLMySQLStore, error) {
	// Carga las variables del archivo .env que va a tener mis credenciales
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Leemos cada variable
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Armamos el DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// db.Ping() es el que abre la puerta a la base de datos
	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("MySQL Conectadoo!")
	return &URLMySQLStore{db: db}, nil
}
