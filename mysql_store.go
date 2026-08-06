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

func (s *URLMySQLStore) Save(u *URL) error {
	// 1. La consulta SQL.
	query := `
        INSERT INTO urls (short_code, original_url, created_at, update_at, access_count) 
        VALUES (?, ?, ?, ?, ?)
    `

	// 2. Ejecutamos la consulta.
	_, err := s.db.Exec(
		query,
		u.ShortCode,
		u.OriginalUrl,
		u.CreateAt,
		u.UpdateAt,
		u.AccessCount,
	)

	// 3. Manejamos los errores
	if err != nil {
		return fmt.Errorf("error al insertar URL en MySQL: %w", err)
	}

	return nil
}

// Obtener una URL
func (s *URLMySQLStore) Get_URL(shortCode string) (*URL, bool) {
	// 1. Preparamos la variable donde vamos a guardar los datos
	var u URL

	// 2. Armamos la consulta
	query := `
        SELECT short_code, original_url, created_at, update_at 
		FROM urls
		WHERE short_code = ? 
    `

	// 3. Ejecutamos la consulta usando QueryRow (porque esperamos 1 sola fila)
	err := s.db.QueryRow(query, shortCode).Scan(
		&u.ShortCode,
		&u.OriginalUrl,
		&u.CreateAt,
		&u.UpdateAt,
	)

	//4. Verificamos
	if err != nil {
		// Este error es especial de Go, basicamente es que no encontro nada.
		if err == sql.ErrNoRows {
			return nil, false
		}

		// Si es otro tipo de error (puede ser que se cayo la red, etc)
		log.Printf("Error al buscar URL en MySQL: %v", err)
		return nil, false
	}

	return &u, true
}
