package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/putabones/RATT/cmdline"
	"github.com/putabones/RATT/structs"
)

var ascii string = `
__________    ___________________________
\______   \  /  _  \__    ___/\__    ___/
 |       _/ /  /_\  \|    |     |    |   
 |    |   \/    |    \    |     |    |   
 |____|_  /\____|__  /____|     |____|   
        \/         \/                    
`

// version string
var version string = "2.0"

// parses user inputs
func parserFunc(t *structs.Target) {
	var parse = argparse.NewParser("RATT", "RATT stands for \"Recon All The Things\", it will perform scans against a target that is as intrusive as you want.\n\nRATT can run in 3 different modes\n   Replay: Replay results from a previous scan\n      CLI: Interactive mode to build and launch scans\n     Live: Immediately launches scans")

	// ip address for manual scans
	var i = parse.String("i", "ip", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "IP address to scan, leave blank for CLI mode",
		Default:  nil,
	})

	// folder where all outputs should be
	var f = parse.String("f", "folder", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Folder to save outputs",
		Default:  "/tmp/",
	})

	// nmap specific string
	// i.e. -Pn -sT -sC -sV
	var o = parse.String("o", "nmap", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Override NMAP Options",
		Default:  "-sT",
	})

	// amount of ports to scan, starts at 1
	var p = parse.Int("p", "ports", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Ports to scan, starts at 1 then up to 65535",
		Default:  200,
	})

	// number of concurrent workers
	var w = parse.Int("w", "workers", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Amount of concurrent workers to spawn",
		Default:  100,
	})

	// hostname
	var n = parse.String("n", "hostname", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Hostname for your target",
		Default:  "NoName",
	})

	// usernme
	var u = parse.String("", "user", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Username for follow on auths",
		Default:  nil,
	})

	// password
	var pass = parse.String("", "pass", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Password for follow on auths",
		Default:  nil,
	})

	// domain
	var d = parse.String("", "domain", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Domain for Windows auths",
		Default:  nil,
	})

	// version
	var v = parse.Flag("v", "version", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Prints the current version",
		Default:  false,
	})

	// check inputs
	err := parse.Parse(os.Args)
	if err != nil {
		fmt.Println(parse.Usage(nil))
		// CHANNGE THIS once the port is parsed right
	} else if *p > 65535 {
		fmt.Println("[-] Ports can't be more than 65535")
		os.Exit(31)
	} else if *v {
		fmt.Println("[i] Version: " + version)
		os.Exit(0)
	} else {
		t.Ip = *i
		t.Amt = *p // Ports
		t.NmapOptions = *o
		t.PortsCap = *w // Workers
		t.Hostname = *n
		t.Folder = *f
		t.Domain = *d
		t.Username = *u
		t.Password = *pass
	}
}

// main
func main() {

	// new target
	var tgt = new(structs.Target)

	// always gotta have ascii art
	fmt.Println(ascii)

	// parse user input
	parserFunc(tgt)

	// check if its Live, CLI, or Read Mode
	if tgt.Ip == "" {
		fmt.Print("[i] CLI Mode\n\n")
		cmdline.StartCLI()
	} else if tgt.Ip != "" {
		fmt.Println("[i] Live Mode")

		// ids linux or windows and empty string slice
		var user = os.Getuid()

		// looks from windows or nix
		switch user {
		case -1:
			fmt.Println("[+] OS: Windows")
		case 0:
			fmt.Println("[+] OS: Linux, User: root")
		default:
			fmt.Println("[-] OS: Linux, User: not root, some of the scans may not work...")
		}

		// launch scan
		tgt.Start()
		tgt.SmbCheck()
		if tgt.NmapOptions != "" {
			tgt.Nmap()
		}
	}
}
