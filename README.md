# RATT

RATT stands for "Recon All The Things", its purpose is to speed up scanning and service enumeration by using Golang's concurrency.  Works on both Linux and Windows, some scans are dependent on packages that are installed.

#### Version: 1.5.1

## RATT Help

### Switches

```
-h|--help                Normal help menu
-i|--ip "<value>"        REQUIRED; IP to be scanned
-f|--folder "<value>"    Folder to save output. Default: /tmp  NOTE: Does not create the folder
-o|--nmap "<value>"      Custom switches for nmap.
-p|--ports <integer>     Number of ports to scan. NOTE: Not the actual port numbers. i.e 100=ports 1-100.
-w|--workers <integer>   Number of scans to take place simultaneously. Default: 100
-n|--hostname "<value>"  Hostname of the device. It will later give the options to save to host file.
-b|--banner              Boolean, if present, indicates you want banners to be grabbed if possible
```

### RATT Help Output

```
usage: RATT [-h|--help] -i|--ip "<value>" [-f|--folder "<value>"] [-o|--nmap
            "<value>"] [-p|--ports <integer>] [-w|--workers <integer>]
            [-n|--hostname "<value>"] [-b|--banner]

            RATT stands for "Recon All The Things", it will perform scans
            against a target that is as intrusive as you want. Version: 1.3

Arguments:

  -h  --help      Print help information
  -i  --ip        IP address to scan
  -f  --folder    Folder to save outputs. Default: /tmp/
  -o  --nmap      Override NMAP Options. Default: -sT
  -p  --ports     Ports to scan, starts at 1 then up to 65535. Default: 200
  -w  --workers   Amount of concurrent workers to spawn. Default: 100
  -n  --hostname  Hostname for your target, will be added to host file if
                  added. Default: NoName
  -b  --banner    Banner grabs on open TCP Ports. Default: false

```
