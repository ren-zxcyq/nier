package main

import (
	"bufio"
	"flag"
	. "fmt"
	"log"
	"os"
	"os/exec" //	Launch SubProcess
	"runtime" //	Identify OS
	"strings"
	//"syscall"	//
)

var cOS string
var targetHost string
var targetPort int
var subdomainEnumeration bool
var sessionTokens string

//	Opens another program in go (os/exec etc): https://stackoverflow.com/a/37123000
//	@TODO	-	go doc os/exec.Cmd

func execCmd(cmd string) string {
	var s []string = strings.Split(cmd, " ")

	out, err := exec.Command(s[0], s[1:]...).Output()
	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("Err in ex")
	}

	var res string = Sprintf("\n%s output is: \n-------------\n%s\n%s\n\n", cmd, out, err) //Sprintf() questionable

	return res
}

//	This will probably be removed
//	Opens another program in go		-	Reading Std Output Stream
//	@TODO	Check	https://gobyexample.com/spawning-processes
func execCmdDontGoForThis(cmd string) {

	var s []string = strings.Split(cmd, " ")

	var res string = Sprintf("\n%s output is: \n-------------\n", s[0]) //Sprintf() questionable
	Print(res)

	//csubprocess := exec.Command("sqlmap", "-u 192.168.1.20/index.php", "--forms", "--tamper=randomcase,space2comment", "--all")
	csubprocess := exec.Command(s[0], s[1:]...)
	cstderr, cerr := csubprocess.StdoutPipe()
	if cerr = csubprocess.Start(); cerr != nil {
		log.Fatal("Err in cerr", cerr)
	}

	scanner := bufio.NewScanner(cstderr)
	//scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		if strings.Contains(m, "[") {
			Print("\n", m)
		} else {
			Print(m)
		}
	}

	csubprocess.Wait() //	Wait for the Process to Exit
}

//	Opens another program in go		-	Reading Std Output Stream
//
//	This:
//	https://github.com/kioopi/extedit/blob/master/extedit.go
//	Led to:
//	https://www.reddit.com/r/golang/comments/2nd4pq/how_can_i_open_an_interactive_subprogram_from/
func execInteractiveCmd(cmd string) {

	var s []string = strings.Split(cmd, " ")

	var res string = Sprintf("\n%s output is: \n-------------\n", s[0]) //Sprintf() questionable
	Print(res)

	//subprocess := exec.Command("sqlmap", "-u 192.168.1.20/index.php", "--forms", "--tamper=randomcase,space2comment", "--all")
	subprocess := exec.Command(s[0], s[1:]...)
	//stdout, suberr := subprocess.StdoutPipe()
	//stderr, suberrerr := subprocess.StderrPipe()

	subprocess.Stdin = os.Stdin
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr

	//	This works on Debian	=>	@TODO - Figure out how to - crossplatform terminate child processes
	// if cOS == "Linux" {
	// 	subprocess.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
	// 	// {
	// 	// 	cmd := exec.Command("/bin/sh", "-c", "watch date > date.txt")
	// 	// 	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// 	// 	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	// 	// }
	// }
	subprocess.Start()

	subprocess.Wait() //	Wait for the Process to Exit
}

