
CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL,
  name TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);

-- Holds the resumes that get tailored for an application
CREATE TABLE IF NOT EXISTS base_resumes  (
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

-- Holds a tailoring of a base resume for a given job
CREATE TABLE IF NOT EXISTS applications (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  base_resume_id INTEGER NOT NULL,

  company TEXT NOT NULL,
  position TEXT NOT NULL,
  description TEXT NOT NULL,
  resume TEXT NOT NULL,

  status TEXT CHECK( status IN ('pending', 'interviewed', 'rejected', 'accepted') ) DEFAULT 'pending' NOT NULL,

  -- User controls
  applied_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

  -- User does not control
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at DATETIME,

  FOREIGN KEY (base_resume_id) REFERENCES base_resumes(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Holds  user sessions creating an application
CREATE TABLE IF NOT EXISTS sessions (
  uuid TEXT PRIMARY KEY,
  base_resume_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,

  -- Input
  company TEXT NOT NULL,
  position TEXT NOT NULL,
  description TEXT NOT NULL,

  -- Output
  resume TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at DATETIME,

  FOREIGN KEY (base_resume_id, user_id) REFERENCES base_resumes(id, user_id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- CREATE TABLE IF NOT EXISTS resumes (
--   id INTEGER PRIMARY KEY,
--   name TEXT NOT NULL,
--   body TEXT NOT NULL
-- );


INSERT INTO users(id, name) VALUES (0, 'admin') ON CONFLICT DO NOTHING;
