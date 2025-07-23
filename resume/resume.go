package resume

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

//go:embed html.template
var rawTemplateHTML string
var htmlTemplate = template.Must(template.New("html").Parse(rawTemplateHTML))

type SectionType string

const (
	SectionTypeDefault    SectionType = ""
	SectionTypeExperience             = "experience"
	SectionTypeEducation              = "education"
	SectionTypeSkills                 = "skills"
	SectionTypeProjects               = "projects"
)

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
	data       map[string]any      `yaml:"-"`
	experience []Experience        `yaml:"experience"`
	education  []Education         `yaml:"education"`
	skills     map[string][]string `yaml:"skills"`
	projects   []Project           `yaml:"projects"`
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
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Technologies []string `yaml:"technologies"`
	Github       string   `yaml:"github"`
	Demo         string   `yaml:"demo"`
	Npm          string   `yaml:"npm"`
}

func (section Section) Education() ([]Education, bool) {
	_, ok := section.data[SectionTypeEducation]
	return section.education, ok
}

func (section Section) Experience() ([]Experience, bool) {
	_, ok := section.data[SectionTypeExperience]
	return section.experience, ok
}

func (section Section) Projects() ([]Project, bool) {
	_, ok := section.data[SectionTypeProjects]
	return section.projects, ok
}

func (section Section) Skills() (map[string][]string, bool) {
	_, ok := section.data[SectionTypeSkills]
	return section.skills, ok
}

// TODO: Optional stylesheets??
func (resume Resume) ToHTML(w io.Writer) error {
	return ResumePage(resume).Render(context.TODO(), w)
	// return htmlTemplate.Execute(w, resume)
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
	if err := yaml.NewDecoder(
		reader,
		yaml.CustomUnmarshaler[Section](func(target *Section, input []byte) error {
			target.data = map[string]any{}
			if err := yaml.Unmarshal(input, &target.data); err != nil {
				return fmt.Errorf(`failed to parse: %w`, err)
			}

			// Extract common fields
			if title, ok := target.data["title"].(string); ok {
				target.Title = title
				delete(target.data, "title")
			}
			// Extract known fields
			if experience, ok := target.data[SectionTypeExperience]; ok {
				if err := decode(experience, &target.experience); err != nil {
					return fmt.Errorf("parsing experience: %w", err)
				}
			}

			if education, ok := target.data[SectionTypeEducation]; ok {
				if err := decode(education, &target.education); err != nil {
					return fmt.Errorf("parsing education: %w", err)
				}
			}

			if skills, ok := target.data[SectionTypeSkills]; ok {
				if err := decode(skills, &target.skills); err != nil {
					return fmt.Errorf("parsing skills: %w", err)
				}
			}

			if projects, ok := target.data[SectionTypeProjects]; ok {
				if err := decode(projects, &target.projects); err != nil {
					return fmt.Errorf("parsing projects: %w", err)
				}
			}

			return nil
		}),
	).Decode(&resume); err != nil {
		return resume, fmt.Errorf("parsing yaml: %w", err)
	}
	return resume, nil
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

func FromJSON(reader io.Reader) (Resume, error) {
	return Resume{}, errors.New("not implemented")
}