//	This will probably be removed
func execSQLMap() {
	/*
		//	WRITES BUT IN A WEIRD WAY

		//csubprocess := exec.Command("sqlmap", "-u 192.168.1.20/index.php", "--forms", "--method=post", "--risk=3", "--level=5", "--tamper=randomcase,space2comment", "--dbs")
		csubprocess := exec.Command("sqlmap", "-u 192.168.1.20/index.php", "--forms", "--tamper=randomcase,space2comment", "--all")
		// Single Wrapped next to the Double Wrapped are the Write attempt
		// //csubprocess := exec.Command("cat")

		// // cstdin, cerr := csubprocess.StdinPipe()
		// // if cerr != nil {
		// // 	Println("cerr")
		// // 	log.Fatal("Err in ex")
		// // }
		cstderr, cerr := csubprocess.StdoutPipe()
		if cerr = csubprocess.Start(); cerr != nil {
			log.Fatal("Err in cerr", cerr)
		}

		scanner := bufio.NewScanner(cstderr)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			m := scanner.Text()
			if strings.Contains(m, "[") {
				Print("\n", m)
			} else {
				Print(m)
			}
		}
	*/

	csubprocess := exec.Command("sqlmap", "-u 192.168.1.20/index.php", "--forms", "--tamper=randomcase,space2comment", "--all")
	// Single Wrapped next to the Double Wrapped are the Write attempt
	// //csubprocess := exec.Command("cat")

	// // cstdin, cerr := csubprocess.StdinPipe()
	// // if cerr != nil {
	// // 	Println("cerr")
	// // 	log.Fatal("Err in ex")
	// // }
	cstderr, cerr := csubprocess.StdoutPipe()
	if cerr = csubprocess.Start(); cerr != nil {
		log.Fatal("Err in cerr", cerr)
	}

	scanner := bufio.NewScanner(cstderr)
	//scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		if strings.Contains(m, "[") {
			Print("\n", m)
		} else {
			Print(m)
		}
	}

	csubprocess.Wait()

	//cmd.Stdout = os.Stdout

	// //csubprocess.Stdout = os.Stdout
	// //csubprocess.Stderr = os.Stderr

	// // if cerr = csubprocess.Start(); cerr != nil {
	// // 	log.Fatal("Err in Start()", cerr)
	// // }

	// // io.WriteString(cstdin, "whoami\n")
	// // //csubprocess.Wait()
	// go func() {
	// 	defer cstdin.Close()
	// 	io.WriteString(cstdin, "values written to stdin are passed to cmd's std in")
	// }()

	// out, err := csubprocess.CombinedOutput()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Printf("%s", out)

	/*
		var res string = Sprintf("\n%s output is: \n-------------\n%s\n%s\n\n", c, cout, cerr) //Sprintf() questionable

		Printf(res)

		os.Stdin
		os.Stdout
	*/
	//var sqlmap string = execCmd("sqlmap --forms --crawl=2 " + targetHost)
	//var sqlmap string = execCmd("sqlmap -u 192.168.1.20/index.php --forms --method=post --risk 3 --level 5 --tamper=space2comment,randomcase --dbs")
	// https://stackoverflow.com/questions/23166468/how-can-i-get-stdin-to-exec-cmd-in-golang
	// Printf(sqlmap)
}

/*
	Identifies & Returns Host OS as per: https://golangbyexample.com/detect-os-golang/
	//	runtime.GOOS
	//	runtime.GOARCH
	//	for a full list of OS & ARCH	->		go tool dist list
*/
func detectOS() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return Sprintf("Windows")
	case "darwin":
		return Sprintf("Mac OS")
	case "linux":
		return Sprintf("Linux")
	default:
		return Sprintf("%s", os)
	}
}

/*	@TODO	check if it works correctly
 *	Check if file Exists	-	as per:	https://golangbyexample.com/check-if-file-or-directory-exists-go/
 *	Basically if it returns nil -> Everything is OK
 */
func checkFile(fileNamePath string) {
	//	of type os.FileInfo
	fileinfo, err := os.Stat(fileNamePath)
	if os.IsNotExist(err) {
		log.Fatal("Error while reading:", fileNamePath, ". File does not exist.")
	}
	log.Println(fileinfo)
	//Println(fileinfo)	//	Printing shows just <nil>
}

/*	@TODO	check if it works correctly
 *	Check if folder Exists	-	as per:	https://golangbyexample.com/check-if-file-or-directory-exists-go/
 */
func checkFolder(folderNamePath string) {
	//
	folderInfo, err := os.Stat(folderNamePath)
	if os.IsNotExist(err) {
		log.Fatal("Folder does not exist")
	}
	log.Println(folderInfo)
	//Println(folderInfo)	//	Printing shows just <nil>
}

