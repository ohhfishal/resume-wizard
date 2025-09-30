
CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL,
  name TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS resumes (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  resume TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  last_used DATETIME,
  deleted_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS application_history (
  -- UUID per state change
  id INTEGER PRIMARY KEY,

  -- Foreign Keys
  user_id INTEGER NOT NULL,
  resume_id INTEGER NOT NULL,

  status TEXT CHECK( status IN ('pending', 'interviewed', 'rejected', 'accepted') ) DEFAULT 'pending' NOT NULL,

  -- Job info
  company TEXT NOT NULL,
  position TEXT NOT NULL,
  description TEXT,

  -- DB managed fields
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

  FOREIGN KEY (resume_id) REFERENCES resumes(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users(id, name) VALUES (0, 'admin') ON CONFLICT DO NOTHING;
