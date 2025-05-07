package cmdline

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/putabones/RATT/cmds"
	"github.com/putabones/RATT/structs"
)

// help func
func help(list []structs.CMD) {
	// tab spaced writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Println("\n[i] Commands\n------------")
	for _, c := range list {
		fmt.Fprintln(w, c.Name+"\t:\t"+c.Help)
	}
	fmt.Fprintln(w, "set\t:\tprint current target variables")
	fmt.Fprintln(w, "set <var>\t:\tset the associated variable")
	fmt.Fprintln(w, "unset <var>\t:\tunset the associated variable")

	// writes data to stdout from buffer
	w.Flush()
	fmt.Println("\nq to quit")
}

// launches the RATT cli
func StartCLI() {
	var tgt = new(structs.Target)

	// loop over the commands
	var cmdList = []structs.CMD{
		cmds.Nmap,
		cmds.Scan,
	}

	// label for loop
cmdlineLoop:
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("h for help >>> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		cmdStrings := strings.Split(cmd, " ")
		switch cmdStrings[0] {
		case "h":
			help(cmdList)
		case "set":
			cmds.Set(tgt, cmdStrings)
		case "unset":
			cmds.Unset(tgt, cmdStrings)
		case "q":
			break cmdlineLoop
		case "scan":
			cmds.Scan.Execute(tgt)
		case "nmap":
			if len(tgt.NmapOptions) == 0 {
				fmt.Println("[!] No nmap switches set!!!")
			} else if len(tgt.Folder) == 0 {
				fmt.Println("[!] No folder set!!! Make sure it exists!!!")
			} else {
				cmds.Nmap.Execute(tgt)
			}
		}
	}
}
