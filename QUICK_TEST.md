# Step-by-Step Testing Guide

Follow these steps in order to test the Auta Metadata Service. Each section takes about 5 minutes.

---

## 🟢 Step 1: Manual Testing (Quickest - 5 minutes)

This tests all 10 endpoints with an automated script.

### 1a: Start the Service

Open Terminal 1:

```bash
cd /path/to/auta

# Start the database
make db-up

# Wait for it to be ready
sleep 3

# Run migrations to create schema
make db-migrate

# Start the metadata service
make run-metadata
```

**Expected Output**:
```
Waiting for PostgreSQL to be ready...
Running migrations...
Starting metadata service on port 8000
```

### 1b: Run the Automated Test

Open Terminal 2:

```bash
cd /path/to/auta

# Run the automated test script
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

[3/10] Testing Get Node
Response: {"node_id":"550e8400...",...}
✓ Get node passed

... (continues for all 10 tests)

=== All 10 Tests Passed! ===
```

### 1c: Verify Results

✅ **If you see "All 10 Tests Passed!"** → Your service is working!  
❌ **If you see errors** → Check troubleshooting section below

---

## 🟢 Step 2: Unit Tests (Comprehensive - 2 minutes)

This tests the database layer with 10 focused unit tests.

### 2a: Setup Test Database

Open Terminal 3:

```bash
cd /path/to/auta

# Create test database on port 5433
make db-test-up

# Wait for it to start
sleep 3
```

**Expected Output**:
```
Setting up test database...
(may see some Docker output)
```

### 2b: Run Unit Tests

```bash
# Run all unit tests
make test-unit
```

**Expected Output**:
```
=== RUN   TestNodeRegisterAndGet
--- PASS: TestNodeRegisterAndGet (0.05s)
=== RUN   TestNodeHeartbeat
--- PASS: TestNodeHeartbeat (0.04s)
=== RUN   TestFileCreateAndGet
--- PASS: TestFileCreateAndGet (0.03s)
... (more tests)
PASS
ok    github.com/anomalyco/auta/internal/metadata    0.524s
```

### 2c: Cleanup Test Database

```bash
# Stop the test database
make db-test-down
```

### 2d: Verify Results

✅ **If you see "PASS" at the end** → All unit tests passed!  
❌ **If you see "FAIL"** → Check troubleshooting section below

---

## 📊 Step 3: Manual Testing with curl (Optional - 10 minutes)

If you want to see exactly what's happening with each endpoint:

### 3a: Test Health Endpoint

```bash
curl http://localhost:8000/health | jq '.'
```

**Expected Response**:
```json
{
  "status": "ok"
}
```

### 3b: Register a Storage Node

```bash
RESPONSE=$(curl -s -X POST http://localhost:8000/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "test_pk_001",
    "hostname": "storage-node-1.local",
    "endpoint": "http://storage-node-1.local:8001",
    "capacity_bytes": 1099511627776
  }')

echo $RESPONSE | jq '.'

# Save the node ID for later
NODE_ID=$(echo $RESPONSE | jq -r '.node_id')
echo "Node ID: $NODE_ID"
```

**Expected Response**:
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "public_key": "test_pk_001",
  "hostname": "storage-node-1.local",
  "endpoint": "http://storage-node-1.local:8001",
  "status": "healthy",
  "capacity_bytes": 1099511627776,
  "used_bytes": 0,
  "created_at": "2026-04-14T12:00:00Z",
  "updated_at": "2026-04-14T12:00:00Z"
}
```

### 3c: Create a File

```bash
FILE_RESPONSE=$(curl -s -X POST http://localhost:8000/files \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "filename": "test-file.bin",
    "original_size": 8388608,
    "chunk_size": 4194304,
    "encryption_alg": "AES-256-GCM",
    "wrapped_file_key": "encrypted_key_material"
  }')

echo $FILE_RESPONSE | jq '.'

