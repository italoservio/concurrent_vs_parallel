package scripts

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

func calculateFileAmount(filename string) (float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	var file_amount float32 = 0.0
	for {
		content, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		if line_amount, err := strconv.ParseFloat(content[1], 32); err == nil {
			file_amount += float32(line_amount)
		}
	}

	return file_amount, nil
}
