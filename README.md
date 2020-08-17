# nier

## Installation
```
go get github.com/ren-zxcyq/nier
go install github.com/ren-zxcyq/nier
```

## Run

Run all checks
```
~/go/src/github.com/ren-zxcyq/nier# go install
# ~/go/bin/nier -host 192.168.1.20 -o ~/Desktop/report -sess "PHPSESSID:g47sf085b2hq5hv3871pdtuo11;SecretCookie:VzuuL2gfLJWNnTSwn2kuLv5wo20vBwpjAGWwLJD2LwDkAJL0ZwplLmR5BQMuLGyuAGOuA2ZmBwR1Amt2AmH1ZGV" -all -cve


## Usage
```
# ~/go/bin/nier -h

        ⣤⡄⠀⠀⣤⢠⢠⠀⠀⠀⠀⣤⠄⠀⢤⡀⠀⠀⠀⢀⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⢀⢀⣤⠀⠀⠀⠀⣒⣒
        ⡇⠙⢦⠀⣿⠀⡄⠀⣀⡀⠀⣿⠀⢀⡼⠃⢠⠀⠀⡘⢻⠀⣀⠀⢀⡀⣰⣀⠀⢀⣀⠀⠁⡀⠀⣀⠀⢨⣄⠠⠤⣤⠇⠿⣢⣀
        ⡇⠀⠈⠳⣿⠀⡇⡜⠀⢹⡆⣿⠚⠙⣆⠀⠀⠀⢀⢃⣸⡇⢸⠀⠀⡇⢸⠀⢰⠁⠈⡃⡄⣵⡆⢐⣖⠈⣷⡾⠇⢨⠀⣽⢡⡌⡆⡧⢺
        ⡇⠀⠀⠀⣿⠀⡇⣷⠊⠁⠀⣿⠀⠀⢹⡀⠐⠀⡘⠉⠀⡷⠸⠀⠀⡃⡈⣶⡎⢶⣴⠇⡇⣿⡇⢸⣿⠀⠋⣴⡇⢸⠀⣿⢘⡅⡇⡇⢸
        ⠓⠀⠀⠐⠛⠐⠓⠈⠓⠒⠃⠛⠂⠀⠘⠃⠀⠀⠃⠀⠀⠓⠂⠓⠂⠃⠃⠈⠚⠀⠉⠚⠁⠙⠀⠘⠛⠀⠂⠉⠘⠈⠃⠈⠓⠐⠃⠃⠘

Usage of /root/go/bin/nier:
  -all
        Execute every type of check. If present, flags [rinj,sqlinj,subdomain] are enabled. If any of the flags [rinj,sqlinj,subdomain] are submitted while flag --all is submitted, they are silently ignored.
  -cve
        Enable Listing of CVEs related to banners discovered.
  -host string
        Identifies target host - i.e. 127.0.0.1 or www.myshop.com or http://myshop.com (default "127.0.0.1")
  -o string
        Output Folder PATH - in format: -o "~/Desktop/report" (default "/root/Desktop/Nier_Automaton_Report")
  -p int
        Target Port (default 80)
  -rinj
        Enable User Controlled Input Injection checking.
  -sess string
        Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2
  -sqlinj
        Enable SQL Injection checking. (SQLMap).
  -subdomain
        Enable Subdomain Enumeration.
  -test
        PoC scenario. i.e. Prioritize "testimonials" during injection detection. Just append "-test" or "--test" to the commandline.
```