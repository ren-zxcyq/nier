package main

import (
	"flag"
	. "fmt"
	"log"
	"os/exec"
)

var targetHost string
var targetPort int
var subdomainEnumeration bool

//	Opens another program in go (os/exec etc): https://stackoverflow.com/a/37123000
//	@TODO	-	go doc os/exec.Cmd
func execCmd() {

	/*	@EXAMPLE
		//	This works but once thre process is Run() -	Need to w8 for the process to end
		//		https://stackoverflow.com/a/15815730
		//
		cmd := exec.Command("nmap", "-sSV -T5 127.0.0.1")
		cmd.Run() // and wait
		//cmd.Start()
		stdout, err := cmd.Output()
		if err != nil {
			Println(err.Error())
			return
		}

		Print(string(stdout))
	*/

	/*	@EXAMPLE
		//	https://stackoverflow.com/a/15815730
		//
		//		//	where out =>	stdout
		//
		out, err := exec.Command("date").Output()
		if err != nil {
			log.Fatal(err)
		}
		Printf("The date is %s\n", out)

	*/
	//	PING WORKS
	//out, err := exec.Command("ping", "www.google.com").Output()
	out, err := exec.Command("ping", targetHost).Output()
	//	NMAP WORKS
	//	@TODO	Decide & Handle OS type
	//out, err := exec.Command("C:/Program Files (x86)/Nmap/nmap", "-T5", "-sSV", targetHost).Output()
	if err != nil {
		log.Fatal(err)
	}
	Printf("Nmap output is \n%s", out)
	//log.Println("log")
}

func main() {
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

	var targetHostPointer = flag.String("host", "127.0.0.1", "Identifies target host - i.e. 127.0.0.1 or www.myshop.com")
	var targetPortPointer = flag.Int("p", 80, "Target Port")
	var subdomainEnumerationPointer = flag.Bool("s", false, "Enable Subdomain Enumeration") ///Disable Subdomain Enumeration - Pass in [true or True] to enable (default false)")

	flag.Parse() //	execute cmd-line parsing

	Println("targethost:", *targetHostPointer)
	Println("targetport:", *targetPortPointer)
	Println("subdomainEnumeration:", *subdomainEnumerationPointer)

	targetHost = *targetHostPointer
	targetPort = *targetPortPointer
	subdomainEnumeration = *subdomainEnumerationPointer

	/*	@EXAMPLE
		textPtr := flag.String("text", "", "Text to parse. (Required)")
		if *textPtr == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	*/
	execCmd()

	/*
		 *	@TODO	test if multiple flagsets can be used at a time

			nmapCmd := flag.NewFlagSet("nmap", flag.ExitOnError)
			nmapEnable := nmapCmd.Bool("enable", false, "enable")
	*/

	//	@TODO	Perform Checks on the flags
	//	@TODO	Assign them to program flags

}
