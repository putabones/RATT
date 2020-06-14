package target

import (
	"fmt"
	"github.com/putabones/RATT/target/helpers"
)

func (t *Target) Enum4linux() {

	var command = " enum4linux -a " + t.Ip + " | tee " + t.Folder + "/" + "enum_" + t.Hostname + "_" + t.Ip + ".out"
	fmt.Println("[+] Running:", command)

	helpers.ExecuteCommand("enum4linux", command)

}
