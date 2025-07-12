package worker

import (
	"fmt"
	"sync"
	"ts-copy/internal/transfer"
)

// ProcessFiles processes a list of files using a worker pool with the specified
// number of workers. Each file is copied to the target machine using Tailscale.
// Returns the number of errors encountered.
func ProcessFiles(files []string, targetMachine string, dryRun bool, maxWorkers int) int {
	jobs := make(chan string, len(files))
	errors := make(chan int, len(files))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, dryRun, targetMachine, errors)
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
	close(errors)

	// Count errors
	errorCount := 0
	for range errors {
		errorCount++
	}

	return errorCount
}

// worker processes files from the jobs channel
func worker(jobs <-chan string, wg *sync.WaitGroup, dryRun bool, targetMachine string, errors chan<- int) {
	defer wg.Done()
	for file := range jobs {
		if err := transfer.CopyFile(file, targetMachine, dryRun); err != nil {
			fmt.Printf("Error: %v\n", err)
			errors <- 1
			continue
		}
	}
}