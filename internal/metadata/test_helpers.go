package metadata

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

// setupTestRepo creates a test repository with a clean database
func setupTestRepo(t *testing.T) *Repository {
	// Connect to test database
	// For now, use an in-memory approach or test database
	// In production, you'd use a test container or test database

	db := openTestDB(t)
	cleanTestDB(t, db)
	createTestSchema(t, db)

	return &Repository{db: db}
}

// openTestDB opens a connection to the test database
func openTestDB(t *testing.T) *sql.DB {
	connStr := "postgres://auta:auta_dev_password@localhost/auta_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("Could not connect to test database: %v. Run 'make test-db-up' to start it", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Skipf("Test database not available: %v. Run 'make test-db-up' to start it", err)
	}

	return db
}

// cleanTestDB truncates all test tables
func cleanTestDB(t *testing.T, db *sql.DB) {
	// Delete in reverse order of foreign keys
	tables := []string{"replicas", "chunks", "files", "nodes"}
	for _, table := range tables {
		if _, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE"); err != nil {
			// Table might not exist yet, that's ok
		}
	}
}

// createTestSchema creates the test database schema
func createTestSchema(t *testing.T, db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS nodes (
		node_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		public_key TEXT NOT NULL UNIQUE,
		hostname TEXT NOT NULL,
		endpoint TEXT NOT NULL UNIQUE,
		status VARCHAR(20) DEFAULT 'healthy' NOT NULL,
		capacity_bytes BIGINT NOT NULL,
		used_bytes BIGINT DEFAULT 0 NOT NULL,
		last_heartbeat_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS files (
		file_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		owner_id UUID NOT NULL,
		filename TEXT NOT NULL,
		original_size BIGINT NOT NULL,
		mime_type TEXT,
		chunk_size INT DEFAULT 4194304 NOT NULL,
		encryption_alg VARCHAR(50) DEFAULT 'AES-256-GCM' NOT NULL,
		wrapped_file_key TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_owner_filename UNIQUE (owner_id, filename)
	);

	CREATE TABLE IF NOT EXISTS chunks (
		chunk_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		file_id UUID NOT NULL REFERENCES files(file_id) ON DELETE CASCADE,
		chunk_index INT NOT NULL,
		chunk_size INT NOT NULL,
		content_hash VARCHAR(64) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_file_chunk_index UNIQUE (file_id, chunk_index)
	);

	CREATE TABLE IF NOT EXISTS replicas (
		replica_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		chunk_id UUID NOT NULL REFERENCES chunks(chunk_id) ON DELETE CASCADE,
		node_id UUID NOT NULL REFERENCES nodes(node_id),
		status VARCHAR(20) DEFAULT 'pending' NOT NULL,
		stored_at TIMESTAMP WITH TIME ZONE,
		verified_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_chunk_node_replica UNIQUE (chunk_id, node_id)
	);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Logf("Warning: could not create schema: %v", err)
	}
}
