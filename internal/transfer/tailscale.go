package transfer

import (
	"fmt"
	"os/exec"
)

// CopyFile copies a single file to the specified Tailscale machine.
// If dryRun is true, it only prints what would be done without executing.
func CopyFile(filePath, targetMachine string, dryRun bool) error {
	destination := targetMachine + ":"
	
	if dryRun {
		fmt.Printf("[DRY RUN] Would copy: %s to %s\n", filePath, destination)
		return nil
	}

	fmt.Printf("Copying: %s to %s\n", filePath, destination)
	cmd := exec.Command("sudo", "tailscale", "cp", filePath, destination)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error copying %s: %v", filePath, err)
	}
	
	return nil
}