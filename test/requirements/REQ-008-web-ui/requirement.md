# REQ-008: Web UI

## Requirement
The system must provide a web-based user interface for monitoring and managing migrations.

## Acceptance Criteria
- [ ] Dashboard shows list of migrations
- [ ] Migration wizard allows creating new migrations
- [ ] Progress page shows real-time migration progress
- [ ] Configuration page allows viewing/editing settings
- [ ] Log viewer shows migration logs
- [ ] UI is responsive and works on mobile
- [ ] UI connects to WebSocket for real-time updates
- [ ] UI handles errors gracefully

## Test Coverage
- `e2e_test.go`: End-to-end tests using browser automation
- Coverage: Must validate all acceptance criteria

## UI Pages

### Dashboard (`/`)
- List of recent migrations
- Quick actions (new migration, view all)
- System status

### Migration Wizard (`/new`)
- Step 1: Select source type
- Step 2: Configure source repository
- Step 3: Configure target repository
- Step 4: Author/branch/tag mapping
- Step 5: Review and start

### Migration Progress (`/migration/:id`)
- Real-time progress bar
- Current step indicator
- Log viewer
- Stop/pause buttons

### Configuration (`/config`)
- View/edit configuration
- Save changes

## Status
- [x] Requirement defined
- [ ] Tests written
- [ ] Implementation complete
- [ ] All tests passing
