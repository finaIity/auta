# Architecture

This document defines the core technical model for Auta's encrypted distributed storage system.

## Node Model

The system is made up of four operational node types plus one control plane service:

### 1. Client

- Runs on the user device
- Encrypts files locally before upload
- Splits encrypted files into chunks
- Decrypts files after download
- Never sends plaintext to the network

### 2. API Gateway

- Accepts upload and download requests
- Authenticates users and nodes
- Coordinates placement decisions
- Returns manifests and chunk locations
- Does not store file contents

### 3. Metadata Service

- Stores manifests, file records, chunk records, node records, and replica mappings
- Tracks which chunks belong to which file
- Tracks where each replica lives
- Verifies that expected replica counts are satisfied

### 4. Storage Node

- Stores encrypted chunk blobs on disk
- Serves chunk reads and writes
- Reports health and capacity
- Verifies chunk hashes before accepting or serving data
- Never sees plaintext content

### 5. Replication Engine

- Monitors node health and replica counts
- Recreates missing replicas when nodes fail
- Balances placement across available nodes
- Initiates repair work after corruption or node loss

## Replication Logic

The default durability model is a replication factor of 3.

### Placement Rules

- Every chunk must exist on at least three distinct nodes
- Replicas should be distributed across different failure domains when possible
- The gateway or replication engine should avoid placing multiple replicas on the same node
- Placement should prefer nodes with enough free capacity and healthy status

### Upload-Time Flow

1. Client encrypts the file and generates chunk hashes.
2. Gateway selects candidate storage nodes for each chunk.
3. Storage nodes accept the chunk write.
4. Metadata service records the replica locations.
5. Chunk is considered durable only after the replica threshold is met.

### Failure Recovery Flow

1. Heartbeat misses or failed reads mark a node offline.
2. Replication engine identifies chunks below the replica threshold.
3. New destination nodes are selected.
4. Missing replicas are rebuilt from healthy copies.
5. Metadata is updated after the new replica is confirmed.

### Consistency Model

- Eventual consistency is acceptable for replica placement updates
- Metadata writes should be atomic at the record level
- Reads should prefer healthy replicas and reject corrupted data

## Chunk Size Strategy

### Default Size

The initial chunk size is 4 MB fixed chunks.

### Why 4 MB

- Large enough to keep metadata overhead manageable
- Small enough to retry efficiently when a node fails mid-upload
- Simple to implement and reason about in the MVP

### Strategy Rules

- Use fixed-size chunks for the first version
- Allow the final chunk to be smaller than the default size
- Store the chunk index, hash, and size in metadata
- Include the chunk order in the manifest so files can be reassembled deterministically

### Future Options

- Adaptive chunk sizing for very large files
- Content-defined chunking for deduplication
- Hybrid strategy for mixed workloads

## Metadata Schema

The metadata service should store enough information to locate and reconstruct a file without exposing plaintext.

### Core Entities

#### files

- file_id
- owner_id
- filename
- original_size
- mime_type
- chunk_size
- encryption_alg
- wrapped_file_key
- created_at
- updated_at

#### chunks

- chunk_id
- file_id
- chunk_index
- chunk_size
- content_hash
- created_at

#### replicas

- replica_id
- chunk_id
- node_id
- status
- stored_at
- verified_at

#### nodes

- node_id
- public_key or auth identifier
- hostname or endpoint
- status
- capacity_bytes
- used_bytes
- last_heartbeat_at

### Manifest Shape

```json
{
  "file_id": "abc123",
  "filename": "photo.jpg",
  "chunk_size": 4194304,
  "encryption_alg": "AES-256-GCM",
  "wrapped_file_key": "encrypted-key-here",
  "chunks": [
    {
      "chunk_id": "c1",
      "chunk_index": 0,
      "content_hash": "sha256...",
      "nodes": ["node1", "node3", "node4"]
    }
  ]
}
```

### Design Notes

- The manifest should allow file reconstruction without querying every node individually
- The manifest should never contain plaintext file contents
- If possible, store only encrypted key material and minimal user-visible metadata

## Encryption Lifecycle

### 1. Key Generation

- Generate a unique symmetric key per file
- Use AES-256-GCM for file encryption
- Generate a secure random initialization vector per encryption operation

### 2. Client-Side Encryption

- Encrypt the file locally before any network transfer
- Split the encrypted payload into chunks
- Compute a SHA-256 hash for each chunk

### 3. Key Wrapping

- Encrypt or wrap the per-file symmetric key for storage in the manifest
- The wrapped key is the only key material that leaves the client
- Future versions may support public-key sharing for multi-user access

### 4. Upload and Storage

- Only encrypted chunks are transmitted to storage nodes
- Storage nodes persist chunk blobs without decrypting them
- Metadata records the chunk hash and replica locations

### 5. Retrieval and Decryption

- Client requests the manifest
- Client downloads the needed chunks from healthy nodes
- Client verifies chunk hashes
- Client reassembles the file
- Client unwraps the file key and decrypts the payload locally

### 6. Rotation and Revocation

- If a file key is rotated, a new manifest version should be created
- If access is revoked, future key-sharing mechanisms should stop issuing wrapped keys to removed users
- Old replicas remain unreadable without the corresponding key material

### Security Properties

- Confidentiality is enforced by client-side encryption
- Integrity is enforced by chunk hashes and AEAD authentication
- Availability is improved through multi-node replication

### Must-add

- chunking strategy
- replication factor
- encryption model
