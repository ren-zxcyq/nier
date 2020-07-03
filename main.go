package main

import (
	"bufio"
	"flag"
	. "fmt"
	"io"
	"log"
	"os"
	"os/exec" //	Launch SubProcess
	"path"
	"path/filepath"
	"runtime" //	Identify OS
	"strings"
	"syscall" //

	"github.com/ren-zxcyq/Nier/nier/handleFolder"
	"github.com/ren-zxcyq/Nier/nier/handlePdf"
)

var cOS string
var targetHost string
var targetPort int
var subdomainEnumeration bool
var outputFolder string
var sessionTokens string

/*
 *	Opens another program in go (os/exec etc): https://stackoverflow.com/a/37123000
 *	@TODO	-	go doc os/exec.Cmd
 */
func execCmd(cmd string) string {
	var s []string = strings.Split(cmd, " ")

	out, err := exec.Command(s[0], s[1:]...).Output()
	if err != nil {
		//Printf("Err in ex", err.Error())
		log.Fatal(err.Error())
	}

	var res string = Sprintf("\n%s output is: \n-------------\n%s\n%s\n\n", cmd, out, err) //Sprintf() questionable

	return res
}

/*
 *	Executes Subprocess interactively - Separates StdOut & StdErr in separate files - Just in case
 *	@TODO	-	verify for sqlmap-shell
 */
func execInteractive(cmd string) {

	var s []string = strings.Split(cmd, " ")

	var res string = Sprintf("\n%s output is: \n-------------\n", s[0]) //Sprintf() questionable
	Print(res)

	//!was NOT commented OUT
	//f, err := os.OpenFile(outputFolder + "/" + s[0] + ".out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Println("Creating FILE")
	var e []string = strings.Split(s[0], "/")
	var toolname string = e[len(e)-1]
	of, err := os.Create(outputFolder + "/" + toolname + "_out")
	ef, err := os.Create(outputFolder + "/" + toolname + "_err")
	if err != nil {
		Printf("error opening file: %v", err)
	}
	defer ef.Close()
	defer of.Close()

	//	FROM STACKOVERFLOW
	//f, _ := os.Create("file")
	//cmd.Stdout = io.MultiWriter(os.Stdout, f)

	subprocess := exec.Command(s[0], s[1:]...)

	subprocess.Stdin = os.Stdin

	//!was NOT commented out
	// redirect output to files
	//subprocess.Stdout = f
	//subprocess.Stderr = f

	subprocess.Stdout = io.MultiWriter(os.Stdout, of)
	subprocess.Stderr = io.MultiWriter(ef)

	subprocess.Start()

	subprocess.Wait() //	Wait for the Process to Exit

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
	Println("INTERACTIVE CMD")
	var s []string = strings.Split(cmd, " ")

	var res string = Sprintf("\n%s output is: \n-------------\n", s[0]) //Sprintf() questionable
	Print(res)
	Println("SHOULLDA OUTPUT")
	//subprocess := exec.Command("sqlmap", "-u 192.168.1.20/index.php", "--forms", "--tamper=randomcase,space2comment", "--all")
	subprocess := exec.Command(s[0], s[1:]...)
	//stdout, suberr := subprocess.StdoutPipe()
	//stderr, suberrerr := subprocess.StderrPipe()

	subprocess.Stdin = os.Stdin
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr
	Println("JUST ASSIGNED subprocess.Stdin = os.Stdin")
	//	This works on Debian	=>	@TODO - Figure out how to - crossplatform terminate child processes
	subprocess.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}

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

func setUpFlags() {
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
	var outputFolderPointer = flag.String("o", "Nier_Automata_Report", "Output Folder PATH RELATIVE to cwd - in format: -o \"./report\"")
	var sessionTokensPointer = flag.String("sess", "", "Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2")

	//	Parse args	-	They return pointers
	flag.Parse() //	execute cmd-line parsing

	//	Show args
	Println("Selected:", "\n-------------")
	Println("Current OS:", cOS)
	Println("targethost:", *targetHostPointer)
	Println("targetport:", *targetPortPointer)
	Println("subdomainEnumeration:", *subdomainEnumerationPointer)
	Println("outputFolder:", *outputFolderPointer)
	Println("sessionTokens:", *sessionTokensPointer)

	targetHost = *targetHostPointer
	targetPort = *targetPortPointer
	subdomainEnumeration = *subdomainEnumerationPointer
	outputFolder = *outputFolderPointer
	cwd, _ := os.Getwd()
	outputFolder = path.Join(cwd, outputFolder)
	sessionTokens = *sessionTokensPointer
	Println("-------------")

	/*	@EXAMPLE
		textPtr := flag.String("text", "", "Text to parse. (Required)")
		if *textPtr == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	*/
}

/*
 *	@TODO - MUST CLEANUP
 */
func runTools() {

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

	var ping string = execCmd("ping -" + pcount + " 1 " + targetHost)
	Printf(ping)

	var nmap string = execCmd("nmap -sSV -T5 -oA" + outputFolder + " " + targetHost)
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
	//	@TODO	Add tools:	httprint, WPScan, WhatWeb, BlindElephant, Wapiti
	//	@TODO	gobuster
	//	@TODO	sqlmap
	//	@TODO	xxs	-	XSSSniper (Don't worry about if it works now)
}

func generateReportFile() {
	Println("Initiating Document Creation Process")
	//handlePdf.CreateDoc(outputFolder)
	handlePdf.CreatePdf(outputFolder)
}

func generateFolder() {
	Println("Initiating Document Creation Process")

	handleFolder.CreateFolder(outputFolder)

}

func test() {

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

	var ping string = execCmd("ping -" + pcount + " 1 " + targetHost)
	Printf(ping)
	var nmapOutFilesUrl string = path.Join(outputFolder, "nmap_1_sSV")

	nmapOutFilesUrl = filepath.ToSlash(nmapOutFilesUrl)
	//nmapOutFilesUrl = strings.Replace(nmapOutFilesUrl, ":", "", -1)

	var nmap string = execCmd("nmap -sSV -T5 -oA " + nmapOutFilesUrl + " " + targetHost)
	Println(nmap)

	var niktoOutFile string = path.Join(outputFolder, "nikto.txt")
	var nikto string = execCmd("nikto -h " + targetHost + " -output " + niktoOutFile)
	Printf(nikto)

	var gobusterFileUrl string = path.Join(outputFolder, "gobuster.txt")
	Println(gobusterFileUrl)
	//

	//	THIS WORKS NORMALLY
	//execInteractiveCmd("/root/go/bin/gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost + " -o " + gobusterFilesUrl)
	//	PEEEERFEEECT	@TODO	test with -o
	execInteractive("/root/go/bin/gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost)
	execInteractive("sqlmap -u " + targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")

	//	?Alt?
	//execCmd("/root/go/bin/gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php -u=" + targetHost + " | tee "+ gobusterFilesUrl)
	//var sqlmapFileUrl string = path.Join(outputFolder, "sqlmap.txt")
	//execInteractiveCmd("sqlmap -u " + targetHost + "/index.php --forms --tamper=randomcase,space2comment --all --output-dir=" + outputFolder)
	//	?tee?
	//execInteractiveCmd("sqlmap -u " + targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")	// 2>&1 | tee " + sqlmapFileUrl)
}

func main() {
	setUpFlags()
	generateFolder()
	generateReportFile()
	test()
}
