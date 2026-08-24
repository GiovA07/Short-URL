# URL Shortener API

Una API REST rapida para acortar URLs, construida en Golang. 

## Requisitos previos

- [Go](https://go.dev/) 1.22 o superior
- [MySQL](https://www.mysql.com/)

##  Instalacion y Uso

1. **Clonar el repositorio**
```bash
git clone [https://github.com/GiovA07/Short-URL.git](https://github.com/GiovA07/Short-URL.git)
cd Short-URL
```

2. **Configurar Base de Datos**

Podes ejecutar el contenido del archivo sql/schema.sql en tu MySQL

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

