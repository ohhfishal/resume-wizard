
CREATE TABLE IF NOT EXISTS resumes (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  body TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS applications (
  resume_id INTEGER NOT NULL,
  company TEXT NOT NULL,
  position TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  status TEXT CHECK( status IN ('pending', 'interviewed', 'rejected', 'accepted') ) DEFAULT 'pending',

  PRIMARY KEY (company, position),
  FOREIGN KEY(resume_id) REFERENCES resumes(id)
);
