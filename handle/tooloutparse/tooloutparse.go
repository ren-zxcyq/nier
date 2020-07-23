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
	// "reflect"
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


func (h *Toolparser) ParseNikto(cmdout string) []string {
	
	var extract []string

	if strings.Contains(cmdout, "No web server found on") && strings.Contains(cmdout, "0 host(s) tested") {
		fmt.Println("Nikto - OK")
	} else {
		fmt.Println("Nikto - FAIL")
	}
	// fmt.Println(cmdout)

	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {
		fmt.Println(v)
		if strings.HasPrefix(v, "+") {
			// fmt.Println("YES", v)
			extract = append(extract,v)
		} else {
			// fmt.Println("NO", v)
			continue
		}
	
	}

	return extract
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

func (h *Toolparser) ParseBanners(cmdout string) []string {
	var extract []string
	var cmdoutlist []string

	cmdoutlist = strings.SplitN(cmdout, "conn-refused", 2)
	cmdoutlist = strings.SplitN(cmdoutlist[1], "Service Info:", 2)

	extract, err := u.StringToLines(cmdoutlist[0])

	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}

	return extract
}

func (h *Toolparser) ParseComments(cmdout string) []string {

	var extract []string
	var cmdoutlist []string
	var tmp, t string
	// var tmplist []string
	
	
	cmdoutlist = strings.SplitN(cmdout, "| http-comments-displayer:", 2)
	cmdoutlist = strings.SplitN(cmdoutlist[1], "MAC Address: ", 2)

	strCont, err := u.StringToLines(cmdoutlist[0])
	
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {	//	k
		// fmt.Println(k,"-",v)
		if strings.Contains(v, "Path:") || strings.Contains(v, "Comment:") || strings.Contains(v,"Line number:") {
			continue
		} else {
			tmp = strings.TrimLeft(v, "|")
			tmp = strings.TrimSpace(tmp)
			t = strings.ToLower(tmp)
			//	@TODO	Check using regex maybe
			if strings.Contains(t, "pass") || strings.Contains(t, "cred") || strings.Contains(t, "u:") || strings.Contains(t, "p:") || strings.Contains(t, "http") || strings.Contains(t, "https") || strings.Contains(t, "@") || strings.Contains(t, "log") || strings.Contains(t, ".com") || strings.Contains(t, "git") || strings.Contains(t, "maybe") || strings.Contains(t, "todo") {
				//	Add to the report just the lines identified by the above filter
				fmt.Println(tmp)
				extract = append(extract,tmp)
			}
		}
	}
	
	return extract
}

func (h *Toolparser) ParseHTTPrint(cmdout string) []string {
	var extract []string
	var cmdoutlist []string
	var tmplist []string
	var line string

	//	guessing
	//	opt1
	// cmdoutlist = strings.SplitN(cmdout, "Derived Signature:", 2)
	// cmdoutlist = strings.SplitN(cmdoutlist[1], "------------------------", 2)
	//	opt2
	// cmdoutlist = strings.SplitN(cmdout, "<!-- Reported signature", 2)
	// cmdoutlist = strings.SplitN(cmdoutlist[1], "-->", 2)

	cmdoutlist = strings.SplitN(cmdout, "<servers>", 2)
	cmdoutlist = strings.SplitN(cmdoutlist[1], ">", 2)

	
	// strCont, err := u.StringToLines(cmdoutlist[0])
	// if err != nil {
	// 	log.Println("Failed while separating lines in formatted tool output")
	// }
	// for _,v := range strCont {



	// 	//extract = append(extract,v)
	// 	fmt.Println("first", v)
	// }

	
	//	opt1
	// cmdoutlist = strings.SplitN(cmdoutlist[1], "<!-- Best Matches", 2)
	// // cmdoutlist = strings.SplitN(cmdoutlist[1], "-->", 2)
	//	opt2
	// cmdoutlist = strings.

	// extract = append(extract, "Guesses:")
	//	version guess ranking
	strCont, err := u.StringToLines(cmdoutlist[1])
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	fmt.Println(len(strCont))
	//strCont = strings.Split(strCont, "<match")


	for k,val := range strCont {
		// extract = append(extract,val)
		// fmt.Println("second", val)
		// fmt.Println("second", val)
		tmplist = strings.Split(val, `"`)

		line = ""
		for kk,vl := range tmplist {
			//	fmt.Println(reflect.TypeOf(vl))
			// fmt.Println("third", vl)
			if kk == 1 {
				line += vl
				line += " - "
				// fmt.Println("1 ######## ", vl, " ######## ", line)	//	fmt.Sprintf("%s\t-\t",vl)		//	strings.Join(vl, " - ")
			} else if kk == 3 {
				line += "Score: "
				// val += strings.Join(t, "/100 - ")
				line += vl
				line += " - "
				// fmt.Println("2 ######## ", vl, " ######## ", line)	//	val = fmt.Sprintf("Score: %s/100\t-\t", vl)		//	strings.Join("Score: ", vl)

			} else if kk == 5 {
				line += "Confidence: "
				line += vl
				// fmt.Println("3 ######## ", vl, " ######## ", line)	// val += fmt.Sprintf("Confidence: %s",vl)	//	strings.Join("Confidence: ", vl)
			}
			
		}

		// fmt.Println("third", line)
		extract = append(extract, line)
		// val = tmplist[1] + " " + tmplist[3] + " " + tmplist[5]
		// fmt.Println(val)
		// extract = append(extract, val)
		if k == 10 {
			break
		}
	}

	return extract
}