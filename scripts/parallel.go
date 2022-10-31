package scripts

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/italoservio/concurrent_vs_parallel/structures"
)

func ParallelExecution(
	dir string,
	dir_entries []fs.DirEntry,
) (float32, error) {
	waiter := new(sync.WaitGroup)
	channel := make(chan structures.Output)

	output_dir, err := creatOutputDir()
	if err != nil {
		return 0, err
	}

	output_file, err := createOutputFile(output_dir, "parallel.csv")
	if err != nil {
		return 0, err
	}

	writeHeaderRow(output_file)

	for _, dir_entry := range dir_entries {
		filename := filepath.Join(dir, dir_entry.Name())

		waiter.Add(1)
		fmt.Println("Green Thread started")
		go sumFileRowsAndWriteFileInThread( // Green Thread: 2kB
			waiter,
			&channel,
			filename,
			output_file,
		)
	}

	go waitToCloseChannel(waiter, &channel) // Thread to lock and wait for other threads

	var total_amount float32 = 0
	for output := range channel {
		fmt.Println("Thread finished")
		total_amount += output.Amount

		fmt.Println("Writing result line...")
		err = writeValueRow(output_file, output.Entry, output.Amount)
		if err != nil {
			return 0, err
		}
	}

	err = writeValueRow(output_file, "total", total_amount)
	if err != nil {
		return 0, err
	}

	return total_amount, nil
}

func sumFileRowsAndWriteFileInThread(
	waiter *sync.WaitGroup,
	channel *chan structures.Output,
	filename string,
	output_file *os.File,
) {
	defer waiter.Done()

	file_amount, err := calculateFileAmount(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	base_filename := filepath.Base(filename)
	fmt.Printf("EOF: %s\n", base_filename)

	*channel <- *structures.NewOutput(base_filename, file_amount)
}

func waitToCloseChannel(waiter *sync.WaitGroup, channel *chan structures.Output) {
	waiter.Wait()
	close(*channel)
}
