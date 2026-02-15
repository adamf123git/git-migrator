# REQ-018: REST API

## Requirement
The system must provide a REST API for managing migrations, viewing status, and configuring the application.

## Acceptance Criteria
- [ ] API server starts on configurable port
- [ ] `GET /api/migrations` returns list of migrations
- [ ] `POST /api/migrations` starts new migration
- [ ] `GET /api/migrations/:id` returns migration status
- [ ] `POST /api/migrations/:id/stop` stops running migration
- [ ] `GET /api/config` returns current configuration
- [ ] `POST /api/config` updates configuration
- [ ] `GET /api/repos/analyze` analyzes source repository
- [ ] All endpoints return proper JSON responses
- [ ] All endpoints return proper HTTP status codes
- [ ] API handles errors gracefully

## Test Coverage
- `api_test.go`: Validates all acceptance criteria
- Coverage: Must be 90%+ for this requirement

## API Specification

### Endpoints

```
GET  /api/health              # Health check
GET  /api/migrations          # List all migrations
POST /api/migrations          # Start new migration
GET  /api/migrations/:id      # Get migration status
POST /api/migrations/:id/stop # Stop migration
GET  /api/config              # Get configuration
POST /api/config              # Update configuration
POST /api/repos/analyze       # Analyze source repository
```

### Response Format

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

### Error Format

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid configuration"
  }
}
```

## Status
- [x] Requirement defined
- [ ] Tests written
- [ ] Implementation complete
- [ ] All tests passing
