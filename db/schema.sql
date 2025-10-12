
CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);

-- Holds the resumes that get tailored for an application
-- CREATE TABLE IF NOT EXISTS base_resumes  (
--   id INTEGER PRIMARY KEY,
--   user_id INTEGER NOT NULL,
--   name TEXT NOT NULL,
--   resume TEXT NOT NULL,
--   created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
--   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
--   last_used DATETIME,
--   deleted_at DATETIME,
--   FOREIGN KEY (user_id) REFERENCES users(id)
-- );
