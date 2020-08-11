// Package injectiondetector is responsible for scraping the target website,
// extracting <form> tags, filter for unique forms, submitting all of them
// and identifying user controlled input which appears on the application pages.
// interface their execution with parsing.
package utilities

import (
	"fmt"
	"strings"
	// "path"
	"path/filepath"
	// "net/http"
	"io/ioutil"
	// "reflect"
	// "bytes"
	"bufio"
	"regexp"
	"os"
	
	"net/url"
	// "resource"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)


var u Utils
var t Agent
var target string
var extforms []extractedform
// var submitedforms []extractedform
// var possibleinjections string
var xssdiscovered string
var inj	map[string]string
var flaggedforms []extractedform

type InjectionHandler struct {
	target			string
	targetport		int
	// targethost		string	//	A concatenation of proto://target:targetport
	sessiontokens	string
	outputFolder	string
	httpprefix		string

	debug			bool	//	@TODO	-	Consider using a similar flag to the other operations
}

type extractedform struct {	//	@TODO	add src
	src				string
	method			string
	action			string
	enctype			string
	
	elements		[]string
	uqelemstring	[]string

	contents		string
	request			string
}

//
func NewInjectionHandler(target string, targetport int, outputFolder string, stokens string) *InjectionHandler {

	//	Create InjectionHandler
	var h InjectionHandler = InjectionHandler{target: target, targetport: targetport, sessiontokens: stokens, outputFolder: outputFolder, httpprefix: "http://"}

	//fmt.Printf("Address of InjectionHandler - %p", &h) //	Prints the address of the Handler
	return &h
}

func (h *InjectionHandler) InjFormCheck() {
	
	fmt.Println("\r\n\r\n[*]\tCrawling URLs provided to extract <form> tags.\r\n-------------")
	if h.debug {
		fmt.Println("********************************************")
		fmt.Println("fROM inj")
		fmt.Println(h.target)
		fmt.Println(h.targetport)

		fmt.Println("Tokens:")
		fmt.Println(h.sessiontokens)
		fmt.Println(h.outputFolder)
		fmt.Println(h.httpprefix)
		fmt.Println("********************************************")	

		fmt.Println("--------------------------------")
		fmt.Println("URLs Combined:")
	}

	var urls []string

	urls = h.combinedURLs()
	
	var urlsduringinj string
	for _,i := range urls {
		urlsduringinj += i + "\r\n"
	}
	var location string = h.outputFolder + "/urls_used_during_detection.txt"
	fmt.Println("\r\n[*]\tWriting URLs used during injection & XSS testing to file:\t" + location + "\r\n")
	u.SaveStringToFile(location, urlsduringinj)

	// /root/Desktop/report/urls_used_during_detection.txt

	for _,url := range urls {
		// fmt.Println(url)
		// os.Exit(1)
		if !strings.Contains(url,"log") || !strings.Contains(url,"forgot") {	//	avoid logging out
			// fmt.Println(url)
			h.injRequestURLi(url)
		} else {
			fmt.Println("[*]== [Skipping URL]",url,"[Reason]: \"log\" is contained in the URL")
		}
	}

	fmt.Println("\r\n[*]\tSubmitting <forms>\r\n")	//-------------")
	
	for i,_ := range extforms {
		// h.handleSubmission(&extforms[i])						//	Check if this actually works
		h.submitform(&extforms[i])
	}

	h.checkforuqstrings(urls)									//	@TODO	Swap within for a method that submits cookies

	fmt.Println("\r\n\r\n\r\n[*]\tPossible Injections\t",len(xssdiscovered),"\r\n\t", xssdiscovered)
	if len(xssdiscovered) > 0 {

		var location string = h.outputFolder + "/form_injection_detection.txt"
		fmt.Println("\r\n[*]\tWriting <form> Submission Results to file:\t" + location + "\r\n")
		u.SaveStringToFile(location, xssdiscovered)
	} else {
		fmt.Println("[*]== [Did not write file: form_injection_detection.txt]")
	}
	//	@TODO
	// fmt.Println("\r\n[*]\tForms Identified as Origins of User Controlled Input\r\n")
	// for _,j := range flaggedforms {
	// 	fmt.Println("i\t",j.contents)
	// }

	// fmt.Println(len(possibleinjections))
	// fmt.Println("\r\n")

}

func (h *InjectionHandler) getwithcookies(urltoget string) string {
	// var target string = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport) + urltoget)

	
	// GET WITH COOKIES urltoget
	rurl,errurl := url.ParseRequestURI(urltoget)	//h.httpprefix +  + urltoget)	//	h.httpprefix + form.action)	//	targethost)
	if errurl != nil {
		fmt.Printf("\r\nparsing error - %s", errurl, "\r\n")
	}

	urlStr := rurl.String()
	
	client := &http.Client{}


	// fmt.Println("Get With COOKIES")
	// fmt.Println(http.MethodGet)
	// fmt.Println(urlStr)
	r,e := http.NewRequest(http.MethodGet, urlStr, nil)	//	URL-encoded payload

	if e != nil {
		fmt.Printf("\r\nERR 1 - %s",e,"\r\n")
	}

	if len(h.sessiontokens) > 0 {
		var tokens []string = u.StringCookiesToList(h.sessiontokens)
		for _,k := range tokens {
			var token []string = u.SeparateCookie(k)
			r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})

		}
	}
	
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("\r\nERR 2 - %s",err,"\r\n")
	}
	
	
	//	Extract Body
	body, _ := ioutil.ReadAll(resp.Body)
	
	return string(body)
}

