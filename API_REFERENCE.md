# Metadata Service API Reference

## Overview

The Metadata Service is the central coordinator for Auta's distributed storage system. It manages:
- Storage node registration and health tracking
- File metadata and encryption information
- Chunk records and locations
- Replica placement and tracking

**Base URL**: `http://localhost:8000`

## Endpoints

### Health Check

#### GET /health

Check service health status.

**Response**: 200 OK
```json
{
  "status": "ok"
}
```

---

### Node Management

#### POST /nodes

Register a new storage node.

**Request**:
```json
{
  "public_key": "pk_node123...",
  "hostname": "storage-node-1.local",
  "endpoint": "http://storage-node-1.local:8001",
  "capacity_bytes": 1099511627776
}
```

**Response**: 201 Created
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "public_key": "pk_node123...",
  "hostname": "storage-node-1.local",
  "endpoint": "http://storage-node-1.local:8001",
  "status": "healthy",
  "capacity_bytes": 1099511627776,
  "used_bytes": 0,
  "last_heartbeat_at": null,
  "created_at": "2026-04-14T12:00:00Z",
  "updated_at": "2026-04-14T12:00:00Z"
}
```

---

#### GET /nodes/{node_id}

Get node details and current capacity.

**Path Parameters**:
- `node_id` (UUID): Node identifier

**Response**: 200 OK
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "public_key": "pk_node123...",
  "hostname": "storage-node-1.local",
  "endpoint": "http://storage-node-1.local:8001",
  "status": "healthy",
  "capacity_bytes": 1099511627776,
  "used_bytes": 104857600,
  "last_heartbeat_at": "2026-04-14T12:05:30Z",
  "created_at": "2026-04-14T12:00:00Z",
  "updated_at": "2026-04-14T12:05:30Z"
}
```

---

#### POST /nodes/{node_id}/heartbeat

Update node heartbeat and health status.

**Path Parameters**:
- `node_id` (UUID): Node identifier

**Request**:
```json
{
  "status": "healthy",
  "used_bytes": 104857600
}
```

**Response**: 200 OK
```json
{
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "public_key": "pk_node123...",
  "hostname": "storage-node-1.local",
  "endpoint": "http://storage-node-1.local:8001",
  "status": "healthy",
  "capacity_bytes": 1099511627776,
  "used_bytes": 104857600,
  "last_heartbeat_at": "2026-04-14T12:05:30Z",
  "created_at": "2026-04-14T12:00:00Z",
  "updated_at": "2026-04-14T12:05:30Z"
}
```

---

### File Management

#### POST /files

Create a new file record with encryption metadata.

**Request**:
```json
{
  "owner_id": "123e4567-e89b-12d3-a456-426614174000",
  "filename": "photo.jpg",
  "original_size": 2097152,
  "mime_type": "image/jpeg",
  "chunk_size": 4194304,
  "encryption_alg": "AES-256-GCM",
  "wrapped_file_key": "encrypted_key_material_base64..."
}
```

**Response**: 201 Created
```json
{
  "file_id": "660e8400-e29b-41d4-a716-446655440001",
  "owner_id": "123e4567-e89b-12d3-a456-426614174000",
  "filename": "photo.jpg",
  "original_size": 2097152,
  "mime_type": "image/jpeg",
  "chunk_size": 4194304,
  "encryption_alg": "AES-256-GCM",
  "wrapped_file_key": "encrypted_key_material_base64...",
  "created_at": "2026-04-14T12:10:00Z",
  "updated_at": "2026-04-14T12:10:00Z"
}
```

---

#### GET /files/{file_id}

Get file metadata as a manifest with chunk locations.

**Path Parameters**:
- `file_id` (UUID): File identifier

**Response**: 200 OK
```json
{
  "file_id": "660e8400-e29b-41d4-a716-446655440001",
  "filename": "photo.jpg",
  "chunk_size": 4194304,
  "encryption_alg": "AES-256-GCM",
  "wrapped_file_key": "encrypted_key_material_base64...",
  "chunks": [
    {
      "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
      "chunk_index": 0,
      "content_hash": "sha256_hash_here",
      "chunk_size": 4194304,
      "nodes": [
        "550e8400-e29b-41d4-a716-446655440000",
        "550e8400-e29b-41d4-a716-446655440001",
        "550e8400-e29b-41d4-a716-446655440002"
      ]
    }
  ]
}
```

---

#### DELETE /files/{file_id}

