# SFConnect Backend Weekly Progress Report

## Week 8: Chatbot Service - Basic NLU & Order Status Retrieval

### Accomplishments

- **Chatbot Service Scaffolded:**
  - Created chatbot-service with standard Go microservice structure.
  - Added Gin HTTP server and `/chatbot/query` endpoint.
- **Order ID Extraction:**
  - Implemented regex-based extraction of order IDs from user messages (supports Vietnamese and English patterns).
  - Unit tests for extraction logic.
- **Order Status Retrieval (Stub):**
  - Endpoint returns canned status for recognized order IDs (integration with order-service to be added).
- **Swagger/OpenAPI:**
  - Documented chatbot API and request/response schemas.
- **Testing:**
  - Unit and handler tests for chatbot endpoint and extraction logic.

### Notes

- Ready for integration with order-service for real order status retrieval.
- All code compiles and tests pass.
