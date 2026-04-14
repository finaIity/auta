package metadata

import (
	"time"

	"github.com/google/uuid"
)

// Node represents a storage node in the cluster
type Node struct {
	NodeID          uuid.UUID  `json:"node_id"`
	PublicKey       string     `json:"public_key"`
	Hostname        string     `json:"hostname"`
	Endpoint        string     `json:"endpoint"`
	Status          string     `json:"status"` // "healthy", "offline", "degraded"
	CapacityBytes   int64      `json:"capacity_bytes"`
	UsedBytes       int64      `json:"used_bytes"`
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// File represents a user's encrypted file
type File struct {
	FileID         uuid.UUID `json:"file_id"`
	OwnerID        uuid.UUID `json:"owner_id"`
	Filename       string    `json:"filename"`
	OriginalSize   int64     `json:"original_size"`
	MimeType       string    `json:"mime_type,omitempty"`
	ChunkSize      int       `json:"chunk_size"`
	EncryptionAlg  string    `json:"encryption_alg"`
	WrappedFileKey string    `json:"wrapped_file_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Chunk represents a chunk of an encrypted file
type Chunk struct {
	ChunkID     uuid.UUID `json:"chunk_id"`
	FileID      uuid.UUID `json:"file_id"`
	ChunkIndex  int       `json:"chunk_index"`
	ChunkSize   int       `json:"chunk_size"`
	ContentHash string    `json:"content_hash"` // SHA-256 hash
	CreatedAt   time.Time `json:"created_at"`
}

// Replica represents a chunk replica on a storage node
type Replica struct {
	ReplicaID  uuid.UUID  `json:"replica_id"`
	ChunkID    uuid.UUID  `json:"chunk_id"`
	NodeID     uuid.UUID  `json:"node_id"`
	Status     string     `json:"status"` // "pending", "stored", "verified", "failed"
	StoredAt   *time.Time `json:"stored_at,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Manifest represents a file with all its chunks and replica locations
type Manifest struct {
	FileID         uuid.UUID   `json:"file_id"`
	Filename       string      `json:"filename"`
	ChunkSize      int         `json:"chunk_size"`
	EncryptionAlg  string      `json:"encryption_alg"`
	WrappedFileKey string      `json:"wrapped_file_key"`
	Chunks         []ChunkInfo `json:"chunks"`
}

// ChunkInfo represents a chunk with its replica locations
type ChunkInfo struct {
	ChunkID     uuid.UUID   `json:"chunk_id"`
	ChunkIndex  int         `json:"chunk_index"`
	ContentHash string      `json:"content_hash"`
	ChunkSize   int         `json:"chunk_size"`
	Nodes       []uuid.UUID `json:"nodes"` // Node IDs where replicas exist
}

// Request/Response types

// RegisterNodeRequest is used to register a new storage node
type RegisterNodeRequest struct {
	PublicKey     string `json:"public_key" validate:"required"`
	Hostname      string `json:"hostname" validate:"required"`
	Endpoint      string `json:"endpoint" validate:"required"`
	CapacityBytes int64  `json:"capacity_bytes" validate:"required,min=1000000"` // At least 1MB
}

// HeartbeatRequest is used to update node health status
type HeartbeatRequest struct {
	Status    string `json:"status" validate:"required,oneof=healthy degraded offline"`
	UsedBytes int64  `json:"used_bytes" validate:"required,min=0"`
}

// CreateFileRequest is used to create a new file record
type CreateFileRequest struct {
	OwnerID        uuid.UUID `json:"owner_id" validate:"required"`
	Filename       string    `json:"filename" validate:"required,max=255"`
	OriginalSize   int64     `json:"original_size" validate:"required,min=1"`
	MimeType       string    `json:"mime_type" validate:"max=100"`
	ChunkSize      int       `json:"chunk_size" validate:"required,min=1000000"` // At least 1MB
	EncryptionAlg  string    `json:"encryption_alg" validate:"required"`
	WrappedFileKey string    `json:"wrapped_file_key" validate:"required"`
}

// CreateChunkRequest is used to create a chunk record
type CreateChunkRequest struct {
	FileID      uuid.UUID `json:"file_id" validate:"required"`
	ChunkIndex  int       `json:"chunk_index" validate:"required,min=0"`
	ChunkSize   int       `json:"chunk_size" validate:"required,min=1"`
	ContentHash string    `json:"content_hash" validate:"required,len=64"` // SHA-256 hex string
}

// CreateReplicaRequest is used to create a replica record
type CreateReplicaRequest struct {
	ChunkID uuid.UUID `json:"chunk_id" validate:"required"`
	NodeID  uuid.UUID `json:"node_id" validate:"required"`
}

// UpdateReplicaStatusRequest is used to update a replica's status
type UpdateReplicaStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending stored verified failed"`
}

// Error response type
type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
}
