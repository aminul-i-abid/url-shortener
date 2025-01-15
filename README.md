# URL Shortener Service

## Features

- Shorten long URLs.
- Retrieve original URLs using the shortened URL.
- Update the original URL for a shortened URL.
- Delete a shortened URL.
- Get statistics for a shortened URL.
- Health check endpoint.
- API documentation with Swagger.
- Dockerized for easy deployment.

---

## API Endpoints

### Base URL

`http://localhost:8080/api/v1`

### Endpoints
![image](https://github.com/user-attachments/assets/05cf1741-9fab-44eb-8bd3-1aa8e15f3b4a)

---

## API Documentation

Swagger is used for API documentation.

- **Generate Swagger Documentation**:

  ```bash
  make generate-swagger
  ```

- **Access Swagger UI**:
  Once the application is running, navigate to `http://localhost:8080/swagger/index.html` to view the API documentation.

---

## Commands

### Makefile Commands

| Command                 | Description                              |
| ----------------------- | ---------------------------------------- |
| `make build`            | Build the Docker image                   |
| `make run`              | Run the application using Docker Compose |
| `make generate-swagger` | Generate Swagger API documentation       |

---
