# REQ-006: Docker Support

## Requirement
The system must run in Docker containers with full functionality, supporting both migration execution and web UI.

## Acceptance Criteria
- [ ] Dockerfile builds successfully
- [ ] Container runs migrations via CLI
- [ ] Container runs web UI on port 8080
- [ ] Volume mounts for repository access work
- [ ] Volume mounts for configuration work
- [ ] Volume mounts for state persistence work
- [ ] Multi-architecture images (amd64, arm64)
- [ ] docker-compose.yml provides complete setup
- [ ] Container handles signals gracefully (SIGTERM)
- [ ] Minimal image size (< 100MB)

## Test Coverage
- `docker_test.go`: Integration tests with Docker
- Coverage: Must validate all acceptance criteria

## Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o git-migrator ./cmd/git-migrator

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
COPY --from=builder /app/git-migrator /usr/local/bin/
EXPOSE 8080
ENTRYPOINT ["git-migrator"]
```

### docker-compose.yml
```yaml
version: '3.8'
services:
  git-migrator:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./repos:/repos
      - ./config:/config
      - ./data:/data
    environment:
      - GIT_MIGRATOR_CONFIG=/config/migrator.yaml
    command: web --port 8080
```

## Status
- [x] Requirement defined
- [ ] Tests written
- [ ] Implementation complete
- [ ] All tests passing
