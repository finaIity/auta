-- Create nodes table
CREATE TABLE nodes (
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

-- Create files table
CREATE TABLE files (
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

-- Create chunks table
CREATE TABLE chunks (
    chunk_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL REFERENCES files(file_id) ON DELETE CASCADE,
    chunk_index INT NOT NULL,
    chunk_size INT NOT NULL,
    content_hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_file_chunk_index UNIQUE (file_id, chunk_index)
);

-- Create replicas table
CREATE TABLE replicas (
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

-- Create indexes for common queries
CREATE INDEX idx_files_owner_id ON files(owner_id);
CREATE INDEX idx_chunks_file_id ON chunks(file_id);
CREATE INDEX idx_replicas_chunk_id ON replicas(chunk_id);
CREATE INDEX idx_replicas_node_id ON replicas(node_id);
CREATE INDEX idx_replicas_status ON replicas(status);
CREATE INDEX idx_nodes_status ON nodes(status);
CREATE INDEX idx_nodes_last_heartbeat ON nodes(last_heartbeat_at);
