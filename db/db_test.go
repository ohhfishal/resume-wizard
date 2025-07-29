package db_test

import (
	assert "github.com/alecthomas/assert/v2"
	"github.com/ohhfishal/resume-wizard/db"
	"github.com/ohhfishal/resume-wizard/resume"
	"testing"
)

func TestInsertResume(t *testing.T) {
	database, err := db.Open(t.Context(), "sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = database.InsertResume(t.Context(), db.InsertResumeParams{
		Name: "test",
		Body: &resume.Resume{},
	})
	assert.NoError(t, err)
}
