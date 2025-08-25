package wizard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/ohhfishal/resume-wizard/resume"
	"log/slog"
	"strings"
	"text/template"
)

// var DefaultModel = anthropic.ModelClaude3_5Haiku20241022
var DefaultModel = anthropic.ModelClaudeOpus4_1_20250805

type TailorPromptArgs struct {
	JobDescription string
	ResumeJSON     string
	// TODO: Include the rest of the job description
}

var tailorTemplate = template.Must(template.New("tailor-tempalte").Parse(promptTemplate))

func (wizard *Wizard) annotateClaude(ctx context.Context, args AnnotationContext) (*resume.Resume, error) {
	content, err := promptFrom(args)
	if err != nil {
		return nil, fmt.Errorf("creating prompt: %w", err)
	}
	message, err := wizard.Claude.client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 2048,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model: DefaultModel,
	})
	if err != nil {
		return nil, fmt.Errorf("making query to model: %w", err)
	}
	// TODO: This array access is not safe and crashes if empty
	wizard.logger.Info("recieved response",
		slog.String("content", fmt.Sprintf("%s", message.Content[0].Text)),
	)

	input := strings.TrimPrefix(message.Content[0].Text, "```json")
	input = strings.TrimPrefix(input, "```")
	input = strings.TrimSpace(input)

	var response promptResponse
	if err = json.NewDecoder(strings.NewReader(input)).Decode(&response); err != nil {
		return nil, fmt.Errorf("model return invalid data: %w", err)
	}

	slog.Info("response from model parsed as json",
		slog.Any("tailored_resume", response.TailoredResume),
		slog.String("explanation", response.Explanation),
		slog.String("error", response.Error),
		slog.String("error_detailed", response.ErrorDetailed),
	)

	if response.Error != "" {
		// TODO: Using the global logger here
		slog.Error("model error", "error_detailed", response.ErrorDetailed, "error", response.Error)
		return nil, fmt.Errorf("model believes there is an error: %s", response.Error)
	}
	return &response.TailoredResume, nil
}

func promptFrom(ctx AnnotationContext) (string, error) {
	if ctx.Base.Resume == nil {
		return "", errors.New("missing field: Base.Resume")
	}

	resumeContent, err := ctx.Base.Resume.JSON()
	if err != nil {
		return "", fmt.Errorf("converting base resume to json: %w", err)
	}

	var writer strings.Builder
	if err := tailorTemplate.Execute(&writer, TailorPromptArgs{
		JobDescription: ctx.Description,
		ResumeJSON:     resumeContent,
	}); err != nil {
		return "", fmt.Errorf("templating: %w", err)
	}
	return writer.String(), nil
}

type promptResponse struct {
	TailoredResume resume.Resume `json:"tailored_resume"`
	Explanation    string        `json:"explanation"`
	Error          string        `json:"error"`
	ErrorDetailed  string        `json:"error_detailed"`
}

const promptTemplate = `
You are an AI assistant tasked with tailoring a resume to a specific job description. Your goal is to create a modified version of the provided JSON resume that highlights the most relevant skills, experiences, and qualifications for the given job.

First, carefully read and analyze the following job description:

<job_description>
{{.JobDescription}}
</job_description>

Now, examine the provided JSON resume:

<json_resume>
{{.ResumeJSON}}
</json_resume>

You may reword, add or remove bullet points as you see fit, but differ to the YAML and responses to questions you ask as the sources of truth.


Present only valid json.

Present your tailored version of the JSON resume using the "tailored_resume" field. Maintain proper JSON formatting and structure in your output.

After the tailored resume, provide a brief explanation of the main changes you made and why, enclosed in "explanation" field.

You also may use the "error" and "error_detailed" fields to describe errors that prevent you from proceeding. "error" will be returned as a golang error so keep the information brief but descriptive.

All your responses must be within a single JSON object.

Remember to focus on creating a targeted resume that highlights the candidate's most relevant qualifications for the specific job described.
`