FILE_ID=$(echo $FILE_RESPONSE | jq -r '.file_id')
echo "File ID: $FILE_ID"
```

### 3d: Create a Chunk

```bash
CHUNK_RESPONSE=$(curl -s -X POST http://localhost:8000/chunks \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": "'$FILE_ID'",
    "chunk_index": 0,
    "chunk_size": 4194304,
    "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }')

echo $CHUNK_RESPONSE | jq '.'

CHUNK_ID=$(echo $CHUNK_RESPONSE | jq -r '.chunk_id')
echo "Chunk ID: $CHUNK_ID"
```

### 3e: Create a Replica

```bash
REPLICA_RESPONSE=$(curl -s -X POST http://localhost:8000/replicas \
  -H "Content-Type: application/json" \
  -d '{
    "chunk_id": "'$CHUNK_ID'",
    "node_id": "'$NODE_ID'"
  }')

echo $REPLICA_RESPONSE | jq '.'
```

### 3f: Get File Manifest (with Replica Locations)

```bash
curl http://localhost:8000/files/$FILE_ID | jq '.'
```

**Expected Response** (file manifest with chunk locations):
```json
{
  "file_id": "660e8400-e29b-41d4-a716-446655440001",
  "filename": "test-file.bin",
  "chunk_size": 4194304,
  "encryption_alg": "AES-256-GCM",
  "wrapped_file_key": "encrypted_key_material",
  "chunks": [
    {
      "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
      "chunk_index": 0,
      "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "chunk_size": 4194304,
      "nodes": ["550e8400-e29b-41d4-a716-446655440000"]
    }
  ]
}
```

### 3g: Delete the File

```bash
curl -X DELETE http://localhost:8000/files/$FILE_ID -v
```

**Expected**: HTTP 204 No Content

---

## 🚨 Troubleshooting

### Issue: "Connection refused" when starting service

**Solution**:
```bash
# Make sure database is running
docker ps | grep postgres

# If not, start it
make db-up

# Wait and try again
sleep 5
make run-metadata
```

### Issue: "psql: command not found"

**Solution**:
```bash
# Install PostgreSQL client tools
brew install postgresql          # macOS
sudo apt install postgresql-client # Ubuntu
```

### Issue: Manual test script says "Could not connect"

**Solution**:
```bash
# Check service is running
curl http://localhost:8000/health

# If not, start it
make run-metadata

# Check database is running
docker ps
```

### Issue: Unit tests fail with database errors

**Solution**:
```bash
# Make sure test database exists
make db-test-down
make db-test-up

# Try again
make test-unit
```

### Issue: Port 8000 already in use

**Solution**:
```bash
# Kill the existing process
pkill -f metadata-service

# Or find what's using it
lsof -i :8000

# Kill that process
kill -9 <PID>

# Start fresh
make run-metadata
```

---

## ✅ Success Criteria

You've successfully tested the Metadata Service when:

- ✅ **Step 1**: `bash scripts/manual-test.sh` shows "All 10 Tests Passed!"
- ✅ **Step 2**: `make test-unit` shows "PASS" with all tests
- ✅ **Step 3** (optional): Manual curl commands return expected JSON responses

---

## 📊 What Was Tested

| Component | Tests | Status |
|-----------|-------|--------|
| Health Check | 1 | ✅ |
| Node Management | 3 | ✅ |
| File Management | 3 | ✅ |
| Chunk Management | 2 | ✅ |
| Replica Management | 2 | ✅ |
| **Total** | **10** | **✅** |

---

## 🎓 What Each Test Validates

### Manual Tests (10 endpoints)
1. Service is running and healthy
2. Nodes can be registered and retrieved
3. Node heartbeat updates work
4. Files can be created and retrieved
5. Chunks can be managed
6. Replicas can be placed on nodes
7. Manifests include verified replica locations
8. Files can be deleted (with cascading deletes)

### Unit Tests (Repository Layer)
1. Node registration and retrieval
2. Node status updates
3. File operations with metadata
4. Chunk creation and ordering
5. Replica status transitions (pending → stored → verified)
6. Manifest generation with verified replicas
7. Cascading deletes preserve integrity
8. Query filtering works correctly

---

## 🎯 Next Steps After Testing

1. ✅ Verify tests pass
2. → Read [IMPLEMENTATION_LOG.md](IMPLEMENTATION_LOG.md) to understand the architecture
3. → (Future) Implement Storage Node (Phase 2)
4. → (Future) Integrate with API Gateway (Phase 3)

---

## 💡 Quick Commands Reference

```bash
# Start service
make db-up && make db-migrate && make run-metadata

# Run manual tests
bash scripts/manual-test.sh

# Run unit tests
make test-unit

# Run all tests
make test

# View help
make help

# Stop service
pkill -f metadata-service

# Stop database
make db-down
```

---

**Ready to test?** Start with [Step 1](#-step-1-manual-testing-quickest---5-minutes) above!
