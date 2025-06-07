Here's a more detailed, step-by-step plan for implementing the SFConnect backend system, expanding on each phase and task.

## Detailed Step-by-Step Implementation Plan: SFConnect Backend System

This plan breaks down the project into weekly sprints, providing a clear roadmap for execution.

### Phase 1: Project Setup and Foundational Services (Weeks 1-3)

**Goal:** Establish the core development environment, project structure, and implement the User Service as the first functional microservice.

---

#### **Week 1: Environment Setup & Project Foundation**

**Objective:** Get the development environment ready, define project structure, and set up basic infrastructure components.

**Tasks:**

1.  **Go Development Environment Setup:**
    * Install Go (if not already installed).
    * Configure your IDE (VS Code with Go extensions, GoLand, etc.).
    * Set up GOPATH and other environment variables.
2.  **Version Control Initialization:**
    * Create a new Git repository (e.g., on GitHub or GitLab).
    * Clone the repository to your local machine.
    * Add a `.gitignore` file to exclude Go binaries, environment files, and IDE-specific files.
    * Make an initial commit: `git commit -m "Initial project setup and .gitignore"`
3.  **Project Directory Structure:**
    * Create a root directory for the project (e.g., `sfconnect-backend`).
    * Inside, create subdirectories for each planned microservice:
        * `./user-service`
        * `./product-service`
        * `./order-service`
        * `./chatbot-service` (optional for now)
    * Create a `.` `/.` `docker` directory for shared Docker configurations.
    * Create a `./docs` directory for project documentation.
4.  **Database (PostgreSQL) Setup:**
    * **Option A (Local Install):** Install PostgreSQL directly on your machine.
    * **Option B (Docker - Recommended):**
        * Create a `docker-compose.yml` file in the root `.` `/` `sfconnect-backend` directory.
        * Define a PostgreSQL service in `docker-compose.yml`:
            ```yaml
            version: '3.8'
            services:
              postgres_user:
                image: postgres:13
                environment:
                  POSTGRES_DB: user_db
                  POSTGRES_USER: user_service_user
                  POSTGRES_PASSWORD: password
                ports:
                  - "5432:5432"
                volumes:
                  - postgres_user_data:/var/lib/postgresql/data
              # ... potentially other postgres services for product/order
            volumes:
              postgres_user_data:
            ```
        * Start the database: `docker-compose up -d postgres_user`
5.  **Caching (Redis) Setup:**
    * Add a Redis service to your `docker-compose.yml`:
        ```yaml
        # ... (under services)
              redis:
                image: redis:6-alpine
                ports:
                  - "6379:6379"
                volumes:
                  - redis_data:/data
        # ... (under volumes)
              redis_data:
        ```
    * Start Redis: `docker-compose up -d redis` (or `docker-compose up -d` for all services).
6.  **Initial Go Modules:**
    * Inside each service directory (e.g., `user-service`), initialize a Go module:
        `cd user-service && go mod init github.com/your-username/sfconnect-backend/user-service`
    * Do this for all services.
7.  **Weekly Progress Report:** Draft an initial weekly report detailing environment setup and planned work.

---

#### **Week 2: User Service - Core Features**

**Objective:** Implement the fundamental user management functionalities, including registration, login, and JWT.

**Tasks:**

1.  **User Service: Database Schema Design:**
    * Define the `users` table schema for PostgreSQL (e.g., `id`, `email`, `password_hash`, `full_name`, `role`, `created_at`, `updated_at`).
    * Consider `uuid` for `id`s for better distributed system compatibility.
    * Write SQL migration scripts (or use a Go ORM with migrations like GORM + goose/migrate) to create the table.
2.  **User Service: Project Structure & Dependencies:**
    * Inside `user-service/`:
        * `main.go`: Entry point.
        * `internal/models`: Go structs representing database entities (e.g., `User`).
        * `internal/repository`: Database interaction logic (e.g., `UserRepository` interface and PostgreSQL implementation).
        * `internal/service`: Business logic (e.g., `UserService` for registration, login).
        * `internal/handler`: HTTP request handling (e.g., `UserHandler`).
        * `pkg/utils`: Utility functions (e.g., password hashing, JWT generation).
    * Install necessary Go libraries:
        * `github.com/gin-gonic/gin` for routing.
        * `github.com/lib/pq` for PostgreSQL driver.
        * `golang.org/x/crypto/bcrypt` for password hashing.
        * `github.com/dgrijalva/jwt-go` for JWT.
        * `github.com/joho/godotenv` for environment variables.