func main() {

	//	Arg Definitions
	/*
		type Flag struct {
			Name     string // name as it appears on command line
			Usage    string // help message
			Value    Value  // value as set
			DefValue string // default value (as text); for usage message
		}
	*/

	/*
	 *	Example 1
	 */
	//	https://gobyexample.com/command-line-flags
	/*	@EXAMPLE
		//variable := flag.String("name", "defaultvalue", "This is the usage string") //flag.String() returns a pointer
		//
		// var svar string
		// flag.StringVar(&svar, "svar", "bar", "a string var")

		//	-h or help automatically generated
		//
	*/
	cOS = detectOS()
	var targetHostPointer = flag.String("host", "127.0.0.1", "Identifies target host - i.e. 127.0.0.1 or www.myshop.com")
	var targetPortPointer = flag.Int("p", 80, "Target Port")
	var subdomainEnumerationPointer = flag.Bool("s", false, "Enable Subdomain Enumeration") ///Disable Subdomain Enumeration - Pass in [true or True] to enable (default false)")

	var sessionTokensPointer = flag.String("sess", "", "Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2")

	//	Parse args	-	They return pointers
	flag.Parse() //	execute cmd-line parsing

	//	Show args
	Println("Selected:", "\n-------------")
	Println("Current OS:", cOS)
	Println("targethost:", *targetHostPointer)
	Println("targetport:", *targetPortPointer)
	Println("subdomainEnumeration:", *subdomainEnumerationPointer)
	Println("sessionTokens:", *sessionTokensPointer)

	targetHost = *targetHostPointer
	targetPort = *targetPortPointer
	subdomainEnumeration = *subdomainEnumerationPointer
	sessionTokens = *sessionTokensPointer
	Println("-------------")

	/*	@EXAMPLE
		textPtr := flag.String("text", "", "Text to parse. (Required)")
		if *textPtr == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	*/
	//execCmdEx()
	//var nmap string = execCmd("nmap", "-T5", "-sSV", targetHost)
	//var ping string = execCmd("ping", targetHost)
	//var nikto string = execCmd("nikto", "-h", targetHost)	//	Breaks when nikto or the requested tool is not installed

	//	Exec
	//	Adjust ping flag
	var pcount string
	if cOS == "Windows" {
		pcount = "n"
	} else if cOS == "Mac OS" {
		pcount = "c"
	} else if cOS == "Linux" {
		pcount = "c"
	} else {
		pcount = "c" //	If none of the 3 use the *nix variation
	}

	var example string = execCmd("asdf")
	Printf(example)

	var ping string = execCmd("ping -" + pcount + " 1 " + targetHost)
	Printf(ping)

	var nmap string = execCmd("nmap -sSV -T5 " + targetHost)
	Printf(nmap)
	//var nikto string = execCmd("nikto -h " + targetHost)
	//Printf(nikto)

	// var gobuster string = execCmd("./gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost)
	// // //	-c 'PHPSESSID:bp78fb8ser34n0hc83v3eu85n6; SecretCookie:VzuuL2gfLJWNnTSwn2kuLv5wo20vBwpjAGWwLJD2LwDkAJL0ZwplLmR5BQMuLGyuAGOuA2ZmBwR1Amt2Awp1ZmH%3D'
	// Printf(gobuster)
	//execCmd("./gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost)

	//THIS WORKS	@TODO: Consider gobuster for non-interactive
	execInteractiveCmd("./gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost)
	//execCmdInt("sqlmap -u " + targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")
	execInteractiveCmd("sqlmap -u " + targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")
	//Println("File & Folder Utilities:", "\n-------------")

	// go doc windows

	//	@EXAMPLE
	//os.Exit(0)	//	Can pass in - 1,2,3,4	-	DEFER'd actions won't be run

	/*
		 *	@TODO	test if multiple flagsets can be used at a time

			nmapCmd := flag.NewFlagSet("nmap", flag.ExitOnError)
			nmapEnable := nmapCmd.Bool("enable", false, "enable")
	*/

	//	@TODO	Perform Checks on the flags
	//	@TODO	Assign them to program flags
	//	@TODO	Add tools:	httprint, WPScan, WhatWeb, BlindElephant
	//	@TODO	gobuster
	//	@TODO	sqlmap
	//	@TODO	xxs
}