func (h *InjectionHandler) getwithcookiesforuqstrings(urltoget string) string {

	rurl,errurl := url.ParseRequestURI(urltoget)	//h.httpprefix +  + urltoget)	//	h.httpprefix + form.action)	//	targethost)
	if errurl != nil {
		fmt.Printf("\r\nparsing error - %s", errurl, "\r\n")
	}

	urlStr := rurl.String()
	
	client := &http.Client{}

	// fmt.Println("Get With COOKIES")
	// fmt.Println(http.MethodGet)
	// fmt.Println(urlStr)
	r,e := http.NewRequest(http.MethodGet, urlStr, nil)	//	URL-encoded payload

	if e != nil {
		fmt.Printf("\r\nERR 3 - %s",e,"\r\n")
	}

	if len(h.sessiontokens) > 0 {
		var tokens []string = u.StringCookiesToList(h.sessiontokens)
		for _,k := range tokens {
			var token []string = u.SeparateCookie(k)
			// fmt.Println("Using Token 2\t-\t",token[0],"\t-\t",token[1])
			r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
		}
	}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("\r\nERR 4 - %s",err,"\r\n")
	}
	

	//	Extract Body
	body, _ := ioutil.ReadAll(resp.Body)
	
	return string(body)
}

// Function checkforuqstrings() iterates over the uqstrings submitted.
// 1) sets up the [*]==== output  which is saved to file - variable: _.
// 2) appends each form that resulted in reflection to a list.
// 3) iterates that list, finds the form, reads an XSS Payload, modifies it, submits it.
// 4) calls selenium to identify any potential <script>alert()</script> boxes.
func (h *InjectionHandler) checkforuqstrings(urls []string) {
	var uqstrings []string = h.getUQstrings()	//	SAVE [] into var

	//	uqstrings []string - contains every submitted value as in - formkey:value

	// Our aim here is to find elements of uqstrings that are reflected.
	// save our results to a map of the structure
	// map {
	// 	uqstring1 = [ url1, url2, url3],
	// 	uqstring3 = [ url1, url3],
	// }
	fmt.Println("creating reflectd map")
	var reflected map[string][]string = make(map[string][]string)
	fmt.Println("uqstrings contain:",uqstrings)


	// populate the map by
	//	get each url
	//	for each url
	//	for each uq in the response
	//	append to map[uq] += url
	fmt.Println("populating reflectd")
	for _,appurl := range urls {
		if !strings.Contains(appurl,"log") || !strings.Contains(appurl,"forgot") {
			// GET URL WITH COOKIES
			// fmt.Println("GET WITH COOKIEs")
			r := h.getwithcookiesforuqstrings(appurl)
			// fmt.Println("\r\narrrrgh is",r)
			for _, marker := range uqstrings {
				if strings.Contains(r, marker) {
					reflected[marker] = append(reflected[marker], appurl)
					fmt.Println("appending to",marker,"-",appurl)
				}
			}
		}
	}

	fmt.Println("ECHOING REFLECTED")
	for m,_ := range reflected {
		fmt.Println("\r\n-")
		fmt.Println(m, "\t-\t", reflected[m])
		fmt.Println("-")
	}

	fmt.Println("[*]\tBegin Injection")

	h.removeelements()

	for uqmapkey,uqmapvalues := range reflected {
		usekey:
			for _,f := range extforms {
				// fmt.Println("\r\nFOR FORM and ",uqmapkey)
				var found bool = false
				findform:
					for _,el := range f.uqelemstring {
						if strings.Contains(el,uqmapkey) {
							// fmt.Println("FOUND FORM")
							found = true
							break findform
						}
					}
				if found {	//	@TODO - Make sure to continue the outer loop once this is done so that we don't check everything again.
					// fmt.Println("found form", f.uqelemstring)
					fmt.Println("\r\n\r\nFor mapkey:",uqmapkey,"\r\n-------------------")

					// Get Payloads. I.E. Either Gen() or Read entire file
					// h.genXSSPayloads()
					// Read Payloads file
					payloadfile, err := os.Open("/root/go/src/github.com/ren-zxcyq/nier/resources/poc-xss-payloads.txt")
					if err != nil {
						fmt.Println(err)
					}
					defer payloadfile.Close()
				
					var scanner = bufio.NewScanner(payloadfile)
					lines := []string{}
				
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}

					// for payload in payloads
					// submit
					for _, line := range lines {
						// convert line from (a0*) to (uqmapkey)
						re := regexp.MustCompile(`(?i)\({1}[a-z0-9]+\){1}`)
						var payload string = re.ReplaceAllString(line,`("`+uqmapkey+`");`)
						// fmt.Println("Convert Line --", line,"\t-\t",payload)
						// fmt.Println("Submit =>\t",payload)	//"$1"))

						// // ?alt?	=>	opt for re right now
						// // line = strings.Replace(line,"(1)","("+uqmapkey+")",1)
						// // fmt.Println(line)


						// build data=url.Values{}
						data := url.Values{}		//	:= url.Values{}
						data = url.Values{}
						for _,el := range f.uqelemstring {
							var tmp []string = nil
							tmp = strings.SplitN(el,":",2)
							if strings.Contains(tmp[1],uqmapkey) {
								// fmt.Println("add payload\t", tmp[0],"-",payload)
								data.Set(tmp[0],payload)
							} else {
								// fmt.Println("add normal\t", tmp[0],"-",tmp[1])
								data.Set(tmp[0],tmp[1])
							}
						}

						// fmt.Println("\r\ndata {} built successfully.\r\n")
						for k,v := range data {
							fmt.Println("\t",k,"\t-\t",v)
						}
						// fmt.Println(reflect.TypeOf(data))
						// fmt.Println("Submitting with payload")
						h.submitWithPayload(f, data)
						

						// selenium return the thing. if it returns break right above line range lines?
						fmt.Println("\r\n\r\n-------------------\r\n",uqmapkey, "\r\n-------------------\r\n")
						for _,location := range uqmapvalues {
							
							searchresult := h.checkForAlertInPage(location, uqmapkey)
							
							// fmt.Println("\r\n[*]\t----- ALERT DISCOVERY",searchresult)
							fmt.Println("\tlocation -",location)
							
							if strings.Contains(searchresult,uqmapkey) {
								xssdiscovered += "\r\n\r\n\r\n[*]\tXSS - Detected at:\t" + location
								xssdiscovered += "\t-\r\n\tForm Location:\t" + f.src
								xssdiscovered += "\r\n\tPayload:\t" + payload
								xssdiscovered += "\r\n\tForm Contents:\r\n-\r\n\t\t" + f.contents + "\r\n-\r\n\r\n"

								// Found at least 1 - ?Jump to the next
								break usekey
							}
							
						}
						h.removeelements()


					}
					break usekey
				}
			}
	}
	// fmt.Println(xssdiscovered)
	// os.Exit(0)
}

