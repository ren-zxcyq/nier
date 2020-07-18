# nier

## Installation
```
go get github.com/ren-zxcyq/nier
cd into the above directory and run build & install
go build
go install
```

## Run
```
If run as user other than root:      (In debian -E specifies using the same ENV structure)
depending on whether u built the binary or not:
sudo -E ~/go/bin/nier -host [TARGET]
sudo -E ~/go/src/github.com/ren-zxcyq/nier/main.go -host [TARGET]

If run as root:
~/go/bin/nier -host [TARGET]
~/go/src/github.com/ren-zxcyq/nier/main.go -host [TARGET]
```


## Usage
```
~/go/src/github.com/ren-zxcyq/nier$ ~/go/bin/nier -h

        ⣤⡄⠀⠀⣤⢠⢠⠀⠀⠀⠀⣤⠄⠀⢤⡀⠀⠀⠀⢀⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⢀⢀⣤⠀⠀⠀⠀⣒⣒
        ⡇⠙⢦⠀⣿⠀⡄⠀⣀⡀⠀⣿⠀⢀⡼⠃⢠⠀⠀⡘⢻⠀⣀⠀⢀⡀⣰⣀⠀⢀⣀⠀⠁⡀⠀⣀⠀⢨⣄⠠⠤⣤⠇⠿⣢⣀
        ⡇⠀⠈⠳⣿⠀⡇⡜⠀⢹⡆⣿⠚⠙⣆⠀⠀⠀⢀⢃⣸⡇⢸⠀⠀⡇⢸⠀⢰⠁⠈⡃⡄⣵⡆⢐⣖⠈⣷⡾⠇⢨⠀⣽⢡⡌⡆⡧⢺
        ⡇⠀⠀⠀⣿⠀⡇⣷⠊⠁⠀⣿⠀⠀⢹⡀⠐⠀⡘⠉⠀⡷⠸⠀⠀⡃⡈⣶⡎⢶⣴⠇⡇⣿⡇⢸⣿⠀⠋⣴⡇⢸⠀⣿⢘⡅⡇⡇⢸
        ⠓⠀⠀⠐⠛⠐⠓⠈⠓⠒⠃⠛⠂⠀⠘⠃⠀⠀⠃⠀⠀⠓⠂⠓⠂⠃⠃⠈⠚⠀⠉⠚⠁⠙⠀⠘⠛⠀⠂⠉⠘⠈⠃⠈⠓⠐⠃⠃⠘

Usage of ~/go/bin/nier:
  -host string
        Identifies target host - i.e. 127.0.0.1 or www.myshop.com (default "127.0.0.1")
  -o string
        Output Folder PATH RELATIVE to cwd - in format: -o "./report" (default "~/Desktop/Nier_Automaton_Report")
  -p int
        Target Port (default 80)
  -s    Enable Subdomain Enumeration
  -sess string
        Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2
```