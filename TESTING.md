# Testing Strategy for Auta Metadata Service

## Overview

Three-tier testing approach:

1. **Manual Testing** (Quick validation) - Start here
2. **Unit Tests** (Fast feedback) - Test individual functions
3. **Integration Tests** (Real database) - Test entire service

---

## Level 1: Manual Testing (5 minutes)

### Quick Start

```bash
# Terminal 1: Start the service
make db-up
make db-migrate
make run-metadata

# Terminal 2: Run tests
make test-manual  # (We'll create this)
```

### Manual Test Script

Run these commands to verify everything works:

```bash
# 1. Health Check
curl http://localhost:8000/health

# 2. Register a Node
curl -X POST http://localhost:8000/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "pk_node_001",
    "hostname": "storage-node-1.local",
    "endpoint": "http://storage-node-1.local:8001",
    "capacity_bytes": 1099511627776
  }'

# Save the node_id from the response, then use it below:
NODE_ID="<copy-node-id-here>"

# 3. Get Node Details
curl http://localhost:8000/nodes/$NODE_ID

# 4. Node Heartbeat
curl -X POST http://localhost:8000/nodes/$NODE_ID/heartbeat \
  -H "Content-Type: application/json" \
  -d '{
    "status": "healthy",
    "used_bytes": 104857600
  }'

# 5. Create File
curl -X POST http://localhost:8000/files \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "filename": "test-file.bin",
    "original_size": 8388608,
    "mime_type": "application/octet-stream",
    "chunk_size": 4194304,
    "encryption_alg": "AES-256-GCM",
    "wrapped_file_key": "encrypted_key_base64_here"
  }'

# Save the file_id from response, then use it below:
FILE_ID="<copy-file-id-here>"

# 6. Get File Manifest
curl http://localhost:8000/files/$FILE_ID

# 7. Create Chunk
curl -X POST http://localhost:8000/chunks \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": "'$FILE_ID'",
    "chunk_index": 0,
    "chunk_size": 4194304,
    "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }'

# Save chunk_id from response:
CHUNK_ID="<copy-chunk-id-here>"

# 8. Create Replica
curl -X POST http://localhost:8000/replicas \
  -H "Content-Type: application/json" \
  -d '{
    "chunk_id": "'$CHUNK_ID'",
    "node_id": "'$NODE_ID'"
  }'

# 9. List Replicas for Chunk
curl http://localhost:8000/replicas/chunk/$CHUNK_ID

# 10. Delete File (cleanup)
curl -X DELETE http://localhost:8000/files/$FILE_ID
```

---

## Level 2: Unit Tests (Repository Pattern)

### What's Tested

The repository layer has comprehensive unit tests covering:

- **Node Operations** (3 tests)
  - `TestNodeRegisterAndGet` - Registration and retrieval
  - `TestNodeHeartbeat` - Status and capacity updates
  - `TestListHealthyNodes` - Filtering and sorting

- **File Operations** (2 tests)
  - `TestFileCreateAndGet` - Creation and retrieval
  - `TestFileDelete` - Cascading delete with chunks and replicas

- **Chunk Operations** (2 tests)
  - `TestChunkCreateAndGet` - Creation and retrieval
  - `TestGetFileChunks` - Fetching all chunks for a file

- **Replica Operations** (3 tests)
  - `TestReplicaCreateAndGet` - Creation and retrieval
  - `TestReplicaStatusUpdate` - Status transitions (pending → stored → verified)
  - `TestGetVerifiedReplicaNodes` - Filtering verified replicas

### Running Unit Tests

```bash
# Run all unit tests
make test-unit

# Run with verbose output
go test -v ./internal/metadata -run "Test"

# Run specific test
go test -v ./internal/metadata -run "TestNodeRegisterAndGet"

# Run with coverage
go test -v -cover ./internal/metadata
```

### Test Files

- `internal/metadata/repository_test.go` - All 10 unit tests
- `internal/metadata/test_helpers.go` - Test setup and helpers

---

## Level 3: Integration Tests (Real Database)

### What's Tested

Full end-to-end flows with real PostgreSQL:

- File upload workflow (node → file → chunks → replicas)
- Manifest generation with verified replicas
- Cascading deletes
- Status transitions
- Replica verification

### Setup

```bash
# Create test database
make db-test-up

# Run integration tests
make test-integration

# Cleanup
make db-test-down
```

### Test Database

- Separate database: `auta_test` (port 5433)
- Automatic cleanup between tests
- Real PostgreSQL instance
- Full schema validation

---

## Quick Testing Checklist

