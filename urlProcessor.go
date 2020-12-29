package main

import (
	"pexChallenge/model"
)

type UrlProcessor interface {
	ProcessURLs(urls []string) ([]model.ImageResult, error)
	Process(url string) (*model.ImageResult, error)
}
