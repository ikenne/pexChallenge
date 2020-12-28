package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

const (
	outputFile           = "myOutput.txt"
	defaultPartitionSize = 100
	defaultNumOfWorkers  = 2
)

var (
	filePath      string
	partitionSize int
	numOfWorkers  int
	results       chan []imageResult
	jobs          chan string
)

func main() {

	flag.StringVar(&filePath, "inputFile", "input.txt", "full path to input file")
	flag.IntVar(&partitionSize, "pSize", defaultPartitionSize, "max lines in a partition")
	flag.IntVar(&numOfWorkers, "workers", defaultNumOfWorkers, "number of workers")
	flag.Parse()

	fp := newFileParser()
	/*
		urls, err := fp.readInputFile(filePath)
		if err != nil {
			fmt.Printf("Error loading input %v", err)
			return
		}

		ip := new(processor)
		ir, err := ip.processURLs(urls)
		if err != nil {
			fmt.Printf("error processing image urls %v", err)
		}
		writeToFile(ir)
	*/

	fmt.Println("partition size", partitionSize)
	partitions, err := fp.partitionInputFile(filePath, partitionSize)
	if err != nil {
		fmt.Printf("Error partitioning input %v", err)
		return
	}

	ip := new(processor)

	/*
		// without goroutines

		var irs []imageResult
		for _, p := range partitions {
			fmt.Println("processing partition file:", p)
			urls, err := fp.readInputFile(p)
			if err != nil {
				fmt.Printf("Error loading partition %s: %v", p, err)
				return
			}

			ir, err := ip.processURLs(urls)
			if err != nil {
				fmt.Printf("error processing image urls %v", err)
				return
			}

			irs = append(irs, ir...)
		}
		writeToFile(irs)
	*/

	results = make(chan []imageResult, len(partitions))
	jobs = make(chan string, len(partitions))
	done := make(chan bool)

	//create output file
	_, err = os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file %s:%v", outputFile, err)
		return
	}

	go writeResults(done)

	//allocate jobs
	go allocateJobs(partitions)

	// create workers
	createWorkers(numOfWorkers, fp, ip)

	<-done
	fmt.Println("Finished processing")
}

// writeToFile outputs the result to a file
func writeToFile(ir []imageResult) error {
	// file, err := os.Create(outputFile)
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, r := range ir {
		if _, err := file.WriteString(r.String()); err != nil {
			return err
		}
	}

	return nil
}

// write results to the output file from results channel
func writeResults(done chan bool) {
	for r := range results {
		writeToFile(r)
	}

	done <- true
}

// adds a job (partition file) to the buffered jobs channel
func allocateJobs(partitions []string) {
	for _, p := range partitions {
		jobs <- p
	}
	close(jobs)
}

// create worker pool - goroutines to process the partitions
func createWorkers(num int, fp *fileParser, ip *processor) {
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go worker(&wg, fp, ip)
	}
	wg.Wait()

	close(results)
}

// the worker - processes the partition (job) and puts output in results channel
func worker(wg *sync.WaitGroup, fp *fileParser, ip *processor) {
	for job := range jobs {
		urls, err := fp.readInputFile(job)
		if err != nil {
			fmt.Printf("Error loading partition %s: %v", job, err)
			return
		}

		ir, err := ip.processURLs(urls)
		if err != nil {
			fmt.Printf("error processing image urls %v", err)
			return
		}

		results <- ir
	}

	wg.Done()
}
