#!/bin/bash

# Manual Testing Script for Auta Metadata Service
# This script runs through all 10 endpoints to verify functionality

set -e

BASE_URL="http://localhost:8000"
COLORS_GREEN='\033[0;32m'
COLORS_BLUE='\033[0;34m'
COLORS_RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${COLORS_BLUE}=== Auta Metadata Service Manual Tests ===${NC}\n"

# Test 1: Health Check
echo -e "${COLORS_BLUE}[1/10] Testing Health Check${NC}"
HEALTH=$(curl -s $BASE_URL/health)
echo "Response: $HEALTH"
if echo "$HEALTH" | grep -q "ok"; then
    echo -e "${COLORS_GREEN}✓ Health check passed${NC}\n"
else
    echo -e "${COLORS_RED}✗ Health check failed${NC}\n"
    exit 1
fi

# Test 2: Register Node
echo -e "${COLORS_BLUE}[2/10] Testing Node Registration${NC}"
NODE_RESPONSE=$(curl -s -X POST $BASE_URL/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "test_pk_001",
    "hostname": "storage-node-1.local",
    "endpoint": "http://storage-node-1.local:8001",
    "capacity_bytes": 1099511627776
  }')
echo "Response: $NODE_RESPONSE" | head -c 100
echo "..."
NODE_ID=$(echo "$NODE_RESPONSE" | grep -o '"node_id":"[^"]*' | cut -d'"' -f4)
if [ ! -z "$NODE_ID" ]; then
    echo -e "${COLORS_GREEN}✓ Node registered: $NODE_ID${NC}\n"
else
    echo -e "${COLORS_RED}✗ Node registration failed${NC}\n"
    exit 1
fi

# Test 3: Get Node
echo -e "${COLORS_BLUE}[3/10] Testing Get Node${NC}"
GET_NODE=$(curl -s $BASE_URL/nodes/$NODE_ID)
echo "Response: $GET_NODE" | head -c 100
echo "..."
if echo "$GET_NODE" | grep -q "$NODE_ID"; then
    echo -e "${COLORS_GREEN}✓ Get node passed${NC}\n"
else
    echo -e "${COLORS_RED}✗ Get node failed${NC}\n"
    exit 1
fi

# Test 4: Node Heartbeat
echo -e "${COLORS_BLUE}[4/10] Testing Node Heartbeat${NC}"
HEARTBEAT=$(curl -s -X POST $BASE_URL/nodes/$NODE_ID/heartbeat \
  -H "Content-Type: application/json" \
  -d '{
    "status": "healthy",
    "used_bytes": 104857600
  }')
echo "Response: $HEARTBEAT" | head -c 100
echo "..."
if echo "$HEARTBEAT" | grep -q "104857600"; then
    echo -e "${COLORS_GREEN}✓ Node heartbeat passed${NC}\n"
else
    echo -e "${COLORS_RED}✗ Node heartbeat failed${NC}\n"
    exit 1
fi

# Test 5: Create File
echo -e "${COLORS_BLUE}[5/10] Testing Create File${NC}"
FILE_RESPONSE=$(curl -s -X POST $BASE_URL/files \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "filename": "test-file-'$(date +%s)'.bin",
    "original_size": 8388608,
    "mime_type": "application/octet-stream",
    "chunk_size": 4194304,
    "encryption_alg": "AES-256-GCM",
    "wrapped_file_key": "encrypted_key_base64_here"
  }')
echo "Response: $FILE_RESPONSE" | head -c 100
echo "..."
FILE_ID=$(echo "$FILE_RESPONSE" | grep -o '"file_id":"[^"]*' | cut -d'"' -f4)
if [ ! -z "$FILE_ID" ]; then
    echo -e "${COLORS_GREEN}✓ File created: $FILE_ID${NC}\n"
else
    echo -e "${COLORS_RED}✗ File creation failed${NC}\n"
    exit 1
fi

# Test 6: Create Chunk
echo -e "${COLORS_BLUE}[6/10] Testing Create Chunk${NC}"
CHUNK_RESPONSE=$(curl -s -X POST $BASE_URL/chunks \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": "'$FILE_ID'",
    "chunk_index": 0,
    "chunk_size": 4194304,
    "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }')
echo "Response: $CHUNK_RESPONSE" | head -c 100
echo "..."
CHUNK_ID=$(echo "$CHUNK_RESPONSE" | grep -o '"chunk_id":"[^"]*' | cut -d'"' -f4)
if [ ! -z "$CHUNK_ID" ]; then
    echo -e "${COLORS_GREEN}✓ Chunk created: $CHUNK_ID${NC}\n"
else
    echo -e "${COLORS_RED}✗ Chunk creation failed${NC}\n"
    exit 1
fi

# Test 7: Get Chunk
echo -e "${COLORS_BLUE}[7/10] Testing Get Chunk${NC}"
GET_CHUNK=$(curl -s $BASE_URL/chunks/$CHUNK_ID)
echo "Response: $GET_CHUNK" | head -c 100
echo "..."
if echo "$GET_CHUNK" | grep -q "$CHUNK_ID"; then
    echo -e "${COLORS_GREEN}✓ Get chunk passed${NC}\n"
else
    echo -e "${COLORS_RED}✗ Get chunk failed${NC}\n"
    exit 1
fi

# Test 8: Create Replica
echo -e "${COLORS_BLUE}[8/10] Testing Create Replica${NC}"
REPLICA_RESPONSE=$(curl -s -X POST $BASE_URL/replicas \
  -H "Content-Type: application/json" \
  -d '{
    "chunk_id": "'$CHUNK_ID'",
    "node_id": "'$NODE_ID'"
  }')
echo "Response: $REPLICA_RESPONSE" | head -c 100
echo "..."
REPLICA_ID=$(echo "$REPLICA_RESPONSE" | grep -o '"replica_id":"[^"]*' | cut -d'"' -f4)
if [ ! -z "$REPLICA_ID" ]; then
    echo -e "${COLORS_GREEN}✓ Replica created: $REPLICA_ID${NC}\n"
else
    echo -e "${COLORS_RED}✗ Replica creation failed${NC}\n"
    exit 1
fi

# Test 9: Get File Manifest
echo -e "${COLORS_BLUE}[9/10] Testing Get File Manifest${NC}"
MANIFEST=$(curl -s $BASE_URL/files/$FILE_ID)
echo "Response: $MANIFEST" | head -c 100
echo "..."
if echo "$MANIFEST" | grep -q "$CHUNK_ID"; then
    echo -e "${COLORS_GREEN}✓ Get file manifest passed${NC}\n"
else
    echo -e "${COLORS_RED}✗ Get file manifest failed${NC}\n"
    exit 1
fi

# Test 10: Delete File
echo -e "${COLORS_BLUE}[10/10] Testing Delete File${NC}"
DELETE_RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE $BASE_URL/files/$FILE_ID)
HTTP_CODE=$(echo "$DELETE_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" = "204" ]; then
    echo -e "${COLORS_GREEN}✓ File deleted successfully${NC}\n"
else
    echo -e "${COLORS_RED}✗ File deletion failed (HTTP $HTTP_CODE)${NC}\n"
    exit 1
fi

echo -e "${COLORS_GREEN}=== All 10 Tests Passed! ===${NC}"
echo ""
echo "Summary:"
echo "✓ Health check"
echo "✓ Node registration"
echo "✓ Get node"
echo "✓ Node heartbeat"
echo "✓ File creation"
echo "✓ Chunk creation"
echo "✓ Get chunk"
echo "✓ Replica creation"
echo "✓ Get file manifest"
echo "✓ File deletion"
