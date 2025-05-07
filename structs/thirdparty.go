package structs

import "os/exec"

// target 3rdparty binary
type ThirdParty struct {
	Nmap bool
	Nxc  bool
}

// method to check if binaries are in the execution path
func (tp ThirdParty) CheckInPath(path string) bool {

	_, err := exec.LookPath(path)
	if err != nil {
		return false
	}
	return true
}
