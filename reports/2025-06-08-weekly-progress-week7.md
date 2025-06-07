# SFConnect Backend Weekly Progress Report

## Week 7: Order Service - Commission Calculation & Inter-Service Communication

### Accomplishments

- **Commission Calculation:**
  - Added `commission` field to the `orders` table and Go model.
  - Implemented commission calculation (10% of total amount) upon order delivery confirmation.
  - Repository and service layers updated to support commission storage and retrieval.
- **Swagger/OpenAPI:**
  - Updated documentation to include commission in the Order schema.
- **Testing:**
  - Unit and integration tests updated to cover commission logic.
- **Inter-Service Communication (Preparation):**
  - Service structure ready for HTTP client integration to product-service for real price lookup (currently mocked).

### Notes

- All code compiles and tests pass.
- Ready for further inter-service communication and chatbot integration in Week 8.
