package main

import (
	. "fmt"
	//
	handleCmdLine "github.com/ren-zxcyq/Nier/nier/handle/cmdline"
	handleExec "github.com/ren-zxcyq/Nier/nier/handle/execution"
	handleFolder "github.com/ren-zxcyq/Nier/nier/handle/folder"
	handlePdf "github.com/ren-zxcyq/Nier/nier/handle/report"
)

type elementsHandler struct {
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
	handlePdf.CreatePdf(hCmd.outputFolder)
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
	ex := handleExec.NewExecHandler(hCmd.configFilePath, hCmd.cOS, hCmd.targetHost, hCmd.targetPort, hCmd.subdomainEnumeration, hCmd.outputFolder, hCmd.sessionTokens, hCmd.tools)
	ex.Exec()
	Println("-------------")
}

func main() {

	cCmd := handleCmdLine.NewCmdlineHandler()

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
}
