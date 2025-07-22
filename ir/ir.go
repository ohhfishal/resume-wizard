package ir

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"github.com/goccy/go-yaml"
)

//go:embed html.template
var rawTemplateHTML string
var htmlTemplate = template.Must(template.New("html").Parse(rawTemplateHTML))

const (
	JSON string = ".json"
	YAML = ".yaml"
	STDIN_BASE = "stdin"
	STDIN_EXT = ""
)

type Resume struct {
}

// TODO: Optional stylesheets??
func (resume Resume) ToHTML(w io.Writer) error {
	return htmlTemplate.Execute(w, resume)
}

func FromFile(file *os.File) (Resume, error) {
	info, err := file.Stat()
	if err != nil {
		return Resume{}, fmt.Errorf("getting file info: %w", err)
	}

	base := filepath.Base(info.Name())
	extension := filepath.Ext(info.Name()) 

	switch {
	case extension == YAML:
		return FromYAML(file)
	case extension == JSON:
		return FromJSON(file)
	case base == STDIN_BASE && extension == STDIN_EXT:
		data, err := io.ReadAll(file)
		if err != nil {
			return Resume{}, fmt.Errorf("reading from stdin: %s", err)
		}

		resume1, err := FromYAML(bytes.NewReader(data))
		if err == nil {
			return resume1, nil
		}

		resume2, err2 := FromJSON(bytes.NewReader(data))
		if err2 == nil {
			return resume2, nil
		}
		return Resume{}, fmt.Errorf(
			"failed to parse from stdin: %s", errors.Join(err, err2),
		)
	default: 
		return Resume{}, fmt.Errorf("unknown file type: %s %s", base, extension)
	}
}

func FromYAML(reader io.Reader) (Resume, error) {
	var resume Resume
	if err := yaml.NewDecoder(reader).Decode(&resume); err != nil {
		return resume, fmt.Errorf("parsing yaml: %w", err)
	}
	return resume, nil
}

func FromJSON(reader io.Reader) (Resume, error) {
	var resume Resume
	if err := json.NewDecoder(reader).Decode(&resume); err != nil {
		return resume, fmt.Errorf("parsing json: %w", err)
	}
	return resume, nil
}
