// Package handleexec is responsible for connecting and running tools needed and
// interface their execution with parsing.
package handleexec

import (
	"fmt"
	// "io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ren-zxcyq/nier/handle/cveretrieval"
	"github.com/ren-zxcyq/nier/handle/injdetect"
	"github.com/ren-zxcyq/nier/handle/spider"
	"github.com/ren-zxcyq/nier/handle/tooloutparse"
	"github.com/ren-zxcyq/nier/utilities"

	// "os/signal"
	// "syscall"
	// "time"
	"sync"

	// "github.com/creack/pty"
	// "golang.org/x/crypto/ssh/terminal"

	"bufio"
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
	runAll               bool
	ucinputInjection     bool
	sqlinjection		 bool
	subdomainEnumeration bool
	outputFolder         string
	cveRetrieval         bool
	sessionTokens        string
	tools                map[string]string
	test                 bool
}

var u utilities.Utils
var wg sync.WaitGroup

// Function NewExecHandler defines a new execHandler struct.
// execHandler.Exec() runs all the tools.
// execHandler is attached with all the internal unexported functions.
func NewExecHandler(installationDir, configPath, os, targetH string, targetP int, runAll, ucinputInjection, sqlinjection, subdomainEnum bool, outFolder string, cveRetrieval bool, sesTokens string, tL map[string]string, test bool) *execHandler {

	//	Create an elementsHandler Object to be passed to the exported execHandler
	var l elementsHandler = elementsHandler{
		installationDir:      installationDir,
		configFilePath:       configPath,
		cOS:                  os,
		targetHost:           targetH,
		targetPort:           targetP,
		runAll:               runAll,
		ucinputInjection:     ucinputInjection,
		sqlinjection:		  sqlinjection,
		subdomainEnumeration: subdomainEnum,
		outputFolder:         outFolder,
		cveRetrieval:         cveRetrieval,
		sessionTokens:        sesTokens,
		tools:                tL,
		test:                 test,
	}

	//	Create execHandler
	var h execHandler = execHandler{state: false, e: l}
	// u = utilities.Utils
	//fmt.Printf("Address of execHandler - %p", &h) //	Prints the address of outputFolderHandler
	return &h
}

// Opens another program in go (os/exec etc):					https://stackoverflow.com/a/37123000
// go doc os/exec.Cmd
// To be used for short subprocesses
func (h *execHandler) execCmd(cmd string) string {
	h.state = true
	var s []string = strings.Split(cmd, " ")

	out, err := exec.Command(s[0], s[1:]...).Output()
	if err != nil {
		//fmt.Printf("Err in ex", err.Error())
		log.Println(err.Error())
	}

	var res string = fmt.Sprintf("\r\n%s output is: \r\n-------------\r\n%s\r\n%s\n\n", cmd, out, err)

	h.state = false
	return res
}

