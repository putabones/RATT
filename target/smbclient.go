package target

import (
	"RATT/target/helpers"
	"fmt"
)

func (t *Target) Smbclient() {

	// setting smbclient command string, then executing
	var command = " smbclient -L //" + t.Ip + " -N -U anonymous " + "| tee " + t.Folder + "/" + "smbclient_" + t.Hostname + "_" + t.Ip + ".out"
	fmt.Println("[+] Running:", command)

	helpers.ExecuteCommand("smbclient", command)

}
