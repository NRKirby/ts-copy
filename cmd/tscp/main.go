package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"ts-copy/internal/discovery"
	"ts-copy/internal/transfer"
	"ts-copy/internal/worker"
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

	// Find target machine before parsing flags
	var targetMachine string
	var filteredArgs []string
	
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Target machine is required as the first argument\n\n")
		flag.Usage()
		os.Exit(1)
	}
	
	targetMachine = os.Args[1]
	
	// Create new args list with target machine removed
	filteredArgs = append([]string{os.Args[0]}, os.Args[2:]...)
	os.Args = filteredArgs
	
	flag.Parse()

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

	// Find matching files
	matchingFiles, err := discovery.FindMatchingFiles(currentDir, extensions)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d matching files\n", len(matchingFiles))

	// Preflight check: verify target machine is online
	if !*dryRun {
		fmt.Printf("Checking if target machine '%s' is online...\n", targetMachine)
		if err := transfer.CheckTargetMachineOnline(targetMachine); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Target machine '%s' is online and ready\n", targetMachine)
	}

	// Process files with worker pool
	const maxWorkers = 5
	errorCount := worker.ProcessFiles(matchingFiles, targetMachine, *dryRun, maxWorkers)
	
	if errorCount > 0 {
		if errorCount == len(matchingFiles) {
			// All files failed, don't print anything positive
			os.Exit(1)
		} else {
			fmt.Printf("Completed with %d error(s)\n", errorCount)
			os.Exit(1)
		}
	} else {
		fmt.Println("All files processed successfully")
	}
}