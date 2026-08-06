CREATE DATABASE shortURL;
USE shortURL;
CREATE TABLE urls (
	ID int AUTO_INCREMENT PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url VARCHAR(2048) NOT NULL,
    created_at DATETIME DEFAULT current_timestamp, # si no mando una fecha de registro se pone la actual
    update_at DATETIME DEFAULT current_timestamp,
    access_count INT DEFAULT 0
)

SELECT * FROM urls;