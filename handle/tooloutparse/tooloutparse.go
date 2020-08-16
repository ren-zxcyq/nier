// Package tooloutparse receives string objects from handleexec package.
// Extracts features and ?-> Act?
package tooloutparse


import (
	// "os"
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

func (h *Toolparser) ParsePing(cmdout string) bool {

	var r bool = true
	if strings.Contains(cmdout, "1 packets transmitted, 1 received, 0% packet loss") {
		fmt.Println("Ping - OK")
		//	Did receive Reply
	} else {
		fmt.Println("Ping - FAIL")
		r = false
		//	Did not receive Reply Host is unreachable
		//	@TODO	-	?Fail Gracefully?
	}
	// fmt.Println(cmdout)
	return r
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
		// //	?DONE?@TODO	-	Connect with Reporting
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
		fmt.Println("NmapVuln - OK")
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
		fmt.Println("\r\n\r\n[*]\tNikto - FAIL")
	} else {
		fmt.Println("\r\n\r\n[*]\tNikto - OK")
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

// ParseGobuster Filters out results that are of (Status: 403)
// Returns an array of lines.
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
	var tmp []string
	var cmdoutlist []string

	cmdoutlist = strings.SplitN(cmdout, "conn-refused", 2)
	cmdoutlist = strings.SplitN(cmdoutlist[1], "Service Info:", 2)

	tmp, err := u.StringToLines(cmdoutlist[0])

	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}

	for _,l := range tmp {
		if len(l) > 0 && (len(strings.TrimSpace(l)) > 0) {
			extract = append(extract,l)
		}
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
	// fmt.Println(len(strCont))
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
		if len(line) > 0 {
			extract = append(extract, line)
		}
		// val = tmplist[1] + " " + tmplist[3] + " " + tmplist[5]
		// fmt.Println(val)
		// extract = append(extract, val)
		if k == 10 {
			break
		}
	}

	return extract
}

func (h *Toolparser) ParseHTTPMethods(cmdout string) []string {
	var extract []string
	var cmdoutlist []string



	cmdoutlist = strings.SplitN(cmdout, "Method - Status\r\n-------------", 2)
	// fmt.Println("AAAAAAAA", cmdoutlist[0]),
	if strings.Contains(cmdoutlist[0], "HTTPS") {
		extract = append(extract, "HTTPS")
	} else {
		extract = append(extract, "HTTP")
	}

	strCont, err := u.StringToLines(cmdoutlist[1])
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	// fmt.Println(len(strCont))
	// strCont = strings.Split(strCont, "<match")
	extract = append(extract, "---")

	for _,val := range strCont {
		// fmt.Println(k,"\t-\t",val)
		if strings.Contains(val, "-") {
			extract = append(extract, val)
		}
	}

	return extract
}

func (h *Toolparser) ParseRobots(cmdout string) []string {
	var extract []string

	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}

	for _,val := range strCont {
		if len(val) > 0 {
			extract = append(extract, val)
		}
	}
	
	return extract
}

// ParseGobuster Filters out results that are of (Status: 403)
// Returns an array of lines.
func (h *Toolparser) ParseGobusterAndSpidersLinks(cmdout string) []string {
	
	var extract []string
	// extract = u.StringToLines(cmdout)
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {
		if len(v) > 0 {
			extract = append(extract,string(v))
		}
	}
	return extract
}

func (h *Toolparser) ParseGobusterSubdomains(cmdout string) []string {

	// Found: chrome.google.com
	// Found: ns1.google.com
	// Found: admin.google.com
	// Found: www.google.com
	// Found: m.google.com
	// Found: support.google.com

	var extract []string
	// extract = u.StringToLines(cmdout)
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {	//	k
		// fmt.Println(k,"-",v)
		if !strings.Contains(v,"Found:") {
			continue
		} else {
			// extract += string(v)
			var tmp []string = strings.SplitN(v,"Found:",2)
			var r string = tmp[1]
			extract = append(extract,strings.TrimSpace(r))
		}
	}
	return extract
}

func (h *Toolparser) ParseWPScanner(cmdout string) []string {
	
	var extract []string


	var tmp []string = strings.SplitN(cmdout,"_______________________________________________________________",3)
	
	// for _,i := range tmp {
	// 	fmt.Println("TMP ",i)
	// }
	cmdout = tmp[2]
	// cmdout
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {
		if len(v) > 0 {
			extract = append(extract, string(v))
		}
	}

	return extract
}

func (h *Toolparser) ParseSeleniumXSS(cmdout string) []string {
	
	var extract []string

	var numberOfAlertsDetected int = strings.Count(cmdout,"[*]")
	
	if numberOfAlertsDetected > 0 {

		extract = append(extract, "Active alert pop-ups contaning the injected strings were detected")
		var tmp []string = strings.SplitN(cmdout,"[*]",numberOfAlertsDetected+1)
		
		// var tmp []string
		for _,j := range tmp {
			// fmt.Println("TMP",i,"-",j)
			if len(j) > 0 {
				var tmps []string
				// extract = append(extract, j)

				// cmdout
				tmps, err := u.StringToLines(j)
				if err != nil {
					log.Println("Failed while separating lines in formatted tool output")
				}
				for _,l := range tmps {
					if len(l) > 0 && (len(strings.TrimSpace(l)) > 0) {
						extract = append(extract, l)
					}
				}
				// extract = append(extract,"\r\n")
			}
		}
		
		// for _,v := range strCont {
		// 	extract = append(extract, string(v))
		// }
	} else {
		extract = append(extract, "No Active alert pop-ups containing the injected strings were detected.")
	}
	return extract
}

func (h *Toolparser) ParseReflectedOutput(cmdout string) []string {

	var extract []string

	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v,"[*]") {
			var tmp []string = strings.SplitN(v,"[*]",2)
			var tmps string = strings.TrimSpace(tmp[1])
			extract = append(extract,tmps)
		} else {
			var t string = strings.TrimSpace(string(v))
			if (len(t) > 0) {
				extract = append(extract,string(t))
			}
		}
	}
	return extract
}

func (h *Toolparser) ParseXSStrikeOutput(cmdout string) []string {

	var extract []string

	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {
		
		// // 	b += u.ConvertFromAnsiToUTF8([]byte(string(r)))
		// var sv string = string(v)	//	strings.TrimSpace(string(vb))
		v = u.StripANSI(v)
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v,`[+]`) || strings.HasPrefix(v,`[++]`) || strings.HasPrefix(v,`[!]`) {
			extract = append(extract,strings.TrimSpace(v))
		}
	}
	return extract
}

func (h *Toolparser) ParseCVEs(cmdout string) []string {
	
	var extract []string
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {
		if len(v) > 0 {
			extract = append(extract,v)
		}
	}
	return extract
}