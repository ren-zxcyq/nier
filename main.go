package main

import (
	"flag"
	. "fmt"
	"log"
	"os"
	"os/exec" //	Launch SubProcess
	"runtime" //	Identify OS
	"strings"
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
		Println("err")
		log.Fatal("Err in ex")
	}

	var res string = Sprintf("\n%s output is: \n-------------\n%s\n%s\n\n", cmd, out, err) //Sprintf() questionable

	return res
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
	//Printf(nmap)
	//Printf(nikto)

	Println("File & Folder Utilities:", "\n-------------")
	checkFile("/c:/Users/blush/Desktop/tools_509_web.txt") //	RETURNS SAME
	checkFile("/c:/Users/blush/Desktop/tools_.txt")        //	AS THIS
	checkFolder("/c:/Users/blush/Desktop/int")

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
