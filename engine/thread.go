package engine

import (
	"fmt"
	"sort"
	"strings"
)

// ThreadVersion identifies a bounded, source-backed structural path. A Thread
// is deliberately not a free-form behavioral summary: every node and edge
// comes from an Atlas snapshot and remains jumpable by its declaration.
const ThreadVersion = "atlas.thread.v1"

type Thread struct {
	Version   string       `json:"version"`
	Root      string       `json:"root"`
	Entry     Declaration  `json:"entry"`
	Nodes     []ThreadNode `json:"nodes"`
	Truncated bool         `json:"truncated,omitempty"`
}

type ThreadNode struct {
	Symbol  Declaration `json:"symbol"`
	Calls   []string    `json:"calls,omitempty"`
	Callers []string    `json:"callers,omitempty"`
}

// BuildThread resolves an exact symbol ID or an unambiguous symbol name, then
// extracts a breadth-first, bounded call path. It never reads source text,
// invents purposes, or treats incomplete graph evidence as execution truth.
func BuildThread(snapshot Snapshot, entry string, cap int) (Thread, error) {
	if snapshot.Version != SnapshotVersion {
		return Thread{}, fmt.Errorf("unsupported Atlas snapshot version %q", snapshot.Version)
	}
	entry = strings.TrimSpace(entry)
	if entry == "" || cap <= 0 {
		return Thread{}, fmt.Errorf("entry and positive cap are required")
	}
	byID := make(map[string]Declaration, len(snapshot.Symbols))
	var matches []Declaration
	for _, symbol := range snapshot.Symbols {
		byID[symbol.ID] = symbol
		if symbol.ID == entry || symbol.Name == entry {
			matches = append(matches, symbol)
		}
	}
	if len(matches) == 0 {
		return Thread{}, fmt.Errorf("entry %q is not in the snapshot", entry)
	}
	if len(matches) > 1 {
		return Thread{}, fmt.Errorf("entry %q is ambiguous; use a declaration ID", entry)
	}
	called, callers := map[string][]string{}, map[string][]string{}
	for _, evidence := range snapshot.Evidence {
		calls := append([]string(nil), evidence.Calls...)
		sort.Strings(calls)
		called[evidence.Symbol.ID] = calls
		for _, target := range calls {
			callers[target] = append(callers[target], evidence.Symbol.ID)
		}
	}
	result := Thread{Version: ThreadVersion, Root: snapshot.Root, Entry: matches[0]}
	seen := map[string]bool{}
	queue := []string{matches[0].ID}
	for len(queue) > 0 && len(result.Nodes) < cap {
		id := queue[0]
		queue = queue[1:]
		if seen[id] {
			continue
		}
		symbol, exists := byID[id]
		if !exists {
			continue
		}
		seen[id] = true
		node := ThreadNode{Symbol: symbol, Callers: append([]string(nil), callers[id]...)}
		sort.Strings(node.Callers)
		for _, next := range called[id] {
			if _, exists := byID[next]; !exists {
				continue
			}
			node.Calls = append(node.Calls, next)
			if !seen[next] {
				queue = append(queue, next)
			}
		}
		result.Nodes = append(result.Nodes, node)
	}
	result.Truncated = len(queue) > 0
	return result, nil
}
