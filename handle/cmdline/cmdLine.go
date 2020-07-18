package handleCmdLine

import (
	"flag"
	. "fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ren-zxcyq/nier/utilities"
)

const INSTALLATION string = "github.com/"

type cmdlineHandler struct {
	InstallationDir      string
	ConfigFilePath       string
	C_OS                 string
	TargetHost           string
	TargetPort           int
	SubdomainEnumeration bool
	OutputFolder         string
	SessionTokens        string
	Tools                map[string]string //	Config File Contents:	map[tool] = location
}

func NewCmdlineHandler() *cmdlineHandler {
	var h cmdlineHandler = cmdlineHandler{}
	h.PrintBanner()
	h.Tools = h.SetUpFlags()
	//Printf("Address of cmdlineHandler - %p", &h) //	Prints the address of cmdlineHandler
	return &h
}

var targetHostPointer = flag.String("host", "127.0.0.1", "Identifies target host - i.e. 127.0.0.1 or www.myshop.com")
var targetPortPointer = flag.Int("p", 80, "Target Port")
var subdomainEnumerationPointer = flag.Bool("s", false, "Enable Subdomain Enumeration") ///Disable Subdomain Enumeration - Pass in [true or True] to enable (default false)")
var outputFolderPointer = flag.String("o", "~/Desktop/Nier_Automaton_Report", "Output Folder PATH RELATIVE to cwd - in format: -o \"./report\"")
var sessionTokensPointer = flag.String("sess", "", "Session Token(s) - in format: -sess PHPSESSID:TOKEN1;JSESSID:TOKEN2")

func (h *cmdlineHandler) PrintBanner() {
	var banner string = "\r\n\t⣤⡄⠀⠀⣤⢠⢠⠀⠀⠀⠀⣤⠄⠀⢤⡀⠀⠀⠀⢀⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⢀⢀⣤⠀⠀⠀⠀⣒⣒"
	banner += "\r\n\t⡇⠙⢦⠀⣿⠀⡄⠀⣀⡀⠀⣿⠀⢀⡼⠃⢠⠀⠀⡘⢻⠀⣀⠀⢀⡀⣰⣀⠀⢀⣀⠀⠁⡀⠀⣀⠀⢨⣄⠠⠤⣤⠇⠿⣢⣀"
	banner += "\r\n\t⡇⠀⠈⠳⣿⠀⡇⡜⠀⢹⡆⣿⠚⠙⣆⠀⠀⠀⢀⢃⣸⡇⢸⠀⠀⡇⢸⠀⢰⠁⠈⡃⡄⣵⡆⢐⣖⠈⣷⡾⠇⢨⠀⣽⢡⡌⡆⡧⢺"
	banner += "\r\n\t⡇⠀⠀⠀⣿⠀⡇⣷⠊⠁⠀⣿⠀⠀⢹⡀⠐⠀⡘⠉⠀⡷⠸⠀⠀⡃⡈⣶⡎⢶⣴⠇⡇⣿⡇⢸⣿⠀⠋⣴⡇⢸⠀⣿⢘⡅⡇⡇⢸"
	banner += "\r\n\t⠓⠀⠀⠐⠛⠐⠓⠈⠓⠒⠃⠛⠂⠀⠘⠃⠀⠀⠃⠀⠀⠓⠂⠓⠂⠃⠃⠈⠚⠀⠉⠚⠁⠙⠀⠘⠛⠀⠂⠉⠘⠈⠃⠈⠓⠐⠃⠃⠘"
	banner += "\r\n\r\n"
	Printf("%s", banner)
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

	//	@STARTING HERE
	//cOS = detectOS()
	var u utilities.Utils
	h.C_OS = u.DetectOS()
	// cwd, _ := os.Getwd()
	// h.InstallationDir = cwd
	// // configFilePath = path.Join(cwd, ".config")
	// h.ConfigFilePath = path.Join(cwd, ".config")

	// Println("YOOO", u.GetGOROOT())
	// Println("YAA", path.Join(u.GetGOPATH(), "src/github.com/ren-zxcyq/nier/"))

	h.InstallationDir = path.Join(u.GetGOPATH(), "src/github.com/ren-zxcyq/nier/")
	h.ConfigFilePath = path.Join(h.InstallationDir, ".config")

	//	Parse args	-	They return pointers
	flag.Parse() //	execute cmd-line parsing

	//	Show args
	Println("\r\nSelected:", "\r\n-------------")
	Println("Installation Dir:", h.InstallationDir)
	Println("Loading Config:", h.ConfigFilePath)
	Println("Current OS:", h.C_OS)
	Println("targethost:", *targetHostPointer)
	Println("targetport:", *targetPortPointer)
	Println("subdomainEnumeration:", *subdomainEnumerationPointer)
	Println("outputFolder:", *outputFolderPointer)
	Println("sessionTokens:", *sessionTokensPointer)

	h.TargetHost = *targetHostPointer
	h.TargetPort = *targetPortPointer
	h.SubdomainEnumeration = *subdomainEnumerationPointer
	//h.OutputFolder = *outputFolderPointer
	//h.OutputFolder = path.Join(cwd, h.OutputFolder)

	//	Determine report file path
	//tmpstr := *outputFolderPointer
	tmpPath, _ := filepath.Abs(*outputFolderPointer)
	h.OutputFolder = tmpPath

	h.SessionTokens = *sessionTokensPointer
	Println("-------------")

	//	Print Contents of the Config File
	u.PrintFileContents(h.ConfigFilePath)

	return h.toolPaths()
}

/*
 *	Reads & extracts - Tool Names & Locations
 */
func (h *cmdlineHandler) toolPaths() map[string]string {
	// Println("\r\nExtracting Utility Location Information from .config\r\n-------------")
	var u utilities.Utils
	var ls []string = u.ReturnLinesFromFile(h.ConfigFilePath)
	//Println(ls)	//	[]

	var toolList map[string]string = make(map[string]string)
	var tool string
	var toolpath string
	var i int = 0
	for i < len(ls) {
		//Println(ls[i], " = becomes =>")
		tool, toolpath = forl(ls[i])
		toolList[tool] = toolpath
		i++
	}
	h.verifyTools(toolList) //	Exits on Error.
	return toolList
}

/*
 *	Reads a line in format	[substring1 = substring2]
 *	Returns [substring1, substring2]
 */
func forl(line string) (string, string) {
	var exp []string
	exp = strings.Split(line, "=")
	//Println("[", exp[0], ", ", exp[1], "]")
	exp[0] = strings.TrimSpace(exp[0])
	exp[1] = strings.TrimSpace(exp[1])
	return exp[0], exp[1]
}

/*
 *	Iterates over a map of format.		map[tool] = location
 *	In case of error exits
 */
func (h *cmdlineHandler) verifyTools(tList map[string]string) {
	// We want to stat a file, and continue if it exists :

	for k, v := range tList {

		_, err := os.Stat(v)

		if err != nil {
			if os.IsNotExist(err) {
				// file does not exist, do something
				Println("Error encountered:", k, "cannot be found at:", v)
				os.Exit(1)
			} else {
				// more serious errors
				Println("Error encountered while attempting to execute", k, "at", v, "\r\n", err)
				os.Exit(1)
			}
		}
	}
	// Printf("Verified that the files exist.\r\n-------------\r\n\r\n")
}
