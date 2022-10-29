package scripts

import (
	"fmt"
	"io/fs"
	"os"
)

func ConcurrentExecution(dir string, dir_entries []fs.DirEntry) (float32, error) {
	var total_amount float32 = 0.0

	for _, dir_entry := range dir_entries {
		filename := fmt.Sprintf("%s/%s", dir, dir_entry.Name())
		sumFileRowsConcurrently(filename, &total_amount)
	}

	return total_amount, nil
}

func sumFileRowsConcurrently(filename string, total_amount *float32) {
	file_amount, err := calculateFileAmount(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	*total_amount += float32(file_amount)

	fmt.Printf("EOF: %s\n", filename)
}