3.  **User Service: Database Connection & Repository:**
    * Implement a function to establish a connection to PostgreSQL.
    * Create `UserRepository` to handle CRUD operations for users.
4.  **User Service: Business Logic (`UserService`):**
    * **Registration:**
        * Hash password using `bcrypt`.
        * Save user to the database.
        * Assign a default role (e.g., `buyer`).
    * **Login:**
        * Retrieve user by email.
        * Compare hashed password using `bcrypt`.
        * Generate JWT token upon successful login, including `user_id` and `role` in claims.
5.  **User Service: API Endpoints (`UserHandler`):**
    * Set up a Gin router.
    * Define routes:
        * `POST /register`
        * `POST /login`
    * Implement handlers to call `UserService` methods and return JSON responses.
6.  **Initial Dockerfile for User Service:**
    * Create a `Dockerfile` in `user-service/` to build a Go executable and run it.
    * Update `docker-compose.yml` to include the `user-service`.
7.  **Testing (Unit Tests):**
    * Write basic unit tests for password hashing/comparison.
    * Write tests for JWT token generation and validation.
    * Test `UserService` methods (registration, login) with mock repositories.
8.  **Weekly Progress Report:** Update and submit your progress report.

---

#### **Week 3: User Service - Authentication, Authorization & Refinement**

**Objective:** Complete User Service with authentication middleware, authorization, and basic profile management.

**Tasks:**

1.  **User Service: JWT Middleware:**
    * Create a Gin middleware function (`AuthMiddleware`) to:
        * Extract JWT token from the `Authorization` header.
        * Validate the token.
        * Extract claims (`user_id`, `role`) and inject them into the Gin context for later use.
2.  **User Service: Authorization Middleware/Helper:**
    * Create an `AuthorizeRole` middleware/helper that takes a required role (`buyer`, `partner`, `admin`) and checks if the authenticated user's role matches.
3.  **User Service: Profile Management:**
    * Implement API endpoints for:
        * `GET /me` (or `GET /users/{id}`): Get authenticated user's profile.
        * `PUT /me` (or `PUT /users/{id}`): Update user's profile (e.g., name, non-sensitive info).
    * Apply `AuthMiddleware` to these routes.
4.  **User Service: Role Management (Basic):**
    * If time permits, implement an `admin`-only endpoint to change a user's role (e.g., `PUT /users/{id}/role`). This will require both authentication and authorization middleware.
5.  **User Service: Error Handling Enhancement:**
    * Refine error responses (e.g., consistent JSON error format with status codes).
    * Implement custom error types where appropriate.
6.  **User Service: Swagger/OpenAPI Documentation:**
    * Integrate `swaggo/swag` into the `user-service`.
    * Add Swagger comments to handlers, models, and routes.
    * Generate Swagger JSON/YAML and potentially set up a Swagger UI.
7.  **Testing (Integration Tests & Coverage):**
    * Write integration tests for API endpoints using a test database.
    * Aim for reasonable test coverage for `user-service`.
8.  **Consolidate Docker Compose:** Ensure `docker-compose.yml` can build and run `user-service`, PostgreSQL, and Redis together.
    * Test running all services with `docker-compose up --build`.
9.  **Weekly Progress Report:** Finalize `Week 3` report, summarizing User Service completion.

### Phase 2: Core Business Logic Services (Weeks 4-7)

**Goal:** Implement the Product Service and Order Service, focusing on inter-service communication and data consistency.

---

#### **Week 4: Product Service - Core Features & Caching**

**Objective:** Develop the Product Service with CRUD operations and Redis caching.

**Tasks:**

1.  **Product Service: Database Schema Design:**
    * Define the `products` table schema for PostgreSQL (e.g., `id`, `name`, `description`, `price`, `image_url`, `partner_id`, `stock`, `created_at`, `updated_at`).
    * Consider `partner_id` as a foreign key to the `users` table (conceptually, not necessarily physically if services are truly decoupled).
2.  **Product Service: Project Structure & Dependencies:**
    * Mirror `user-service` structure (`models`, `repository`, `service`, `handler`, `utils`).
    * Install Gin, PostgreSQL driver, Redis client (`github.com/go-redis/redis/v8`).
3.  **Product Service: Database Connection & Repository:**
    * Implement connection to PostgreSQL.
    * Create `ProductRepository` for CRUD operations.
