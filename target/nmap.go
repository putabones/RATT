package target

import (
	"RATT/target/helpers"
	"strconv"
	"strings"
)

// nmap method
func (t *Target) Nmap() {

	var str []string

	// combines the open ports
	for p := range t.Tcpopen {
		str = append(str, strconv.Itoa(t.Tcpopen[p]))
	}

	// set the nmap command to a string
	var command = "nmap " + t.NmapOptions + " -p " + strings.Join(str, ",") + " " + t.Ip + " -oA " + t.Folder + "/" + t.Hostname + "_" + t.Ip

	helpers.ExecuteCommand("nmap", command)
}
