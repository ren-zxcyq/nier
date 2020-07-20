package main

import (
	"fmt"
	// "strconv"

	//
	handlecmdline "github.com/ren-zxcyq/nier/handle/cmdline"
	handleexec "github.com/ren-zxcyq/nier/handle/execution"
	handlefolder "github.com/ren-zxcyq/nier/handle/folder"
	handlepdf "github.com/ren-zxcyq/nier/handle/report"
	// "github.com/ren-zxcyq/nier/utilities"
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
	fmt.Println("\r\nInitiating Document Creation Process\t-\t", hCmd.outputFolder, "\r\n-------------")
	//handlepdf.CreateDoc(outputFolder)
	fmt.Println(hCmd.installationDir, hCmd.outputFolder)
	handlepdf.CreatePdf(hCmd.installationDir, hCmd.outputFolder)
	fmt.Println("-------------")
}

func generateFolder() {
	fmt.Println("\r\nInitiating Folder Creation Process\t-\t", hCmd.outputFolder, "\r\n-------------")
	handlefolder.CreateFolder(hCmd.outputFolder)
	// fmt.Println("-------------")
}

func runTools() {
	fmt.Println("\r\nInitiating Exec\r\n-------------")
	ex := handleexec.NewExecHandler(hCmd.installationDir, hCmd.configFilePath, hCmd.cOS, hCmd.targetHost, hCmd.targetPort, hCmd.subdomainEnumeration, hCmd.outputFolder, hCmd.sessionTokens, hCmd.tools)
	ex.Exec()
	fmt.Println("-------------")
}

func main() {
	//	@Main
	cCmd := handlecmdline.NewCmdlineHandler()

	hCmd.installationDir = cCmd.InstallationDir
	hCmd.configFilePath = cCmd.ConfigFilePath
	hCmd.cOS = cCmd.C_OS
	hCmd.targetHost = cCmd.TargetHost
	hCmd.targetPort = cCmd.TargetPort
	hCmd.subdomainEnumeration = cCmd.SubdomainEnumeration
	hCmd.outputFolder = cCmd.OutputFolder
	hCmd.sessionTokens = cCmd.SessionTokens
	hCmd.tools = cCmd.Tools

	generateFolder()
	//exec
	runTools()
	generateReportFile()

	//testHttpMethods()
	// testPdf()
}

// /*
//  *	Need to make sure that -host contains http://
//  */
// func testHttpMethods() {
// 	//	NEEDS hCmd Assignments to happen before running.
// 	// host := "http://192.168.1.20"
// 	// port := 80
// 	// // var tar string = host + ":" + string(port)
// 	// var tar string = host + ":" + strconv.Itoa(port)
// 	// var u utilities.Utils
// 	// u.EncodingTest()
// 	// fmt.Println("-------------")
// 	// var a utilities.Agent
// 	// //a.Robots("http://www.google.com")
// 	// //a.Head("http://192.168.1.20")
// 	// //a.OptionsRequest("http://192.168.1.20")
// 	// //a.OptionsVerify("http://192.168.1.20")
// 	// a.OptionsVerify(tar)
// 	// fmt.Println("-------------")
// 	fmt.Println("\r\nInitiating HTTP Methods Checking\r\n-------------")
// 	var a utilities.Agent
// 	a.OptionsRequest(hCmd.targetHost + ":" + strconv.Itoa(hCmd.targetPort))
// 	a.OptionsVerify(hCmd.targetHost + ":" + strconv.Itoa(hCmd.targetPort))
// 	fmt.Println("-------------")
// }

func testPdf() {
	generateFolder()
	generateReportFile()
}
