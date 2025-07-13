package transfer

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// TailscaleStatus represents the JSON response from 'tailscale status --json'
type TailscaleStatus struct {
	Self struct {
		HostName string `json:"HostName"`
	} `json:"Self"`
	Peer map[string]struct {
		HostName string `json:"HostName"`
		Online   bool   `json:"Online"`
	} `json:"Peer"`
}

// CheckTargetMachineOnline verifies that the target machine is online via Tailscale
func CheckTargetMachineOnline(targetMachine string) error {
	cmd := exec.Command("tailscale", "status", "--json")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get tailscale status: %v", err)
	}

	var status TailscaleStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return fmt.Errorf("failed to parse tailscale status: %v", err)
	}

	// Check if target machine matches our own hostname
	if strings.EqualFold(status.Self.HostName, targetMachine) {
		return fmt.Errorf("cannot copy to self (target machine '%s' is the current machine)", targetMachine)
	}

	// Look for target machine in peers
	for _, peer := range status.Peer {
		if strings.EqualFold(peer.HostName, targetMachine) {
			if !peer.Online {
				return fmt.Errorf("target machine '%s' is offline", targetMachine)
			}
			return nil
		}
	}

	return fmt.Errorf("target machine '%s' not found in tailscale network", targetMachine)
}

// CopyFile copies a single file to the specified Tailscale machine.
// If dryRun is true, it only prints what would be done without executing.
func CopyFile(filePath, targetMachine string, dryRun bool) error {
	destination := targetMachine + ":"
	
	if dryRun {
		fmt.Printf("[DRY RUN] Would copy: %s to %s\n", filePath, destination)
		return nil
	}

	fmt.Printf("Copying: %s to %s\n", filePath, destination)
	cmd := exec.Command("tailscale", "file", "cp", filePath, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error copying %s: %s", filePath, string(output))
	}
	
	return nil
}