package main

import (
	"RATT/target"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

// TODO
// 	- add dirb and wfuzz methods
// 	- eventually add switches for dirb, wfuzz, etc.
// 		- wfuzz -c -w /usr/share/wordlists/dirb/common.txt http://10.10.10.180/FUZZ/web.config
// 		- dirb http://10.10.10.180 /usr/share/wordlists/dirb/small.txt -x /usr/share/wordlists/dirb/extensions_common.txt
// 	- add showmount listing method
// 		- sudo showmount -e 10.10.10.180
// 	- add on disk save location
// 	- add "|tee" to save outputs individually
// 		- smbclient
// 		- enum4linux

// parses user inputs
func parserFunc(t *target.Target) {
	var parse = argparse.NewParser("RATT", "RATT stands for \"Recon All The Things\", it will perform "+
		"scans against a target that is as intrusive as you want. Version: 1.3")

	var i = parse.String("i", "ip", &argparse.Options{
		Required: true,
		Validate: nil,
		Help:     "IP address to scan",
		Default:  nil,
	})

	var f = parse.String("f", "folder", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "Folder to save outputs",
		Default:  "/tmp/",
	})

	// can add multiple -n to append to the list
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
		Help:     "Hostname for your target, will be added to host file if added",
		Default:  "NoName",
	})

	// check inputs
	err := parse.Parse(os.Args)
	if err != nil {
		fmt.Println(parse.Usage(nil))
	} else if *i == "" {
		fmt.Println("[-] Need an IP with \"-i\"")
		fmt.Println("[-] Try \"-h\" or \"--help\" for arguments")
	} else if *p > 65535 {
		fmt.Println("[-] Ports can't be more than 65535")
	} else {
		t.Ip = *i
		t.Amt = *p
		t.NmapOptions = *o
		t.PortsCap = *w
		t.Hostname = *n
		t.Folder = *f
	}
}

// main
func main() {

	// new target
	var tgt = new(target.Target)

	// parse user input
	parserFunc(tgt)

	if tgt.Ip != "" {
		// ids linux or windows and empty string slice
		var user = os.Getuid()

		// looks from windows or nix
		switch user {
		case -1:
			fmt.Println("\n[+] OS: Windows")
		case 0:
			fmt.Println("\n[+] OS: Linux, User: root")
		default:
			fmt.Println("\n[-] OS: Linux, User: not root, some of the scans may not work...")
		}

		// launch scan
		tgt.Start()
	}
}
