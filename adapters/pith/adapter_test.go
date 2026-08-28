package pithadapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotAdaptsPithWithoutSourceBodies(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package demo\nfunc B() {}\nfunc A() { B() }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if s.Provider != "atlas" || len(s.Symbols) != 2 || s.Symbols[0].ID == "" {
		t.Fatalf("snapshot=%#v", s)
	}
}

func TestMethodReceiversHaveDistinctStableIdentities(t *testing.T) {
	root := t.TempDir()
	body := []byte("package demo\ntype A struct{}\ntype B struct{}\nfunc (A) Scope() {}\nfunc (B) Scope() {}\n")
	if err := os.WriteFile(filepath.Join(root, "main.go"), body, 0o600); err != nil {
		t.Fatal(err)
	}
	s, err := Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, d := range s.Symbols {
		if seen[d.ID] {
			t.Fatalf("duplicate %q", d.ID)
		}
		seen[d.ID] = true
	}
}
