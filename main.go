package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Extensions        []string `yaml:"extensions"`
	TargetTsMachine   string   `yaml:"targetTsMachine"`
}

func loadConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	configPath := filepath.Join(usr.HomeDir, ".ts-copy", "config.yaml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s\n\nPlease create a config file with the following format:\n\nextensions:\n  - \".mp3\"\n  - \".flac\"\n  - \".wav\"\ntargetTsMachine: \"my-tailscale-machine\"", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.TargetTsMachine == "" {
		return nil, fmt.Errorf("targetTsMachine is required in config file but was not found or is empty")
	}

	if len(config.Extensions) == 0 {
		return nil, fmt.Errorf("extensions array is required in config file but was not found or is empty")
	}

	return &config, nil
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Show what would be copied without executing commands")
	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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
			for _, configExt := range config.Extensions {
				if ext == strings.ToLower(configExt) {
					audioFiles = append(audioFiles, path)
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

	fmt.Printf("Found %d audio files\n", len(audioFiles))

	const maxWorkers = 5
	jobs := make(chan string, len(audioFiles))
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, *dryRun, config.TargetTsMachine)
	}

	for _, file := range audioFiles {
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