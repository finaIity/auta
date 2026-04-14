package metadata

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

// Service represents the metadata service
type Service struct {
	db *sql.DB
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

	return &Service{db: db}, nil
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

// Handlers

func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Service) handleCreateFile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleGetFile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleCreateChunk(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleGetChunk(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleRegisterNode(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleGetNode(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleNodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleCreateReplica(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Service) handleGetChunkReplicas(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}
