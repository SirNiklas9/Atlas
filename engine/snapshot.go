// Package engine defines Atlas's provider-neutral structural-evidence core.
// Parsers, language servers, and repository indexers adapt into Declaration
// and Relation; consumers receive a compact, deterministic snapshot with no
// source bodies or parser-specific state.
package engine

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const SnapshotVersion = "atlas.snapshot.v1"

type Declaration struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	File string `json:"file"`
	Line int    `json:"line"`
}

type Relation struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Evidence struct {
	Symbol Declaration `json:"symbol"`
	Calls  []string    `json:"calls,omitempty"`
}

type Snapshot struct {
	Provider string        `json:"provider"`
	Version  string        `json:"version"`
	Root     string        `json:"root"`
	Symbols  []Declaration `json:"symbols"`
	Evidence []Evidence    `json:"evidence,omitempty"`
}

// Build is the engine boundary shared by every Atlas provider.  IDs are
// stable across a declaration's line moves; only name, kind, or file changes
// mint a new identity. Relations are restricted to included symbols, so an
// indexer cannot leak an unbounded graph into context.
func Build(provider, root string, declarations []Declaration, relations []Relation, cap int) (Snapshot, error) {
	provider, root = strings.TrimSpace(provider), strings.TrimSpace(root)
	if provider == "" || root == "" {
		return Snapshot{}, errors.New("provider and root are required")
	}
	if cap <= 0 {
		return Snapshot{}, errors.New("symbol cap must be positive")
	}
	byID := map[string]Declaration{}
	for _, d := range declarations {
		d.Name, d.Kind, d.File = strings.TrimSpace(d.Name), strings.TrimSpace(d.Kind), strings.TrimSpace(d.File)
		if d.Name == "" || d.Kind == "" || d.File == "" || d.Line <= 0 {
			return Snapshot{}, fmt.Errorf("invalid declaration %q", d.Name)
		}
		d.File = filepath.ToSlash(filepath.Clean(d.File))
		if d.ID == "" {
			d.ID = d.File + ":" + d.Kind + ":" + d.Name
		}
		if prior, exists := byID[d.ID]; exists && prior != d {
			return Snapshot{}, fmt.Errorf("conflicting declaration identity %q", d.ID)
		}
		byID[d.ID] = d
	}
	symbols := make([]Declaration, 0, len(byID))
	for _, d := range byID {
		symbols = append(symbols, d)
	}
	sort.Slice(symbols, func(i, j int) bool {
		if symbols[i].File != symbols[j].File {
			return symbols[i].File < symbols[j].File
		}
		if symbols[i].Line != symbols[j].Line {
			return symbols[i].Line < symbols[j].Line
		}
		return symbols[i].ID < symbols[j].ID
	})
	if len(symbols) > cap {
		symbols = symbols[:cap]
	}
	allowed := map[string]bool{}
	for _, d := range symbols {
		allowed[d.ID] = true
	}
	calls := map[string]map[string]bool{}
	for _, r := range relations {
		if r.From == r.To || !allowed[r.From] || !allowed[r.To] {
			continue
		}
		if calls[r.From] == nil {
			calls[r.From] = map[string]bool{}
		}
		calls[r.From][r.To] = true
	}
	evidence := make([]Evidence, 0, len(calls))
	for _, d := range symbols {
		if len(calls[d.ID]) == 0 {
			continue
		}
		row := Evidence{Symbol: d}
		for to := range calls[d.ID] {
			row.Calls = append(row.Calls, to)
		}
		sort.Strings(row.Calls)
		evidence = append(evidence, row)
	}
	return Snapshot{Provider: provider, Version: SnapshotVersion, Root: root, Symbols: symbols, Evidence: evidence}, nil
}
