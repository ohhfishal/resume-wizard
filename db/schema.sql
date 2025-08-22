
CREATE TABLE IF NOT EXISTS resumes (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  body TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS applications (
  resume_id INTEGER NOT NULL,
  company TEXT NOT NULL,
  position TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  status TEXT CHECK( status IN ('pending', 'interviewed', 'rejected', 'accepted') ) DEFAULT 'pending' NOT NULL,

  PRIMARY KEY (company, position),
  FOREIGN KEY(resume_id) REFERENCES resumes(id)
);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL,
  name TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
);

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

-- TODO: Hard coding examples
INSERT INTO users(id, name) VALUES (0, 'admin') ON CONFLICT DO NOTHING;

INSERT INTO base_resumes (user_id, name, resume) VALUES (0, 'Example Resume', '{}') ON CONFLICT DO NOTHING;

