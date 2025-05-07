package cmds

import (
	"github.com/putabones/RATT/structs"
)

var Nmap = structs.CMD{
	Name: "nmap",
	Help: "start nmap scan, !!!make sure your switches are set!!!",
	Execute: func(t *structs.Target) {

		// call nmap method from target struct
		t.Nmap()
	},
}