// @TODO CHECK getUQstrings
func (h *InjectionHandler) getUQstrings() []string {
	var list []string
	for _,f := range extforms {
		for _,k := range f.uqelemstring {
			var t []string = strings.Split(k,":")
			if t[1] != "" {
				list = append(list,t[1])
			}
		}
	}
	return list
}

func (h *InjectionHandler) injRequestURLi(url string) {

	//	Check & add if not present - http://
	// target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
	// if h.debug {
	// 	fmt.Println("==============\t",target," - ",url,"\t==============")
	// }
	// target = target + url
	// // fmt.Println("==============\t\t",url)

	var r string

	r = h.getwithcookies(url)	//	target
	
	if h.debug {
		fmt.Println("getwithcookies() WAS RUN", len(r))
	}

	var tmpforms []extractedform

	if strings.Contains(r, "Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking") {	//	or
		fmt.Println("[+]\tUpgrading to HTTPS\t",url)
		h.httpprefix = "https://"


	} else {								//	@TODO	consider checking for another error
		// fmt.Println("-------------")
		if h.debug {
			fmt.Println("[+]\tContinue HTTP\t",url)
		}

		h.httpprefix = "http://"
		// fmt.Println(url)
		// fmt.Println(r[0:10])
		if len(r) == 0 {
			// fmt.Println("r LENGTH = 0")
		} else if len(r) >= 1 {
			// fmt.Println(r[0:1])					//	print 2 lines from each response
		}

		//	If HTML response contains a form -> pass it to the parser
		if strings.Contains(r, "<form") {	//	or
			tmpforms = h.extractForms(r, url)	//	h.httpprefix + h.target + 
			uniquenesscheck(tmpforms)	//	, extforms
		}
	}

}

