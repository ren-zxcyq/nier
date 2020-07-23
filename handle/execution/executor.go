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
	"path/filepath"
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
	var a utilities.Agent
	a.OptionsRequest(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))
	a.OptionsVerify(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))
	// fmt.Println("-------------")
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

	//	Ping
	var ping string = h.execCmd(h.e.tools["ping"] + " -" + pcount + " 1 " + h.e.targetHost)
	//fmt.Printf(ping)
	toolparser.ParsePing(ping)

	
	//	Web Server Fingerprinting

	//	Nmap
	var nmapOutFilesURL string = path.Join(h.e.outputFolder, "nmap_1_sSV")
	nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
	//nmapOutFilesURL = strings.Replace(nmapOutFilesURL, ":", "", -1)

	// var nmap string = h.execCmd(h.e.tools["nmap"] + " -sSV -T5 -oA " + nmapOutFilesURL + " " + h.e.targetHost)
	h.execCmd(h.e.tools["nmap"] + " -Pn -sSV -T5 -oA " + nmapOutFilesURL + " " + h.e.targetHost)
	/*	@PRV
	//	
	// toolparser.ParseNmapSV(nmap)
	// // fmt.Println(nmap)

	// fmt.Println("Initiating nmap --script=vuln scanning.","-------------")
	// nmapOutFilesURL string = path.Join(h.e.outputFolder, "/nmap-vuln.nmap")
	// nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
	// var nmapvuln string = h.execCmd(h.e.tools["nmap"] + " --script=vuln -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-vuln")) + " " + h.e.targetHost)
	*/
	//	nmap -p- -Pn -vv -sTV -T5 --script=banner -oA /root/Desktkop/nmap-banner- 192.168.1.20
	h.execCmd(h.e.tools["nmap"] + " -Pn -p- -vv -sTV -T5 --script=banner -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-banners")) + " " + h.e.targetHost)
	h.execCmd(h.e.tools["httprint"] + " -P0 -s /usr/share/httprint/signatures.txt -ox " + filepath.ToSlash(path.Join(h.e.outputFolder, "/httprint-srv-version")) + " -h " + h.e.targetHost)
	h.execCmd(h.e.tools["nmap"] + " -Pn --script=vuln -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-vuln")) + " " + h.e.targetHost)
	h.execCmd(h.e.tools["nmap"] + " -Pn -p" + strconv.Itoa(h.e.targetPort) + " --script=http-comments-displayer -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-comments")) + " " + h.e.targetHost)

	//	HTTP Methods
	h.checkHTTPMethods()
	/*
	//tester := utilities.NewHTTPShandler()
	//tester.TestHTTPS(h.e.targetHost)
	*/
	//var niktoOutFile string = path.Join(h.e.outputFolder, "nikto.txt")
	h.execCmd(h.e.tools["nikto"] + " -h " + h.e.targetHost + " -output " + filepath.ToSlash(path.Join(h.e.outputFolder, "nikto.txt")))


	/*	@RMV_COMMENT
	// THIS WORKS NORMALLY
	execInteractiveCmd("/root/go/bin/gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost + " -o " + gobusterFilesUrl)
	// @TODO	test with -o
	*/
	/*	@RMV_COMMENT
	// h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/wordlists/dirb/common.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + h.e.targetHost)
	*/
	//	h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -o "+ filepath.ToSlash(path.Join(h.e.outputFolder, "gobuster-URLs")) + " -u=" + h.e.targetHost)
	
	//	-s, --statuscodes string            Positive status codes (will be overwritten with statuscodesblacklist if set) (default "200,204,301,302,307,401,403")
	/*	@RMV_COMMENT
		var niktoOutFile string = path.Join(h.e.outputFolder, "nikto.txt")
		var nikto string = h.execCmd(h.e.tools["nikto"] + " -h " + h.e.targetHost + " -output " + niktoOutFile)
		// fmt.Printf(nikto)
		toolparser.ParseNikto(nikto)

		h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + h.e.targetHost)
		h.execInteractive(h.e.tools["sqlmap"] + " -u " + h.e.targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")
	*/
}
