package cmds

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/putabones/RATT/structs"
)

// print current settings
func print(t *structs.Target) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Ports = Amt
	// Workers = PortsCap

	// variables for pretty columns
	var ipLen string
	var hnLen string
	var tcpOpenLen string
	amtLen := strings.Repeat("-", 5)
	portsCapLen := strings.Repeat("-", 7)
	var nmapLen string
	var folderLen string
	var domainLen string
	var userLen string
	var passLen string

	// making tcp open ports string
	var portsJoined string
	if len(t.Tcpopen) > 0 {
		strPorts := []string{}
		for _, p := range t.Tcpopen {
			strPorts = append(strPorts, strconv.Itoa(p))
		}

		portsJoined = strings.Join(strPorts, ",")
	}

	// getting the - length for pretty printing
	if len(t.Ip) > 2 {
		ipLen = strings.Repeat("-", len(t.Ip))
	} else {
		ipLen = strings.Repeat("-", 2)
	}
	if len(t.Hostname) > 8 {
		hnLen = strings.Repeat("-", len(t.Hostname))
	} else {
		hnLen = strings.Repeat("-", 8)
	}
	if len(portsJoined) > 14 {
		tcpOpenLen = strings.Repeat("-", len(portsJoined))
	} else {
		tcpOpenLen = strings.Repeat("-", 14)
	}
	if len(t.NmapOptions) > 4 {
		nmapLen = strings.Repeat("-", len(t.NmapOptions))
	} else {
		nmapLen = strings.Repeat("-", 4)
	}
	if len(t.Folder) > 6 {
		folderLen = strings.Repeat("-", len(t.Folder))
	} else {
		folderLen = strings.Repeat("-", 6)
	}

	// auths
	if len(t.Domain) > 5 {
		domainLen = strings.Repeat("-", len(t.Domain))
	} else {
		domainLen = strings.Repeat("-", 5)
	}
	if len(t.Username) > 4 {
		userLen = strings.Repeat("-", len(t.Username))
	} else {
		userLen = strings.Repeat("-", 4)
	}
	if len(t.Password) > 4 {
		passLen = strings.Repeat("-", len(t.Password))
	} else {
		passLen = strings.Repeat("-", 4)
	}

	// adding table to buffer
	fmt.Fprintln(w, "\nip\thostname\tTCP Open Ports\tports\tworkers\tnmap\tfolder\tdomain\tuser\tpass")
	fmt.Fprintln(w, ipLen+"\t"+hnLen+"\t"+tcpOpenLen+"\t"+amtLen+"\t"+portsCapLen+"\t"+nmapLen+"\t"+folderLen+"\t"+domainLen+"\t"+userLen+"\t"+passLen)

	// prints if there is ports
	if portsJoined != "" {
		fmt.Fprintln(w, t.Ip+"\t"+t.Hostname+"\t"+portsJoined+"\t"+strconv.Itoa(t.Amt)+"\t"+strconv.Itoa(t.PortsCap)+"\t"+t.NmapOptions+"\t"+t.Folder+"\t"+t.Domain+"\t"+t.Username+"\t"+t.Password+"\n")
	} else {
		fmt.Fprintln(w, t.Ip+"\t"+t.Hostname+"\t"+"\t"+strconv.Itoa(t.Amt)+"\t"+strconv.Itoa(t.PortsCap)+"\t"+t.NmapOptions+"\t"+t.Folder+"\t"+t.Domain+"\t"+t.Username+"\t"+t.Password+"\n")
	}

	// writes the output from the buffer
	w.Flush()
}

// set function
func Set(t *structs.Target, v []string) {
	if len(v) == 1 {
		print(t)
	} else if len(v) == 2 {
		fmt.Printf("[!] Need the value for the variable\n[i] CMD: %s %s\n", v[0], v[1])
	} else if len(v) == 3 {
		switch v[1] {
		case "ip":
			t.Ip = v[2]
		case "hostname":
			t.Hostname = v[2]
		case "ports":
			t.Amt, _ = strconv.Atoi(v[2])
		case "workers":
			t.PortsCap, _ = strconv.Atoi(v[2])
		case "nmap":
			t.NmapOptions = v[2]
		case "folder":
			t.Folder = v[2]
		case "domain":
			t.Domain = v[2]
		case "user":
			t.Username = v[2]
		case "pass":
			t.Password = v[2]
		default:
			fmt.Printf("[!] Bad set CMD: %s\n", strings.Join(v[:], " "))
		}
	} else if len(v) >= 3 && v[1] == "nmap" {
		t.NmapOptions = strings.Join(v[2:], " ")
	} else if len(v) >= 3 && v[1] == "pass" {
		t.Password = strings.Join(v[2:], " ")
	} else {
		fmt.Printf("[!] Bad set CMD: %s\n", strings.Join(v[:], " "))
	}
}

// write yaml config

// read yaml config
