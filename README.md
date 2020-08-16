# nier

## Installation
```
go get github.com/ren-zxcyq/nier
go install github.com/ren-zxcyq/nier
```

## Run

If run as user other than root:
```
                                    (In debian -E specifies using the same ENV structure)
                                    (depending on whether u built the binary or not)

sudo -E ~/go/bin/nier -host [TARGET]
sudo -E ~/go/src/github.com/ren-zxcyq/nier/main.go -host [TARGET]
```

If run as root:
```
~/go/bin/nier -host [TARGET]
~/go/src/github.com/ren-zxcyq/nier/main.go -host [TARGET]
```


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
        Execute every type of check. If present, flags [inj,subdomain] are enabled. If any of the flags [inj,subdomain] are submitted while flag --all is submitted, they are silently ignored.
  -cve
        Enable Listing of CVEs related to banners discovered.
  -host string
        Identifies target host - i.e. 127.0.0.1 or www.myshop.com or http://myshop.com (default "127.0.0.1")
  -inj
        Enable User Controlled Input Injection checking.
  -o string
        Output Folder PATH - in format: -o "~/Desktop/report" (default "/root/Desktop/Nier_Automaton_Report")
  -p int
        Target Port (default 80)
  -sess string
        Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2
  -subdomain
        Enable Subdomain Enumeration.
  -test
        PoC scenario. i.e. Prioritize "testimonials" during injection detection. Just append "-test" or "--test" to the commandline.
```