func uniquenesscheck(tocheck []extractedform) {
	if len(extforms) > 0 {
		for _,v := range tocheck {
			if !isstrinforms(v.contents) {
				//	Appending
				extforms = append(extforms, v)
			} else {
				//	Not Appending
			}
		}

	} else {
		extforms = append(extforms, tocheck[0])
		uniquenesscheck(tocheck[1:])
	}
}

func isstrinforms(str string) bool {

	var r bool
	for _,v := range extforms {

		if str == v.contents {
			r = true
		}
	}

	return r
}

// func (h *InjectionHandler) handleSubmission(f *extractedform) {
// 	h.submitform(f)
// }

var scounter, fcounter int
// Submits the given form using the autogenerated values for the fields.
// Returns the fields submitted
func (h *InjectionHandler) submitform(form *extractedform) {
	
	data := url.Values{}

	if len(form.elements) > 0 {

		// fmt.Println("IN submitform() - IF LEN > 0")

		var tmp []string
		var uq string
		for _,v := range form.elements {
			tmp = strings.SplitN(v,":",2)
			uq = u.UniqueStringAlphaNum()
			// fmt.Println("APPENDING", tmp[0] + ":" + uq)


			if tmp[0] != "submit" {
				form.uqelemstring = append(form.uqelemstring, tmp[0] + ":" + uq)				//	REMOVED UQ FROM HERE
				data.Set(tmp[0],uq)
			} else {
				// fmt.Println("CONTAINS submit")
				form.uqelemstring = append(form.uqelemstring, tmp[0] + ":")
				data.Set(tmp[0],"")
			}
		}
		scounter += 1
	} else {
		fcounter += 1

		//	<form> contains no elements
		//	Notify & move on to the next

		return
	}
	if strings.Contains(form.contents, "password") {	//	@FILTER
		// fmt.Println("IN submitform() - contains pass")
		return
	}

	//	@TODO	HERE	-	remove this
	if !strings.Contains(form.contents, "testimonial") {
		return
	}

	if h.debug {
		// fmt.Println("Submitting:\r\n-------")
		// fmt.Println("\tSRC\t-", form.src)
		// fmt.Println("\tMET\t-", form.method)
		// fmt.Println("\tACT\t-", form.action)
		// fmt.Println("\tENC\t-", form.enctype)
		// fmt.Println("\tELM\t-", form.elements)
		// fmt.Println("\tUQL\t-", form.uqelemstring)
		// fmt.Println("\tCON\t-", form.contents)
		// fmt.Println("-------")
	}

	/*	@HERE		REMOVE form.action = _
	rurl,_ := url.ParseRequestURI("http://127.0.0.1:9988")
	*/
	// fmt.Println("form.action", form.action)
	// form.action = "http://127.0.0.1:9988/post-testimonial.php"
	rurl,_ := url.ParseRequestURI(form.action)	//	h.httpprefix + form.action)	//	targethost)
	urlStr := rurl.String()
	// fmt.Println("form.action", urlStr)

	client := &http.Client{}	//	@TODO	Consider how HTTPS will be handled. This has to be changed here.
	
	// r,_ := http.NewRequest(form.method, urlStr, strings.NewReader(data.Encode()))	//	URL-encoded payload
	
	r,_ := http.NewRequest(form.method, urlStr, strings.NewReader(data.Encode()))	//	URL-encoded payload


	//	From now on it will be submitted.
	
	//	Setting HTTP Headers.
	r.Header.Set("Referer", form.action)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if len(h.sessiontokens) > 0 {
		var tokens []string = u.StringCookiesToList(h.sessiontokens)
		for _,k := range tokens {
			var token []string = u.SeparateCookie(k)
			if h.debug {
				// fmt.Println("Using Token 1\t-\t",token[0],"\t-\t",token[1])
			}
			r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
		}
	}

	//	Right before submitting the request
	//	Parse the Submitted Request
	//	Save it on the extractedform struct
	pp,_ := url.Parse(r.URL.String())
	var requestString string = r.Method + " " + pp.Path + " " + r.Proto + "\r\n"
	
	// requestString += "Host: " + 
	// p,_ := url.Parse(r.URL.String())
	// p,_ := url.ParseRequestURI(r.URL.String())
	// pp := p.Path.String()
	requestString += "Host: " + string(r.Host) + "\r\n"
	for ii,vvv := range r.Header {
		// requestString += ii + " " + vvv + "\r\n"
		var val string
		for _,jj := range vvv {
			val += jj
		}
		// val = strings.Trim(val," ")
		requestString += fmt.Sprintf("%s: %s\r\n",ii, val)
	}
	requestString += "\r\n"
	requestString += data.Encode() + "\r\n"

	if h.debug {
		// fmt.Println("^^^^^^^^^^^^^^^^^^^^")
		fmt.Println(requestString)
		fmt.Println(form.request)
		// fmt.Println("^^^^^^^^^^^^^^^^^^^^")

		// for k, vv := range r.Header {
		// 	n := copy(sv, vv)
		// 	h2[k] = sv[:n:n]
		// 	sv = sv[n:]
		// }
	}
	form.request = requestString


	resp, _ := client.Do(r)	//	resp, _
	if (h.debug) {
		fmt.Println("\r\nRESPONSE:\t",resp.Status,"\r\n/////////////////////////////////////////////////")
	}
	fmt.Println(resp.Status)
	if strings.Contains(resp.Status, "404") {
		fmt.Println(form.request)
	}
}

