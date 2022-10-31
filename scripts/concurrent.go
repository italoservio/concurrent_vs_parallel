package scripts

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/italoservio/concurrent_vs_parallel/structures"
)

type ThreadResponse struct {
	total float32
	rows  []structures.Output
}

func ConcurrentExecution(
	dir string,
	dir_entries []fs.DirEntry,
) (float32, error) {
	waiter := new(sync.WaitGroup)
	channel := make(chan ThreadResponse)

	waiter.Add(1)

	fmt.Println("Green Thread started")
	go calculateAmountInDirectoryEntries(
		waiter,
		&channel,
		dir,
		dir_entries,
	) // Green Thread: 2kB

	output_dir, err := creatOutputDir()
	if err != nil {
		return 0, err
	}

	output_file, err := createOutputFile(output_dir, "concurrent.csv")
	if err != nil {
		return 0, err
	}

	writeHeaderRow(output_file)

	thread_response := <-channel
	waiter.Wait() // Waiter after the channel was empty

	fmt.Println("Writing result lines...")
	for _, output := range thread_response.rows {
		err := writeValueRow(output_file, output.Entry, output.Amount)
		if err != nil {
			return 0, err
		}
	}

	err = writeValueRow(output_file, "total", thread_response.total)
	if err != nil {
		return 0, err
	}

	return thread_response.total, nil
}

func calculateAmountInDirectoryEntries(
	waiter *sync.WaitGroup,
	channel *chan ThreadResponse,
	dir string,
	dir_entries []fs.DirEntry,
) {
	defer close(*channel)
	defer waiter.Done()

	rows := []structures.Output{}
	var total_amount float32 = 0.0

	for _, dir_entry := range dir_entries {
		filename := filepath.Join(dir, dir_entry.Name())
		base_filename := filepath.Base(filename)
		file_amount := sumFileRows(filename)

		rows = append(rows, *structures.NewOutput(
			filepath.Base(base_filename),
			file_amount,
		))

		total_amount += file_amount
		fmt.Printf("EOF: %s\n", base_filename)
	}

	fmt.Println("Finished files reading")
	*channel <- ThreadResponse{total: total_amount, rows: rows}
}
