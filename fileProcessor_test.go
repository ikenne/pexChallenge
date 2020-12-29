package main

import (
	"pexChallenge/mocks"
	"pexChallenge/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testInputFile = "fileParser_test.txt"
	maxLines      = 5
)

var inputLines = []string{
	"http://i.imgur.com/FApqk3D.jpg",
	"http://i.imgur.com/TKLs9lo.jpg",
	"https://i.redd.it/d8021b5i2moy.jpg",
	"https://i.redd.it/4m5yk8gjrtzy.jpg",
	"https://i.redd.it/xae65ypfqycy.jpg",
	"http://i.imgur.com/lcEUZHv.jpg",
	"https://i.redd.it/1nlgrn49x7ry.jpg",
	"http://i.imgur.com/M3NOzLC.jpg",
	"https://i.redd.it/w5q6gldnvcuy.jpg",
	"https://i.redd.it/s5viyluv421z.jpg",
}

var partitionFiles = []string{
	"fileParser_test_p1.txt",
	"fileParser_test_p2.txt",
}

func TestFileRead(t *testing.T) {
	fp := newFileProcessor()
	r, err := fp.readInputFile(testInputFile)
	assert.Nil(t, err)
	assert.Equal(t, inputLines, r)
}

func TestInputPartition(t *testing.T) {
	fp := newFileProcessor()
	r, err := fp.partitionInputFile(testInputFile, maxLines)
	assert.Nil(t, err)
	assert.Equal(t, partitionFiles, r)

	// remove temp partition files
	cleanUpPartitionFiles(partitionFiles)
}

func TestProcess(t *testing.T) {
	url := "http://i.imgur.com/FApqk3D.jpg"
	ir := model.ImageResult{
		URL: url,
		Colors: [3]model.RGB{
			{0, 0, 255},
			{255, 0, 0},
			{0, 0, 255},
		},
	}

	p := new(mocks.UrlProcessor)
	p.On("Process", url).Return(
		&ir, nil)

	fp := new(fileProcessor)
	fp.UrlProcessor = p
	r, err := fp.Process(url)

	assert.Nil(t, err)
	assert.Equal(t, ir, *r)
}
