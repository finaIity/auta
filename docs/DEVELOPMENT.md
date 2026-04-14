# Auta Development Guide

This directory contains the Auta distributed storage system implementation.

## Project Structure

```
├── cmd/                      # Service entry points
│   ├── metadata-service/     # Metadata coordination service
│   ├── storage-node/         # Storage node service
│   ├── api-gateway/          # API gateway
│   ├── client/               # Client CLI/library
│   └── replication-engine/   # Replication engine
├── internal/                 # Internal packages
│   ├── metadata/             # Metadata service implementation
│   ├── storage/              # Storage node implementation
│   ├── api/                  # API gateway implementation
│   ├── client/               # Client implementation
│   └── replication/          # Replication engine implementation
├── pkg/                      # Public packages
│   ├── crypto/               # Encryption utilities
│   └── chunking/             # File chunking logic
├── migrations/               # Database migrations
├── docker-compose.yml        # Local development environment
└── Makefile                  # Development commands
```

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL client tools (psql)

## Quick Start

### 1. Start PostgreSQL

```bash
make db-up
```

### 2. Run Database Migrations

```bash
make db-migrate
```

### 3. Build and Run Metadata Service

```bash
make run-metadata
```

The metadata service will start on `http://localhost:8000`. Check health at `http://localhost:8000/health`.

## Development Commands

```bash
make help              # Show all available commands
make build             # Build all services
make test              # Run tests
make clean             # Clean build artifacts
make db-down           # Stop PostgreSQL
```

## Database Schema

The metadata service uses PostgreSQL with the following core entities:

- **nodes** - Storage nodes in the network
- **files** - User files with metadata
- **chunks** - Encrypted file chunks
- **replicas** - Chunk replicas on nodes

See `migrations/001_initial_schema.sql` for full schema details.

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Files
- `POST /files` - Create file metadata
- `GET /files/{file_id}` - Get file metadata
- `DELETE /files/{file_id}` - Delete file

### Chunks
- `POST /chunks` - Create chunk record
- `GET /chunks/{chunk_id}` - Get chunk metadata

### Nodes
- `POST /nodes` - Register storage node
- `GET /nodes/{node_id}` - Get node info
- `POST /nodes/{node_id}/heartbeat` - Node heartbeat

### Replicas
- `POST /replicas` - Create replica record
- `GET /replicas/chunk/{chunk_id}` - List replicas for chunk

## Next Steps

1. ✅ Set up project structure
2. ✅ Configure database
3. ⏳ Implement metadata service handlers
4. ⏳ Build storage node service
5. ⏳ Build API gateway
6. ⏳ Build client layer
7. ⏳ Build replication engine
