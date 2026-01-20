package structure

import (
	"os"
	"path/filepath"
	"sort"

	"craft/internal/workflow"
)

const (
	PitchFile = "pitch.md"
	CardsDir  = "cards"
)

// CraftDir is exported for convenience, sourced from workflow package.
var CraftDir = workflow.CraftDir

// PitchPath returns the path to the pitch file.
func PitchPath() string {
	return filepath.Join(CraftDir, PitchFile)
}

// CardsDirPath returns the path to the cards directory.
func CardsDirPath() string {
	return filepath.Join(CraftDir, CardsDir)
}

// EnsureStructureDir creates .craft/cards/ if it doesn't exist.
func EnsureStructureDir() error {
	return os.MkdirAll(CardsDirPath(), 0755)
}

// HasPitch returns true if .craft/pitch.md exists.
func HasPitch() bool {
	info, err := os.Stat(PitchPath())
	return err == nil && !info.IsDir()
}

// HasCards returns true if any cards exist in .craft/cards/.
func HasCards() bool {
	cards, err := ListCards()
	return err == nil && len(cards) > 0
}

// ListCards returns paths to all card files in .craft/cards/, sorted.
func ListCards() ([]string, error) {
	entries, err := os.ReadDir(CardsDirPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cards []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
			cards = append(cards, filepath.Join(CardsDirPath(), e.Name()))
		}
	}

	sort.Strings(cards)
	return cards, nil
}

// ListStructure returns the pitch path (if exists) and card paths.
func ListStructure() (pitch string, cards []string, err error) {
	if HasPitch() {
		pitch = PitchPath()
	}

	cards, err = ListCards()
	if err != nil {
		return "", nil, err
	}

	return pitch, cards, nil
}
