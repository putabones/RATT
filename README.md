# RATT

RATT stands for "Recon All The Things", its purpose is to speed up scanning and service enumeration by using Golang's concurrency.  Works on both Linux and Windows, some scans are dependent on packages that are installed.

#### Version: 2.0

## RATT Help

### Switches

```
__________    ___________________________
\______   \  /  _  \__    ___/\__    ___/
 |       _/ /  /_\  \|    |     |    |   
 |    |   \/    |    \    |     |    |   
 |____|_  /\____|__  /____|     |____|   
        \/         \/

usage: RATT [-h|--help] [-i|--ip "<value>"] [-f|--folder "<value>"] [-o|--nmap
            "<value>"] [-p|--ports <integer>] [-w|--workers <integer>]        
            [-n|--hostname "<value>"] [--user "<value>"] [--pass "<value>"]   
            [--domain "<value>"] [-v|--version]

            RATT stands for "Recon All The Things", it will perform scans     
            against a target that is as intrusive as you want.

RATT can run in 3 different modes
   Replay: Replay results from a previous scan
      CLI: Interactive mode to build and launch scans
     Live: Immediately launches scans

Arguments:

  -h  --help      Print help information
  -i  --ip        IP address to scan, leave blank for CLI mode
  -f  --folder    Folder to save outputs. Default: /tmp/
  -o  --nmap      Override NMAP Options. Default: -sT
  -p  --ports     Ports to scan, starts at 1 then up to 65535. Default: 200
  -w  --workers   Amount of concurrent workers to spawn. Default: 100
  -n  --hostname  Hostname for your target. Default: NoName
      --user      Username for follow on auths
      --pass      Password for follow on auths
      --domain    Domain for Windows auths
  -v  --version   Prints the current version. Default: false
  ```
  