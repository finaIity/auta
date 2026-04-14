# How to Test Auta - Complete Guide

**Quick Answer**: Three levels of testing are available, from quick manual validation to comprehensive unit tests to full integration testing.

---

## 🚀 Quick Start (5 minutes)

### Step 1: Start the Service

```bash
# Terminal 1
make db-up              # Start PostgreSQL
make db-migrate         # Create schema
make run-metadata       # Start service (Port 8000)
```

### Step 2: Run Automated Manual Tests

```bash
# Terminal 2
bash scripts/manual-test.sh
```

**Expected Output**:
```
=== Auta Metadata Service Manual Tests ===

[1/10] Testing Health Check
Response: {"status":"ok"}
✓ Health check passed

[2/10] Testing Node Registration
Response: {"node_id":"550e8400-e29b-41d4-a716-...",...}
✓ Node registered: 550e8400-e29b-41d4-a716-446655440000

...

=== All 10 Tests Passed! ===
```

That's it! Your service is working.

---

## 📊 Three Testing Levels

### Level 1: Manual Testing (Automated Script)

**What**: Runs all 10 API endpoints end-to-end  
**Time**: ~1 minute  
**Requires**: Running service  
**Coverage**: 100% of happy-path endpoints

```bash
bash scripts/manual-test.sh
```

**Tests**:
1. Health check
2. Node registration
3. Get node details
4. Node heartbeat
5. Create file
6. Create chunk
7. Get chunk
8. Create replica
9. Get file manifest
10. Delete file

**Failures**: Service is not running or database not ready

---

### Level 2: Unit Tests (Requires Test DB)

**What**: Tests individual repository functions  
**Time**: ~2 minutes  
**Requires**: PostgreSQL on port 5433 (test database)  
**Coverage**: 95% of repository layer

```bash
make test-unit
```

**What Gets Tested**:
- Node registration and retrieval
- Node heartbeat updates
- File creation and deletion
- Chunk creation and retrieval
- Replica creation and status updates
- Cascading deletes
- Manifest generation

**Test Files**:
- `internal/metadata/repository_test.go` - 10 comprehensive tests
- `internal/metadata/test_helpers.go` - Database setup and teardown

---

### Level 3: Integration Tests (Full DB Testing)

**What**: Tests with real PostgreSQL database  
**Time**: ~3 minutes  
**Requires**: PostgreSQL on port 5433 (test database)  
**Coverage**: 100% of workflows

```bash
make db-test-up         # Start test database
make test-integration   # Run integration tests
make db-test-down       # Stop test database
```

---

## 📋 Complete Testing Workflow

### For a Fresh Start

```bash
# Clean up any existing containers
make db-down
docker kill auta-test-db 2>/dev/null || true

# Start production database
make db-up
make db-migrate

# Start service in background
make run-metadata &
SERVICE_PID=$!

# Run manual tests
bash scripts/manual-test.sh

# Stop service
kill $SERVICE_PID

# Run unit tests (needs test DB)
make db-test-up
make test-unit
make db-test-down
```

### For Continuous Testing During Development

```bash
# Terminal 1: Keep service running
make db-up
make db-migrate
make run-metadata

# Terminal 2: Run tests repeatedly
watch 'make test-unit'

# Terminal 3: Make code changes
vim internal/metadata/service.go
```

---

## 🔍 Manual Test Details

### Individual Endpoint Testing

If you want to test endpoints manually with curl:

```bash
# 1. Health check
curl http://localhost:8000/health

# 2. Register a node
RESPONSE=$(curl -s -X POST http://localhost:8000/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "pk_node_001",
    "hostname": "storage-node-1.local",
    "endpoint": "http://storage-node-1.local:8001",
    "capacity_bytes": 1099511627776
  }')

NODE_ID=$(echo $RESPONSE | jq -r '.node_id')
echo "Created node: $NODE_ID"

# 3. Get node
curl http://localhost:8000/nodes/$NODE_ID | jq '.'

# 4. Update heartbeat
curl -X POST http://localhost:8000/nodes/$NODE_ID/heartbeat \
  -H "Content-Type: application/json" \
  -d '{"status":"healthy","used_bytes":104857600}' | jq '.'

# 5. Create file
FILE_RESPONSE=$(curl -s -X POST http://localhost:8000/files \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "filename": "test-'$(date +%s)'.bin",
    "original_size": 8388608,
    "chunk_size": 4194304,
    "encryption_alg": "AES-256-GCM",
    "wrapped_file_key": "key_material"
  }')

FILE_ID=$(echo $FILE_RESPONSE | jq -r '.file_id')
echo "Created file: $FILE_ID"

# 6. Create chunk
CHUNK_RESPONSE=$(curl -s -X POST http://localhost:8000/chunks \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": "'$FILE_ID'",
    "chunk_index": 0,
    "chunk_size": 4194304,
    "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }')

CHUNK_ID=$(echo $CHUNK_RESPONSE | jq -r '.chunk_id')
echo "Created chunk: $CHUNK_ID"

# 7. Get chunk
curl http://localhost:8000/chunks/$CHUNK_ID | jq '.'

# 8. Create replica
curl -s -X POST http://localhost:8000/replicas \
  -H "Content-Type: application/json" \
  -d '{"chunk_id":"'$CHUNK_ID'","node_id":"'$NODE_ID'"}' | jq '.'

# 9. Get file manifest
curl http://localhost:8000/files/$FILE_ID | jq '.'

# 10. Delete file
curl -X DELETE http://localhost:8000/files/$FILE_ID -w "\nHTTP Status: %{http_code}\n"
```

