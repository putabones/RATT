package Target

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"net"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// target struct
type Target struct {
	Ip          string
	Hostname    string
	Tcpopen     []int
	Amt         int
	PortsCap    int
	Results     chan int
	NmapOptions string
	Folder      string
}

// nmap method
func (t Target) nmap() {
	var str []string

	// combines the open ports
	for p := range t.Tcpopen {
		str = append(str, strconv.Itoa(t.Tcpopen[p]))
	}

	// set the nmap command to a string
	var command = "nmap " + t.NmapOptions + " -p " + strings.Join(str, ",") + " " + t.Ip + " -oA " + t.Folder + "/" + t.Hostname + "_" + t.Ip

	// nmap command execution
	var cmd = exec.Command("bash", "-c", command)
	fmt.Println("[+] nmap command:", cmd.String())

	stderrstdout, err := cmd.CombinedOutput()
	if err != nil{
		fmt.Println("[-] Broke on nmap, check command; Error:", err)
	} else {
		fmt.Println("[+] nmap output:\n", string(stderrstdout))
	}
}

// port check method
func (t Target) portCheck(ports, results chan int, b *pb.ProgressBar) {
	// start workers
	for i := 0; i < cap(ports); i++ {
		go t.worker(ports, results, b)
	}

	// loads port into ports channel
	go func() {
		for i := 1; i <= t.Amt; i++ {
			// increment progress bar
			ports <- i
		}
	}()
}

// smb checks
func (t Target) smbCheck(c chan string) {
	// setting smbclient command string, then executing
	var command = " smbclient -L //" + t.Ip + " -N -U anonymous " + "| tee " +  t.Folder + "/" + "smbclient_" + t.Hostname + "_" + t.Ip + ".out"
	fmt.Println("[+] Running:", command)
	var cmd = exec.Command("bash", "-c", command)
	var stderrstdout, err = cmd.CombinedOutput()
	if err != nil{
		fmt.Println("[-] smbclient error:", err)
	}

	// print the output from smbclient
	fmt.Println(string(stderrstdout))

	command = " enum4linux -a " + t.Ip + " | tee " + t.Folder + "/" + "enum_" + t.Hostname + "_" + t.Ip + ".out"
	fmt.Println("[+] Running:", command)
	cmd = exec.Command("bash", "-c", command)
	stderrstdout, err = cmd.CombinedOutput()
	if err != nil{
		fmt.Println("[-] enum4linux error:", err)
	}

	// print the output from enum4linux
	fmt.Println(string(stderrstdout))

	// send signal to channel that scan is done
	c <- "SMB"
	close(c)
}

// start scanning
func (t Target) Start() {
	var ports = make(chan int, t.PortsCap) // channel to hold port numbers to be scanned
	var results = make(chan int)           // channel to hold open ports
	var smb = make(chan string)            // channel for smb scan
	var ans string                         // smb scan choice
	var bar = pb.StartNew(t.Amt)           // progress bar
	var start = time.Now().UTC()           // start time from scan

	fmt.Println("[+] IP:", t.Ip)
	fmt.Println("[+] Hostname:", t.Hostname)
	fmt.Println("[+] NMAP Options:", t.NmapOptions)
	fmt.Println("[+] Amount of Ports:", t.Amt)
	fmt.Println("[+] Workers Setup:", t.PortsCap)
	fmt.Println()

	// start the port check
	t.portCheck(ports, results, bar)

	// append results to slice
	for i := 0; i < t.Amt; i++ {
		p := <- results
		if p != 0 {
			t.Tcpopen = append(t.Tcpopen, p)
		}
	}

	// sort the ports Low to High
	sort.Ints(t.Tcpopen)

	// close bar
	bar.Finish()

	// end time of scan
	var end = time.Now().UTC()

	// print results
	for p := range t.Tcpopen {
		fmt.Println("[+] Open:", t.Tcpopen[p])
	}

	// elapsed time
	var duration = end.Sub(start)
	fmt.Println("Scan Time:", duration.Truncate(time.Millisecond))

	// port 445 check
	for p := range t.Tcpopen {
		if t.Tcpopen[p] == 445 {
			fmt.Println()
			fmt.Print("[+] Port 445 open, do you want to run SMB Checks? Y/N: ")
			fmt.Scanln(&ans)
			if strings.ToUpper(ans) == "Y" {
				go t.smbCheck(smb)
			} else {
				fmt.Println("[-] Not scanning 445")
				close(smb)
			}
		}
	}

	// close channels and bar
	close(ports)
	close(results)

	t.nmap()

	// close channel catch in case there is no answer to smb check
	if ans == "Y" {
		for s := range smb {
			fmt.Printf("[+] %v Checks Complete\n", s)
		}
	} else {
		close(smb)
	}
}

// worker method
func (t Target) worker(ports, results chan int, b * pb.ProgressBar){
	// loop for 100 ports per Go routine
	for p := range ports{
		// progress bar
		b.Increment()

		// convert and cat the IP and Port as a string
		var sock = t.Ip + ":" + strconv.Itoa(p)

		// creates connection
		var conn, err = net.DialTimeout("tcp", sock, time.Second * 2)

		// either adds the port or a 0
		if err == nil {
			results <- p
			_ = conn.Close()
		} else {
			results <- 0
		}
		continue
	}
}