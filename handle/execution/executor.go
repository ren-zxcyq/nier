// Package handleexec is responsible for connecting and running tools needed and
// interface their execution with parsing.
package handleexec

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	// "path/filepath"
	"strconv"
	"strings"

	"github.com/ren-zxcyq/nier/handle/tooloutparse"
	"github.com/ren-zxcyq/nier/utilities"
)

type execHandler struct {
	state bool            //	States whether it is running	1:running	0:not-running
	e     elementsHandler //	Receives Main.hCmd
}

type elementsHandler struct {
	installationDir      string
	configFilePath       string
	cOS                  string
	targetHost           string
	targetPort           int
	subdomainEnumeration bool
	outputFolder         string
	sessionTokens        string
	tools                map[string]string
}

//
// NewExecHandler defines a new execHandler struct.
// execHandler.Exec() runs all the tools.
// execHandler is attached with all the internal unexported functions.
//
func NewExecHandler(installationDir, configPath, os, targetH string, targetP int, subdomainEnum bool, outFolder, sesTokens string, tL map[string]string) *execHandler {

	//	Create an elementsHandler Object to be passed to the exported execHandler
	var l elementsHandler = elementsHandler{
		installationDir:      installationDir,
		configFilePath:       configPath,
		cOS:                  os,
		targetHost:           targetH,
		targetPort:           targetP,
		subdomainEnumeration: subdomainEnum,
		outputFolder:         outFolder,
		sessionTokens:        sesTokens,
		tools:                tL,
	}

	//	Create execHandler
	var h execHandler = execHandler{state: false, e: l}

	//fmt.Printf("Address of execHandler - %p", &h) //	Prints the address of outputFolderHandler
	return &h
}

/*
 *	Opens another program in go (os/exec etc): https://stackoverflow.com/a/37123000
 *	@TODO	-	go doc os/exec.Cmd
 */
func (h *execHandler) execCmd(cmd string) string {
	h.state = true
	var s []string = strings.Split(cmd, " ")

	out, err := exec.Command(s[0], s[1:]...).Output()
	if err != nil {
		//fmt.Printf("Err in ex", err.Error())
		log.Println(err.Error())
	}

	var res string = fmt.Sprintf("\r\n%s output is: \r\n-------------\r\n%s\r\n%s\n\n", cmd, out, err) //fmt.Sprintf() questionable

	h.state = false
	return res
}

/*
 *	Executes Subprocess interactively - Separates StdOut & StdErr in separate files - Just in case
 *	@TODO	-	verify for sqlmap-shell
 */
func (h *execHandler) execInteractive(cmd string) {

	h.state = true
	var s []string = strings.Split(cmd, " ")

	var res string = fmt.Sprintf("\n%s output is: \n-------------\n", s[0]) //fmt.Sprintf() questionable
	fmt.Print(res)

	//!was NOT commented OUT
	//f, err := os.OpenFile(outputFolder + "/" + s[0] + ".out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//fmt.Println("Creating FILE")
	var e []string = strings.Split(s[0], "/")
	var toolname string = e[len(e)-1]
	of, err := os.Create(h.e.outputFolder + "/" + toolname + "_out")
	ef, err := os.Create(h.e.outputFolder + "/" + toolname + "_err")
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer ef.Close()
	defer of.Close()

	//	?alt?
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
	h.state = false
}

