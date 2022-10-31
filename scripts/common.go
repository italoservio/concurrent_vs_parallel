package scripts

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func creatOutputDir() (string, error) {
	fmt.Println("Verifying if output directory exists...")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	output_dir := filepath.Join(pwd, "output")

	if _, err := os.Stat(output_dir); os.IsNotExist(err) {
		fmt.Println("Creating output directory...")
		if err := os.Mkdir(output_dir, os.ModePerm); err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	return output_dir, nil
}

func createOutputFile(output_dir string, name string) (*os.File, error) {
	fmt.Println("Creating output file...")
	output_file, err := os.Create(
		filepath.Join(
			output_dir,
			name,
		),
	)
	if err != nil {
		return &os.File{}, err
	}

	return output_file, nil
}

func writeHeaderRow(output_file *os.File) error {
	fmt.Println("Writing file header...")
	_, err := output_file.WriteString("entry;amount")
	if err != nil {
		return err
	}
	return nil
}

func writeValueRow(output_file *os.File, entry string, amount float32) error {
	_, err := output_file.WriteString(fmt.Sprintf("\n%s;%v", entry, amount))
	if err != nil {
		return err
	}
	return nil
}
