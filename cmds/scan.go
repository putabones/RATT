package cmds

import (
	"github.com/putabones/RATT/structs"
)

var Scan = structs.CMD{
	Name: "scan",
	Help: "tcp scan your target",
	Execute: func(t *structs.Target) {

		// calling the target start method
		t.Start()
	},
}