/*
 *	Need to make sure that -host contains http://
 *
 *	1) OPTIONS request
 *	2) Verify Each Method
 *
 *	Upgrade to HTTPS upon encountering HTTPS reply
 *
 */
 func (h *execHandler) checkHTTPMethods() {
	//	NEEDS hCmd Assignments to happen before running.
	// host := "http://192.168.1.20"
	// port := 80
	// // var tar string = host + ":" + string(port)
	// var tar string = host + ":" + strconv.Itoa(port)
	// var u utilities.Utils
	// u.EncodingTest()
	// fmt.Println("-------------")
	// var a utilities.Agent
	// //a.Robots("http://www.google.com")
	// //a.Head("http://192.168.1.20")
	// //a.OptionsRequest("http://192.168.1.20")
	// //a.OptionsVerify("http://192.168.1.20")
	// a.OptionsVerify(tar)
	// fmt.Println("-------------")
	fmt.Println("\r\nInitiating HTTP Methods Checking\r\n-------------")
	var t utilities.Agent
	var results string
	var or string

	or = t.OptionsRequest(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))
	// fmt.Println(or)	//	Write this to file

	if strings.Contains(or, "Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking") {
		fmt.Println("[+]\tUpgrading to HTTPS")
		results = "-------------\r\nProtocol:\tHTTPS\r\n-------------\r\n\r\n"
		tsec := utilities.NewHTTPShandler()
		// var ors string
		or = tsec.OptionsRequest(h.e.targetHost)
		// fmt.Println("HTTPS Response Contents:\r\n", tsec.OptionsRequest(target))
		results += "Options Request - Response:\r\n-------------\r\n" + or + "-------------\r\nHTTPS Method - Status\r\n-------------\r\n"



		// fmt.Println(or)	//	Write this to file
		var optionsRes []string
		optionsRes = tsec.OptionsVerify(h.e.targetHost)

		for _,opt := range optionsRes {
			// fmt.Println(opt)	//	Write this to file
			results += opt + "\r\n"
		}
		// fmt.Println("HTTPS test\r\n",h.RequestMethodStatus("OPTIONS", target))
		// fmt.Println("_______________________", h.Robots(target))
		// fmt.Println("_______________________", h.Head(target))
		// fmt.Println("_______________________")
		// h.OptionsVerify(target)
		/*
		tester := utilities.NewHTTPShandler()
		tester.TestHTTPS(h.e.targetHost)	//	@TODO	consider if targetPort
		tester.Robots(h.e.targetHost)		//	@		should be added
		fmt.Println("[+]\tClosing HTTPS Testing")
		*/
	} else {								//	@TODO	consider checking for another error
		// fmt.Println("[+]\tContinue Performing HTTP Testing")
		results = "-------------\r\nProtocol:\tHTTP\r\n-------------\r\n\r\n"
		results += "Options Request - Response:\r\n-------------\r\n" + or + "-------------\r\nHTTP Method - Status\r\n-------------\r\n"

		var optionsRes []string
		// fmt.Println(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))
		// os.Exit(0)
		optionsRes = t.OptionsVerify(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))
		
		for _,opt := range optionsRes {
			// fmt.Println(opt)	//	Write this to file
			results += opt + "\r\n"
		}
	}
	// fmt.Println("-------------")
	// fmt.Println(or)
	// fmt.Println("-------------")
	// fmt.Println(results)
	// fmt.Println("-------------")

	//	Save to File
	var u utilities.Utils
	u.SaveStringToFile(h.e.outputFolder + "/httptesting.txt", results)

	// os.Exit(0)
}

func (h *execHandler) getRobots() {
	var t utilities.Agent
	var r string
	r = "Retrieved over HTTP\r\n"
	r += t.Robots(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))

	if strings.Contains(r,"server gave HTTP response to HTTPS client") {	//"Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:") {
		tsec := utilities.NewHTTPShandler()
		r = "Retrieved over HTTPS\r\n"
		r += tsec.Robots(h.e.targetHost)
	}
	// fmt.Println(r)
	var u utilities.Utils
	u.SaveStringToFile(h.e.outputFolder + "/getrobots.txt", r)
}

/*
 *	Uses methods:	execCmd & execInteractive	essentially main execution happens here
 */
func (h *execHandler) Exec() {

	toolparser := tooloutparse.NewToolparser()
	//	Adjust ping flag
	var pcount string
	if h.e.cOS == "Windows" {
		pcount = "n"
	} else if h.e.cOS == "Mac OS" {
		pcount = "c"
	} else if h.e.cOS == "Linux" {
		pcount = "c"
	} else {
		pcount = "c" //	If none of the 3 use the *nix variation
	}

	var u utilities.Utils

	//	Ping
	var ping string = h.execCmd(h.e.tools["ping"] + " -" + pcount + " 1 " + u.Trimurlsuffixhttp(h.e.targetHost))
	//fmt.Printf(ping)
	toolparser.ParsePing(ping)

	// h.runTools()
	h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/wordlists/dirb/common.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -o "+ "/root/Desktop/report/gobuster-URLs" + " -u=" + h.e.targetHost)
	h.injectionTest()
}

