package shaper

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"craft/internal/structure"
)

// ShapeCLIShaper invokes the shape-cli tool for structure generation.
type ShapeCLIShaper struct{}

func (s *ShapeCLIShaper) Name() string {
	return NameShapeCLI
}

func (s *ShapeCLIShaper) Available() bool {
	_, err := exec.LookPath("shape")
	return err == nil
}

func (s *ShapeCLIShaper) Shape(req ShapeRequest) (ShapeResult, error) {
	// Ensure structure directory exists
	if err := structure.EnsureStructureDir(); err != nil {
		return ShapeResult{}, fmt.Errorf("failed to create structure dir: %w", err)
	}

	// Build command arguments
	args := []string{
		"generate",
		"--intent", req.Intent,
		"--output", structure.CraftDir,
	}

	// Add notes if present
	if len(req.Notes) > 0 {
		notes := strings.Join(req.Notes, "; ")
		args = append(args, "--notes", notes)
	}

	cmd := exec.Command("shape", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return ShapeResult{}, fmt.Errorf("shape-cli failed: %s", errMsg)
	}

	// Parse output for created files
	result := ShapeResult{
		Shaper: NameShapeCLI,
	}

	// Look for created file paths in output
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasSuffix(line, ".md") {
			// Determine if it's a pitch or card
			if strings.Contains(line, "pitch") {
				result.PitchPath = line
			} else if strings.Contains(line, "cards/") || strings.Contains(line, "cards\\") {
				result.CardPaths = append(result.CardPaths, line)
			}
		}
	}

	// If output parsing didn't find files, check the filesystem
	if result.PitchPath == "" && structure.HasPitch() {
		result.PitchPath = structure.PitchPath()
	}

	if len(result.CardPaths) == 0 {
		cards, _ := structure.ListCards()
		result.CardPaths = cards
	}

	// Normalize paths
	if result.PitchPath != "" {
		result.PitchPath = filepath.Clean(result.PitchPath)
	}
	for i, p := range result.CardPaths {
		result.CardPaths[i] = filepath.Clean(p)
	}

	return result, nil
}