### Before You Start

```bash
# Check Go is installed
go version

# Check Docker is running
docker version

# Check psql is available
psql --version
```

### Testing Sequence

**Phase 1: Manual Testing (5 mins)**
```bash
# Start the service
make db-up
make db-migrate
make run-metadata

# In another terminal, run automated manual tests
bash scripts/manual-test.sh
```

**Phase 2: Unit Testing (2 mins)**
```bash
# Run unit tests
make test-unit

# Expected: All 10 tests pass
```

**Phase 3: Integration Testing (3 mins)**
```bash
# Setup test database
make db-test-up

# Run integration tests
make test-integration

# Cleanup
make db-test-down
```

---

## Expected Results

### Manual Tests (10 endpoints)

```
✓ Health check
✓ Node registration
✓ Get node
✓ Node heartbeat
✓ File creation
✓ Chunk creation
✓ Get chunk
✓ Replica creation
✓ Get file manifest
✓ File deletion
```

### Unit Tests

```
✓ TestNodeRegisterAndGet
✓ TestNodeHeartbeat
✓ TestFileCreateAndGet
✓ TestFileDelete
✓ TestChunkCreateAndGet
✓ TestReplicaCreateAndGet
✓ TestReplicaStatusUpdate
✓ TestGetVerifiedReplicaNodes
✓ ... (all passing)
```

### Integration Tests

All workflows complete successfully with real database.

---

## Troubleshooting

### "Database connection refused"

```bash
# Make sure PostgreSQL is running
make db-up

# Wait for it to be ready
sleep 5

# Try migrating again
make db-migrate
```

### "Port 5432 already in use"

```bash
# Kill existing container
docker kill auta-postgres 2>/dev/null || true

# Start fresh
make db-down
make db-up
```

### "psql: command not found"

Install PostgreSQL tools:
```bash
# macOS
brew install postgresql

# Ubuntu/Debian
sudo apt-get install postgresql-client

# Fedora
sudo dnf install postgresql
```

### Tests timeout or fail

```bash
# Check PostgreSQL is healthy
psql -h localhost -U auta -d auta -c "SELECT 1"

# Increase timeout in tests if needed
go test -v -timeout 30s ./...
```

---

## Coverage Analysis

### Current Coverage

Expected coverage by component:

- Repository layer: ~95% (20/20 methods tested)
- Service handlers: ~70% (basic functionality, error paths in progress)
- Models: ~100% (pure data types)

### Running Coverage Reports

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out

# Show coverage by function
go tool cover -func=coverage.out
```

---

## Performance Testing

### Load Testing Endpoints

```bash
# Install Apache Bench (ab)
brew install ab                    # macOS
sudo apt-get install apache2-utils # Ubuntu

# Test endpoint performance
ab -n 1000 -c 10 http://localhost:8000/health

# Test with POST
ab -n 100 -c 5 -p request.json -T application/json http://localhost:8000/files
```

### Expected Performance

- Health check: <1ms
- Node operations: <5ms
- File operations: <10ms
- Chunk operations: <5ms
- Replica operations: <5ms

---

## Next Testing Phase

After initial testing passes, add:

- [ ] Error case testing (invalid inputs)
- [ ] Concurrent operation testing
- [ ] Database constraint testing
- [ ] Performance benchmarks
- [ ] Stress testing
- [ ] Security testing

---

## Tips for Developers

### Run tests while coding

```bash
# Watch mode (requires entr)
find . -name "*.go" | entr make test-unit
```

### Debug failing tests

```bash
# Run with extra logging
go test -v -run TestName ./internal/metadata

# Run with race detector
go test -race ./internal/metadata

# Run with coverage
go test -cover -run TestName ./internal/metadata
```

### Keep tests organized

```bash
# Group related tests
# ✓ TestNodeXxx
# ✓ TestFileXxx
# ✓ TestChunkXxx
# ✓ TestReplicaXxx
```

---

## Testing Best Practices

1. **Always clean up**: Tests should not affect each other
2. **Use meaningful names**: Test names describe what they test
3. **One assertion focus**: Each test validates one behavior
4. **Test the contract**: Test public API, not implementation
5. **Keep tests fast**: Avoid slow operations in unit tests
6. **Document weird cases**: Add comments for non-obvious tests
7. **Use table-driven tests**: For multiple similar test cases

---

## Continuous Integration

When ready, add CI/CD:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_PASSWORD: postgres
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: go test -v ./...
```

---

## Resources

- [Go Testing Docs](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Test Coverage](https://blog.golang.org/coverage)


