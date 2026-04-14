package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository provides data access methods for metadata operations
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Node operations

// RegisterNode creates a new storage node record
func (r *Repository) RegisterNode(ctx context.Context, node *Node) error {
	query := `
		INSERT INTO nodes (node_id, public_key, hostname, endpoint, status, capacity_bytes, used_bytes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		node.NodeID, node.PublicKey, node.Hostname, node.Endpoint, node.Status,
		node.CapacityBytes, node.UsedBytes, node.CreatedAt, node.UpdatedAt,
	)
	return err
}

// GetNode retrieves a node by ID
func (r *Repository) GetNode(ctx context.Context, nodeID uuid.UUID) (*Node, error) {
	query := `
		SELECT node_id, public_key, hostname, endpoint, status, capacity_bytes, used_bytes, last_heartbeat_at, created_at, updated_at
		FROM nodes WHERE node_id = $1
	`
	node := &Node{}
	err := r.db.QueryRowContext(ctx, query, nodeID).Scan(
		&node.NodeID, &node.PublicKey, &node.Hostname, &node.Endpoint, &node.Status,
		&node.CapacityBytes, &node.UsedBytes, &node.LastHeartbeatAt, &node.CreatedAt, &node.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("node not found")
	}
	return node, err
}

// ListHealthyNodes retrieves all healthy nodes
func (r *Repository) ListHealthyNodes(ctx context.Context) ([]*Node, error) {
	query := `
		SELECT node_id, public_key, hostname, endpoint, status, capacity_bytes, used_bytes, last_heartbeat_at, created_at, updated_at
		FROM nodes WHERE status = 'healthy'
		ORDER BY used_bytes ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*Node
	for rows.Next() {
		node := &Node{}
		if err := rows.Scan(
			&node.NodeID, &node.PublicKey, &node.Hostname, &node.Endpoint, &node.Status,
			&node.CapacityBytes, &node.UsedBytes, &node.LastHeartbeatAt, &node.CreatedAt, &node.UpdatedAt,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, rows.Err()
}

// UpdateNodeHeartbeat updates node heartbeat and status
func (r *Repository) UpdateNodeHeartbeat(ctx context.Context, nodeID uuid.UUID, status string, usedBytes int64) error {
	query := `
		UPDATE nodes SET status = $1, used_bytes = $2, last_heartbeat_at = $3, updated_at = $4
		WHERE node_id = $5
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, status, usedBytes, now, now, nodeID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("node not found")
	}
	return nil
}

// File operations

// CreateFile creates a new file record
func (r *Repository) CreateFile(ctx context.Context, file *File) error {
	query := `
		INSERT INTO files (file_id, owner_id, filename, original_size, mime_type, chunk_size, encryption_alg, wrapped_file_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		file.FileID, file.OwnerID, file.Filename, file.OriginalSize, file.MimeType,
		file.ChunkSize, file.EncryptionAlg, file.WrappedFileKey, file.CreatedAt, file.UpdatedAt,
	)
	return err
}

// GetFile retrieves a file by ID
func (r *Repository) GetFile(ctx context.Context, fileID uuid.UUID) (*File, error) {
	query := `
		SELECT file_id, owner_id, filename, original_size, mime_type, chunk_size, encryption_alg, wrapped_file_key, created_at, updated_at
		FROM files WHERE file_id = $1
	`
	file := &File{}
	err := r.db.QueryRowContext(ctx, query, fileID).Scan(
		&file.FileID, &file.OwnerID, &file.Filename, &file.OriginalSize, &file.MimeType,
		&file.ChunkSize, &file.EncryptionAlg, &file.WrappedFileKey, &file.CreatedAt, &file.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not found")
	}
	return file, err
}

// DeleteFile deletes a file and all its chunks and replicas
func (r *Repository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete replicas for all chunks of this file
	_, err = tx.ExecContext(ctx, `
		DELETE FROM replicas WHERE chunk_id IN (
			SELECT chunk_id FROM chunks WHERE file_id = $1
		)
	`, fileID)
	if err != nil {
		return err
	}

	// Delete chunks
	_, err = tx.ExecContext(ctx, `DELETE FROM chunks WHERE file_id = $1`, fileID)
	if err != nil {
		return err
	}

	// Delete file
	_, err = tx.ExecContext(ctx, `DELETE FROM files WHERE file_id = $1`, fileID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Chunk operations

// CreateChunk creates a new chunk record
func (r *Repository) CreateChunk(ctx context.Context, chunk *Chunk) error {
	query := `
		INSERT INTO chunks (chunk_id, file_id, chunk_index, chunk_size, content_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		chunk.ChunkID, chunk.FileID, chunk.ChunkIndex, chunk.ChunkSize, chunk.ContentHash, chunk.CreatedAt,
	)
	return err
}

// GetChunk retrieves a chunk by ID
func (r *Repository) GetChunk(ctx context.Context, chunkID uuid.UUID) (*Chunk, error) {
	query := `
		SELECT chunk_id, file_id, chunk_index, chunk_size, content_hash, created_at
		FROM chunks WHERE chunk_id = $1
	`
	chunk := &Chunk{}
	err := r.db.QueryRowContext(ctx, query, chunkID).Scan(
		&chunk.ChunkID, &chunk.FileID, &chunk.ChunkIndex, &chunk.ChunkSize, &chunk.ContentHash, &chunk.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("chunk not found")
	}
	return chunk, err
}

// GetFileChunks retrieves all chunks for a file
func (r *Repository) GetFileChunks(ctx context.Context, fileID uuid.UUID) ([]*Chunk, error) {
	query := `
		SELECT chunk_id, file_id, chunk_index, chunk_size, content_hash, created_at
		FROM chunks WHERE file_id = $1
		ORDER BY chunk_index ASC
	`
	rows, err := r.db.QueryContext(ctx, query, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []*Chunk
	for rows.Next() {
		chunk := &Chunk{}
		if err := rows.Scan(
			&chunk.ChunkID, &chunk.FileID, &chunk.ChunkIndex, &chunk.ChunkSize, &chunk.ContentHash, &chunk.CreatedAt,
		); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, rows.Err()
}

// Replica operations

// CreateReplica creates a new replica record
func (r *Repository) CreateReplica(ctx context.Context, replica *Replica) error {
	query := `
		INSERT INTO replicas (replica_id, chunk_id, node_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		replica.ReplicaID, replica.ChunkID, replica.NodeID, replica.Status, replica.CreatedAt, replica.UpdatedAt,
	)
	return err
}

// GetReplica retrieves a replica by ID
func (r *Repository) GetReplica(ctx context.Context, replicaID uuid.UUID) (*Replica, error) {
	query := `
		SELECT replica_id, chunk_id, node_id, status, stored_at, verified_at, created_at, updated_at
		FROM replicas WHERE replica_id = $1
	`
	replica := &Replica{}
	err := r.db.QueryRowContext(ctx, query, replicaID).Scan(
		&replica.ReplicaID, &replica.ChunkID, &replica.NodeID, &replica.Status, &replica.StoredAt, &replica.VerifiedAt, &replica.CreatedAt, &replica.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("replica not found")
	}
	return replica, err
}

// GetChunkReplicas retrieves all replicas for a chunk
func (r *Repository) GetChunkReplicas(ctx context.Context, chunkID uuid.UUID) ([]*Replica, error) {
	query := `
		SELECT replica_id, chunk_id, node_id, status, stored_at, verified_at, created_at, updated_at
		FROM replicas WHERE chunk_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, chunkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replicas []*Replica
	for rows.Next() {
		replica := &Replica{}
		if err := rows.Scan(
			&replica.ReplicaID, &replica.ChunkID, &replica.NodeID, &replica.Status, &replica.StoredAt, &replica.VerifiedAt, &replica.CreatedAt, &replica.UpdatedAt,
		); err != nil {
			return nil, err
		}
		replicas = append(replicas, replica)
	}
	return replicas, rows.Err()
}

// UpdateReplicaStatus updates the status of a replica
func (r *Repository) UpdateReplicaStatus(ctx context.Context, replicaID uuid.UUID, status string) error {
	now := time.Now()
	query := `
		UPDATE replicas SET status = $1, updated_at = $2`

	// Update timestamp based on status
	if status == "stored" {
		query += `, stored_at = $2`
	} else if status == "verified" {
		query += `, verified_at = $2`
	}

	query += ` WHERE replica_id = $3`

	result, err := r.db.ExecContext(ctx, query, status, now, replicaID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("replica not found")
	}
	return nil
}

// CountChunkReplicas returns the number of verified replicas for a chunk
func (r *Repository) CountChunkReplicas(ctx context.Context, chunkID uuid.UUID, status string) (int, error) {
	query := `SELECT COUNT(*) FROM replicas WHERE chunk_id = $1 AND status = $2`
	var count int
	err := r.db.QueryRowContext(ctx, query, chunkID, status).Scan(&count)
	return count, err
}

// GetVerifiedReplicaNodes returns node IDs of verified replicas for a chunk
func (r *Repository) GetVerifiedReplicaNodes(ctx context.Context, chunkID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT node_id FROM replicas 
		WHERE chunk_id = $1 AND status = 'verified'
	`
	rows, err := r.db.QueryContext(ctx, query, chunkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodeIDs []uuid.UUID
	for rows.Next() {
		var nodeID uuid.UUID
		if err := rows.Scan(&nodeID); err != nil {
			return nil, err
		}
		nodeIDs = append(nodeIDs, nodeID)
	}
	return nodeIDs, rows.Err()
}
