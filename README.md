# Microservicio en Golang con MongoDB

## 🧩 Objetivo

Este proyecto consiste en el desarrollo de un microservicio en Go (Golang) que implementa operaciones CRUD sobre una entidad seleccionada. Utiliza MongoDB como base de datos, pruebas unitarias y de integración, contenerización con Docker y orquestación con Docker Compose.


## 📦 Estructura del Proyecto

- `controllers/` – Endpoints REST.
- `services/` – Lógica de negocio.
- `repositories/` – Acceso a datos con `mongo-go-driver`.
- `config/` – Configuración y manejo de variables de entorno.
- `main.go` – Punto de entrada del microservicio.
- `Dockerfile` – Multi-Stage Build para la aplicación.
- `docker-compose.yml` – Orquestación de servicios.
- `tests/` – Pruebas unitarias e integración.
- `postman/` – Colección exportada de Postman.


## 🛠️ Requerimientos Técnicos

### 1. CRUD sobre una entidad
- CRUD completo con arquitectura modular.
- Uso de `mongo-go-driver`.

### 2. MongoDB en Docker
- MongoDB expuesto en el puerto `27017`.
- Volumen persistente para no perder datos.

### 3. Pruebas
- Pruebas unitarias con `Testify` o `GoMock`.
- Pruebas de integración completas.
- Cobertura con `go test -cover`.

### 4. Documentación con Postman
- Colección de endpoints con pruebas básicas incluidas.

### 5. Dockerización
- Dockerfile con Multi-Stage Build.
- Uso de `.dockerignore` y buenas prácticas.

### 6. Docker Compose
- Definición de servicios y redes en `docker-compose.yml`.

### 7. Publicación
- Imagen publicada en Docker Hub bajo el formato sugerido.

### 8. Diagrama de Infraestructura
- Diagrama visual detallado (componentes, redes, volúmenes, seguridad).


### Ejecutar pruebas y generar cobertura:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
