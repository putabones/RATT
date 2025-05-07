package structs

import (
	"fmt"
	"net"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// third party var
var tp ThirdParty
var ThirdPartyChecks = ThirdParty{
	Nmap: tp.CheckInPath("nmap"),
	Nxc:  tp.CheckInPath("nxc"),
}

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
	Domain      string
	Username    string
	Password    string
}

// port check method
func (t *Target) portCheck(ports, results chan int, b *pb.ProgressBar) {
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
func (t *Target) smbCheck(c chan string) {

	if ThirdPartyChecks.Nxc {
		t.NxcSmbAuth()
	} else {
		fmt.Println("[!] Smbclient or Bash not found in path")
	}

	// send signal to channel that scan is done
	c <- "SMB"
	close(c)
}

// start scanning
func (t *Target) Start() {

	// Return true/false for binaries needed
	//checkBinary()

	var ports = make(chan int, t.PortsCap) // channel to hold port numbers to be scanned
	var results = make(chan int)           // channel to hold open ports
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
		p := <-results
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

	// close channels
	close(ports)
	close(results)
}

// worker method
func (t *Target) worker(ports, results chan int, b *pb.ProgressBar) {
	// loop for 100 ports per Go routine
	for p := range ports {
		// progress bar
		b.Increment()

		// convert and cat the IP and Port as a string
		var sock = t.Ip + ":" + strconv.Itoa(p)

		// creates connection
		var conn, err = net.DialTimeout("tcp", sock, time.Second*2)

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

// nmap method
func (t *Target) Nmap() {
	// do nmap
	if ThirdPartyChecks.Nmap {
		var str []string

		// combines the open ports
		for p := range t.Tcpopen {
			str = append(str, strconv.Itoa(t.Tcpopen[p]))
		}

		// excute command
		var cmd = exec.Command(
			"nmap",
			t.NmapOptions,
			"-p", strings.Join(str, ","),
			t.Ip,
			"-oA", t.Folder+"/"+t.Hostname+"_"+t.Ip,
		)
		stderrstdout, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("[-] Broke on nmap, check command; Error:", err)
		}
		fmt.Println("[+] nmap output:\n", string(stderrstdout))
	} else {
		fmt.Println("[!] Nmap not found in path")
	}
}

// smb checks
func (t *Target) SmbCheck() {
	// vars
	var ans string              // smb scan choice
	var smb = make(chan string) // channel for smb scan

	// port 445 check
	for _, p := range t.Tcpopen {
		if p == 445 {
			fmt.Print("\n[+] Port 445 open, do you want to run SMB Checks? Y/N: ")
			fmt.Scanln(&ans)
			if strings.ToUpper(ans) == "Y" {
				go t.smbCheck(smb)
			} else {
				fmt.Println("[-] Not scanning 445")
			}
		}
	}

	// Checks the SMB channel, closes once checks are done, if you choose not to scan, or 445 wasnt open
	if ans == "Y" {
		for s := range smb {
			fmt.Printf("[+] %v Checks Complete\n", s)
		}
	} else {
		close(smb)
	}

}

// nxc smb auth check
// example nxc smb --log /home/covid/Documents/files/htb/escapetwo/scans/nxc_rose_smb.out --dns-server 10.10.11.51 -u rose -p KxEPkKe6R8su -d sequel.htb DC01.sequel.htb
func (t *Target) NxcSmbAuth() {
	// do nxc auth
	if ThirdPartyChecks.Nxc {
		var str []string

		// combines the open ports
		for p := range t.Tcpopen {
			str = append(str, strconv.Itoa(t.Tcpopen[p]))
		}

		// excute command
		var cmd = exec.Command(
			"nxc",
			"smb",
			"-u", t.Username,
			"-p", t.Password,
			"-d", t.Domain,
			"-log", t.Folder+"/"+t.Hostname+"_"+t.Ip+"_"+t.Domain+"_"+t.Username+"nxcSmbAuth.log",
			t.Ip,
			t.NmapOptions,
			"-p", strings.Join(str, ","),
			t.Ip,
			"-oA", t.Folder+"/"+t.Hostname+"_"+t.Ip,
		)
		stderrstdout, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("[-] Broke on nxc, check command; Error:", err)
		}
		fmt.Println("[+] nxc output:\n", string(stderrstdout))
	} else {
		fmt.Println("[!] Netexec not found in path")
	}

}