func (h *execHandler) execSubP(filesprefix string, cmd string) {

	var s []string = strings.Split(cmd, " ")

	var res string = fmt.Sprintf("\n%s output is: \n-------------\n", s[0])
	fmt.Print(res)

	var prwll bool = false
	if strings.Contains(filesprefix,"sqlmap") {
		fmt.Println("SET TO TRUE")
		prwll = true
	}

	of, oferr := os.OpenFile(h.e.outputFolder+"/"+filesprefix+"_out",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if oferr != nil {
		fmt.Printf("Error Opening file for %s: %v", filesprefix, oferr)
	}

	defer of.Close()

	subprocess := exec.Command(s[0], s[1:]...)

	////////////////////////////////////////////////////////////////////////
	// inpipe, err := subprocess.StdinPipe()
	// if err != nil {
	// 	fmt.Println("ERROR WHILE OPENING STDIN", filesprefix)
	// }
	// inpipe = os.Stdin
	// defer inpipe.Close()
	
	subprocess.Stdin = os.Stdin
	////////////////////////////////////////////////////////////////////////

	outpipe, operr := subprocess.StdoutPipe()
	if operr != nil {
		fmt.Println("Could not open StdoutPipe")
	}

	if starterror := subprocess.Start(); starterror != nil {
		fmt.Println("Error while starting:", filesprefix)
	} 
	
	done := make(chan struct{})
	go func() {
		defer outpipe.Close()

		stdoutscanner := bufio.NewScanner(outpipe)

		for stdoutscanner.Scan() {
			var tmp string = stdoutscanner.Text()
			if prwll {
				fmt.Print(tmp+"\r\n")
			} else {
				fmt.Println(tmp)
			}
			of.WriteString(tmp + "\r\n")
		}
		// outpipe = io.MultiWriter(os.Stdout, of)

		done <- struct{}{}
	}()

	<-done
	fmt.Println("[*] - Wait was reached",filesprefix, "-", cmd)
	if err := subprocess.Wait(); err != nil {
		fmt.Println("Wait returned error:", err)
	}
	fmt.Println("[*] - DONE READING", cmd)

	//	Kill Sub-Process Group after a specific time
	// subprocess.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// start := time.Now()

	// var mins time.Duration = time.Duration(2) * time.Minute
	// if killaftermins > 0 {
	// 	mins = time.Duration(killaftermins) * time.Minute
	// }
	// time.AfterFunc(mins, func() {
	// 	if killaftermins > 0 {
	// 		syscall.Kill(-subprocess.Process.Pid, syscall.SIGKILL)
	// 		fmt.Println("Killed Subprocess:", subprocess.Process.Pid)
	// 	} else {
	// 		fmt.Println("Subprocess Left Running:", subprocess.Process.Pid)
	// 	}
	// })

	// // subprocess.Run()	//	change for Start() and Wait()
	// err := subprocess.Start()
	// if err != nil {
	// 	fmt.Println("Error starting subprocess", err)
	// }
	// <-done
	// err := subprocess.Wait()
	// if err != nil {
	// 	fmt.Println("Error while w8ing for subprocess,", err)
	// }

	// // fmt.Println("Waiting for ", subprocess.Process.Pid)
	// // subprocess.Wait() //	Wait for the Process to Exit

	// fmt.Printf("pid=%d duration=%s\n", subprocess.Process.Pid, time.Since(start))
	// // err := cmd.Run()
}

// 1) OPTIONS request
// 2) Verify Each Method
//
// Upgrade to HTTPS upon encountering HTTPS reply
func (h *execHandler) checkHTTPMethods() {
	//	NEEDS hCmd Assignments to happen before running.
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

		for _, opt := range optionsRes {
			// fmt.Println(opt)	//	Write this to file
			results += opt + "\r\n"
		}
	} else { //	@TODO	consider checking for another error
		// fmt.Println("[+]\tContinue Performing HTTP Testing")
		results = "-------------\r\nProtocol:\tHTTP\r\n-------------\r\n\r\n"
		results += "Options Request - Response:\r\n-------------\r\n" + or + "-------------\r\nHTTP Method - Status\r\n-------------\r\n"

		var optionsRes []string
		// fmt.Println(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))
		// os.Exit(0)
		optionsRes = t.OptionsVerify(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))

		for _, opt := range optionsRes {
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
	// var u utilities.Utils
	u.SaveStringToFile(h.e.outputFolder+"/httptesting.txt", results)

}

func (h *execHandler) getRobots() {
	var t utilities.Agent
	var r string
	r = "Retrieved over HTTP\r\n"
	r += t.Robots(h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))

	if strings.Contains(r, "server gave HTTP response to HTTPS client") { //"Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:") {
		tsec := utilities.NewHTTPShandler()
		r = "Retrieved over HTTPS\r\n"
		r += tsec.Robots(h.e.targetHost)
	}
	// fmt.Println(r)
	// var u utilities.Utils
	u.SaveStringToFile(h.e.outputFolder+"/getrobots.txt", r)
}

