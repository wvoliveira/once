CREATE TABLE
  IF NOT EXISTS content (
    id TEXT PRIMARY KEY,
    file_id TEXT,
    file_name TEXT,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL
  );


CREATE INDEX IF NOT EXISTS idx_content_file_id ON content (file_id);
CREATE INDEX IF NOT EXISTS idx_content_expires_at ON content (expires_at);
