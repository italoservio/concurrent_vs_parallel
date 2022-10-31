package scripts

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func SequentialExecution(
	dir string,
	dir_entries []fs.DirEntry,
) (float32, error) {
	var total_amount float32 = 0.0

	output_dir, err := creatOutputDir()
	if err != nil {
		return 0, err
	}

	output_file, err := createOutputFile(output_dir, "sequential.csv")
	if err != nil {
		return 0, err
	}

	writeHeaderRow(output_file)

	for _, dir_entry := range dir_entries {
		filename := filepath.Join(dir, dir_entry.Name())
		file_amount := sumFileRows(filename)
		base_filename := filepath.Base(filename)

		writeValueRow(output_file, base_filename, file_amount)
		if err != nil {
			return 0, err
		}

		total_amount += file_amount

		fmt.Printf("EOF: %s\n", base_filename)
	}

	writeValueRow(output_file, "total", total_amount)
	if err != nil {
		return 0, err
	}

	return total_amount, nil
}

func sumFileRows(
	filename string,
) float32 {
	file_amount, err := calculateFileAmount(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return float32(file_amount)
}