4.  **Product Service: Redis Integration:**
    * Implement a caching layer (e.g., `ProductCache` interface).
    * Use Redis client to store and retrieve product data (e.g., product details by ID, product list).
    * Implement cache-aside pattern: check cache first, if not found, fetch from DB, then cache it. Invalidate cache on updates/deletes.
5.  **Product Service: Business Logic (`ProductService`):**
    * Implement methods for creating, updating, getting, listing, and deleting products.
    * Integrate caching logic here.
6.  **Product Service: API Endpoints (`ProductHandler`):**
    * Define routes for product CRUD.
    * Apply authentication middleware (from `user-service`'s JWT or a shared library).
    * Implement authorization for creating/updating/deleting products (only `partner` or `admin` roles).
7.  **Testing (Unit & Integration):**
    * Write tests for product CRUD operations.
    * Crucially, test Redis caching logic (cache hit/miss, invalidation).
8.  **Swagger Documentation:** Document Product Service APIs.
9.  **Weekly Progress Report.**

---

#### **Week 5: Order Service - Order Creation & Basic Flow**

**Objective:** Develop the Order Service with order creation functionality and initial status management.

**Tasks:**

1.  **Order Service: Database Schema Design:**
    * Define `orders` table (e.g., `id`, `user_id`, `total_amount`, `status`, `created_at`, `updated_at`).
    * Define `order_items` table (e.g., `id`, `order_id`, `product_id`, `quantity`, `unit_price`).
    * Define `order_status_history` (optional but good for tracking: `order_id`, `status`, `timestamp`).
2.  **Order Service: Project Structure & Dependencies:**
    * Mirror previous service structures.
    * Install Gin, PostgreSQL driver.
3.  **Order Service: Database Connection & Repository:**
    * Implement connection to PostgreSQL.
    * Create `OrderRepository` for order and order item CRUD.
4.  **Order Service: Business Logic (`OrderService`):**
    * **Create Order:**
        * Accept product IDs and quantities.
        * (Conceptual) Interact with Product Service to verify product existence and retrieve prices. *For now, simplify by assuming product data is valid or use mock data.*
        * Calculate total amount.
        * Create `order` and `order_item` records in a transaction.
        * Set initial status (e.g., "Pending").
5.  **Order Service: API Endpoints (`OrderHandler`):**
    * Define routes:
        * `POST /orders` (for buyers to create orders).
    * Apply authentication middleware.
6.  **Testing (Unit & Integration):**
    * Test order creation, ensuring atomic operations for order and order items.
7.  **Swagger Documentation:** Document Order Service APIs.
8.  **Weekly Progress Report.**

---

#### **Week 6: Order Service - Status Updates & Partner Flow**

**Objective:** Enhance Order Service with partner-specific functionalities and detailed status management.

**Tasks:**

1.  **Order Service: Status Update Logic:**
    * Define possible order statuses (e.g., "Pending", "Processing", "ReadyForDelivery", "Shipped", "Delivered", "Canceled").
    * Implement methods to update order status. Ensure valid status transitions.
2.  **Order Service: Partner Interaction:**
    * **Assign Order (Conceptual):** Implement a placeholder for assigning orders to partners (e.g., an `admin`- or `system`-triggered update). The requirement states SendoFarm processes and transfers, so this might be an internal API or a direct update.
    * **Partner Confirmation:**
        * API endpoint `PUT /orders/{id}/confirm-ready` (for partners).
        * Requires `partner` role authorization.
        * Updates status to "ReadyForDelivery" or similar.
3.  **Order Service: Buyer Confirmation:**
    * API endpoint `PUT /orders/{id}/confirm-delivery` (for buyers).
    * Requires `buyer` role authorization.
    * Updates status to "Delivered".
4.  **Order Service: Listing/Viewing Orders:**
    * `GET /orders`: List all orders (e.g., for admin/SendoFarm).
    * `GET /orders/my-orders`: List orders for the authenticated user (buyer/partner).
    * `GET /orders/{id}`: Get order details.
    * Apply appropriate authentication and authorization.
5.  **Testing:**
    * Test all status transition APIs with different roles.
    * Test order listing/details for different user types.
6.  **Swagger Documentation:** Update Swagger for new endpoints.
7.  **Weekly Progress Report.**

---

#### **Week 7: Order Service - Commission Calculation & Inter-Service Communication**

**Objective:** Finalize Order Service with commission calculation and refine inter-service communication.

**Tasks:**

1.  **Order Service: Commission Calculation:**
    * Upon `Delivered` status confirmation by the buyer, implement logic to calculate commission for the assigned partner.
    * Store commission (e.g., in a separate `commissions` table or as a field in `orders`).
    * Define commission rate (e.g., a configurable constant).
2.  **Inter-service Communication Strategy (Refinement):**
    * **Option 1 (Direct HTTP Calls - Simplest for Internship):** If Product Service needs to be queried by Order Service, implement HTTP clients (e.g., using `net/http` or a library like `resty`) within Order Service to call Product Service APIs.
    * **Option 2 (Message Broker - More Robust, if time permits):** Briefly explore using a message broker (like RabbitMQ or Kafka) for asynchronous communication, e.g., for order events. (Likely out of scope for initial internship but good to understand).
    * *Focus for this plan:* Stick to direct HTTP calls for simplicity, but acknowledge the alternative.
    * **Error Handling for Inter-Service Calls:** Implement robust error handling (retries, circuit breakers - simplified) for calls between services.
3.  **Testing:**
    * Test commission calculation thoroughly.
    * Test inter-service communication (mocking the external service if needed for unit tests, or using integration tests with running services).
4.  **Database Optimization (Initial Pass):**
    * Review SQL queries in all services for potential performance issues.
    * Add indexes to frequently queried columns (e.g., `email` in users, `product_id` in order items, `user_id` in orders, `status` in orders).
5.  **Weekly Progress Report.**

### Phase 3: Advanced Features & Refinement (Weeks 8-10)

**Goal:** Implement optional features (Chatbot), refine existing services, and prepare for comprehensive testing.

---

#### **Week 8: Chatbot Service - Basic NLU & Order Status Retrieval**

**Objective:** Implement the Chatbot Service with basic natural language understanding and integration with Order Service.

**Tasks:**

1.  **Chatbot Service: Project Structure & Dependencies:**
    * Standard service structure.
    * No complex NLP libraries needed initially; focus on pattern matching.
2.  **Chatbot Service: API Endpoint:**
    * `POST /chatbot/query` (takes a text `message` from the user).
3.  **Chatbot Service: Basic NLU (Order ID Extraction):**
    * Use regular expressions to extract order IDs from the user's question (e.g., "Đơn hàng #123 của tôi", "order 123", "mã đơn hàng 456").
4.  **Chatbot Service: Order Service Integration:**
    * Create an HTTP client within Chatbot Service to call the Order Service's `GET /orders/{id}` API.
    * Handle cases where the order ID is not found or the Order Service returns an error.
5.  **Chatbot Service: Response Generation:**
    * Based on the data retrieved from Order Service, construct a natural language response in Vietnamese (e.g., "Đơn hàng #123 hiện đang vận chuyển. Dự kiến sẽ được giao trước ngày 10/06/2025.").
6.  **Testing:**
    * Test order ID extraction with various input phrases.
    * Test the full chatbot flow, mocking the Order Service response.
7.  **Swagger Documentation:** Document Chatbot Service APIs.
8.  **Weekly Progress Report.**

---

#### **Week 9: System Refinement & SOLID Principles**

**Objective:** Review and refactor existing code to adhere to SOLID principles, clean code, and Go best practices.

**Tasks:**

1.  **Code Review & Refactoring (All Services):**
    * **Single Responsibility Principle (SRP):** Ensure each function/module has one clear responsibility.
    * **Open/Closed Principle (OCP):** Design modules to be open for extension, but closed for modification.
    * **Liskov Substitution Principle (LSP):** Ensure derived types are substitutable for their base types (less common explicitly in Go, more about interfaces).
    * **Interface Segregation Principle (ISP):** Create smaller, client-specific interfaces rather than large, monolithic ones.
    * **Dependency Inversion Principle (DIP):** Depend on abstractions (interfaces) rather than concrete implementations (e.g., inject `UserRepository` interface into `UserService`).
    * **Clean Code:**
        * Meaningful variable and function names.
        * Consistent formatting (use `gofmt`).
        * Avoid magic numbers/strings.
        * Reduce code duplication.
2.  **Error Handling Review:**
    * Ensure consistent and robust error handling across all services.
    * Implement custom error types where semantic meaning is important.
3.  **Configuration Management:**
    * Consolidate environment variable loading (e.g., using `godotenv` in `main.go` of each service).
    * Externalize all configuration parameters.
4.  **Logging:**
    * Implement basic structured logging (e.g., using `log` package or `zap`/`logrus` if desired) for debugging and monitoring.
5.  **Testing Coverage Improvement:**
    * Identify areas with low test coverage and add more unit/integration tests.
    * Ensure critical paths are well-tested.
6.  **Weekly Progress Report.**

---

#### **Week 10: Performance & DevOps Readiness**

**Objective:** Conduct initial performance checks, ensure Docker readiness, and consider CI/CD.

**Tasks:**

1.  **Database Indexing & Query Optimization:**
    * Run `EXPLAIN ANALYZE` on critical database queries to identify bottlenecks.
    * Add or optimize indexes based on query patterns.
2.  **Caching Strategy Review:**
    * Review Redis caching implementation. Are all frequently read data cached? Is invalidation handled correctly?
    * Consider cache eviction policies if relevant.
3.  **Docker Readiness:**
    * Review and optimize Dockerfiles for each service (multi-stage builds for smaller images).
    * Ensure `docker-compose.yml` can build and run all services seamlessly.
    * Test network communication between Docker containers for services.
4.  **CI/CD Pipeline (Optional but Highly Recommended):**
    * **If using GitHub:** Set up a `.github/workflows/main.yml` for GitHub Actions.
    * **If using GitLab:** Set up a `.gitlab-ci.yml`.
    * Define stages: `build`, `test`, `deploy` (placeholder for now).
    * Configure it to run `go test ./...` on every push to `main` branch.
5.  **Basic Load Testing:**
    * Use a simple tool like `ab` (ApacheBench) or Postman's collection runner to hit a few endpoints on each service to get a rough idea of performance.
6.  **Weekly Progress Report.**

### Phase 4: Deployment & Documentation (Weeks 11-12)

**Goal:** Finalize all technical and soft skill requirements, prepare for final presentation.

---

#### **Week 11: Comprehensive Documentation & Reporting**

**Objective:** Compile all technical documentation, preparing for the final report.

**Tasks:**

1.  **Technical Design Document (TDD) Creation:**
    * Start writing the main technical report (refer to `A practical guide to writing technical specs` reference).
    * **Introduction:** Project overview, problem statement.
    * **Goals & Objectives:** What was achieved.
    * **Architecture:** Detailed microservice architecture diagram.
        * Explain each service's role.
        * Describe data flow between services.
    * **Technology Stack:** List all technologies and why they were chosen.
    * **Detailed Design (for each service):**
        * API Endpoints (using Swagger output as reference).
        * Database Schema.
        * Key business logic flows.
        * Authentication & Authorization mechanisms.
        * Caching strategy.
        * Error handling strategy.
    * **Deployment Strategy:** Explain Docker setup, `docker-compose`, and potential Kubernetes considerations.
    * **Testing Approach:** Describe unit and integration testing.
    * **Future Considerations/Improvements:** What could be done next.
2.  **API Documentation Review:**
    * Ensure all Swagger documentation is up-to-date and accurate for all services.
    * Verify that it's easy to understand and use.
3.  **Weekly Progress Report.**

---

#### **Week 12: Final Review, Presentation & Handover**

**Objective:** Finalize the project, prepare the presentation, and ensure all requirements are met.

**Tasks:**

1.  **Final Code Review:**
    * Conduct a final self-review of all code.
    * Ensure all comments are clear and sufficient.
    * Check for any remaining TODOs.
2.  **Comprehensive Testing:**
    * Perform a final round of end-to-end testing of the entire system.
    * Verify all user flows (buyer registration -> product search -> order creation -> partner confirmation -> buyer delivery confirmation).
    * Test chatbot integration.
3.  **Refine Technical Design Document:**
    * Proofread the TDD for clarity, completeness, and accuracy.
    * Ensure it meets the `Yêu cầu kỹ năng mềm` and `Tiêu chí đánh giá` criteria.
4.  **Prepare Final Presentation:**
    * Create slides summarizing the project: problem, solution, architecture, key features, technical challenges, achievements, and future work.
    * Prepare a live demonstration of the working system.
    * Practice explaining technical decisions and the system's architecture.
5.  **Project Repository Clean-up:**
    * Ensure the `README.md` in the root repository is comprehensive, including:
        * How to set up and run the project locally (using `docker-compose`).
        * API endpoint examples.
        * Links to Swagger UI for each service.
6.  **Final Weekly Progress Report.**
7.  **Handover:** Prepare for the final discussion and handover to the mentor.

This detailed plan provides a robust framework for successfully completing the SFConnect backend internship project. Remember to communicate regularly with your mentor and adapt the plan as needed based on feedback and challenges encountered.