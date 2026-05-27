package database

import (
	"strings"
	"testing"
)

func TestLegacyTableRenames(t *testing.T) {
	renames := legacyTableRenames()
	if len(renames) != 8 {
		t.Fatalf("unexpected rename count: got %d want %d", len(renames), 8)
	}

	seenLegacy := make(map[string]struct{}, len(renames))
	seenCurrent := make(map[string]struct{}, len(renames))

	for _, rename := range renames {
		if rename.legacy == "" || rename.current == "" {
			t.Fatalf("rename entry must not be empty: %+v", rename)
		}
		if !strings.HasPrefix(rename.current, "epay_") {
			t.Fatalf("prefixed table name expected, got %s", rename.current)
		}
		if _, ok := seenLegacy[rename.legacy]; ok {
			t.Fatalf("duplicate legacy table rename entry: %s", rename.legacy)
		}
		if _, ok := seenCurrent[rename.current]; ok {
			t.Fatalf("duplicate current table rename entry: %s", rename.current)
		}

		seenLegacy[rename.legacy] = struct{}{}
		seenCurrent[rename.current] = struct{}{}
	}
}
