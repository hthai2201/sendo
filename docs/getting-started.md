# SFConnect Backend – Getting Started Guide

This guide will help you set up and run the entire SFConnect backend system locally using Docker Compose. All services (user, product, order, chatbot) are containerized and orchestrated with PostgreSQL and Redis.

---

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop) (with Docker Compose)
- [Go](https://golang.org/) (if you want to run/test services outside Docker)

---

## 1. Clone the Repository

```sh
git clone <your-repo-url>
cd sendo
```

---

## 2. Environment Variables

Each service has a `.env` file in its directory. These are pre-filled with sensible defaults for local development. You can adjust secrets and DB names as needed.

---

## 3. Build & Start All Services

From the project root:

```sh
docker compose up --build
```

- This will build and start all services, PostgreSQL, and Redis.
- Services will be available on the following ports by default:
  - user-service: http://localhost:8081
  - product-service: http://localhost:8082
  - order-service: http://localhost:8083
  - chatbot-service: http://localhost:8084

---

## 4. API Documentation

Each service exposes Swagger/OpenAPI docs:

- user-service: http://localhost:8081/swagger/index.html
- product-service: http://localhost:8082/swagger/index.html
- order-service: http://localhost:8083/swagger/index.html
- chatbot-service: http://localhost:8084/swagger/index.html

---

## 5. Running Tests

To run all tests for a service (from the service directory):

```sh
go test ./...
```

---

## 6. Useful Commands

- Stop all containers:
  ```sh
  docker compose down
  ```
- Rebuild a single service:
  ```sh
  docker compose build <service-name>
  ```
- View logs for a service:
  ```sh
  docker compose logs -f <service-name>
  ```

---

## 7. Load Testing

- See the README or docs for instructions on using `ab` (ApacheBench) or Postman for load testing.

---

## 8. Troubleshooting

- Ensure ports 8081–8084, 5432 (Postgres), and 6379 (Redis) are free.
- If you change `.env` files, restart the affected service(s).
- For database issues, check the `migrations/` folder in each service for schema setup.

---

## 9. Further Information

- See `docs/technical-design-document.md` for architecture and design details.
- See `docs/implementation-plan-and-milestone-requirements.md` for the full project plan.
- For questions, see the handover checklist or contact the maintainer.

---

**You're ready to start developing and testing the SFConnect backend!**
