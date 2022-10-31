package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/italoservio/concurrent_vs_parallel/scripts"
)

const (
	CONCURRENT string = "concurrent"
	PARALLEL   string = "parallel"
)

func main() {
	start_time := time.Now()
	algorithm := flag.String(
		"algorithm",
		"concurrent",
		"[concurrent, parallel, sequential]",
	)
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dir := filepath.Join(pwd, "sheets")
	dir_entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var script func(
		dir string,
		dir_entries []fs.DirEntry,
	) (float32, error)

	switch *algorithm {
	case CONCURRENT:
		script = scripts.ConcurrentExecution
	case PARALLEL:
		script = scripts.ParallelExecution
	default:
		script = scripts.SequentialExecution
	}

	amount, err := script(dir, dir_entries)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Execution time: %s\n", time.Since(start_time))
	fmt.Printf("Amount: %.2f USD\n", amount)
}