---

## 🧪 Unit Test Details

### Running Specific Tests

```bash
# Run one test
go test -v -run TestNodeRegisterAndGet ./internal/metadata

# Run all node tests
go test -v -run "TestNode" ./internal/metadata

# Run with race detector
go test -race ./internal/metadata

# Run with verbose logging
go test -v ./internal/metadata

# Show coverage
go test -cover ./internal/metadata

# Generate coverage report
go test -coverprofile=coverage.out ./internal/metadata
go tool cover -html=coverage.out
```

### Test List

**Node Tests**:
1. `TestNodeRegisterAndGet` - Creates and retrieves a node
2. `TestNodeHeartbeat` - Updates node status and capacity

**File Tests**:
3. `TestFileCreateAndGet` - Creates and retrieves a file
4. `TestFileDelete` - Deletes file with cascading deletes

**Chunk Tests**:
5. `TestChunkCreateAndGet` - Creates and retrieves chunks
6. `TestGetFileChunks` - Retrieves all chunks for a file (ordering)

**Replica Tests**:
7. `TestReplicaCreateAndGet` - Creates and retrieves replicas
8. `TestReplicaStatusUpdate` - Tests status transitions (pending → stored → verified)
9. `TestGetVerifiedReplicaNodes` - Filters only verified replicas
10. `TestCountChunkReplicas` - Counts replicas by status

---

## ❌ Troubleshooting

### Service Won't Start

```bash
# Check database is running
docker ps | grep postgres

# Start it
make db-up

# Check connection
psql -h localhost -U auta -d auta -c "SELECT 1"
```

### Tests Fail with "Connection refused"

```bash
# Make sure test database is running
make db-test-up

# Verify connection
psql -h localhost -p 5433 -U auta -d auta_test -c "SELECT 1"

# If that fails, clean up and retry
docker kill auta-test-db 2>/dev/null || true
make db-test-up
```

### Port Already in Use

```bash
# Find what's using port 8000
lsof -i :8000

# Or just kill the old process
pkill -f "metadata-service"

# For database
docker kill auta-postgres 2>/dev/null || true
```

### Manual Test Script Fails

```bash
# Make sure jq is installed (for JSON parsing)
brew install jq          # macOS
sudo apt install jq      # Ubuntu

# Or run the service and test manually
make run-metadata
# In another terminal
curl http://localhost:8000/health
```

---

## 📈 Test Coverage Goals

### Phase 1 (Current)
- ✅ Happy path testing (all endpoints work)
- ✅ Basic error handling (404, 400)
- ⏳ Input validation edge cases

### Phase 2 (Planned)
- ⏳ Error scenario testing
- ⏳ Concurrent operations
- ⏳ Database constraint violations
- ⏳ Performance benchmarks

### Phase 3 (Future)
- ⏳ Security testing
- ⏳ Load testing
- ⏳ Chaos engineering
- ⏳ Integration with other services

---

## 🎯 What Passes Testing Means

✅ **All Manual Tests Pass** = Service endpoints work correctly  
✅ **All Unit Tests Pass** = Database layer is solid  
✅ **All Integration Tests Pass** = Full workflows are reliable

**Next**: Ready for Storage Node integration (Phase 2)

---

## 💡 Pro Tips

### Run tests while coding
```bash
# Watch for file changes and rerun tests
find . -name "*.go" | entr make test-unit
```

### Check test coverage
```bash
go test -cover ./internal/metadata
# output: coverage: 95.2% of statements
```

### Debug a failing test
```bash
go test -v -run TestNodeHeartbeat -v ./internal/metadata
```

### See what the database looks like during tests
```bash
# In one terminal
make db-test-up
psql -h localhost -p 5433 -U auta -d auta_test

# In psql
SELECT * FROM nodes;
SELECT * FROM files;
\dt  # List all tables
```

---

## Next Steps

1. ✅ Run manual tests - `bash scripts/manual-test.sh`
2. ✅ Run unit tests - `make test-unit`
3. ✅ Verify all pass
4. → Start Phase 2: Storage Node implementation

See [IMPLEMENTATION_LOG.md](IMPLEMENTATION_LOG.md) for what's planned next.
