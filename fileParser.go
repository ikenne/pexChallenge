package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	delimiter = '\n'
)

// ensure interface compliance
// var _ parser.PayloadParser = (*FileParser)(nil)

type fileParser struct {
}

func newFileParser() *fileParser {
	return new(fileParser)
}

// readInputFile parses the input file
func (fp *fileParser) readInputFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	lines := make([]string, 0)

	var line string
	for {
		line, err = reader.ReadString(delimiter)
		if line != "" && line != string(delimiter) {
			x := strings.TrimSpace(strings.TrimSuffix(line, string(delimiter)))

			lines = append(lines, x)

		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf("Error loading input: %v\n", err)
		return nil, err
	}

	fmt.Println("number of lines parsed ", len(lines))
	return lines, nil
}
func (fp *fileParser) partitionInputFile(inputFile string, maxLines int) ([]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	partitionFiles := make([]string, 0)

	ext := filepath.Ext(inputFile)
	partitionFilePrefix := strings.TrimSuffix(inputFile, ext) + "_p"

	partitionCnt := 1
	partitionFile := partitionFilePrefix + fmt.Sprintf("%v", partitionCnt) + ext
	pFile, err := os.Create(partitionFile)
	if err != nil {
		return nil, err
	}

	partitionFiles = append(partitionFiles, partitionFile)

	lineCnt := 0
	lines := make([]string, 0)
	var line string

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	for {
		line, err = reader.ReadString(delimiter)
		if line != "" && line != string(delimiter) {
			x := strings.TrimSpace(strings.TrimSuffix(line, string(delimiter)))
			lineCnt++
			if lineCnt > 1 && (lineCnt%maxLines) == 1 {
				// write lines to partition file
				for _, l := range lines {
					l = fmt.Sprintf("%s\n", l)
					if _, err := pFile.WriteString(l); err != nil {
						return nil, err
					}
				}
				// close the partition file
				pFile.Close()

				// new partition file
				partitionCnt++
				partitionFile = partitionFilePrefix + fmt.Sprintf("%v", partitionCnt) + ext
				pFile, err = os.Create(partitionFile)
				if err != nil {
					return nil, err
				}
				partitionFiles = append(partitionFiles, partitionFile)
				lines = make([]string, 0)
			}

			lines = append(lines, x)
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf("Error loading input: %v\n", err)
		return nil, err
	}

	if len(lines) > 0 && pFile != nil {
		for _, l := range lines {
			l = fmt.Sprintf("%s\n", l)
			if _, err := pFile.WriteString(l); err != nil {
				return nil, err
			}
		}

		pFile.Close()
	}

	fmt.Println("number of lines parsed ", lineCnt)
	return partitionFiles, nil
}
