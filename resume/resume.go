package resume

import (
	"bytes"
	"context"
	"database/sql/driver"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	Version      string       `json:"version" yaml:"version"`
	Title        string       `json:"title" yaml:"title"`
	Summary      string       `json:"summary" yaml:"summary"`
	PersonalInfo PersonalInfo `json:"personalInfo" yaml:"personalInfo"`
	Sections     []Section    `json:"sections" yaml:"sections"`
}

type PersonalInfo struct {
	Name        string `json:"name" yaml:"name"`
	Email       string `json:"email" yaml:"email"`
	LinkedIn    string `json:"linkedin" yaml:"linkedin"`
	Github      string `json:"github" yaml:"github"`
	Portfolio   string `json:"portfolio" yaml:"portfolio"`
	Citizenship string `json:"citizenship" yaml:"citizenship"`
	Phone       string `json:"phone" yaml:"phone"`
	// Links map[string]string `yaml:"linksi"`
}

type Section struct {
	Title      string              `json:"title" yaml:"title"`
	Experience []Experience        `json:"experience,omitempty" yaml:"experience,omitempty"`
	Education  []Education         `json:"education,omitempty" yaml:"education,omitempty"`
	Skills     map[string][]string `json:"skills,omitempty" yaml:"skills,omitempty"`
	Projects   []Project           `json:"projects,omitempty" yaml:"projects,omitempty"`
}

type Experience struct {
	Title            string   `yaml:"title" json:"title"`
	Company          string   `yaml:"company" json:"company"`
	Duration         string   `yaml:"duration" json:"duration"`
	Responsibilities []string `yaml:"responsibilities" json:"responsibilities"`
}

type Education struct {
	Degree             string   `yaml:"degree" json:"degree"`
	Institution        string   `yaml:"institution" json:"institution"`
	Location           string   `yaml:"location" json:"location"`
	Duration           string   `yaml:"duration" json:"duration"`
	GPA                string   `yaml:"gpa" json:"gpa"`
	Focus              string   `yaml:"focus" json:"focus"`
	RelevantCoursework []string `yaml:"relevantCoursework" json:"relevantCoursework"`
}

type Project struct {
	Name         string   `yaml:"title" json:"name"`
	Technologies []string `yaml:"technologies" json:"technologies"`
	Duration     string   `yaml:"duration" json:"duration"`
	Description  string   `yaml:"description" json:"description"`
	Github       string   `yaml:"github" json:"github"`
	Demo         string   `yaml:"demo" json:"demo"`
	Npm          string   `yaml:"npm" json:"npm"`
}

// TODO: Optional stylesheets??
func (resume Resume) ToHTML(w io.Writer) error {
	return ResumePage(resume).Render(context.TODO(), w)
	// return htmlTemplate.Execute(w, resume)
}

func (resume Resume) YAML() (string, error) {
	var writer strings.Builder
	if err := resume.ToYAML(&writer); err != nil {
		return "", err
	}
	return writer.String(), nil
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

func FromFiles(files ...*os.File) (Resume, error) {
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
		yaml.DisallowUnknownField(),
	).Decode(resume); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
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

// Implement methods to be used in a database
func (r Resume) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Resume) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Resume", value)
	}

	return json.Unmarshal(bytes, r)
}
