package main

import (
	. "fmt"
	"strconv"

	//
	handleCmdLine "github.com/ren-zxcyq/Nier/nier/handle/cmdline"
	handleExec "github.com/ren-zxcyq/Nier/nier/handle/execution"
	handleFolder "github.com/ren-zxcyq/Nier/nier/handle/folder"
	handlePdf "github.com/ren-zxcyq/Nier/nier/handle/report"
	"github.com/ren-zxcyq/Nier/nier/utilities"
)

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

var hCmd elementsHandler

func generateReportFile() {
	Println("\r\nInitiating Document Creation Process\r\n-------------")
	//handlePdf.CreateDoc(outputFolder)
	handlePdf.CreatePdf(hCmd.installationDir, hCmd.outputFolder)
	Println("-------------")
}

func generateFolder() {
	Println("\r\nInitiating Folder Creation Process\r\n-------------")
	Println(hCmd.outputFolder)
	handleFolder.CreateFolder(hCmd.outputFolder)
	Println("-------------")
}

func runTools() {
	Println("\r\nInitiating Exec\r\n-------------")
	ex := handleExec.NewExecHandler(hCmd.installationDir, hCmd.configFilePath, hCmd.cOS, hCmd.targetHost, hCmd.targetPort, hCmd.subdomainEnumeration, hCmd.outputFolder, hCmd.sessionTokens, hCmd.tools)
	ex.Exec()
	Println("-------------")
}

func main() {
	//	@Main
	cCmd := handleCmdLine.NewCmdlineHandler()

	hCmd.installationDir = cCmd.InstallationDir
	hCmd.configFilePath = cCmd.ConfigFilePath
	hCmd.cOS = cCmd.C_OS
	hCmd.targetHost = cCmd.TargetHost
	hCmd.targetPort = cCmd.TargetPort
	hCmd.subdomainEnumeration = cCmd.SubdomainEnumeration
	hCmd.outputFolder = cCmd.OutputFolder
	hCmd.sessionTokens = cCmd.SessionTokens
	hCmd.tools = cCmd.Tools

	/*generateFolder()
	//exec
	runTools()
	generateReportFile()*/

	//testHttpMethods()
	testPdf()
}

/*
 *	Need to make sure that -host contains http://
 */
func testHttpMethods() {
	//	NEEDS hCmd Assignments to happen before running.
	// host := "http://192.168.1.20"
	// port := 80
	// // var tar string = host + ":" + string(port)
	// var tar string = host + ":" + strconv.Itoa(port)
	// var u utilities.Utils
	// u.EncodingTest()
	// Println("-------------")
	// var a utilities.Agent
	// //a.Robots("http://www.google.com")
	// //a.Head("http://192.168.1.20")
	// //a.OptionsRequest("http://192.168.1.20")
	// //a.OptionsVerify("http://192.168.1.20")
	// a.OptionsVerify(tar)
	// Println("-------------")
	Println("\r\nInitiating HTTP Methods Checking\r\n-------------")
	var a utilities.Agent
	a.OptionsRequest(hCmd.targetHost + ":" + strconv.Itoa(hCmd.targetPort))
	a.OptionsVerify(hCmd.targetHost + ":" + strconv.Itoa(hCmd.targetPort))
	Println("-------------")
}

func testPdf() {
	generateFolder()
	generateReportFile()
}

