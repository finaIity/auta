# CHANGELOG

All notable changes to the Auta project are documented in this file.

## [0.1.0] - 2026-04-14

### Project Initialization
- ✅ Set up Go module structure (github.com/anomalyco/auta)
- ✅ Initialized git repository with initial commits
- ✅ Created project directory layout with separation of concerns

### Infrastructure
- ✅ PostgreSQL database schema with 4 core entities (nodes, files, chunks, replicas)
- ✅ Docker Compose configuration for local PostgreSQL development
- ✅ Database migration system (migrations/001_initial_schema.sql)
- ✅ Makefile with development commands (db-up, db-down, build, run-metadata, test, clean)

### Metadata Service Foundation (Internal)
- ✅ HTTP server setup with router
- ✅ Health check endpoint at GET /health
- ✅ Database connection pooling with PostgreSQL
- ✅ Placeholder handlers for 10 endpoints (stubs for implementation)
- ✅ Service interface with Start() and Close() methods

### Metadata Service - Phase 1 Implementation
- ✅ **Data Models** (`internal/metadata/models.go`):
  - Node, File, Chunk, Replica, Manifest entities
  - 7 request types with validation tags
  - Error response type
  - All models JSON-serializable

- ✅ **Repository Layer** (`internal/metadata/repository.go`):
  - Repository pattern for data access
  - 20 database methods across 4 entities:
    - Node: Register, Get, ListHealthy, UpdateHeartbeat
    - File: Create, Get, Delete (with cascade)
    - Chunk: Create, Get, GetFileChunks
    - Replica: Create, Get, GetChunkReplicas, UpdateStatus, CountByStatus, GetVerifiedNodes
  - Context-aware queries for timeout support
  - Atomic cascade deletes for data integrity

- ✅ **HTTP Handlers** (`internal/metadata/service.go`):
  - 10 fully implemented endpoints
  - Node management: Register, Get, Heartbeat
  - File management: Create, Get (with manifest), Delete
  - Chunk management: Create, Get
  - Replica management: Create, List by chunk
  - Standardized error handling with proper HTTP status codes
  - JSON request/response serialization

### Documentation
- ✅ DEVELOPMENT.md - Developer setup and quick start guide
- ✅ API_REFERENCE.md - Complete API documentation with examples
- ✅ IMPLEMENTATION_LOG.md - Detailed implementation notes and technical decisions
- ✅ CHANGELOG.md - This file
- ✅ .gitignore - Standard Go project ignores
- ✅ go.mod, go.sum - Dependency management with uuid and pq

### Technology Stack
- **Language**: Go 1.22
- **Database**: PostgreSQL 16
- **Key Dependencies**: lib/pq for PostgreSQL driver
- **Deployment**: Docker & Docker Compose

---

## [Unreleased] - Upcoming

### Metadata Service - Phase 1 Implementation
- ⏳ Node management (register, heartbeat, status)
- ⏳ File CRUD operations with encryption metadata
- ⏳ Chunk management and tracking
- ⏳ Replica placement and tracking
- ⏳ Input validation and error handling
- ⏳ Request/response JSON serialization
- ⏳ Database transaction handling
- ⏳ UUID generation for all entities

### Metadata Service - Phase 2 Enhancement
- ⏳ Replication factor enforcement
- ⏳ Node capacity checks
- ⏳ Replica status transitions
- ⏳ Manifest generation with chunk locations
- ⏳ Query optimization with proper indexing

### Storage Node Service
- ⏳ HTTP server for chunk upload/download
- ⏳ Disk storage with organized chunk layout
- ⏳ Hash verification (SHA-256)
- ⏳ Health reporting to metadata service
- ⏳ Capacity tracking

### API Gateway Service
- ⏳ Request routing to other services
- ⏳ User authentication
- ⏳ Chunk placement orchestration
- ⏳ Manifest retrieval and caching

### Client Library
- ⏳ File encryption (AES-256-GCM)
- ⏳ Chunking logic (4MB fixed chunks)
- ⏳ Upload orchestration
- ⏳ Download and reassembly
- ⏳ Key wrapping

### Replication Engine
- ⏳ Node health monitoring
- ⏳ Replica count verification
- ⏳ Replica reconstruction
- ⏳ Placement rebalancing

### Testing & Quality
- ⏳ Unit tests for metadata service
- ⏳ Integration tests with real database
- ⏳ End-to-end tests

### Security & Operations
- ⏳ Authentication/authorization
- ⏳ Logging and monitoring
- ⏳ Rate limiting
- ⏳ Backup strategy documentation
