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

Status: In Active Development
Current Version: MVP v0.1

- [ ] Encrypted chunk engine
- [ ] Multi-node replication
- [ ] Erasure coding
- [ ] S3 compatibility

## Architecture Overview

The platform is organized into five layers:

- Client layer for encryption, chunking, upload, and download
- API gateway for authentication, orchestration, and routing
- Metadata coordination layer for manifests and placement tracking
- Storage node layer for chunk persistence and retrieval
- Replication and recovery layer for fault tolerance and self-healing

For the full technical specification, see [architecture.md](architecture.md).
