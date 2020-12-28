package main

import (
	"fmt"
	"image"
	"net/http"

	"github.com/EdlinOrg/prominentcolor"
)

type processor struct {
}

// processURLs processes the input URLs
func (p processor) processURLs(urls []string) ([]imageResult, error) {
	result := make([]imageResult, 0)
	for _, url := range urls {
		r, err := p.process(url)
		if err != nil {
			// return nil, err
		}
		result = append(result, *r)
	}

	return result, nil
}

// process - gets the prominent colours from the URL
func (p processor) process(url string) (*imageResult, error) {

	r := new(imageResult)
	r.URL = url

	resp, err := http.Get(url)
	if err != nil {
		r.errMsg = err.Error()
		return r, err
	}
	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Printf("Failed to process image, url %s "+
			"format %s: %v\n", url, format, err)
		r.errMsg = err.Error()
		return r, err
	}

	colours, err := prominentcolor.Kmeans(img)
	if err != nil {
		fmt.Println("Failed to process image", err)
		r.errMsg = err.Error()
		return r, err
	}

	// fmt.Println("Dominant colours:")
	// r.Colors = make([]string, 0)
	// for _, c := range colours {
	// 	fmt.Println("#" + c.AsString())
	// 	r.Colors = append(r.Colors, "#"+c.AsString())
	// }

	r.Colors = [3]rgb{}
	for i := 0; i < 3; i++ {
		c := colours[i]
		r.Colors[i][0] = byte(c.Color.R)
		r.Colors[i][1] = byte(c.Color.G)
		r.Colors[i][2] = byte(c.Color.B)
	}

	return r, nil
}
