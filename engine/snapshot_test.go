package engine

import "testing"

func TestBuildNormalizesAndBoundsStructuralEvidence(t *testing.T) {
	s, err := Build("fixture", "C:/repo", []Declaration{{Name: "A", Kind: "func", File: "a.go", Line: 5}, {Name: "B", Kind: "func", File: "a.go", Line: 9}}, []Relation{{From: "a.go:func:A", To: "a.go:func:B"}, {From: "a.go:func:A", To: "outside"}}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if s.Version != SnapshotVersion || len(s.Symbols) != 2 || len(s.Evidence) != 1 || len(s.Evidence[0].Calls) != 1 {
		t.Fatalf("snapshot = %#v", s)
	}
	if s.Symbols[0].ID != "a.go:func:A" {
		t.Fatalf("stable ID = %q", s.Symbols[0].ID)
	}
}

func TestBuildRejectsIdentityConflict(t *testing.T) {
	_, err := Build("p", "r", []Declaration{{ID: "x", Name: "A", Kind: "func", File: "a.go", Line: 1}, {ID: "x", Name: "A", Kind: "func", File: "a.go", Line: 2}}, nil, 10)
	if err == nil {
		t.Fatal("expected identity conflict")
	}
}

func TestBuildThreadFollowsBoundedStructuralCalls(t *testing.T) {
	s, err := Build("fixture", "C:/repo", []Declaration{
		{Name: "Apply", Kind: "func", File: "links.go", Line: 1},
		{Name: "Prepare", Kind: "func", File: "links.go", Line: 10},
		{Name: "Verify", Kind: "func", File: "links.go", Line: 20},
	}, []Relation{
		{From: "links.go:func:Apply", To: "links.go:func:Prepare"},
		{From: "links.go:func:Prepare", To: "links.go:func:Verify"},
	}, 10)
	if err != nil {
		t.Fatal(err)
	}
	thread, err := BuildThread(s, "Apply", 2)
	if err != nil {
		t.Fatal(err)
	}
	if thread.Version != ThreadVersion || len(thread.Nodes) != 2 || thread.Nodes[0].Symbol.Name != "Apply" || thread.Nodes[1].Symbol.Name != "Prepare" {
		t.Fatalf("thread=%#v", thread)
	}
	if !thread.Truncated || len(thread.Nodes[1].Callers) != 1 || thread.Nodes[1].Callers[0] != "links.go:func:Apply" {
		t.Fatalf("thread=%#v", thread)
	}
}

func TestBuildThreadRejectsAmbiguousName(t *testing.T) {
	s, err := Build("fixture", "C:/repo", []Declaration{{Name: "Apply", Kind: "func", File: "a.go", Line: 1}, {Name: "Apply", Kind: "func", File: "b.go", Line: 1}}, nil, 10)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BuildThread(s, "Apply", 10); err == nil {
		t.Fatal("expected ambiguous name")
	}
}
