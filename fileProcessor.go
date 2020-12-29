package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pexChallenge/model"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
)

const (
	delimiter = '\n'
)

// var _ UrlProcessor = (*fileParser)(nil)
var _ UrlProcessor = (*processor)(nil)

type fileProcessor struct {
	UrlProcessor
}

func newFileProcessor() *fileProcessor {
	return &fileProcessor{
		UrlProcessor: new(processor),
	}
}

// readInputFile parses the input file
func (fp *fileProcessor) readInputFile(filePath string) ([]string, error) {
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

func (fp *fileProcessor) partitionInputFile(inputFile string, maxLines int) ([]string, error) {
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

	fmt.Println("total number of lines parsed ", lineCnt)
	return partitionFiles, nil
}

type processor struct {
}

// processURLs processes the input URLs
func (p processor) ProcessURLs(urls []string) ([]model.ImageResult, error) {
	result := make([]model.ImageResult, 0)
	for _, url := range urls {
		r, err := p.Process(url)
		if err != nil {
			// return nil, err
		}
		result = append(result, *r)
	}

	return result, nil
}

// process - gets the prominent colours from the URL
func (p processor) Process(url string) (*model.ImageResult, error) {

	r := new(model.ImageResult)
	r.URL = url

	resp, err := http.Get(url)
	if err != nil {
		r.ErrMsg = err.Error()
		return r, err
	}
	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Printf("Failed to process image, url %s "+
			"format %s: %v\n", url, format, err)
		r.ErrMsg = err.Error()
		return r, err
	}

	colours, err := prominentcolor.Kmeans(img)
	if err != nil {
		fmt.Println("Failed to process image", err)
		r.ErrMsg = err.Error()
		return r, err
	}

	// fmt.Println("Dominant colours:")
	// r.Colors = make([]string, 0)
	// for _, c := range colours {
	// 	fmt.Println("#" + c.AsString())
	// 	r.Colors = append(r.Colors, "#"+c.AsString())
	// }

	r.Colors = [3]model.RGB{}
	for i := 0; i < 3; i++ {
		c := colours[i]
		r.Colors[i][0] = byte(c.Color.R)
		r.Colors[i][1] = byte(c.Color.G)
		r.Colors[i][2] = byte(c.Color.B)
	}

	return r, nil
}
