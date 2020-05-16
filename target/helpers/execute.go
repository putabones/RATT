package helpers

import (
	"fmt"
	"os/exec"
)

func ExecuteCommand(program, command string) {

	//  command execution
	var cmd = exec.Command("bash", "-c", command)
	fmt.Println("[+]"+program+" command:", cmd.String())

	stderrstdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("[-] Broke on "+program+", check command; Error:", err)
	} else {
		fmt.Println("[+] "+program+" output:\n", string(stderrstdout))
	}

}
