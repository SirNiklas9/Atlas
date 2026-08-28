// Command atlas exports bounded, provider-neutral structural evidence.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	pithadapter "github.com/SirNiklas9/atlas/adapters/pith"
	"github.com/SirNiklas9/atlas/engine"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "atlas:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("atlas", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	root := fs.String("root", ".", "repository root to map")
	output := fs.String("output", "", "snapshot path (default stdout)")
	snapshotPath := fs.String("snapshot", "", "Atlas snapshot path for thread")
	symbol := fs.String("symbol", "", "entry symbol or declaration ID for thread")
	cap := fs.Int("cap", 32, "maximum Thread nodes")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("usage: atlas [flags] snapshot|thread")
	}
	w := out
	var file *os.File
	var err error
	if *output != "" {
		if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil {
			return err
		}
		file, err = os.Create(*output)
		if err != nil {
			return err
		}
		defer file.Close()
		w = file
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	switch fs.Arg(0) {
	case "snapshot":
		s, err := pithadapter.Snapshot(*root)
		if err != nil {
			return err
		}
		return e.Encode(s)
	case "thread":
		if *snapshotPath == "" {
			return fmt.Errorf("--snapshot is required for thread")
		}
		data, err := os.ReadFile(*snapshotPath)
		if err != nil {
			return err
		}
		var snapshot engine.Snapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return err
		}
		thread, err := engine.BuildThread(snapshot, *symbol, *cap)
		if err != nil {
			return err
		}
		return e.Encode(thread)
	default:
		return fmt.Errorf("usage: atlas [flags] snapshot|thread")
	}
}
