package metadata

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestNodeRegisterAndGet tests node registration and retrieval
func TestNodeRegisterAndGet(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// Create a test node
	node := &Node{
		NodeID:        uuid.New(),
		PublicKey:     "test_pk_001",
		Hostname:      "test-host",
		Endpoint:      "http://test:8001",
		Status:        "healthy",
		CapacityBytes: 1000000000,
		UsedBytes:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Register the node
	err := repo.RegisterNode(ctx, node)
	if err != nil {
		t.Fatalf("RegisterNode failed: %v", err)
	}

	// Retrieve the node
	retrieved, err := repo.GetNode(ctx, node.NodeID)
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}

	// Verify fields
	if retrieved.NodeID != node.NodeID {
		t.Errorf("NodeID mismatch: got %v, want %v", retrieved.NodeID, node.NodeID)
	}
	if retrieved.PublicKey != node.PublicKey {
		t.Errorf("PublicKey mismatch: got %v, want %v", retrieved.PublicKey, node.PublicKey)
	}
	if retrieved.Status != node.Status {
		t.Errorf("Status mismatch: got %v, want %v", retrieved.Status, node.Status)
	}
}

// TestNodeHeartbeat tests heartbeat update
func TestNodeHeartbeat(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// Create and register a node
	node := &Node{
		NodeID:        uuid.New(),
		PublicKey:     "test_pk_002",
		Hostname:      "test-host-2",
		Endpoint:      "http://test2:8001",
		Status:        "healthy",
		CapacityBytes: 1000000000,
		UsedBytes:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.RegisterNode(ctx, node)

	// Update heartbeat
	err := repo.UpdateNodeHeartbeat(ctx, node.NodeID, "degraded", 500000000)
	if err != nil {
		t.Fatalf("UpdateNodeHeartbeat failed: %v", err)
	}

	// Retrieve and verify
	updated, err := repo.GetNode(ctx, node.NodeID)
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}

	if updated.Status != "degraded" {
		t.Errorf("Status not updated: got %v, want degraded", updated.Status)
	}
	if updated.UsedBytes != 500000000 {
		t.Errorf("UsedBytes not updated: got %v, want 500000000", updated.UsedBytes)
	}
	if updated.LastHeartbeatAt == nil {
		t.Error("LastHeartbeatAt not set")
	}
}

// TestFileCreateAndGet tests file creation and retrieval
func TestFileCreateAndGet(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	file := &File{
		FileID:         uuid.New(),
		OwnerID:        uuid.New(),
		Filename:       "test-file.bin",
		OriginalSize:   8388608,
		MimeType:       "application/octet-stream",
		ChunkSize:      4194304,
		EncryptionAlg:  "AES-256-GCM",
		WrappedFileKey: "encrypted_key_123",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create file
	err := repo.CreateFile(ctx, file)
	if err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}

	// Retrieve file
	retrieved, err := repo.GetFile(ctx, file.FileID)
	if err != nil {
		t.Fatalf("GetFile failed: %v", err)
	}

	if retrieved.Filename != file.Filename {
		t.Errorf("Filename mismatch: got %v, want %v", retrieved.Filename, file.Filename)
	}
	if retrieved.OriginalSize != file.OriginalSize {
		t.Errorf("OriginalSize mismatch: got %v, want %v", retrieved.OriginalSize, file.OriginalSize)
	}
}

