package cmds

import (
	"fmt"
	"strings"

	"github.com/putabones/RATT/structs"
)

func Unset(t *structs.Target, v []string) {
	if len(v) == 2 {
		switch v[1] {
		case "ip":
			t.Ip = ""
		case "hostname":
			t.Hostname = ""
		case "ports":
			t.Amt = 0
		case "workers":
			t.PortsCap = 0
		case "nmap":
			t.NmapOptions = ""
		case "folder":
			t.Folder = ""
		default:
			fmt.Printf("[!] Bad set CMD: %s\n", strings.Join(v[:], " "))
		}
	} else if len(v) != 2 {
		fmt.Printf("[!] Bad set CMD: %s\n", strings.Join(v[:], " "))
	}
}
