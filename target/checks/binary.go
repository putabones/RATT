package checks

import "os/exec"

// Checks if binary is in env path
func CheckInPath(path string) bool {

	_, err := exec.LookPath(path)
	if err != nil {
		return false
	}
	return true
}
