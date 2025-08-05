package resume_test

import (
	"strings"
	"testing"

	"github.com/ohhfishal/resume-wizard/resume"
)

func TestYAML(t *testing.T) {
	tests := []struct {
		Body string
		OK   bool
	}{
		{
			Body: "version: ok",
			OK:   true,
		},
		{
			Body: "field: unknown",
			OK:   false,
		},
		{
			Body: "version: bad\nfoo bar",
			OK:   false,
		},
	}

	for _, test := range tests {
		var dummy resume.Resume
		err := resume.FromYAML(strings.NewReader(test.Body), &dummy)
		if test.OK && err != nil {
			t.Fatal(err)
		}
		if !test.OK && err == nil {
			t.Fatal("No error when expected one")
		}
	}

}
