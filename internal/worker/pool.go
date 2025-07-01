package worker

import (
	"sync"
	"ts-copy/internal/transfer"
)

// ProcessFiles processes a list of files using a worker pool with the specified
// number of workers. Each file is copied to the target machine using Tailscale.
func ProcessFiles(files []string, targetMachine string, dryRun bool, maxWorkers int) {
	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, dryRun, targetMachine)
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
}

// worker processes files from the jobs channel
func worker(jobs <-chan string, wg *sync.WaitGroup, dryRun bool, targetMachine string) {
	defer wg.Done()
	for file := range jobs {
		if err := transfer.CopyFile(file, targetMachine, dryRun); err != nil {
			// Error is already printed by CopyFile, just continue
			continue
		}
	}
}