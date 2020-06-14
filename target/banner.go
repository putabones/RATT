package target

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

// loops over the open tcp ports, then banner grabs them
func (t *Target) BannerGrab() {
	for _, p := range t.Tcpopen {
		var sock = t.Ip + ":" + strconv.Itoa(p)
		fmt.Println("[+] Trying BanenrGrab on:", sock)

		// choose DialTimeout for any pesky connections
		var conn, err = net.Dial("tcp", sock)
		if err != nil {
			fmt.Println("[-] BannerGrab Error: ", err, "; on: ", sock)
		} else {

			// timer for pesky protocols
			var timer = time.NewTimer(2 * time.Second)

			// Have to launch as go function in order to implement timer
			go func() {

				// Added the GET request for HTTP, does not appear to mess with other protocols
				var _, err = fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
				if err != nil {
					fmt.Println("[-] BannerGrab Format Conn Error: ", err)
				}

				// Create byte array, then read the banner in
				var b = make([]byte, 4096)
				n, err := conn.Read(b)
				if err != nil {
					fmt.Println("[-] BannerGrab Read Error: ", err)
				} else {
					fmt.Println("[+] Banner for: "+sock+" ; Size:", n)

					// Choose to write to file vice straight printing to screen due to web page sizes
					var fname = t.Folder + "/banner_" + sock + ".out"
					err = ioutil.WriteFile(fname, b[:n], 0644)
					if err != nil {
						fmt.Println("[-] BannerGrab Error Writing")
						fmt.Println(string(b[:n]))
					} else {
						fmt.Println("[+] Wrote Banner to:", fname)
					}
				}
			}()
			// start the timer, then close the connection
			<-timer.C
			_ = conn.Close()
		}
	}
}
