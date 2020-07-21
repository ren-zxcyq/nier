package tooloutparse

/*
 *	Receive string objects from handleexec package. Extract features ?-> Act?
 *
 */

import (
	"fmt"
	"strings"
	"github.com/ren-zxcyq/nier/utilities"
	"log"
)

var u utilities.Utils

type Toolparser struct{}

func NewToolparser() *Toolparser {
	var h Toolparser = Toolparser{}
	return &h
}

func (h *Toolparser) ParsePing(cmdout string) {

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

func (h *Toolparser) ParseNmapSV(cmdout string) string {
	var extract []string
	if strings.Contains(cmdout, "Host is up") {
		// //	Nmap was successful.	-	Extract Features
		// fmt.Println("NmapSV - OK")
		extract = strings.Split(cmdout, "ports")
		// // fmt.Println(extract[1])
		extract = strings.Split(extract[1], "Service detection performed.")
		// // fmt.Println(extract[0])
		// //	@TODO	-	Connect with Reporting
		return extract[0]
	} else {
		//	Nmap did not run smoothly. "Host is up" was not part of Stdout
		// fmt.Println("NmapSV - FAIL")
		//	@TODO	-	?Fail Gracefully?
		return fmt.Sprintln("NmapSV - FAIL")
	}
	// fmt.Println(cmdout)
}

func (h *Toolparser) ParseNmapVuln(cmdout string) string {
	var extract []string
	if strings.Contains(cmdout, "Host is up") {
		//	Nmap was successful.	-	Extract Features
		// fmt.Println("NmapVuln - OK")
		extract = strings.Split(cmdout, "ports")
		// // fmt.Println(extract[1])
		// extract = strings.Split(extract[1], "Service detection performed.")
		// // fmt.Println(extract[0])
		// //	@TODO	-	Connect with Reporting
		return extract[1]
	} else {
		//	Nmap did not run smoothly. "Host is up" was not part of Stdout
		// fmt.Println("NmapSV - FAIL")
		//	@TODO	-	?Fail Gracefully?
		return fmt.Sprintln("NmapVuln - FAIL")
	}
	// fmt.Println(cmdout)
}


func (h *Toolparser) ParseNikto(cmdout string) {
	if strings.Contains(cmdout, "No web server found on") && strings.Contains(cmdout, "0 host(s) tested") {
		fmt.Println("Nikto - OK")
	} else {
		fmt.Println("Nikto - FAIL")
	}
	fmt.Println(cmdout)
}

//	ParseGobuster Filters out results that are of (Status: 403)
//	Returns an array of lines.
func (h *Toolparser) ParseGobuster(cmdout string) []string {
	
	var extract []string
	// extract = u.StringToLines(cmdout)
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {	//	k
		// fmt.Println(k,"-",v)
		if strings.Contains(v,"Status: 403") {
			continue
		} else {
			// extract += string(v)
			extract = append(extract,string(v))
		}
	}
	return extract
}