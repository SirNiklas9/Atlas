package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/SirNiklas9/atlas/engine"
)

func TestAtlasCommandExportsSnapshot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte("package demo\nfunc A() {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"--root", root, "snapshot"}, &out); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out.Bytes(), []byte(`"atlas.snapshot.v1"`)) {
		t.Fatalf("output=%s", out.String())
	}
}

func TestAtlasCommandBuildsThreadFromSnapshot(t *testing.T) {
	s, err := engine.Build("fixture", "C:/repo", []engine.Declaration{{Name: "Apply", Kind: "func", File: "a.go", Line: 1}, {Name: "Verify", Kind: "func", File: "a.go", Line: 2}}, []engine.Relation{{From: "a.go:func:Apply", To: "a.go:func:Verify"}}, 10)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "atlas.json")
	data, _ := json.Marshal(s)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"--snapshot", path, "--symbol", "Apply", "thread"}, &out); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out.Bytes(), []byte(`"atlas.thread.v1"`)) {
		t.Fatalf("output=%s", out.String())
	}
}
