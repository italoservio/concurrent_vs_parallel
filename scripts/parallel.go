package scripts

import (
	"fmt"
	"io/fs"
	"os"
	"sync"
)

func ParallelExecution(dir string, dir_entries []fs.DirEntry) (float32, error) {
	waiter := new(sync.WaitGroup)
	channel := make(chan float32)

	for _, dir_entry := range dir_entries {
		filename := fmt.Sprintf("%s/%s", dir, dir_entry.Name())

		waiter.Add(1)
		go sumFileRowsInThread(waiter, &channel, filename) // 2kB
	}

	go waitToCloseChannel(waiter, &channel)

	var total_amount float32 = 0
	for file_amount := range channel {
		total_amount += file_amount
	}

	return total_amount, nil
}

func sumFileRowsInThread(waiter *sync.WaitGroup, channel *chan float32, filename string) {
	defer waiter.Done()
	file_amount, err := calculateFileAmount(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	*channel <- file_amount

	fmt.Printf("EOF: %s\n", filename)
}

func waitToCloseChannel(waiter *sync.WaitGroup, channel *chan float32) {
	waiter.Wait()
	close(*channel)
}
