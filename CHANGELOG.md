# Changelog
## [2.0] - 07 May 2025
- Lots of code clen up and refactoring, too many to list
- Created new mode: CLI, made the original Live Mode
- Created structs folder and moved all methods related to each struct to their respctive files, vice it being separate as before
- Removed smbclient and enum4linux
- Created the cmd struct to make it easier to add furture commands later
- The target struct is where all command syntax will take place
- Added basic NXC SMB auth check
- Added domain creds to parsing and target struct

// TODO
// 	- add dirb and wfuzz methods
// 	- eventually add switches for dirb, wfuzz, etc.
// 		- wfuzz -c -w /usr/share/wordlists/dirb/common.txt http://10.10.10.180/FUZZ/web.config
// 		- dirb http://10.10.10.180 /usr/share/wordlists/dirb/small.txt -x /usr/share/wordlists/dirb/extensions_common.txt
// 	- add showmount listing method
// 		- sudo showmount -e 10.10.10.180
//  - add bloodhound
//  - add config read
//  - add config create
//  - add write config
//  - add cmdline
//  - add dns lookup capabilities
//    - dns server
//    - /etc/host file
//  - add capability to do specific ports and ranges
//  - add UDP capability
//  - add capability to re-read previous scans
//  - add logging capabilities plus printing to the screen
