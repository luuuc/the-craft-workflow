package structure

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTest(t *testing.T) func() {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	return func() {
		os.Chdir(origDir)
	}
}

func TestPitchPath(t *testing.T) {
	want := filepath.Join(".craft", "pitch.md")
	if got := PitchPath(); got != want {
		t.Errorf("PitchPath() = %q, want %q", got, want)
	}
}

func TestCardsDirPath(t *testing.T) {
	want := filepath.Join(".craft", "cards")
	if got := CardsDirPath(); got != want {
		t.Errorf("CardsDirPath() = %q, want %q", got, want)
	}
}

func TestEnsureStructureDir(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	if err := EnsureStructureDir(); err != nil {
		t.Fatalf("EnsureStructureDir() error = %v", err)
	}

	info, err := os.Stat(CardsDirPath())
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if !info.IsDir() {
		t.Error("cards path should be a directory")
	}
}

func TestHasPitch(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// No pitch yet
	if HasPitch() {
		t.Error("HasPitch() = true, want false (no pitch)")
	}

	// Create pitch
	os.MkdirAll(CraftDir, 0755)
	os.WriteFile(PitchPath(), []byte("# Pitch"), 0644)

	if !HasPitch() {
		t.Error("HasPitch() = false, want true (pitch exists)")
	}
}

func TestHasCards(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// No cards yet
	if HasCards() {
		t.Error("HasCards() = true, want false (no cards)")
	}

	// Create cards directory with a card
	EnsureStructureDir()
	os.WriteFile(filepath.Join(CardsDirPath(), "01-first.md"), []byte("# Card"), 0644)

	if !HasCards() {
		t.Error("HasCards() = false, want true (card exists)")
	}
}

func TestListCards(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// No cards directory
	cards, err := ListCards()
	if err != nil {
		t.Fatalf("ListCards() error = %v", err)
	}
	if len(cards) != 0 {
		t.Errorf("ListCards() = %v, want empty", cards)
	}

	// Create cards
	EnsureStructureDir()
	os.WriteFile(filepath.Join(CardsDirPath(), "02-second.md"), []byte("# Card 2"), 0644)
	os.WriteFile(filepath.Join(CardsDirPath(), "01-first.md"), []byte("# Card 1"), 0644)
	os.WriteFile(filepath.Join(CardsDirPath(), "03-third.md"), []byte("# Card 3"), 0644)
	os.WriteFile(filepath.Join(CardsDirPath(), "readme.txt"), []byte("ignore"), 0644) // non-md file

	cards, err = ListCards()
	if err != nil {
		t.Fatalf("ListCards() error = %v", err)
	}

	want := []string{
		filepath.Join(CardsDirPath(), "01-first.md"),
		filepath.Join(CardsDirPath(), "02-second.md"),
		filepath.Join(CardsDirPath(), "03-third.md"),
	}

	if len(cards) != len(want) {
		t.Fatalf("ListCards() = %v, want %v", cards, want)
	}

	for i, c := range cards {
		if c != want[i] {
			t.Errorf("ListCards()[%d] = %q, want %q", i, c, want[i])
		}
	}
}

func TestListStructure(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Nothing exists
	pitch, cards, err := ListStructure()
	if err != nil {
		t.Fatalf("ListStructure() error = %v", err)
	}
	if pitch != "" || len(cards) != 0 {
		t.Errorf("ListStructure() = (%q, %v), want empty", pitch, cards)
	}

	// Create pitch and cards
	os.MkdirAll(CraftDir, 0755)
	os.WriteFile(PitchPath(), []byte("# Pitch"), 0644)
	EnsureStructureDir()
	os.WriteFile(filepath.Join(CardsDirPath(), "01-first.md"), []byte("# Card"), 0644)

	pitch, cards, err = ListStructure()
	if err != nil {
		t.Fatalf("ListStructure() error = %v", err)
	}

	if pitch != PitchPath() {
		t.Errorf("pitch = %q, want %q", pitch, PitchPath())
	}
	if len(cards) != 1 {
		t.Errorf("cards = %v, want 1 card", cards)
	}
}
