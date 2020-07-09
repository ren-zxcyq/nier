package handleExec

import (
	. "fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

	//Printf("Address of execHandler - %p", &h) //	Prints the address of outputFolderHandler
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
		//Printf("Err in ex", err.Error())
		log.Fatal(err.Error())
	}

	var res string = Sprintf("\r\n%s output is: \r\n-------------\r\n%s\r\n%s\n\n", cmd, out, err) //Sprintf() questionable

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

	var res string = Sprintf("\n%s output is: \n-------------\n", s[0]) //Sprintf() questionable
	Print(res)

	//!was NOT commented OUT
	//f, err := os.OpenFile(outputFolder + "/" + s[0] + ".out", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//Println("Creating FILE")
	var e []string = strings.Split(s[0], "/")
	var toolname string = e[len(e)-1]
	of, err := os.Create(h.e.outputFolder + "/" + toolname + "_out")
	ef, err := os.Create(h.e.outputFolder + "/" + toolname + "_err")
	if err != nil {
		Printf("error opening file: %v", err)
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
 *	Uses methods:	execCmd & execInteractive	essentially main execution happens here
 */
func (h *execHandler) Exec() {

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

	var ping string = h.execCmd(h.e.tools["ping"] + " -" + pcount + " 1 " + h.e.targetHost)
	Printf(ping)
	var nmapOutFilesUrl string = path.Join(h.e.outputFolder, "nmap_1_sSV")

	nmapOutFilesUrl = filepath.ToSlash(nmapOutFilesUrl)
	//nmapOutFilesUrl = strings.Replace(nmapOutFilesUrl, ":", "", -1)

	var nmap string = h.execCmd(h.e.tools["nmap"] + " -sSV -T5 -oA " + nmapOutFilesUrl + " " + h.e.targetHost)
	Println(nmap)

	var niktoOutFile string = path.Join(h.e.outputFolder, "nikto.txt")
	var nikto string = h.execCmd(h.e.tools["nikto"] + " -h " + h.e.targetHost + " -output " + niktoOutFile)
	Printf(nikto)

	//	example create file
	// var gobusterFileUrl string = path.Join(h.e.outputFolder, "gobuster.txt")
	// Println(gobusterFileUrl)
	//

	//	THIS WORKS NORMALLY
	//execInteractiveCmd("/root/go/bin/gobuster dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + targetHost + " -o " + gobusterFilesUrl)
	//	@TODO	test with -o
	h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -u=" + h.e.targetHost)
	h.execInteractive(h.e.tools["sqlmap"] + " -u " + h.e.targetHost + "/index.php --forms --tamper=randomcase,space2comment --all")

}