// Uses methods:	execCmd & execInteractive	essentially main execution happens here
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

	if pcount != "c" {
		fmt.Println("[*] This tool is currently designed to run on a *nix host.")
		os.Exit(1)
	}

	// var u utilities.Utils

	//	Ping
	var ping string = h.execCmd(h.e.tools["ping"] + " -" + pcount + " 1 " + u.Trimurlprefixhttp(h.e.targetHost))
	//fmt.Printf(ping)
	if !toolparser.ParsePing(ping) {
		fmt.Println("[*] Host Unreachable.")
		os.Exit(1)
	}

	h.runTools()
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
	var injectionhandler *injdetect.InjectionHandler = injdetect.NewInjectionHandler(h.e.targetHost, h.e.targetPort, h.e.installationDir, h.e.outputFolder, h.e.sessionTokens, h.e.test)
	injectionhandler.InjFormCheck()
}

func (h *execHandler) relativeUrlSpider() {
	var relurlspider *spider.RelativeLinkSpider = spider.NewRelativeLinkSpider(h.e.targetHost, h.e.targetPort, h.e.outputFolder, h.e.sessionTokens)
	// var results string =
	relurlspider.ReqURLs()
	// u.SaveStringToFile(h.e.outputFolder + "/links_gobuster_and_relspider.txt", results)

}

func (h *execHandler) appspider() {
	var appspider *spider.AppSpider = spider.NewAppSpider(h.e.targetHost, h.e.targetPort, h.e.outputFolder)
	appspider.Prepare()
	fmt.Println("\r\n\r\n[*]\tSpider Launched towards the Application\r\n")
	h.execSubP("gospider", h.e.tools["gospider"]+" -S "+filepath.ToSlash(path.Join(h.e.outputFolder, "/prespiderlinks.txt"))+" --depth 2 --no-redirect -t 50 -c 3 --cookie \""+h.e.sessionTokens+"\" --blacklist \"log\"")
	// also consider: -a, --other-source (i.e. find URLs from 3d party (Archive.org, CommonCrawl.org, VirusTotal.com))
	// h.execInteractive(h.e.tools["gospider"] + " -S " + filepath.ToSlash(path.Join(h.e.outputFolder, "/prespiderlinks.txt")) + " --depth 0 --no-redirect -t 50 -c 3 --cookie \"" + h.e.sessionTokens + "\" --blacklist -o " + filepath.ToSlash(path.Join(h.e.outputFolder, "/prespiderlinks.txt")))
	appspider.Organize()
}

// func xsstrike() runs xsstrike using the URL list built in injectiondection.go.InjFormCheck() & used during the XSS injection detection process.
func (h *execHandler) xsstrike() {

	//	Pops nano but needs to be tested.	--headers
	h.execSubP("xsstrike", h.e.tools["python3"]+" "+h.e.tools["xsstrike"]+" -u "+h.e.targetHost+":"+strconv.Itoa(h.e.targetPort)+" --crawl -t 10 --seeds "+h.e.outputFolder+"/urls_used_during_detection.txt") // + " --log-file " + h.e.outputFolder + "/xsstrike.txt")//	/root/testurls.txt")

}

func (h *execHandler) wpscan() {
	// @TODO	-	check if prefix can be enforced during cmdline config parsing
	var httprefix string
	if !strings.HasPrefix(h.e.targetHost, "http://") || !strings.HasPrefix(h.e.targetHost, "https://") {
		httprefix = "http://"
	}
	h.execSubP("wpscan", h.e.tools["wpscan"]+" --no-update -e --url "+httprefix+h.e.targetHost+":"+strconv.Itoa(h.e.targetPort))
}

