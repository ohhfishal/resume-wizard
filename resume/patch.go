package resume

import (
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-yaml"
)


type Patch struct {
	PersonalInfo PersonalInfo      `yaml:"personalInfo"`
	Override     map[string]string `yaml:"override"`
}

func (resume *Resume) ApplyPatch(reader io.Reader) error {
	var patch Patch
	if err := yaml.NewDecoder(reader).Decode(&patch); err != nil {
		return fmt.Errorf("parsing patch yaml: %w", err)
	}
	resume.PersonalInfo = patch.PersonalInfo

	var buffer strings.Builder
	if err := resume.ToYAML(&buffer); err != nil {
		return fmt.Errorf("resume form is invalid: %w", err)
	}

	content := buffer.String()
	for key, value := range patch.Override {
		content = strings.ReplaceAll(content, "$" + key, value)
	}

	if err := FromYAML(strings.NewReader(content), resume); err != nil {
		return fmt.Errorf("applying new content: %w", err)
	}
	return nil
}