// TestChunkCreateAndGet tests chunk creation and retrieval
func TestChunkCreateAndGet(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// First create a file
	file := &File{
		FileID:         uuid.New(),
		OwnerID:        uuid.New(),
		Filename:       "test-chunks.bin",
		OriginalSize:   8388608,
		ChunkSize:      4194304,
		EncryptionAlg:  "AES-256-GCM",
		WrappedFileKey: "key",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	repo.CreateFile(ctx, file)

	// Create chunks
	chunk := &Chunk{
		ChunkID:     uuid.New(),
		FileID:      file.FileID,
		ChunkIndex:  0,
		ChunkSize:   4194304,
		ContentHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		CreatedAt:   time.Now(),
	}

	err := repo.CreateChunk(ctx, chunk)
	if err != nil {
		t.Fatalf("CreateChunk failed: %v", err)
	}

	// Retrieve chunk
	retrieved, err := repo.GetChunk(ctx, chunk.ChunkID)
	if err != nil {
		t.Fatalf("GetChunk failed: %v", err)
	}

	if retrieved.ChunkIndex != chunk.ChunkIndex {
		t.Errorf("ChunkIndex mismatch: got %v, want %v", retrieved.ChunkIndex, chunk.ChunkIndex)
	}
	if retrieved.ContentHash != chunk.ContentHash {
		t.Errorf("ContentHash mismatch: got %v, want %v", retrieved.ContentHash, chunk.ContentHash)
	}
}

// TestReplicaCreateAndGet tests replica creation and retrieval
func TestReplicaCreateAndGet(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// Create node, file, and chunk
	node := &Node{
		NodeID:        uuid.New(),
		PublicKey:     "pk_003",
		Hostname:      "host3",
		Endpoint:      "http://host3:8001",
		Status:        "healthy",
		CapacityBytes: 1000000000,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.RegisterNode(ctx, node)

	file := &File{
		FileID:         uuid.New(),
		OwnerID:        uuid.New(),
		Filename:       "test.bin",
		OriginalSize:   4194304,
		ChunkSize:      4194304,
		EncryptionAlg:  "AES-256-GCM",
		WrappedFileKey: "key",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	repo.CreateFile(ctx, file)

	chunk := &Chunk{
		ChunkID:     uuid.New(),
		FileID:      file.FileID,
		ChunkIndex:  0,
		ChunkSize:   4194304,
		ContentHash: "hash",
		CreatedAt:   time.Now(),
	}
	repo.CreateChunk(ctx, chunk)

	// Create replica
	replica := &Replica{
		ReplicaID: uuid.New(),
		ChunkID:   chunk.ChunkID,
		NodeID:    node.NodeID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateReplica(ctx, replica)
	if err != nil {
		t.Fatalf("CreateReplica failed: %v", err)
	}

	// Retrieve replica
	retrieved, err := repo.GetReplica(ctx, replica.ReplicaID)
	if err != nil {
		t.Fatalf("GetReplica failed: %v", err)
	}

	if retrieved.Status != replica.Status {
		t.Errorf("Status mismatch: got %v, want %v", retrieved.Status, replica.Status)
	}
	if retrieved.NodeID != node.NodeID {
		t.Errorf("NodeID mismatch: got %v, want %v", retrieved.NodeID, node.NodeID)
	}
}

// TestReplicaStatusUpdate tests replica status transitions
func TestReplicaStatusUpdate(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// Create node, file, chunk, and replica
	node := &Node{
		NodeID:        uuid.New(),
		PublicKey:     "pk_004",
		Hostname:      "host4",
		Endpoint:      "http://host4:8001",
		Status:        "healthy",
		CapacityBytes: 1000000000,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.RegisterNode(ctx, node)

	file := &File{
		FileID:         uuid.New(),
		OwnerID:        uuid.New(),
		Filename:       "test.bin",
		OriginalSize:   4194304,
		ChunkSize:      4194304,
		EncryptionAlg:  "AES-256-GCM",
		WrappedFileKey: "key",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	repo.CreateFile(ctx, file)

	chunk := &Chunk{
		ChunkID:     uuid.New(),
		FileID:      file.FileID,
		ChunkIndex:  0,
		ChunkSize:   4194304,
		ContentHash: "hash",
		CreatedAt:   time.Now(),
	}
	repo.CreateChunk(ctx, chunk)

	replica := &Replica{
		ReplicaID: uuid.New(),
		ChunkID:   chunk.ChunkID,
		NodeID:    node.NodeID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.CreateReplica(ctx, replica)

	// Update status: pending -> stored
	err := repo.UpdateReplicaStatus(ctx, replica.ReplicaID, "stored")
	if err != nil {
		t.Fatalf("UpdateReplicaStatus to stored failed: %v", err)
	}

	stored, _ := repo.GetReplica(ctx, replica.ReplicaID)
	if stored.Status != "stored" {
		t.Errorf("Status not updated to stored: got %v", stored.Status)
	}
	if stored.StoredAt == nil {
		t.Error("StoredAt timestamp not set")
	}

	// Update status: stored -> verified
	err = repo.UpdateReplicaStatus(ctx, replica.ReplicaID, "verified")
	if err != nil {
		t.Fatalf("UpdateReplicaStatus to verified failed: %v", err)
	}

	verified, _ := repo.GetReplica(ctx, replica.ReplicaID)
	if verified.Status != "verified" {
		t.Errorf("Status not updated to verified: got %v", verified.Status)
	}
	if verified.VerifiedAt == nil {
		t.Error("VerifiedAt timestamp not set")
	}
}

// TestFileDelete tests cascading delete
func TestFileDelete(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// Create node
	node := &Node{
		NodeID:        uuid.New(),
		PublicKey:     "pk_005",
		Hostname:      "host5",
		Endpoint:      "http://host5:8001",
		Status:        "healthy",
		CapacityBytes: 1000000000,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	repo.RegisterNode(ctx, node)

	// Create file
	file := &File{
		FileID:         uuid.New(),
		OwnerID:        uuid.New(),
		Filename:       "delete-test.bin",
		OriginalSize:   4194304,
		ChunkSize:      4194304,
		EncryptionAlg:  "AES-256-GCM",
		WrappedFileKey: "key",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	repo.CreateFile(ctx, file)

	// Create chunk
	chunk := &Chunk{
		ChunkID:     uuid.New(),
		FileID:      file.FileID,
		ChunkIndex:  0,
		ChunkSize:   4194304,
		ContentHash: "hash",
		CreatedAt:   time.Now(),
	}
	repo.CreateChunk(ctx, chunk)

	// Create replica
	replica := &Replica{
		ReplicaID: uuid.New(),
		ChunkID:   chunk.ChunkID,
		NodeID:    node.NodeID,
		Status:    "verified",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.CreateReplica(ctx, replica)

	// Delete file (should cascade)
	err := repo.DeleteFile(ctx, file.FileID)
	if err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	// Verify file is gone
	_, err = repo.GetFile(ctx, file.FileID)
	if err == nil {
		t.Error("File should not exist after delete")
	}

	// Verify chunks are gone
	_, err = repo.GetChunk(ctx, chunk.ChunkID)
	if err == nil {
		t.Error("Chunk should not exist after file delete")
	}

	// Verify replicas are gone
	_, err = repo.GetReplica(ctx, replica.ReplicaID)
	if err == nil {
		t.Error("Replica should not exist after file delete")
	}
}

// TestGetVerifiedReplicaNodes tests getting verified replica node IDs
func TestGetVerifiedReplicaNodes(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()

	ctx := context.Background()

	// Create 3 nodes
	nodes := make([]uuid.UUID, 3)
	for i := 0; i < 3; i++ {
		node := &Node{
			NodeID:        uuid.New(),
			PublicKey:     "pk_" + string(rune(i)),
			Hostname:      "host",
			Endpoint:      "http://host:8001",
			Status:        "healthy",
			CapacityBytes: 1000000000,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		repo.RegisterNode(ctx, node)
		nodes[i] = node.NodeID
	}

	// Create file and chunk
	file := &File{
		FileID:         uuid.New(),
		OwnerID:        uuid.New(),
		Filename:       "test.bin",
		OriginalSize:   4194304,
		ChunkSize:      4194304,
		EncryptionAlg:  "AES-256-GCM",
		WrappedFileKey: "key",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	repo.CreateFile(ctx, file)

	chunk := &Chunk{
		ChunkID:     uuid.New(),
		FileID:      file.FileID,
		ChunkIndex:  0,
		ChunkSize:   4194304,
		ContentHash: "hash",
		CreatedAt:   time.Now(),
	}
	repo.CreateChunk(ctx, chunk)

	// Create replicas on all nodes - only 2 verified
	for i := 0; i < 3; i++ {
		status := "verified"
		if i == 2 {
			status = "pending"
		}

		replica := &Replica{
			ReplicaID: uuid.New(),
			ChunkID:   chunk.ChunkID,
			NodeID:    nodes[i],
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		repo.CreateReplica(ctx, replica)
	}

	// Get verified replicas
	verifiedNodes, err := repo.GetVerifiedReplicaNodes(ctx, chunk.ChunkID)
	if err != nil {
		t.Fatalf("GetVerifiedReplicaNodes failed: %v", err)
	}

	if len(verifiedNodes) != 2 {
		t.Errorf("Expected 2 verified nodes, got %v", len(verifiedNodes))
	}

	// Verify the right nodes are returned
	for _, nodeID := range verifiedNodes {
		if nodeID == nodes[2] {
			t.Error("Pending replica should not be included in verified nodes")
		}
	}
}
