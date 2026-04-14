# Implementation Log

This document tracks all implementations, changes, and technical decisions made during Auta development.

## Session 1: Metadata Service - Phase 1 Implementation (2026-04-14)

### Overview
Implemented the core metadata service with HTTP handlers, database access layer, and complete request/response handling for all 10 API endpoints.

### Files Created/Modified

#### 1. `internal/metadata/models.go` - Data Models
**What**: Comprehensive type definitions for all metadata entities
**Why**: Standardize data structures across handlers and database layer
**Changes**:
- Defined 4 core models: `Node`, `File`, `Chunk`, `Replica`
- Defined 1 manifest model: `Manifest` with `ChunkInfo`
- Created 7 request types: `RegisterNodeRequest`, `HeartbeatRequest`, `CreateFileRequest`, `CreateChunkRequest`, `CreateReplicaRequest`, `UpdateReplicaStatusRequest`
- Created error response type: `ErrorResponse`
- All models include JSON serialization tags for API responses

**Key Design Decisions**:
- Used `uuid.UUID` for all IDs (security and collision resistance)
- Included timestamps (CreatedAt, UpdatedAt) for audit trails
- Status fields use string enums (e.g., "healthy", "offline", "degraded")
- Replica status follows state machine: pending → stored → verified (or failed)
- Manifest includes only verified replicas for download operations

#### 2. `internal/metadata/repository.go` - Database Access Layer
**What**: Repository pattern implementation for all database operations
**Why**: Separation of concerns, testability, and consistency

**Methods Implemented** (20 total):

**Node Operations** (4):
- `RegisterNode(ctx, node)` - Insert new storage node
- `GetNode(ctx, nodeID)` - Retrieve node by ID
- `ListHealthyNodes(ctx)` - Get all healthy nodes sorted by usage
- `UpdateNodeHeartbeat(ctx, nodeID, status, usedBytes)` - Update node health and capacity

**File Operations** (4):
- `CreateFile(ctx, file)` - Create file record
- `GetFile(ctx, fileID)` - Retrieve file with metadata
- `DeleteFile(ctx, fileID)` - Cascading delete (file → chunks → replicas)
- (Future: ListUserFiles, UpdateFile)

**Chunk Operations** (3):
- `CreateChunk(ctx, chunk)` - Create chunk record
- `GetChunk(ctx, chunkID)` - Retrieve chunk
- `GetFileChunks(ctx, fileID)` - Get all chunks for file (ordered by index)

**Replica Operations** (7):
- `CreateReplica(ctx, replica)` - Create replica record
- `GetReplica(ctx, replicaID)` - Retrieve replica
- `GetChunkReplicas(ctx, chunkID)` - Get all replicas for chunk
- `UpdateReplicaStatus(ctx, replicaID, status)` - Update status with timestamp
- `CountChunkReplicas(ctx, chunkID, status)` - Count replicas by status
- `GetVerifiedReplicaNodes(ctx, chunkID)` - Get node IDs of verified replicas
- (Future: GetReplicasForNode, ListFailedReplicas)

**Key Technical Details**:
- Contextual query execution (`ExecContext`, `QueryContext`) for timeout support
- Row cascade deletion for file cleanup
- Status-aware timestamp updates (stored_at, verified_at)
- Sorted queries (healthy nodes by usage, chunks by index)
- Transaction-based cascading deletes

#### 3. `internal/metadata/service.go` - HTTP Service & Handlers
**What**: Main service with HTTP handler implementations
**Why**: Connect database layer to HTTP API endpoints

**Architecture**:
```
Service
├── db: *sql.DB (connection pool)
└── repo: *Repository (data access)
```

**Helper Methods** (2):
- `respondJSON(w, data, statusCode)` - Serialize response and write with correct content-type
- `respondError(w, error, statusCode)` - Standardized error responses

**Handlers Implemented** (10):

1. **Health Check**
   - `handleHealth()` - GET /health
   - Returns `{"status":"ok"}` with 200 OK

2. **Node Management**
   - `handleRegisterNode()` - POST /nodes
     - Validates request (publicKey, hostname, endpoint, capacity)
     - Generates UUID for node
     - Returns created node
   - `handleGetNode()` - GET /nodes/{node_id}
     - Parses UUID from path parameter
     - Returns node details
   - `handleNodeHeartbeat()` - POST /nodes/{node_id}/heartbeat
     - Updates status and used_bytes
     - Returns updated node

3. **File Management**
   - `handleCreateFile()` - POST /files
     - Validates encryption metadata
     - Generates file UUID
     - Returns created file
   - `handleGetFile()` - GET /files/{file_id}
     - Retrieves file metadata
     - Assembles manifest with all chunks
     - Includes verified replica node IDs
     - Returns complete manifest for download
   - `handleDeleteFile()` - DELETE /files/{file_id}
     - Cascading delete via repository
     - Returns 204 No Content

