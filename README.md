# Auta

Encrypted, decentralized, self-hosted cloud storage mesh for sovereign distributed file storage.
Auta distributes encrypted file chunks across multiple self-hosted nodes, enabling fault-tolerant
private cloud storage without relying on centralized providers.

## Core

- Client-side AES-256 encryption
- Distributed chunk replication
- Multi-node storage mesh
- Automatic node failure recovery
- Zero-knowledge architecture

## Repository Status

Status: In Active Development, contributions welcome!

- [x] Phase 1: Metadata Service - HTTP API & Database Layer
- [ ] Phase 2: Storage Node Service - Chunk Upload/Download/Verification
- [ ] Phase 3: API Gateway - Request Routing & Authentication
- [ ] Phase 4: Client Library - Encryption & Chunking Logic
- [ ] Phase 5: Replication Engine - Node Health & Replica Management
- [ ] Encrypted Chunk Engine - AES-256-GCM with Key Wrapping
- [ ] Multi-Node Replication - 3-replica minimum enforcement
- [ ] Erasure Coding - Data redundancy with reduced storage overhead
- [ ] S3 Compatibility - S3-like API for object storag

## Architecture Overview

The platform is organized into five layers:

- Client layer for encryption, chunking, upload, and download
- API gateway for authentication, orchestration, and routing
- Metadata coordination layer for manifests and placement tracking
- Storage node layer for chunk persistence and retrieval
- Replication and recovery layer for fault tolerance and self-healing

For the full technical specification, see [architecture.md](architecture.md).
For testing, all necessary guides are found in [docs/](/docs/)
