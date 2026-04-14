package metadata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Service represents the metadata service
type Service struct {
	db   *sql.DB
	repo *Repository
}

// NewService creates a new metadata service instance
func NewService(dbURL string) (*Service, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Service{
		db:   db,
		repo: NewRepository(db),
	}, nil
}

// Start starts the metadata service server
func (s *Service) Start(port string) error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.handleHealth)

	// File endpoints
	mux.HandleFunc("POST /files", s.handleCreateFile)
	mux.HandleFunc("GET /files/{file_id}", s.handleGetFile)
	mux.HandleFunc("DELETE /files/{file_id}", s.handleDeleteFile)

	// Chunk endpoints
	mux.HandleFunc("POST /chunks", s.handleCreateChunk)
	mux.HandleFunc("GET /chunks/{chunk_id}", s.handleGetChunk)

	// Node endpoints
	mux.HandleFunc("POST /nodes", s.handleRegisterNode)
	mux.HandleFunc("GET /nodes/{node_id}", s.handleGetNode)
	mux.HandleFunc("POST /nodes/{node_id}/heartbeat", s.handleNodeHeartbeat)

	// Replica endpoints
	mux.HandleFunc("POST /replicas", s.handleCreateReplica)
	mux.HandleFunc("GET /replicas/chunk/{chunk_id}", s.handleGetChunkReplicas)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return server.ListenAndServe()
}

// Close closes the database connection
func (s *Service) Close() error {
	return s.db.Close()
}

// Helper functions

func (s *Service) respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Service) respondError(w http.ResponseWriter, error string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:      error,
		StatusCode: statusCode,
	})
}

// Handlers

func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.respondJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (s *Service) handleRegisterNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	node := &Node{
		NodeID:        uuid.New(),
		PublicKey:     req.PublicKey,
		Hostname:      req.Hostname,
		Endpoint:      req.Endpoint,
		Status:        "healthy",
		CapacityBytes: req.CapacityBytes,
		UsedBytes:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.RegisterNode(r.Context(), node); err != nil {
		s.respondError(w, fmt.Sprintf("failed to register node: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, node, http.StatusCreated)
}

func (s *Service) handleGetNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeIDStr := r.PathValue("node_id")
	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		s.respondError(w, "invalid node_id", http.StatusBadRequest)
		return
	}

	node, err := s.repo.GetNode(r.Context(), nodeID)
	if err != nil {
		s.respondError(w, "node not found", http.StatusNotFound)
		return
	}

	s.respondJSON(w, node, http.StatusOK)
}

func (s *Service) handleNodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeIDStr := r.PathValue("node_id")
	nodeID, err := uuid.Parse(nodeIDStr)
	if err != nil {
		s.respondError(w, "invalid node_id", http.StatusBadRequest)
		return
	}

	var req HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.repo.UpdateNodeHeartbeat(r.Context(), nodeID, req.Status, req.UsedBytes); err != nil {
		s.respondError(w, fmt.Sprintf("failed to update heartbeat: %v", err), http.StatusInternalServerError)
		return
	}

	node, err := s.repo.GetNode(r.Context(), nodeID)
	if err != nil {
		s.respondError(w, "node not found", http.StatusNotFound)
		return
	}

	s.respondJSON(w, node, http.StatusOK)
}

func (s *Service) handleCreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	file := &File{
		FileID:         uuid.New(),
		OwnerID:        req.OwnerID,
		Filename:       req.Filename,
		OriginalSize:   req.OriginalSize,
		MimeType:       req.MimeType,
		ChunkSize:      req.ChunkSize,
		EncryptionAlg:  req.EncryptionAlg,
		WrappedFileKey: req.WrappedFileKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateFile(r.Context(), file); err != nil {
		s.respondError(w, fmt.Sprintf("failed to create file: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, file, http.StatusCreated)
}

func (s *Service) handleGetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileIDStr := r.PathValue("file_id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		s.respondError(w, "invalid file_id", http.StatusBadRequest)
		return
	}

	file, err := s.repo.GetFile(r.Context(), fileID)
	if err != nil {
		s.respondError(w, "file not found", http.StatusNotFound)
		return
	}

	// Build manifest with chunks and replicas
	chunks, err := s.repo.GetFileChunks(r.Context(), fileID)
	if err != nil {
		s.respondError(w, fmt.Sprintf("failed to retrieve chunks: %v", err), http.StatusInternalServerError)
		return
	}

	chunkInfos := make([]ChunkInfo, len(chunks))
	for i, chunk := range chunks {
		nodes, err := s.repo.GetVerifiedReplicaNodes(r.Context(), chunk.ChunkID)
		if err != nil {
			s.respondError(w, fmt.Sprintf("failed to retrieve replicas: %v", err), http.StatusInternalServerError)
			return
		}
		chunkInfos[i] = ChunkInfo{
			ChunkID:     chunk.ChunkID,
			ChunkIndex:  chunk.ChunkIndex,
			ContentHash: chunk.ContentHash,
			ChunkSize:   chunk.ChunkSize,
			Nodes:       nodes,
		}
	}

	manifest := &Manifest{
		FileID:         file.FileID,
		Filename:       file.Filename,
		ChunkSize:      file.ChunkSize,
		EncryptionAlg:  file.EncryptionAlg,
		WrappedFileKey: file.WrappedFileKey,
		Chunks:         chunkInfos,
	}

	s.respondJSON(w, manifest, http.StatusOK)
}

func (s *Service) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileIDStr := r.PathValue("file_id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		s.respondError(w, "invalid file_id", http.StatusBadRequest)
		return
	}

	if err := s.repo.DeleteFile(r.Context(), fileID); err != nil {
		s.respondError(w, fmt.Sprintf("failed to delete file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) handleCreateChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateChunkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	chunk := &Chunk{
		ChunkID:     uuid.New(),
		FileID:      req.FileID,
		ChunkIndex:  req.ChunkIndex,
		ChunkSize:   req.ChunkSize,
		ContentHash: req.ContentHash,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateChunk(r.Context(), chunk); err != nil {
		s.respondError(w, fmt.Sprintf("failed to create chunk: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, chunk, http.StatusCreated)
}

func (s *Service) handleGetChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chunkIDStr := r.PathValue("chunk_id")
	chunkID, err := uuid.Parse(chunkIDStr)
	if err != nil {
		s.respondError(w, "invalid chunk_id", http.StatusBadRequest)
		return
	}

	chunk, err := s.repo.GetChunk(r.Context(), chunkID)
	if err != nil {
		s.respondError(w, "chunk not found", http.StatusNotFound)
		return
	}

	s.respondJSON(w, chunk, http.StatusOK)
}

func (s *Service) handleCreateReplica(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateReplicaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	replica := &Replica{
		ReplicaID: uuid.New(),
		ChunkID:   req.ChunkID,
		NodeID:    req.NodeID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateReplica(r.Context(), replica); err != nil {
		s.respondError(w, fmt.Sprintf("failed to create replica: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, replica, http.StatusCreated)
}

func (s *Service) handleGetChunkReplicas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chunkIDStr := r.PathValue("chunk_id")
	chunkID, err := uuid.Parse(chunkIDStr)
	if err != nil {
		s.respondError(w, "invalid chunk_id", http.StatusBadRequest)
		return
	}

	replicas, err := s.repo.GetChunkReplicas(r.Context(), chunkID)
	if err != nil {
		s.respondError(w, fmt.Sprintf("failed to retrieve replicas: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, replicas, http.StatusOK)
}
