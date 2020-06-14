package structs

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
