# Simple HTTP Proxy

A lightweight HTTP proxy server built in Go for forwarding requests and storing responses.

---

## Features

Forward HTTP requests to any URL with custom headers.
Returns proxied response status, headers, and body length in JSON.
Stores requests and responses in memory with unique IDs.
Built-in Swagger API documentation.

### Getting Started
1. Prerequisites
Docker & Docker Compose
Go 1.22+ (for local development)

2. Setup
Copy .env.example to .env and configure environment variables if needed.
Build and start the container:
make build
make up

3. Visit the Swagger UI:
[Docker & Docker Compose
Go 1.22+ (for local development)](http://localhost:3333/swagger/index.html)

### Stopping the Server
make down

### Viewing Logs
make logs

### Project Structure
main.go: Entry point for the application. Starts the HTTP server.
/internal: Contains internal modules not intended for external import.
/internal/app: Initialization code (config, logging, etc.).
/internal/cache: Caching logic and in-memory storage.
/internal/config: Environment and configuration parsing.
/internal/domain: Core business logic independent of infrastructure.
/internal/handler: HTTP request handlers for the proxy endpoints.
/internal/service: Implements core application functionality.
/pkg: Reusable packages for external use.


### Usage Example
Send a POST request to /proxy:

POST /proxy
Content-Type: application/json

{
  "method": "GET",
  "url": "https://example.com",
  "headers": {
    "User-Agent": "GoProxy"
  }
}


Response:
{
  "id": "uuid-generated-id",
  "status": 200,
  "headers": {
    "Content-Type": "text/html"
  },
  "length": 1234
}



### Docker
Dockerfile:

FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o proxy-server main.go

EXPOSE 3333

CMD ["./proxy-server"]



docker-compose.yml:

version: "3.9"
services:
  proxy:
    build: .
    ports:
      - "3333:3333"


### Dependencies
chi – HTTP router
zerolog – Fast structured logging
Swagger – API documentation
uuid – Unique request IDs


### License
MIT License
