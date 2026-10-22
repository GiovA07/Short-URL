# URL Shortener API

Una API RESTful de alto rendimiento y completa para el acortamiento de URLs, desarrollada en **Go** y persistida en **MySQL**. El proyecto sigue los principios de **Clean Architecture**, garantizando un código desacoplado, mantenible y listo para producción.

## Requisitos previos

![Go Version](https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge)
![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)


## Características Principales

- **Clean Architecture (Repository Pattern):** Desacoplamiento total entre los *handlers* HTTP y la capa de almacenamiento (`URLRepository`).
- **Enrutamiento Moderno:** Utiliza las características avanzadas de routing nativas de Go 1.22+ (`http.ServeMux`).
- **Validaciones & Seguridad:** Validacion de URLs con `net/url` y protección contra colisiones en enlaces personalizados.
- **Custom Short Codes:** Soporte para generar alias personalizados o códigos aleatorios.
- **Métricas e Integración:** Rastreo de contador de accesos en tiempo real e interfaz web minimalista integrada (`HTML/CSS/JS`).
- **Unit Testing:** Cobertura de pruebas automáticas con `httptest`.

---


##  Instalacion y Uso

1. **Clonar el repositorio**
```bash
git clone [https://github.com/GiovA07/Short-URL.git](https://github.com/GiovA07/Short-URL.git)
cd Short-URL
```

2. **Configurar Base de Datos**

Podes ejecutar el contenido del archivo sql/schema.sql en tu MySQL
```sql
CREATE DATABASE IF NOT EXISTS shortURL;
USE shortURL;

CREATE TABLE IF NOT EXISTS urls (
    id INT AUTO_INCREMENT PRIMARY KEY,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    access_count INT DEFAULT 0
);
```

3. **Configurar variables de entorno**
Crea un archivo .env en la raíz del proyecto utilizando en este ejemplo:
```bash
DB_USER=root
DB_PASSWORD=tu_contraseña
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=shortURL
```

4. **Instalar dependencias e iniciar el servidor**
```bash
go mod tidy
go run .
```

> El servidor se iniciará en http://localhost:8080.

## Endpoints de la API

| Metodo  | Ruta | Descripcion |
| ------------- |:-------------:|:-------------:|
| POST     | /shorten                   |Crea una nueva URL acortada    |
| GET      | /shorten/{shortCode}       |Obtiene los detalles de una URL sin redireccionar     |
| PUT      | /shorten/{shortCode}       |Actualiza la URL de destino de un shortCode existente.     |
| DELETE   | /shorten/{shortCode}       |Elimina una URL del sistema     |
| GET      | /shorten/{shortCode}/stats |evuelve las estadisticas de uso     |
| GET      | /{shortCode}               |Redirecciona el ususario a la URL original     |

## Tests
Para ejecutar la suite de pruebas unitarias en memoria:
```bash
go test -v
```