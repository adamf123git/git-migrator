# REQ-019: WebSocket Progress

## Requirement
The system must provide real-time progress updates via WebSocket for running migrations.

## Acceptance Criteria
- [ ] WebSocket endpoint available at `/ws/progress/:id`
- [ ] Clients receive progress events in real-time
- [ ] Progress events include: status, percentage, current step, errors
- [ ] Connection handles client disconnect gracefully
- [ ] Multiple clients can connect to same migration
- [ ] Invalid migration ID returns appropriate error
- [ ] WebSocket closes when migration completes

## Test Coverage
- `websocket_test.go`: Validates all acceptance criteria
- Coverage: Must be 90%+ for this requirement

## WebSocket Protocol

### Connection
```
ws://localhost:8080/ws/progress/<migration-id>
```

### Message Format (Server → Client)
```json
{
  "type": "progress",
  "data": {
    "migrationId": "uuid",
    "status": "running",
    "percentage": 45,
    "currentStep": "Processing commit 450/1000",
    "totalCommits": 1000,
    "processedCommits": 450,
    "errors": []
  }
}
```

### Event Types
- `progress` - Progress update
- `completed` - Migration completed successfully
- `failed` - Migration failed
- `error` - Error occurred

## Status
- [x] Requirement defined
- [ ] Tests written
- [ ] Implementation complete
- [ ] All tests passing