func (h *InjectionHandler) combinedURLs() []string {

	//	Filter Spider Output
	var combinedSpiderURLs string = h.outputFolder + "/spider_URLs.list"

	var contents = u.ReturnFileContentsStr(combinedSpiderURLs)
	res := h.ParseSpiderURLContents(contents)

	var r []string
	for _,k := range res {
		r = append(r,k)
	}

	return r
}

func (h *InjectionHandler) ParseSpiderURLContents(cmdout string) []string {
	var extract []string
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		fmt.Println("Failed while separating lines in fromatted tool output")
	}
	for _,v := range strCont {
		if strings.HasPrefix(v, "[url] - ") {
			v = v[8:]
		} else {
			fmt.Println("NO '[url] - ' prefix\t>", v,"<")
			continue
		}
		if strings.Contains(v,"[code-4") {
			fmt.Println("[ignored] - Response code of family 4xx\t>", v,"<")
			continue
		}
		if strings.Contains(v,"[code-5") {
			fmt.Println("[ignored] - Response code of family 5xx\t>", v,"<")
			continue
		}
		if strings.HasPrefix(v, "[code-") {
			v = v[13:]
			extract = append(extract,v)
		}
	}
	return extract
}

func (h *InjectionHandler) oldCombinedMethodAsINGobusterAndRelLinks() []string {

	//	Filter Gobuster output
	var gobusterOutURL string = h.outputFolder + "/links_gobuster_and_rel.txt"

	var gobuster string = u.ReturnFileContentsStr(gobusterOutURL)
	res := h.parseGobuster(gobuster)
	if h.debug {
		fmt.Println("******hhhhhhhhhh\r\n",res, "\r\nhhhhhhhhhh******")
	}

	var r []string
	for _,k := range res {
		// fmt.Println(i,"\t---\t","http://192.168.1.20" + k)
		// r = append(r, "http://192.168.1.20" + k)
		r = append(r, k)
	}
	if h.debug {
		fmt.Println("r is:\r\n", r)
	}
	return r

}

