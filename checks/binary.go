package checks

import "os/exec"

// Checks if binary is in env path
func CheckInPath(path string) (string, bool, error) {

	binaryPath, err := exec.LookPath(path)
	if err != nil {
		return path, false, err
	}
	return binaryPath, true, err
}