func (h *execHandler) injectionTest() {
	//	Procedure Outline:
		//	Get all URLs
		//	Filter for FORMS	<form
		//	Filter for params
		//	Generate Unique Items
		//	Submit
		//	Get all URLs
		//	Filter for Unique Items
	
	//	Get all URLs
	var injectionhandler *utilities.InjectionHandler = utilities.NewInjectionHandler()
	injectionhandler.InjURLsi()
}

func (h *execHandler) runTools() {
	nmapOutFilesURL := path.Join(h.e.outputFolder, "nmap_1_sSV")
	fmt.Println(nmapOutFilesURL)
	/*
	//	Web Server Fingerprinting
	var nmapOutFilesURL string = path.Join(h.e.outputFolder, "nmap_1_sSV")
	nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
	h.execCmd(h.e.tools["nmap"] + " -Pn -sSV -T5 -oA " + nmapOutFilesURL + " " + h.e.targetHost)
	h.execCmd(h.e.tools["nmap"] + " -Pn -p- -vv -sTV -T5 --script=banner -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-banners")) + " " + h.e.targetHost)
	h.execCmd(h.e.tools["httprint"] + " -P0 -s /usr/share/httprint/signatures.txt -ox " + filepath.ToSlash(path.Join(h.e.outputFolder, "/httprint-srv-version")) + " -h " + h.e.targetHost)
	h.checkHTTPMethods()
		
	//	Application Comments
	h.getRobots()
	h.execCmd(h.e.tools["nmap"] + " -Pn -p" + strconv.Itoa(h.e.targetPort) + " --script=http-comments-displayer -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-comments")) + " " + h.e.targetHost)
	//	@Uncoment
	h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/wordlists/dirb/common.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -o "+ filepath.ToSlash(path.Join(h.e.outputFolder, "gobuster-URLs")) + " -u=" + h.e.targetHost)
	// h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -o "+ filepath.ToSlash(path.Join(h.e.outputFolder, "gobuster-URLs")) + " -u=" + h.e.targetHost)
	
	//	Vulnerability Testing
	h.execCmd(h.e.tools["nmap"] + " -Pn --script=vuln -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-vuln")) + " " + h.e.targetHost)
	//var niktoOutFile string = path.Join(h.e.outputFolder, "nikto.txt")
	h.execCmd(h.e.tools["nikto"] + " -h " + h.e.targetHost + " -output " + filepath.ToSlash(path.Join(h.e.outputFolder, "nikto.txt")))
	
	//	@Uncoment
	// h.execInteractive(h.e.tools["sqlmap"] + " -u " + h.e.targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")
	*/
	
	/*	@PRV
		//	Nmap
		var nmapOutFilesURL string = path.Join(h.e.outputFolder, "nmap_1_sSV")
		nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
		//nmapOutFilesURL = strings.Replace(nmapOutFilesURL, ":", "", -1)
		
		// var nmap string = h.execCmd(h.e.tools["nmap"] + " -sSV -T5 -oA " + nmapOutFilesURL + " " + h.e.targetHost)
		h.execCmd(h.e.tools["nmap"] + " -Pn -sSV -T5 -oA " + nmapOutFilesURL + " " + h.e.targetHost)
		//	
		// toolparser.ParseNmapSV(nmap)
		// // fmt.Println(nmap)
		
		// fmt.Println("Initiating nmap --script=vuln scanning.","-------------")
		// nmapOutFilesURL string = path.Join(h.e.outputFolder, "/nmap-vuln.nmap")
		// nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
		// var nmapvuln string = h.execCmd(h.e.tools["nmap"] + " --script=vuln -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-vuln")) + " " + h.e.targetHost)
	*/
	/*
		execInteractiveCmd("/root/go/bin/gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost + " -o " + gobusterFilesUrl)
		//	-s, --statuscodes string            Positive status codes (will be overwritten with statuscodesblacklist if set) (default "200,204,301,302,307,401,403")
		// @TODO	test with -o
		// h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/wordlists/dirb/common.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + h.e.targetHost)
	*/
	/*
		var niktoOutFile string = path.Join(h.e.outputFolder, "nikto.txt")
		var nikto string = h.execCmd(h.e.tools["nikto"] + " -h " + h.e.targetHost + " -output " + niktoOutFile)
		// fmt.Printf(nikto)
		toolparser.ParseNikto(nikto)
	
		h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + h.e.targetHost)
		h.execInteractive(h.e.tools["sqlmap"] + " -u " + h.e.targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")
	*/
}