// ParseGobuster Filters out results that are of (Status: 403)
// Returns an array of lines.
func (h *InjectionHandler) parseGobuster(cmdout string) []string {

	var extract []string
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		fmt.Println("Failed while separating lines in formatted tool output")
	}

	for _,v := range strCont {
		if !(len(v) == 0) {
			extract = append(extract, v)
		}
	}
	return u.Uniquestrslice(extract)
}


func (h *InjectionHandler) extractForms(r string, r_url string) []extractedform {

	var forms []extractedform

	//	Ignored Elements
	//	-	label			//	Just text
	//	-	<fieldset>		//	Groups related items

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader((r)))
    doc.Find("form").Each(func(i int, form *goquery.Selection) {

		var f extractedform

		f.src = r_url

		// fmt.Println("---------FOUND FORM---------")
		formhtml,_ := goquery.OuterHtml(form)
		f.contents = formhtml

		//	Get <form> Attributes
		action, okaction := form.Attr("action")
		if okaction {	//	basically set form.action to be used by submitform()
			var t []string
			if len(action) != 1 && strings.Contains(action, "#") {
				t = strings.SplitN(action,"#",2)
				action = t[0]
			} else if len(action) == 1 {
				action = ""
			}

			//	@TODO	Consider using native go
			if strings.HasPrefix(action, "http://") || strings.HasPrefix(action, "https://") {
				f.action = action
			} else {
				if h.debug {
					// fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", f.src)
					// fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", action)
				}
				
				if strings.HasPrefix(f.src, "http://") || strings.HasPrefix(f.src, "https://") {

					tmp1, _ := filepath.Split(f.src)
					if h.debug {
						// var tq int = len(tmp)
						// fmt.Println("*************************************", tmp1)
						// fmt.Println("*************************************", tmp2)
					}
					f.action = tmp1 + action
				} else {
					f.action = f.src + action
				}
			}
		} else {
			f.action = f.src
		}

		method, okmethod := form.Attr("method")
		if okmethod {
			f.method = strings.ToUpper(method)
		}

		enctype, okenctype := form.Attr("enctype")
		if okenctype {
			f.enctype = enctype
		}
		
		//	Get <form> Elements

		//	-	<output>
		form.Find("output").Each(func(j int, output *goquery.Selection) {
			//	Skip current <form> gracefully
			// fmt.Println("Found <output> - Notify & Go to Next <form>")
			
			// @TODO	-	Notify & Go to next form

			var outputTag string = `<output`
			outputname,okoutputname := output.Attr("name")
			if okoutputname {
				outputTag += ` name="` + outputname + `"`
			}
			outputfor,okoutputfor := output.Attr("for")
			if okoutputfor {
				outputTag += ` for="` + outputfor + `"`
			}
			outputTag += `>`


			//	Notes:
			//		-	For Now Ignore the entire form if <output> is included.
			//			it is only used to display output.
			//
			//		-	https://developer.mozilla.org/en-US/docs/Web/HTML/Element/output
			//			https://www.w3schools.com/tags/tag_output.asp
			//
			//		-	Format:
								// <form action="/action_page.php" oninput="x.value=parseInt(a.value)+parseInt(b.value)">
								// 	0
								// 	<input type="range" id="a" name="a" value="50">
								// 	100 +
								// 	<input type="number" id="b" name="b" value="50">
								// 	=
								// 	<output name="x" for="a b"></output>
								// 	<input type="submit">
								// </form>
			// 		-	represents the result of a calculation (like one performed by a script)
		})

		//	-	<button>
        form.Find("button").Each(func(j int, b *goquery.Selection) {
			// fmt.Println("\tFOUND <button>")
			var buttonTag string = "<button"
			btype, oktype := b.Attr("type")
			if oktype {
				if strings.Contains(btype, "submit") {
					// fmt.Println("\t\t<button.Type - submit",)
					buttonTag += ` type="submit" `
				} else {
					// fmt.Println("\t\t<button.Type -", btype)
					buttonTag += ` type="` + btype + `" `
				}
			}
			bname, okname := b.Attr("name")
			if okname {
				// fmt.Println("\t\t<button.Name -", bname)
				buttonTag += ` name="` + bname + `"`
				f.elements = append(f.elements, bname + `:` + `BBBBB` + `&`)
			}
			bonclick, okbonclick := b.Attr("onclick")
			if okbonclick {
				buttonTag += ` onclick="` + bonclick + `"`
			}
			buttonTag += `>`

			// htmlb,_ := goquery.OuterHtml(b)
			// fmt.Println("\t\tORIGINAL:\t", htmlb)
			// fmt.Println("\t\tEXTRACTED:\t", buttonTag)
			
			
			//	Notes:
			//		-	Format:
								// <button type="button" onclick="alert('Hello World!')">Click Me!</button>
		})


		//	-	<input>
		form.Find("input").Each(func(j int, in *goquery.Selection) {
			// fmt.Println("\tFOUND <input>")
			var inputTag string = "<input"
			intype, oktype := in.Attr("type")
			if oktype {
				if strings.Contains(intype, "text") {			//	?alt?	just do: inputTag += ` type="` + intype + `"`
					// fmt.Println("\t\t<input.Type - text")
					inputTag += ` type="text" `
				} else if strings.Contains(intype, "range") {
					inputTag += ` type="range"`
				} else if strings.Contains(intype, "number") {
					// fmt.Println("\t\t<input.Type -", intype)
					inputTag += ` type="number"`
				}
			}
			inname, okname := in.Attr("name")
			if okname {
				// fmt.Println("\t\t<input.Name -", inname)
				inputTag += ` name="` + inname + `"`
				// f.contents += inname + `:` + `AAAAA` + `&`
				f.elements = append(f.elements, inname + `:` + `BBBBB` + `&`)
			}

			inlist,okinlist := in.Attr("list")
			if okinlist {
				//	Need to find <datalist id="inlist"> & nested <option>s
				form.Find("datalist").Each(func(k int, datalist *goquery.Selection) {
					dlistid,okdlistid := datalist.Attr("id")
					if okdlistid {
						if dlistid == inlist {
							//	Find Nested <option> tags
							datalist.Find("option").Each(func(o int, opt *goquery.Selection) {
								var optTag string = `<option`
								optval, okoptval := opt.Attr("value")
								if okoptval {
									optTag += ` value="` + optval + `"`
								}
								optTag += `>`
				
								// opthtml,_ := goquery.OuterHtml(opt)
								// fmt.Println("\t\tORIGINAL:\t", opthtml)
								// fmt.Println("\t\tEXTRACTED:\t", optTag)
							})
							//	?Act?
						}
					}
					//	Notes:
					//		-	specifies a list of pre-defined options for an <input> element.
					//				users will see a drop-down list.
					//		-	Format:
											// <form action="/action_page.php">
											// 	<input list="browsers" name="browser">
											// 	<datalist id="browsers">
											// 		<option value="Internet Explorer">
											// 		<option value="Firefox">
											// 		<option value="Chrome">
											// 		<option value="Opera">
											// 		<option value="Safari">
											// 	</datalist>
											// 	<input type="submit">
											// </form>

				})
			}

			inputTag += ">"

			// htmli,_ := goquery.OuterHtml(in)
			// fmt.Println("\t\tORIGINAL:\t", htmli)
			// fmt.Println("\t\tEXTRACTED:\t", inputTag)

			//	Notes:
			//		-	If the type attribute is omitted, the input field gets the default type: "text".
			
		})

		//	-	<textarea>
		form.Find("textarea").Each(func(j int, txtarea *goquery.Selection) {
			// fmt.Println("\tFOUND <textarea>")
			var txtareaTag string = `<textarea`

			txtareaname, oktxtareaname := txtarea.Attr("name")
			if oktxtareaname {
				txtareaTag += ` name="` + txtareaname +`"` 
			}
			txtareaTag += `>`

			// htmltxtarea,_ := goquery.OuterHtml(txtarea)
			// fmt.Println("\t\tORIGINAL:\t", htmltxtarea)
			// fmt.Println("\t\tEXTRACTED:\t", txtareaTag)
			f.elements = append(f.elements, txtareaname + `:` + `BBBBB` + `&`)
			
			// fmt.Println("GOT ]]]]]]]]]]]]]]]]]]]]]]]]]]]]")
			// fmt.Println(txtareaname + `:` + `BBBBB` + `&`)
			// // f.elements = append(f.elements, txtareaname + `:` + `BBBBB` + `&`)
			// fmt.Println("GOT ]]]]]]]]]]]]]]]]]]]]]]]]]]]]")

			//	Notes:
			//		-	Format:
									// <textarea name="message" rows="10" cols="30">
									// The cat was playing in the garden.
									// </textarea>

			//		-	rows & cols			attributes		don't affect what is sent
		})

		//	-	<select>	//	drop-down list
		form.Find("select").Each(func(j int, sl *goquery.Selection) {
			// fmt.Println("\tFOUND <select>")

			var slTag string = "<select"
			slnameval,okslnameval := sl.Attr("name")
			if okslnameval {
				slTag += ` name="` + slnameval + `"`
			}
			_,okslmultiple := sl.Attr("multiple")
			if okslmultiple {
				slTag += ` multiple `
			}
			slTag += `>`

			// htmlsl,_ := goquery.OuterHtml(sl)
			// fmt.Println("\t\tORIGINAL:\t", htmlsl)
			// fmt.Println("\t\tEXTRACTED:\t", slTag)


			sl.Find("option").Each(func(o int, opt *goquery.Selection) {
				var optTag string = `<option`

				optval, okoptval := opt.Attr("value")
				if okoptval {
					optTag += ` value="` + optval + `"`
				}
				optTag += `>`

				// opthtml,_ := goquery.OuterHtml(opt)
				// fmt.Println("\t\tORIGINAL:\t", opthtml)
				// fmt.Println("\t\tEXTRACTED:\t", optTag)
	
			})
			//	Notes:
			//		-	Format:
								// <form action="/action_page.php">
								// 	<label for="cars">Choose a car:</label>
								// 	<select id="cars" name="cars">
								// 		<option value="volvo">Volvo</option>
								// 		<option value="saab">Saab</option>
								// 		<option value="fiat">Fiat</option>
								// 		<option value="audi">Audi</option>
								// 	</select>
								// 	<input type="submit">
								// </form>
			//		-	Define a default option:			selected	attribute
								//	<option value="fiat" selected>Fiat</option>

			//		-	Select Multiple:					multiple	attribute
								//	<select id="cars" name="cars" size="4" multiple>
								//	Results in
								//	cars=volvo&cars=audi
		})

		//	Perform Soft Check:
		


		// fmt.Println("[+]==============\r\n",f,"\r\n[+]==============\r\n")
		// fmt.Println("-------------\r\n-------------\r\n",len(forms),"\r\n-------------\r\n-------------\r\n")
		f = *(h.formsoftCheck(&f))
		if !(&f == nil) {
			// fmt.Println("\r\n!@#$\t-\tAPPENDING FORM - ", f.contents)
			forms = append(forms, f)
		}	// else {
		// 	fmt.Println("\r\n!@#$\t-\tAPPENDING FORM NOT - ", f.contents)
		// }
	})
	if h.debug {
		// fmt.Println("=====================\r\n-------------\r\n",len(forms),"\r\n-------------\r\n=====================\r\n")
	}
	return forms
}

