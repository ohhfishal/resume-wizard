package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type WizardCmd struct {
	File *os.File `arg:"" help:"File to read"`
}

func (cmd *WizardCmd) Run(logger *slog.Logger) error {
	defer cmd.File.Close()
	metadata, err := parse(cmd.File)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	logger.Info("done",
		slog.Any("mapping", metadata.Frequencies),
		slog.Any("words", metadata.Words),
	)

	return nil
}

type Metadata struct {
	Frequencies map[string]int
	Words       []string
}

func (metadata *Metadata) Add(word string) {
	metadata.Words = append(metadata.Words, word)

	normalized := normalize(word)
	if len(normalized) == 0 || slices.Contains(stopWords, normalized) {
		return
	}
	if normalized == "we" {
		panic("WE?")
	}

	if _, ok := metadata.Frequencies[normalized]; !ok {
		metadata.Frequencies[normalized] = 0
	}
	metadata.Frequencies[normalized] += 1
}

func normalize(word string) string {
	return strings.TrimSpace(strings.ToLower(word))
}

func parse(reader io.Reader) (Metadata, error) {
	var word strings.Builder
	queue := bufio.NewReader(reader)
	metadata := Metadata{
		Frequencies: map[string]int{},
	}

	for {
		char, _, err := queue.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return metadata, fmt.Errorf("reading: %w", err)
		}

		switch char {
		case '/', ' ', '.', '\n', ',', '!', '?', '(', ')':
			metadata.Add(word.String())
			word.Reset()
		default:
			word.WriteRune(char)
		}
	}
	return metadata, nil
}
