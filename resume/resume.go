package resume

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type SectionType string

const (
	SectionTypeDefault    SectionType = ""
	SectionTypeExperience             = "experience"
	SectionTypeEducation              = "education"
	SectionTypeSkills                 = "skills"
	SectionTypeProjects               = "projects"
)

var redactedPersonalInfo = PersonalInfo{
	Name:   "REDACTED_NAME",
	Email:  "email@email.com",
	Github: "https://github.com/",
}

const (
	JSON       string = ".json"
	YAML              = ".yaml"
	STDIN_BASE        = "stdin"
	STDIN_EXT         = ""
)

type Resume struct {
	Version      string       `yaml:"version"`
	Title        string       `yaml:"title"`
	Summary      string       `yaml:"summary"`
	PersonalInfo PersonalInfo `yaml:"personalInfo"`
	Sections     []Section    `yaml:"sections"`
}

type PersonalInfo struct {
	Name      string `yaml:"name"`
	Email     string `yaml:"email"`
	LinkedIn  string `yaml:"linkedin"`
	Github    string `yaml:"github"`
	Portfolio string `yaml:"portfolio"`
	// Links map[string]string `yaml:"linksi"`
}

type Section struct {
	Title      string              `yaml:"title"`
	Experience []Experience        `yaml:"experience,omitempty"`
	Education  []Education         `yaml:"education,omitempty"`
	Skills     map[string][]string `yaml:"skills,omitempty"`
	Projects   []Project           `yaml:"projects,omitempty"`
}

type Experience struct {
	Title            string   `yaml:"title"`
	Company          string   `yaml:"company"`
	Duration         string   `yaml:"duration"`
	Responsibilities []string `yaml:"responsibilities"`
}

type Education struct {
	Degree             string   `yaml:"degree"`
	Institution        string   `yaml:"institution"`
	Location           string   `yaml:"location"`
	Duration           string   `yaml:"duration"`
	GPA                string   `yaml:"gpa"`
	Focus              string   `yaml:"focus"`
	RelevantCoursework []string `yaml:"relevantCoursework"`
}

type Project struct {
	Name         string   `yaml:"title"`
	Technologies []string `yaml:"technologies"`
	Duration     string   `yaml:"duration"`
	Description  string   `yaml:"description"`
	Github       string   `yaml:"github"`
	Demo         string   `yaml:"demo"`
	Npm          string   `yaml:"npm"`
}

// TODO: Optional stylesheets??
func (resume Resume) ToHTML(w io.Writer) error {
	return ResumePage(resume).Render(context.TODO(), w)
	// return htmlTemplate.Execute(w, resume)
}

func (resume Resume) ToYAML(w io.Writer) error {
	encoder := yaml.NewEncoder(
		w,
	)
	defer encoder.Close()
	return encoder.Encode(resume)
}

func (resume Resume) ToJSON(w io.Writer) error {
	return json.NewEncoder(
		w,
	).Encode(resume)
}

func FromFiles(files []*os.File) (Resume, error) {
	if len(files) == 0 {
		return Resume{}, errors.New("must include at least one file")
	}
	var resume Resume
	for _, file := range files {
		if file == nil {
			return Resume{}, errors.New("provided nil file")
		}

		if err := fromFile(file, &resume); err != nil {
			return Resume{}, fmt.Errorf(`applying file "%s": %w`, file.Name(), err)
		}
	}

	return resume, nil
}

func fromFile(file *os.File, resume *Resume) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}

	base := filepath.Base(info.Name())
	extension := filepath.Ext(info.Name())

	switch {
	case extension == JSON:
		if err := FromJSON(file, resume); err != nil {
			return fmt.Errorf("converting from yaml: %w", err)
		}
	case extension == YAML:
		if err := FromYAML(file, resume); err != nil {
			return fmt.Errorf("converting from yaml: %w", err)
		}
	case base == STDIN_BASE && extension == STDIN_EXT:
		data, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("reading from stdin: %s", err)
		}

		if err := FromYAML(bytes.NewReader(data), resume); err != nil {
			return fmt.Errorf("converting from yaml: %w", err)
		}
	default:
		return fmt.Errorf("unknown file type: %s %s", base, extension)
	}
	return nil
}

func FromJSON(reader io.Reader, resume *Resume) error {
	if err := json.NewDecoder(
		reader,
	).Decode(resume); err != nil {
		return fmt.Errorf("parsing json: %w", err)
	}
	return nil
}

func FromYAML(reader io.Reader, resume *Resume) error {
	if err := yaml.NewDecoder(
		reader,
	).Decode(resume); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
	}
	return nil
}

func decode[T any](input any, output *T) error {
	data, err := yaml.Marshal(input)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, output); err != nil {
		return err
	}

	return nil
}

func (resume *Resume) HidePersonalInfo() {
	resume.PersonalInfo = redactedPersonalInfo

	for i, section := range resume.Sections {
		for j, _ := range section.Education {
			resume.Sections[i].Education[j].Institution = fmt.Sprintf("$REDACTED_INSTITUTION_%d", j)
		}
		for j, _ := range section.Experience {
			resume.Sections[i].Experience[j].Company = fmt.Sprintf("$REDACTED_COMPANY_%d", j)
		}
	}
}
