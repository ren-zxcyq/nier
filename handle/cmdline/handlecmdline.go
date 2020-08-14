// Package handlecmdline processes command line variables and populates a cmdlineHandler struct
// The struct is populated and passed to main so that these features can be accessible in other
// parts of the project.
package handlecmdline

import (
	"flag"
	"fmt"
	"os"
	"path"
	// "path/filepath"
	"strings"

	"github.com/ren-zxcyq/nier/utilities"
)

type cmdlineHandler struct {
	InstallationDir      string
	ConfigFilePath       string
	C_OS                 string
	TargetHost           string
	TargetPort           int
	RunAll				 bool
	Ucinputinjection	 bool
	SubdomainEnumeration bool
	OutputFolder         string
	SessionTokens        string
	Tools                map[string]string //	Config File Contents:	map[tool] = location
	Test				 bool	//	PoC scenario. i.e. Prioritize "testimonials".
}

func NewCmdlineHandler() *cmdlineHandler {
	var h cmdlineHandler = cmdlineHandler{}
	h.PrintBanner()
	h.Tools = h.SetUpFlags()
	//fmt.Printf("Address of cmdlineHandler - %p", &h) //	Prints the address of cmdlineHandler
	return &h
}

var runallPointer = flag.Bool("all",false, "Execute every type of check. If present, flags [inj,subdomain] are enabled. If any of the flags [inj,subdomain] are submitted while flag --all is submitted, they are silently ignored.")
var targetHostPointer = flag.String("host", "127.0.0.1", "Identifies target host - i.e. 127.0.0.1 or www.myshop.com or http://myshop.com")
var targetPortPointer = flag.Int("p", 80, "Target Port")
var ucinputinjectionPointer = flag.Bool("inj",false, "Enable User Controlled Input Injection checking.")
var subdomainEnumerationPointer = flag.Bool("subdomain", false, "Enable Subdomain Enumeration.") ///Disable Subdomain Enumeration - Pass in [true or True] to enable (default false)")
var outputFolderPointer = flag.String("o", os.Getenv("HOME") + "/Desktop/Nier_Automaton_Report", "Output Folder PATH - in format: -o \"~/Desktop/report\"")
var sessionTokensPointer = flag.String("sess", "", "Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2")
var testPointer = flag.Bool("test", false, "PoC scenario. i.e. Prioritize \"testimonials\" during injection detection. Just append \"-test\" or \"--test\" to the commandline.")

func (h *cmdlineHandler) PrintBanner() {
	var banner string = "\r\n\t⣤⡄⠀⠀⣤⢠⢠⠀⠀⠀⠀⣤⠄⠀⢤⡀⠀⠀⠀⢀⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⢀⢀⣤⠀⠀⠀⠀⣒⣒"
	banner += "\r\n\t⡇⠙⢦⠀⣿⠀⡄⠀⣀⡀⠀⣿⠀⢀⡼⠃⢠⠀⠀⡘⢻⠀⣀⠀⢀⡀⣰⣀⠀⢀⣀⠀⠁⡀⠀⣀⠀⢨⣄⠠⠤⣤⠇⠿⣢⣀"
	banner += "\r\n\t⡇⠀⠈⠳⣿⠀⡇⡜⠀⢹⡆⣿⠚⠙⣆⠀⠀⠀⢀⢃⣸⡇⢸⠀⠀⡇⢸⠀⢰⠁⠈⡃⡄⣵⡆⢐⣖⠈⣷⡾⠇⢨⠀⣽⢡⡌⡆⡧⢺"
	banner += "\r\n\t⡇⠀⠀⠀⣿⠀⡇⣷⠊⠁⠀⣿⠀⠀⢹⡀⠐⠀⡘⠉⠀⡷⠸⠀⠀⡃⡈⣶⡎⢶⣴⠇⡇⣿⡇⢸⣿⠀⠋⣴⡇⢸⠀⣿⢘⡅⡇⡇⢸"
	banner += "\r\n\t⠓⠀⠀⠐⠛⠐⠓⠈⠓⠒⠃⠛⠂⠀⠘⠃⠀⠀⠃⠀⠀⠓⠂⠓⠂⠃⠃⠈⠚⠀⠉⠚⠁⠙⠀⠘⠛⠀⠂⠉⠘⠈⠃⠈⠓⠐⠃⠃⠘"
	banner += "\r\n\r\n"
	fmt.Printf("%s", banner)
}

