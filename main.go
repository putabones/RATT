package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/cheggaaa/pb/v3"
	"net"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// target struct
type Target struct {
	ip          string
	hostname    string
	tcpopen     []int
	amt         int
	portsCap    int
	results     chan int
	nmapOptions string
}

// worker method
func (t Target) worker(ports, results chan int){
	// loop for 100 ports per Go routine
	for p := range ports{

		// convert and cat the IP and Port as a string
		var sock = t.ip + ":" + strconv.Itoa(p)

		// creates connection
		var _, err = net.Dial("tcp", sock)

		// either adds the port or a 0
		if err == nil {
			results <- p
		} else {
			results <- 0
		}
		continue
	}
}


// port check method
func (t Target) portCheck(ports, results chan int, b *pb.ProgressBar) {
	// start workers
	for i := 0; i < cap(ports); i++ {
		go t.worker(ports, results)
	}

	// loads port into ports channel
	go func() {
		for i := 1; i <= t.amt; i++ {
			// increment progress bar
			b.Increment()
			ports <- i
		}
	}()
}

// nmap method
func (t Target) nmap() {
	// ids linux or windows and empty string slice
	var user = os.Getuid()
	var str []string

	// looks from windows or nix
	switch user {
	case -1:
		fmt.Println("\n[+] Windows")
	case 0:
		fmt.Println("\n[+] root")
	default:
		fmt.Println("\n[-] not root, may not work...")
	}

	// combines the open ports
	for p := range t.tcpopen {
		str = append(str, strconv.Itoa(t.tcpopen[p]))
	}

	// set the nmap command to a string
	var command = "nmap " + t.nmapOptions + " -p " + strings.Join(str, ",") + " " + t.ip

	// nmap command execution
	var cmd = exec.Command("bash", "-c", command)
	fmt.Println("[+] nmap command\n", cmd.String())

	stderrstdout, err := cmd.CombinedOutput()
	if err != nil{
		fmt.Println("[-] Broke on nmap, Error:", err)
	} else {
		fmt.Println("[+] nmap output\n", string(stderrstdout))
	}
}

// start scanning
func (t Target) start() {
	var ports = make(chan int, t.portsCap) // channel to hold port numbers to be scanned
	var results = make(chan int)           // channel to hold open ports
	var bar = pb.StartNew(t.amt)           // progress bar
	var start = time.Now().UTC()           // start time from scan

	// start the port check
	t.portCheck(ports, results, bar)

	// append results to slice
	for i := 0; i < t.amt; i++ {
		p := <- results
		if p != 0 {
			t.tcpopen = append(t.tcpopen, p)
		}
	}

	// sort the ports Low to High
	sort.Ints(t.tcpopen)

	// close bar
	bar.Finish()

	// end time of scan
	var end = time.Now().UTC()

	// print results
	for p := range t.tcpopen {
		fmt.Println("[+] Open:", t.tcpopen[p])
	}

	// elapsed time
	var duration = end.Sub(start)
	fmt.Println("Scan Time:", duration.Truncate(time.Millisecond))

	// port 445 check
	for p := range t.tcpopen {
		if t.tcpopen[p] == 445 {
			go t.smbCheck()
		}
	}


	// close channels and bar
	close(ports)
	close(results)

	t.nmap()

}

// smb checks
func (t Target) smbCheck() {
	fmt.Println("[+] SMB Was Open")

	// setting smbclient command string, then executing
	var command = " smbclient -L //" + t.ip + " -N -U anonymous"
	fmt.Println("[+] Running:", command)
	var cmd = exec.Command("bash", "-c", command)
	var stderrstdout, err = cmd.CombinedOutput()
	if err != nil{
		fmt.Println("[-] smbclient error:", err)
	}

	// print the output from smbclient
	fmt.Println(string(stderrstdout))

	command = " enum4linux -a " + t.ip
	fmt.Println("[+] Running:", command)
	cmd = exec.Command("bash", "-c", command)
	stderrstdout, err = cmd.CombinedOutput()
	if err != nil{
		fmt.Println("[-] enum4linux error:", err)
	}

	// print the output from enum4linux
	fmt.Println(string(stderrstdout))

}

// TODO
// 	- parse if a range of IPs is given
// 	- add dirb and wfuzz methods
// 	- eventually add switches for dirb, wfuzz, etc.
// 		- wfuzz -c -w /usr/share/wordlists/dirb/common.txt http://10.10.10.180/FUZZ/web.config
// 		- dirb http://10.10.10.180 /usr/share/wordlists/dirb/small.txt -x /usr/share/wordlists/dirb/extensions_common.txt
// 	- add smbclient, enum4linux, and showmount listing method
// 		- sudo showmount -e 10.10.10.180

// main
func main() {

	var parse = argparse.NewParser("RATT", "RATT stands for \"Recon All The Things\", it will perform scans against a target that is as intrusive as you want")
	var tgt = new(Target) // new object of target

	var i = parse.String("i", "ip", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "IP address to scan",
		Default:  nil,
	})

	// need to parse the list of ips
	var r = parse.String("r", "range", &argparse.Options{
		Required: false,
		Validate: nil,
		Help:     "IP Range to scan, i.e. 192.168.0.0/24 or 192.168.0.1-255",
		Default:  nil,
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
		Default:  nil,
	})

	err := parse.Parse(os.Args)
	if err != nil {
		fmt.Println(parse.Usage(nil))
	} else if *i == "" && *r == "" {
		fmt.Println("[-] Need either \"-i\" or \"-r")
		fmt.Println("[-] Try \"-h\" or \"--help\" for arguments")
	} else if *i != "" && *r != "" {
		fmt.Println("[-] IP Set To:", *i)
		fmt.Println("[-] IP Range Set To:", *r)
		fmt.Println("[-] Need either \"-i\" or \"-r\", but not both")
		fmt.Println("[-] Try \"-h\" or \"--help\" for arguments")
	} else if *p > 65535 {
		fmt.Println("[-]")
	} else {
		tgt.ip = *i
		tgt.amt = *p
		tgt.portsCap = *w
		tgt.hostname = *n
		tgt.nmapOptions = *o
		fmt.Println("[+] IP:", tgt.ip)
		fmt.Println("[+] Hostname:", tgt.hostname)
		fmt.Println("[+] NMAP Options:", tgt.nmapOptions)
		fmt.Println("[+] Amount of Ports:", tgt.amt)
		fmt.Println("[+] Workers Setup:", tgt.portsCap)
		fmt.Println()

		tgt.start()
	}
}