4. **Chunk Management**
   - `handleCreateChunk()` - POST /chunks
     - Validates SHA-256 hash (64 char hex)
     - Links chunk to file via file_id
     - Returns created chunk
   - `handleGetChunk()` - GET /chunks/{chunk_id}
     - Returns chunk metadata

5. **Replica Management**
   - `handleCreateReplica()` - POST /replicas
     - Creates replica record with "pending" status
     - Returns created replica
   - `handleGetChunkReplicas()` - GET /replicas/chunk/{chunk_id}
     - Returns all replicas for chunk (ordered by creation time)

**Error Handling**:
- HTTP method validation (405 Method Not Allowed)
- JSON decode errors (400 Bad Request)
- UUID parse errors (400 Bad Request)
- Not found errors (404)
- Database errors (500)
- Structured error responses with status code

**Key Implementation Patterns**:
- HTTP 1.1.16+ path parameters using `r.PathValue()`
- Context propagation through all database calls
- Generated UUIDs for all entity IDs
- Timestamp generation at service layer (immutable at DB layer)
- JSON encoding for all responses

### Database Schema
Already defined in `migrations/001_initial_schema.sql`:
- 4 tables with proper foreign keys and constraints
- Unique constraints for duplicate prevention (filename per owner)
- Partial indexes for common queries (status, node_id, created_at)
- Cascade deletes for data integrity

### Testing

**Build Verification**:
```bash
go build -o bin/metadata-service ./cmd/metadata-service
# ✅ Compiles successfully with no errors
```

**Next Testing Phase**:
- Unit tests for repository operations
- Integration tests with real PostgreSQL
- End-to-end tests of upload/download flow
- Load testing for concurrent operations

### API Specification Summary

| Method | Path | Request | Response | Status |
|--------|------|---------|----------|--------|
| GET | /health | - | `{status: "ok"}` | 200 |
| POST | /nodes | RegisterNodeRequest | Node | 201 |
| GET | /nodes/{node_id} | - | Node | 200 |
| POST | /nodes/{node_id}/heartbeat | HeartbeatRequest | Node | 200 |
| POST | /files | CreateFileRequest | File | 201 |
| GET | /files/{file_id} | - | Manifest | 200 |
| DELETE | /files/{file_id} | - | - | 204 |
| POST | /chunks | CreateChunkRequest | Chunk | 201 |
| GET | /chunks/{chunk_id} | - | Chunk | 200 |
| POST | /replicas | CreateReplicaRequest | Replica | 201 |
| GET | /replicas/chunk/{chunk_id} | - | Replica[] | 200 |

### Known Limitations & Future Work

1. **No Input Validation**
   - Should add struct tags and validation middleware (go-playground/validator)
   - Content-hash validation (must be valid SHA-256 hex)

2. **No Authentication/Authorization**
   - All endpoints currently public
   - Should add JWT tokens for user identity
   - Should restrict file access to owner

3. **No Replication Factor Enforcement**
   - API accepts replicas but doesn't verify 3-replica minimum
   - Replication engine (future component) will handle this

4. **Limited Query Capabilities**
   - No filtering by owner or date range
   - No pagination for large result sets

5. **No Transactions for Multi-Step Operations**
   - File creation + chunk creation + replica creation should be atomic
   - Need distributed transaction support for cross-service operations

### Performance Considerations

**Indexes**:
- Covered queries for most common operations
- Node queries sorted by used_bytes for efficient placement
- Chunks ordered by index for deterministic reassembly

**Optimization Opportunities**:
- Add caching layer for frequently accessed files/manifests
- Batch replica creation/update operations
- Connection pooling already enabled (db/sql)
- Prepared statement caching (db/sql default)

### Security Notes

- No rate limiting on API endpoints
- No HTTPS/TLS termination (handled by reverse proxy)
- Wrapped keys stored but never unwrapped (zero-knowledge)
- Node API keys not yet implemented
- Chunk hashes can be compared with plaintext to detect corruption

---

## Next Session: Storage Node Implementation

**Expected Tasks**:
1. Create storage node HTTP server
2. Implement chunk upload/download logic
3. Add disk storage with organized directory structure
4. Implement SHA-256 hash verification
5. Add health reporting to metadata service
6. Create storage node CLI for management

---

## Glossary

- **UUID**: Universally Unique Identifier (128-bit)
- **Repository Pattern**: Data access abstraction layer
- **Context**: Go's context.Context for cancellation and timeouts
- **Manifest**: File metadata + chunk locations for reconstruction
- **Replica**: Single copy of a chunk on a storage node
- **Zero-Knowledge**: System cannot access plaintext data
