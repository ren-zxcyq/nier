package main

import (
	"flag"
	"fmt"
	. "fmt"
	"log"
	"os/exec" //	Launch SubProcess
	"runtime" //	Identify OS
	"strings"
)

var targetHost string
var targetPort int
var subdomainEnumeration bool
var sessionTokens string

//	Opens another program in go (os/exec etc): https://stackoverflow.com/a/37123000
//	@TODO	-	go doc os/exec.Cmd
func execCmdEx() {

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
	//out, err := exec.Command("ping", targetHost).Output()
	//	NMAP WORKS
	//	@TODO	Decide & Handle OS type
	//out, err := exec.Command("C:/Program Files (x86)/Nmap/nmap", "-T5", "-sSV", targetHost).Output()
	out, err := exec.Command("nmap", "-T5", "-sSV", targetHost).Output()
	if err != nil {
		log.Fatal(err)
	}
	Printf("Nmap output is \n%s", out)
	//log.Println("log")
}

// Trims trailing given character
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func execCmd(cmd string, arg ...string) string {
	/*	@EXAMPLE
		// var l int = len(arg)
		// var i int = 0
		var argstr string
		//for i; i < l; i++ {
		//	argstr += Fprintf("%s,", arg[i])
		//}
		//for _, i := range arg {

		fmt.Println("HIT")
		if len(arg) > 1 {
			argstr = "\"" + strings.Join(arg, "\", ")
		} else {
			argstr = arg[0]
		}

		fmt.Println(argstr)
		out, err := exec.Command(cmd, argstr).Output()
		if err != nil {
			log.Fatal(err)
		}
		Printf("\n%s output is: \n%s", cmd, out)
	*/
	Println("Here")
	var argstr []string
	var l int = len(arg)
	var i int = 0
	if l > 1 {
		for i < l {
			Println("l", arg[i])
			if i == l-1 {
				argstr = append(argstr, "\""+arg[i]+"\"")
			} else {
				argstr = append(argstr, "\""+arg[i]+"\", ")
			}
			//argstr += "\"" + arg[i] + "\", "
			i++
		}
		//		argstr = "\"" + strings.Join(arg, "\", ")

		// Trim trailing , from last item in the array
		//argstr = TrimSuffix(argstr, ",")
	} else {
		argstr = append(argstr, arg[0])
	}
	Println("Here")
	fmt.Println(argstr)

	out, err := exec.Command(cmd, argstr[0:]...).Output()
	if err != nil {
		Println("err")
		log.Fatal(err)
	}
	Println("Here")

	var res string = Sprintf("\n%s output is: \n-------------\n%s\n%s\n\n", cmd, out, err) //Sprintf() questionable
	Println("Here")

	//cmd.StdoutPipe()

	return res
}

func ex(cmd string) string {
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
	var targetHostPointer = flag.String("host", "127.0.0.1", "Identifies target host - i.e. 127.0.0.1 or www.myshop.com")
	var targetPortPointer = flag.Int("p", 80, "Target Port")
	var subdomainEnumerationPointer = flag.Bool("s", false, "Enable Subdomain Enumeration") ///Disable Subdomain Enumeration - Pass in [true or True] to enable (default false)")

	var sessionTokensPointer = flag.String("sess", "", "Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2")

	//	Parse args	-	They return pointers
	flag.Parse() //	execute cmd-line parsing

	//	Show args
	Println("Selected:", "\n-------------")
	Println("Current OS:", detectOS())
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

	var ping string = ex("ping -c 1 " + targetHost)
	Printf(ping)
	//Printf(nmap)
	//Printf(nikto)

	/*
		 *	@TODO	test if multiple flagsets can be used at a time

			nmapCmd := flag.NewFlagSet("nmap", flag.ExitOnError)
			nmapEnable := nmapCmd.Bool("enable", false, "enable")
	*/

	//	@TODO	Perform Checks on the flags
	//	@TODO	Assign them to program flags

}