func (h *cmdlineHandler) SetUpFlags() map[string]string {
	// fHandler := newFolderHandler(outputFolder)

	// var res bool
	// if res = fHandler.makeFolder(); res == false {
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

	//cOS = detectOS()
	var u utilities.Utils
	h.C_OS = u.DetectOS()
	// cwd, _ := os.Getwd()
	// h.InstallationDir = cwd
	// // configFilePath = path.Join(cwd, ".config")
	// h.ConfigFilePath = path.Join(cwd, ".config")

	// fmt.Println("YOOO", u.GetGOROOT())
	// fmt.Println("YAA", path.Join(u.GetGOPATH(), "src/github.com/ren-zxcyq/nier/"))

	h.InstallationDir = path.Join(u.GetGOPATH(), "src/github.com/ren-zxcyq/nier/")
	h.ConfigFilePath = path.Join(h.InstallationDir, ".config")

	//	Parse args	-	They return pointers
	flag.Parse() //	execute cmd-line parsing

	h.TargetHost = *targetHostPointer
	h.TargetPort = *targetPortPointer
	h.RunAll = *runallPointer
	if *runallPointer == true {
		h.RunAll = true
		h.Ucinputinjection = true
		h.SubdomainEnumeration = true
	} else {
		h.RunAll = *runallPointer
		h.Ucinputinjection = *ucinputinjectionPointer
		h.SubdomainEnumeration = *subdomainEnumerationPointer
	}

	//h.OutputFolder = *outputFolderPointer
	// h.OutputFolder = path.Join(cwd, h.OutputFolder)

	//	Determine report file path
	// tmpstr := *outputFolderPointer
	// fmt.Println("*outputFolderPointer:", *outputFolderPointer)
	//tmpPath, _ := filepath.Abs(tmpstr)	//*outputFolderPointer)
	h.OutputFolder = *outputFolderPointer	//tmpPath
	

	h.SessionTokens = *sessionTokensPointer

	h.Test = *testPointer	//	h.isFlagPassed("testPointer")

	// fmt.Println("TESTING ", h.Test)
	// os.Exit(1)
	
	//	Show args
	fmt.Println("\r\nSelected:", "\r\n-------------")
	fmt.Println("Installation Dir:", h.InstallationDir)
	fmt.Println("Loading Config:", h.ConfigFilePath)
	fmt.Println("Current OS:", h.C_OS)
	fmt.Println("Target Host:", *targetHostPointer)
	fmt.Println("Target Port:", *targetPortPointer)
	fmt.Println("Perform All Checks:", *runallPointer)
	fmt.Println("User Controlled Input Injection:", *ucinputinjectionPointer)
	fmt.Println("Subdomain Enumeration:", *subdomainEnumerationPointer)
	fmt.Println("Output Folder:", *outputFolderPointer)
	fmt.Println("Session Tokens:", *sessionTokensPointer)
	fmt.Println("Test:", *testPointer)



	fmt.Println("-------------")

	//	Print Contents of the Config File
	u.PrintFileContents(h.ConfigFilePath)

	return h.toolPaths()
}

// Reads & extracts - Tool Names & Locations
func (h *cmdlineHandler) toolPaths() map[string]string {
	// fmt.Println("\r\nExtracting Utility Location Information from .config\r\n-------------")
	var u utilities.Utils
	var ls []string = u.ReturnLinesFromFile(h.ConfigFilePath)
	//fmt.Println(ls)	//	[]

	var toolList map[string]string = make(map[string]string)
	var tool string
	var toolpath string
	var i int = 0
	for i < len(ls) {
		//fmt.Println(ls[i], " = becomes =>")
		tool, toolpath = forl(ls[i])
		toolList[tool] = toolpath
		i++
	}
	h.verifyTools(toolList) //	Exits on Error.
	return toolList
}

// Reads a line in format	[substring1 = substring2]
// Returns [substring1, substring2]
func forl(line string) (string, string) {
	var exp []string
	exp = strings.Split(line, "=")
	//fmt.Println("[", exp[0], ", ", exp[1], "]")
	exp[0] = strings.TrimSpace(exp[0])
	exp[1] = strings.TrimSpace(exp[1])
	return exp[0], exp[1]
}

// Iterates over a map of format.		map[tool] = location
// In case of error exits
func (h *cmdlineHandler) verifyTools(tList map[string]string) {
	// We want to stat a file, and continue if it exists :

	for k, v := range tList {

		_, err := os.Stat(v)

		if err != nil {
			if os.IsNotExist(err) {
				// file does not exist, do something
				fmt.Println("Error encountered:", k, "cannot be found at:", v)
				os.Exit(1)
			} else {
				// more serious errors
				fmt.Println("Error encountered while attempting to execute", k, "at", v, "\r\n", err)
				os.Exit(1)
			}
		}
	}
	// fmt.Printf("Verified that the files exist.\r\n-------------\r\n\r\n")
}

// func (h *cmdlineHandler) isFlagPassed(name string) bool {
//     found := false
//     flag.Visit(func(f *flag.Flag) {
//         if f.Name == name {
//             found = true
//         }
//     })
//     return found
// }