func (h *InjectionHandler) formsoftCheck(f *extractedform) *extractedform {
	var keywords []string = []string{"delete", "remove", "edit", "log"}
	for _,word := range keywords {
		if strings.Contains(f.contents, word) {
			*f = extractedform{}
			break
		} else {
		}
	}
		return f
}

func (h *InjectionHandler) submitWithPayload(form extractedform, data url.Values) {
	//@TODO	-	IF LEN(FORM.ELEMENTS) > 0	->	create data
	//@TODO	-	IF form.contents contains "password"	->	return

	rurl,_ := url.ParseRequestURI(form.action)	//	h.httpprefix + form.action)	//	targethost)
	urlStr := rurl.String()
	fmt.Println("form.action", urlStr)

	client := &http.Client{}	//	@TODO	Consider how HTTPS will be handled. This has to be changed here.

	r,_ := http.NewRequest(form.method, urlStr, strings.NewReader(data.Encode()))	//	URL-encoded payload

	//	Setting HTTP Headers.
	r.Header.Set("Referer", form.action)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if len(h.sessiontokens) > 0 {
		var tokens []string = u.StringCookiesToList(h.sessiontokens)
		for _,k := range tokens {
			var token []string = u.SeparateCookie(k)
			if h.debug {
				// fmt.Println("Using Token 1\t-\t",token[0],"\t-\t",token[1])
			}
			r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
		}
	}

	resp, _ := client.Do(r)
	if (h.debug) {
		fmt.Println("\r\nRESPONSE:\t",resp.Status,"\r\n/////////////////////////////////////////////////")
	}
	fmt.Println(resp.Status)
}

// // @TODO
// h.genXSSPayloads()