package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Show what would be copied without executing commands")
	flag.Parse()

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Searching for audio files in: %s\n", currentDir)

	var audioFiles []string
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".mp3" || ext == ".flac" || ext == ".wav" {
				audioFiles = append(audioFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d audio files\n", len(audioFiles))

	const maxWorkers = 5
	jobs := make(chan string, len(audioFiles))
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, *dryRun)
	}

	for _, file := range audioFiles {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	fmt.Println("All files processed")
}

func worker(jobs <-chan string, wg *sync.WaitGroup, dryRun bool) {
	defer wg.Done()
	for file := range jobs {
		if dryRun {
			fmt.Printf("[DRY RUN] Would copy: %s\n", file)
		} else {
			fmt.Printf("Copying: %s\n", file)
			cmd := exec.Command("sudo", "tailscale", "cp", file, "my-tailscale-node:")
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error copying %s: %v\n", file, err)
			}
		}
	}
}