func (h *execHandler) runTools() {

	//	Port Scan - Service Discovery
	var nmapOutFilesURL string = path.Join(h.e.outputFolder, "nmap_1_sSV")
	nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
	h.execSubP("nmap_1_sSV", h.e.tools["nmap"]+" -Pn -sSV -T5 -oA "+nmapOutFilesURL+" "+u.Trimurlprefixhttp(h.e.targetHost))
	h.execSubP("nmap-banners", h.e.tools["nmap"]+" -Pn -p- -vv -sTV -T5 --script=banner -oA "+filepath.ToSlash(path.Join(h.e.outputFolder, "nmap-banners"))+" "+u.Trimurlprefixhttp(h.e.targetHost))

	//	CVE Retrieval
	if h.e.cveRetrieval {
		cveretriever := cveretrieval.NewCVERetriever(h.e.outputFolder)
		cveretriever.Retrieve()
	}

	//	Web Server Fingerprinting
	h.execSubP("httprint", h.e.tools["httprint"]+" -P0 -s /usr/share/httprint/signatures.txt -ox "+filepath.ToSlash(path.Join(h.e.outputFolder, "httprint-srv-version"))+" -h "+h.e.targetHost)
	h.checkHTTPMethods()

	//	Application Comments
	h.getRobots()
	h.execSubP("nmap-comments", h.e.tools["nmap"]+" -Pn -p"+strconv.Itoa(h.e.targetPort)+" --script=http-comments-displayer -oA "+filepath.ToSlash(path.Join(h.e.outputFolder, "nmap-comments"))+" "+u.Trimurlprefixhttp(h.e.targetHost))

	//	Content Discovery
	//	Subdomain Enumeration
	if h.e.subdomainEnumeration {
		h.execSubP("gobuster_subd", h.e.tools["gobuster"]+" dns -d "+u.Trimurlprefixhttp(h.e.targetHost)+" -w /usr/share/amass/wordlists/subdomains-top1mil-5000.txt -o "+filepath.ToSlash(path.Join(h.e.outputFolder, "gobuster-Subdomains")))
	}

	//	Subdirectory & File Enumeration. Brute-Force & Crawl.
	// Rel & App Spiders need to run after gobuster directory discovery has been run. Otherwise execution fails.
	// rel spider reads file OUTPUTFOLDER/gobuster-URLs -> generates another file
	// which is read by appspider -> generating OUTPUTFOLDER/gospider_out
	// and additional filtered and extracted elements [forms, URLs, subdomains etc.] under OUTPUTFOLDER/gospider_URLs.list etc.
	h.execSubP("gobuster_dir", h.e.tools["gobuster"]+" dir -w /usr/share/wordlists/dirb/common.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -o "+filepath.ToSlash(path.Join(h.e.outputFolder, "gobuster-URLs"))+" -u="+u.Trimurlprefixhttp(h.e.targetHost))
	// h.execInteractive(h.e.tools["gobuster"] + " dir -w /usr/share/dirbuster/wordlists/directory-list-2.3-medium.txt -l -t 50 -x .php,.html,.ini,.py,.java,.sh,.js,.git -o "+ filepath.ToSlash(path.Join(h.e.outputFolder, "gobuster-URLs")) + " -u=" + h.e.targetHost)
	h.relativeUrlSpider()
	h.appspider()

	//	User Controll Input Reflection & XSS Injection
	if h.e.ucinputInjection {
		h.injectionTest()
		h.xsstrike()
	}

	//	Vulnerability Testing
	h.execSubP("nmap-vuln", h.e.tools["nmap"]+" -Pn --script=vuln -oA "+filepath.ToSlash(path.Join(h.e.outputFolder, "nmap-vuln"))+" "+u.Trimurlprefixhttp(h.e.targetHost))
	h.execSubP("nikto", h.e.tools["nikto"]+" -h "+h.e.targetHost+" -Tuning x 6 -output "+filepath.ToSlash(path.Join(h.e.outputFolder, "nikto.txt")))
	h.execSubP("whatweb", h.e.tools["whatweb"]+" -a4 -v "+h.e.targetHost+" --log-verbose "+filepath.ToSlash(path.Join(h.e.outputFolder, "whatweb-out.txt")))
	h.wpscan()

	if h.e.sqlinjection {
		h.execSubP("sqlmap", h.e.tools["sqlmap"]+" -u "+h.e.targetHost+" --forms --tamper=randomcase,space2comment --all --fresh-queries --answers=\"follow=Y\" --batch --disable-coloring --output-dir="+filepath.ToSlash(path.Join(h.e.outputFolder, "sqlmap_out")))
	}

	//	Screenshot URLs
	// ./gowitness file --source=/root/Desktop/urls.txt --threads=4 -d=/root/Desktop/report --disable-db report generate
	//	h.execInteractive(h.e.tools["gowitness"] + " file --source=" + h.e.outputFolder + "prespiderlinks.txt --threads=4 -d=" + h.e.outputFolder + " --disable-db report generate")
}
