package tooloutparse

/*
 *	Receive string objects from handleexec package. Extract features ?-> Act?
 *
 */

import (
	"fmt"
	"strings"
)

type toolparser struct{}

func NewToolparser() *toolparser {
	var h toolparser = toolparser{}
	return &h
}

func (h *toolparser) ParsePing(cmdout string) {

	if strings.Contains(cmdout, "1 packets transmitted, 1 received, 0% packet loss") {
		fmt.Println("Ping - OK")
		//	Did receive Reply
	} else {
		fmt.Println("Ping - FAIL")
		//	Did not receive Reply Host is unreachable
		//	@TODO	-	?Fail Gracefully?
	}
	fmt.Println(cmdout)
}

func (h *toolparser) ParseNmapSV(cmdout string) {
	var extract []string
	if strings.Contains(cmdout, "Host is up") {
		//	Nmap was successful.	-	Extract Features
		fmt.Println("NmapSV - OK")
		extract = strings.Split(cmdout, "ports")
		// fmt.Println(extract[1])
		extract = strings.Split(extract[1], "Service detection performed.")
		fmt.Println(extract[0])
		//	@TODO	-	Connect with Reporting

	} else {
		//	Nmap did not run smoothly. "Host is up" was not part of Stdout
		fmt.Println("NmapSV - FAIL")
		//	@TODO	-	?Fail Gracefully?
	}
	fmt.Println(cmdout)
}

func (h *toolparser) ParseNikto(cmdout string) {
	if strings.Contains(cmdout, "No web server found on") && strings.Contains(cmdout, "0 host(s) tested") {
		fmt.Println("Nikto - OK")
	} else {
		fmt.Println("Nikto - FAIL")
	}
	fmt.Println(cmdout)
}
