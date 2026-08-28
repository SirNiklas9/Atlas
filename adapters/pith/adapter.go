// Package pithadapter adapts the pre-existing Pith tree-sitter parser to the
// neutral Atlas engine. It is intentionally the only place Atlas mentions
// Pith, making replacement by an LSP, compiler index, or another parser a
// provider swap rather than a Lexicon change.
package pithadapter

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SirNiklas9/atlas/engine"
	"pith"
)

const DefaultCap = 2048

func Snapshot(root string) (engine.Snapshot, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return engine.Snapshot{}, err
	}
	packages, err := pith.Map(abs)
	if err != nil {
		return engine.Snapshot{}, err
	}
	entries := make([]pith.Entry, 0)
	for _, pkg := range packages {
		r, err := pith.Gather(filepath.Join(abs, filepath.FromSlash(pkg.Dir)), "")
		if err == nil {
			entries = append(entries, r.All...)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].File != entries[j].File {
			return entries[i].File < entries[j].File
		}
		return entries[i].Line < entries[j].Line
	})
	decls := make([]engine.Declaration, 0, len(entries))
	byLocation := map[string]string{}
	for _, entry := range entries {
		d := declaration(abs, entry)
		decls = append(decls, d)
		byLocation[location(entry.File, entry.Line)] = d.ID
	}
	relations := make([]engine.Relation, 0)
	for _, entry := range entries {
		from := byLocation[location(entry.File, entry.Line)]
		end := entry.Line + strings.Count(entry.Source, "\n")
		used, _, err := pith.UsesClosure(entry.File, "uses:project", entry.Line, end)
		if err != nil {
			continue
		}
		for _, u := range used {
			if to, ok := byLocation[location(u.File, u.Line)]; ok && to != from {
				relations = append(relations, engine.Relation{From: from, To: to})
			}
		}
	}
	return engine.Build("atlas", abs, decls, relations, DefaultCap)
}

func declaration(root string, entry pith.Entry) engine.Declaration {
	rel, err := filepath.Rel(root, entry.File)
	if err != nil {
		rel = entry.File
	}
	file := filepath.ToSlash(rel)
	identity := entry.Name
	// Methods can legitimately share a name within one file. Their receiver is
	// structural identity, unlike the mutable line number, so preserve it in
	// the opaque ID while keeping Name friendly for consumers.
	if entry.Kind == "method" && entry.Recv != "" {
		identity = entry.Recv + "." + entry.Name
	}
	return engine.Declaration{ID: fmt.Sprintf("%s:%s:%s", file, entry.Kind, identity), Name: entry.Name, Kind: entry.Kind, File: file, Line: entry.Line}
}
func location(file string, line int) string { return filepath.Clean(file) + fmt.Sprintf(":%d", line) }
