package target

import (
	"RATT/target/helpers"
	"fmt"
)

func (t *Target) Enum4linux() {

	var command = " enum4linux -a " + t.Ip + " | tee " + t.Folder + "/" + "enum_" + t.Hostname + "_" + t.Ip + ".out"
	fmt.Println("[+] Running:", command)

	helpers.ExecuteCommand("enum4linux", command)

}
