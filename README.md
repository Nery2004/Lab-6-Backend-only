# LaLigaTracker Backend

Este es el backend desarrollado en Go para el proyecto LaLigaTracker. Proporciona una API para interactuar con la página `LaLigaTracker.html`, permitiendo gestionar y consultar información relacionada con La Liga.

## Tecnologías Utilizadas

- **Go**: Lenguaje de programación principal.
- **Docker**: Para contenedorización y despliegue.
- **Gorilla Mux**: Router HTTP para manejar las rutas de la API.
- **Base de datos**: (Especificar si usa PostgreSQL, MySQL, etc.)

## Instalación

Para ejecutar el backend en un entorno de desarrollo local, sigue estos pasos:

### 1. Clonar el repositorio
```sh
$ git clone https://github.com/Nery2004/Lab-6-Backend-only.git
$ cd Lab-6-Backend-only
```

### 2. Instalar dependencias
```sh
$ go mod tidy
```

### 3. Configurar variables de entorno
Crea un archivo `.env` en la raíz del proyecto y agrega lo siguiente:
```ini
DB_HOST=your_database_host
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=your_database_name
PORT=8080
```
(Sustituye los valores según tu configuración.)

### 4. Ejecutar el servidor
```sh
$ go run main.go
```
El servidor se ejecutará en `http://localhost:8080` (según la variable `PORT`).

## Uso con Docker

### Construir la imagen
```sh
$ docker build -t laliga-tracker-backend .
```

### Ejecutar el contenedor
```sh
$ docker run -p 8080:8080 --env-file .env laliga-tracker-backend
```
Esto levantará el backend en el puerto `8080`.

## Endpoints Principales

| Método | Endpoint           | Descripción                        |
|---------|-------------------|--------------------------------|
| GET     | `/teams`          | Obtiene la lista de equipos. |
| GET     | `/matches`        | Obtiene los partidos recientes. |
| POST    | `/register`       | Registra un nuevo usuario. |
| POST    | `/login`          | Inicia sesión de usuario. |

(Asegúrate de actualizar la tabla con los endpoints correctos)

## Contribuciones
Si deseas contribuir, por favor sigue estos pasos:

1. Crea un *fork* del repositorio.
2. Crea una nueva rama (`git checkout -b feature-nueva`).
3. Realiza los cambios y haz un *commit* (`git commit -m 'Agregada nueva funcionalidad'`).
4. Sube los cambios (`git push origin feature-nueva`).
5. Abre un *Pull Request*.

---
**Autor:** Nery2004
Imagenes:

![{D6AD95EB-D493-424D-B804-4440A89B0CD2}](https://github.com/user-attachments/assets/3eae895f-f749-4efa-a37d-ec15cadd868d)

![{0A6C0F40-BF11-448E-B443-223AC4EC86A4}](https://github.com/user-attachments/assets/973112df-0feb-475d-9cc5-578382c4ff99)

![{75A9CA51-5E72-4357-B9BA-A8FB86AE51AD}](https://github.com/user-attachments/assets/1f4b5c98-552c-4d22-b124-c05aaf0efab0)

![{ADD8B347-58F8-426E-B7A2-AAC8D0CB5505}](https://github.com/user-attachments/assets/1b813651-0f1b-42fe-8506-77a1b92bc362)
