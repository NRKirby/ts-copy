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

type extensionFlag []string

func (e *extensionFlag) String() string {
	return strings.Join(*e, ", ")
}

func (e *extensionFlag) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func main() {
	var extensions extensionFlag
	dryRun := flag.Bool("dry-run", false, "Show what would be copied without executing commands")
	flag.Var(&extensions, "ext", "File extension to copy (repeatable)")
	flag.Var(&extensions, "e", "File extension to copy (repeatable, short form)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <target-machine> [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <target-machine>    Name of the Tailscale machine to copy files to\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s my-server --ext .mp3 --ext .flac\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s my-server -e .pdf -e .docx --dry-run\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s my-server -e .zip\n", os.Args[0])
	}

	flag.Parse()

	// Validate target machine
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Target machine is required as the first argument\n\n")
		flag.Usage()
		os.Exit(1)
	}

	targetMachine := flag.Args()[0]

	// Validate extensions
	if len(extensions) == 0 {
		fmt.Fprintf(os.Stderr, "Error: At least one file extension must be specified using --ext or -e flag\n\n")
		flag.Usage()
		os.Exit(1)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Searching for files in: %s\n", currentDir)

	var matchingFiles []string
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			for _, configExt := range extensions {
				if ext == strings.ToLower(configExt) {
					matchingFiles = append(matchingFiles, path)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d matching files\n", len(matchingFiles))

	const maxWorkers = 5
	jobs := make(chan string, len(matchingFiles))
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, *dryRun, targetMachine)
	}

	for _, file := range matchingFiles {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	fmt.Println("All files processed")
}

func worker(jobs <-chan string, wg *sync.WaitGroup, dryRun bool, targetMachine string) {
	defer wg.Done()
	for file := range jobs {
		destination := targetMachine + ":"
		if dryRun {
			fmt.Printf("[DRY RUN] Would copy: %s to %s\n", file, destination)
		} else {
			fmt.Printf("Copying: %s to %s\n", file, destination)
			cmd := exec.Command("sudo", "tailscale", "cp", file, destination)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error copying %s: %v\n", file, err)
			}
		}
	}
}