Delete file and all associated chunks and replicas.

**Path Parameters**:
- `file_id` (UUID): File identifier

**Response**: 204 No Content

---

### Chunk Management

#### POST /chunks

Create a chunk record linked to a file.

**Request**:
```json
{
  "file_id": "660e8400-e29b-41d4-a716-446655440001",
  "chunk_index": 0,
  "chunk_size": 4194304,
  "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
}
```

**Response**: 201 Created
```json
{
  "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
  "file_id": "660e8400-e29b-41d4-a716-446655440001",
  "chunk_index": 0,
  "chunk_size": 4194304,
  "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "created_at": "2026-04-14T12:10:00Z"
}
```

---

#### GET /chunks/{chunk_id}

Get chunk metadata.

**Path Parameters**:
- `chunk_id` (UUID): Chunk identifier

**Response**: 200 OK
```json
{
  "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
  "file_id": "660e8400-e29b-41d4-a716-446655440001",
  "chunk_index": 0,
  "chunk_size": 4194304,
  "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "created_at": "2026-04-14T12:10:00Z"
}
```

---

### Replica Management

#### POST /replicas

Create a replica record for a chunk on a node.

**Request**:
```json
{
  "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
  "node_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response**: 201 Created
```json
{
  "replica_id": "880e8400-e29b-41d4-a716-446655440003",
  "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
  "node_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "stored_at": null,
  "verified_at": null,
  "created_at": "2026-04-14T12:10:00Z",
  "updated_at": "2026-04-14T12:10:00Z"
}
```

---

#### GET /replicas/chunk/{chunk_id}

Get all replicas for a chunk.

**Path Parameters**:
- `chunk_id` (UUID): Chunk identifier

**Response**: 200 OK
```json
[
  {
    "replica_id": "880e8400-e29b-41d4-a716-446655440003",
    "chunk_id": "770e8400-e29b-41d4-a716-446655440002",
    "node_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "verified",
    "stored_at": "2026-04-14T12:10:10Z",
    "verified_at": "2026-04-14T12:10:15Z",
    "created_at": "2026-04-14T12:10:00Z",
    "updated_at": "2026-04-14T12:10:15Z"
  }
]
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": "error_code",
  "status_code": 400,
  "message": "Detailed error description"
}
```

### Common Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid input or parsing error
- `404 Not Found` - Resource not found
- `405 Method Not Allowed` - Wrong HTTP method
- `500 Internal Server Error` - Database or server error

---

## Data Types

### Node Status
- `healthy` - Node is operational and accepting new replicas
- `degraded` - Node is operational but under stress
- `offline` - Node is not responding

### Replica Status
- `pending` - Replica created, waiting for upload
- `stored` - Replica data received on node
- `verified` - Replica hash verified and ready
- `failed` - Replica upload/verification failed

---

## Usage Examples

### Register a Node
```bash
curl -X POST http://localhost:8000/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "pk_node1",
    "hostname": "storage1.local",
    "endpoint": "http://storage1.local:8001",
    "capacity_bytes": 1099511627776
  }'
```

### Create and Upload a File
```bash
# 1. Create file record
FILE_ID=$(curl -X POST http://localhost:8000/files \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "filename": "data.bin",
    "original_size": 8388608,
    "chunk_size": 4194304,
    "encryption_alg": "AES-256-GCM",
    "wrapped_file_key": "key..."
  }' | jq -r '.file_id')

# 2. Create chunk records
CHUNK_ID=$(curl -X POST http://localhost:8000/chunks \
  -H "Content-Type: application/json" \
  -d '{
    "file_id": "'$FILE_ID'",
    "chunk_index": 0,
    "chunk_size": 4194304,
    "content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }' | jq -r '.chunk_id')

# 3. Create replica records
curl -X POST http://localhost:8000/replicas \
  -H "Content-Type: application/json" \
  -d '{
    "chunk_id": "'$CHUNK_ID'",
    "node_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

### Retrieve File Manifest
```bash
curl http://localhost:8000/files/660e8400-e29b-41d4-a716-446655440001 \
  -H "Content-Type: application/json" | jq '.'
```

---

## Future Enhancements

- [ ] Input validation middleware
- [ ] Authentication/authorization
- [ ] Rate limiting
- [ ] Pagination for list operations
- [ ] Filtering by owner/date
- [ ] Batch operations
- [ ] GraphQL interface
- [ ] OpenAPI/Swagger documentation
