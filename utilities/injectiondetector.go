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
	// "os"
	
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
var possibleinjections string

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

	for _,url := range urls {
		if !strings.Contains(url,"log") {	//	avoid logging out
			// fmt.Println(url)
			h.injRequestURLi(url)
		} else {
			fmt.Println("[*]== [Skipping URL]",url,"[Reason]: \"log\" is contained in the URL")
		}
	}
	

	fmt.Println("\r\n[*]\tSubmitting <forms>\r\n")	//-------------")
	
	// if h.debug {
	// 	fmt.Println("\r\n--------------------------------")
	// 	fmt.Println("Handle Submission:\r\n/////////////////////////////////////////////////")
	// }
	// var target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))

	for i,_ := range extforms {
		// if strings.Contains(extforms[i].src,"testimonial") {
		// 	fmt.Println("AAAAAAAAAA\r\nAAAAAAAAAAAAAA\r\nAAAAAAAAAAAAAA\r\nAAA")
		// 	fmt.Println(extforms[i].src)
		h.handleSubmission(&extforms[i])						//	Check if this actually works
		// fmt.Println("AAAAAAAAAA\r\nAAAAAAAAAAAAAA\r\nAAAAAAAAAAAAAA\r\nAAA")

		// } else {
		// 	fmt.Println("NOPE")
		// }

	}

	h.checkforuqstrings(urls)									//	@TODO	Swap within for a method that submits cookies
/*	//	@UNCOMMENT
	// fmt.Println("--------------------------------")
	// fmt.Println("Checking application for strings submitted")
	// h.checkforuqstrings(urls)	//"http://" + h.target + h.targetport + urls)
	// for q := 20; q > 0; q-- {
	// 	fmt.Println(q)
	// 	h.checkforuqstrings([]string{"/"})		//	Comment this out
	// }
	// fmt.Println("--------------------------------")
*/
	if len(possibleinjections) > 0 {

		var location string = h.outputFolder + "/form_injection_detection.txt"
		fmt.Println("\r\n[*]\tWriting <form> Submission Results to file:\t" + location + "\r\n")
		u.SaveStringToFile(location, possibleinjections)
	}

	// fmt.Println(")))))))))))))))))))))))))(((((((((((((((((((((((((((((",len(possibleinjections))
	fmt.Println("\r\n")
}

func (h *InjectionHandler) getwithcookies(urltoget string) string {
	// var target string = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport) + urltoget)

	
	// fmt.Println("GET WITH COOKIES", urltoget)
	rurl,errurl := url.ParseRequestURI(urltoget)	//h.httpprefix +  + urltoget)	//	h.httpprefix + form.action)	//	targethost)
	if errurl != nil {
		fmt.Printf("\r\nparsing error - %s", errurl, "\r\n")
	}

	urlStr := rurl.String()
	
	client := &http.Client{}

	// fmt.Println("URLSTR",urlStr)
	// fmt.Println("URLSTR",urlStr)
	// fmt.Println("URLSTR",urlStr)


	// fmt.Println("Get With COOKIES")
	// fmt.Println(http.MethodGet)
	// fmt.Println(urlStr)
	r,e := http.NewRequest(http.MethodGet, urlStr, nil)	//	URL-encoded payload

	if e != nil {
		fmt.Printf("\r\nEEEEEEERRRRR 1 - %s",e,"\r\n")
	}


	var tokens []string = u.StringCookiesToList(h.sessiontokens)
	for _,k := range tokens {
		var token []string = u.SeparateCookie(k)
		r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})

	}

	
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("\r\nEEEEEEERRRRR 2 - %s",err,"\r\n")
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
		fmt.Printf("\r\nEEEEEEERRRRR 3 - %s",e,"\r\n")
	}


	var tokens []string = u.StringCookiesToList(h.sessiontokens)
	for _,k := range tokens {
		var token []string = u.SeparateCookie(k)
		// fmt.Println("Using Token 2\t-\t",token[0],"\t-\t",token[1])
		r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
	}

	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("\r\nEEEEEEERRRRR 4 - %s",err,"\r\n")
	}
	

	//	Extract Body
	body, _ := ioutil.ReadAll(resp.Body)
	
	return string(body)
}

func (h *InjectionHandler) checkforuqstrings(urls []string) {
	
	var uqstrings []string = h.getUQstrings()	//	SAVE [] into var

	if h.debug {
		fmt.Println("uqstrings", uqstrings)
	}
	//	Check 10 times 
	for n:=0; n < 1; n++ {
		for _,v := range urls {
			
			if !strings.Contains(v,"log") {	//	avoid logging out

				//	Check & add if not present - http://
				target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
				// fmt.Println("==============\t",target," - ",v,"\t==============")
				target = target + v
				// fmt.Println("==============\t\t",url)

				
				
					var r string = h.getwithcookiesforuqstrings(target)

					for _,i := range uqstrings {	//	uqstrindex






						// fmt.Println("i\t\t\t",i)
						// fmt.Println(i, "\t-",len(r))
		
						// // fmt.Println(reflect.TypeOf(i))
						// // if strings.Contains(r,`<a href="index.php">Home</a>`) {
						if strings.Contains(r, i) {
							if h.debug {
								fmt.Println("[*]")
								fmt.Println("[*]== FOUND: user controlled string - Potential injection")
								fmt.Println("[*]======= at:", target, "\t-\tResponse Length:", len(r))
								fmt.Println("[*]======= Found Sequence:",i)
								fmt.Println("[*]")
							}
							possibleinjections += "[*]" + "\r\n"
							possibleinjections += "[*]== FOUND: user controlled string - Potential injection" + "\r\n"
							possibleinjections += "[*]======= at: " + target + "\t-\tResponse Length: " + strconv.Itoa(len(r)) + "\r\n"
							possibleinjections += "[*]======= Found Sequence: " + i + "\r\n"
							possibleinjections += "[*]" + "\r\n"
							
							for _,f := range extforms {



								// yes


								for _,val := range f.uqelemstring {
									var tmp []string = strings.SplitN(val,":",2)
									// fmt.Println(tmp[0],tmp[1])
									if tmp[1] == i {
										// uqstrings = append(uqstrings[:uqstrindex], uqstrings[uqstrindex+1:])
										// // fmt.Println("[*]---------------------- Submitted.Param - Submitted.Sequence", tmp[0], "-", tmp[1])
										if h.debug {
											// fmt.Println("[*]======= Submitted.Location:",f.src)
											// fmt.Println("[*]======= Submitted.Parameter:",tmp[0])
											// fmt.Println("[*]======= Submitted.Sequence: ",tmp[1])
											// fmt.Println("[*]======= Form Responsible:\r\n",f.contents)	//	@TODO - Add form.contents
											// fmt.Println("[*]======= Request Submitted:\r\n" + f.request)	//	@TODO - Add Warning?										
											// fmt.Println("[*]=======")
										}
										// fmt.Println("==============\r\n=======\r\n",target,"\t", f.src, "\t",i, "\r\n=======\r\n=======\r\n")
										
										fmt.Println("[*]== Found: [user controlled string]: ",i,"| [URL]:", target, "| Submitted [Location-Param]: [" + f.src,"-",tmp[0] +"]")
										possibleinjections += "[*]======= Submitted.Location:" + f.src + "\r\n"
										possibleinjections += "[*]======= Submitted.Parameter:" + tmp[0] + "\r\n"
										possibleinjections += "[*]======= Submitted.Sequence: " + tmp[1] + "\r\n"
										possibleinjections += "[*]======= Form Responsible:\r\n" + f.contents + "\r\n"
										possibleinjections += "[*]======= Request Submitted:\r\n" + f.request
										possibleinjections += "[*]=======" + "\r\n"
									}
								}
							}
							// var tmp []string
							// var uq string
							// for _,v := range form.elements {
							// 	tmp = strings.SplitN(v,":",2)
							// 	uq = u.UniqueStringAlphaNum()
							// 	// fmt.Println("APPENDING", tmp[0] + ":" + uq)
							// 	form.uqelemstring = append(form.uqelemstring, tmp[0] + ":" + uq)
					
								//	@HERE
							if h.debug {
								// fmt.Println("[*]\r\n\r\n")
							}
							possibleinjections += "[*]\r\n\r\n\r\n\r\n"
							// fmt.Print(possibleinjections)
						}

						// } else {
						// 	fmt.Println("== Checking for:",i,":\t-\tNope")
						// }
					}
				
			
			}
		}
	}
}

//	@TODO CHECK getUQstrings
func (h *InjectionHandler) getUQstrings() []string {
	var list []string
	for _,f := range extforms {
		for _,k := range f.uqelemstring {
			var t []string = strings.Split(k,":")
			list = append(list,t[1])
		}
	}
	return list
}

func (h *InjectionHandler) urlunique(str string) bool {
	var r bool = true
	for _,form := range extforms {
		if strings.Contains(form.action, str) {
			r = false
		}
	}
	return r
}

func (h *InjectionHandler) injRequestURLi(url string) {

	//	Check & add if not present - http://
	target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
	if h.debug {
		fmt.Println("==============\t",target," - ",url,"\t==============")
	}
	target = target + url
	// fmt.Println("==============\t\t",url)

	var r string

	r = h.getwithcookies(target)
	
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
			tmpforms = h.extractForms(r, h.httpprefix + h.target + url)
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

func (h *InjectionHandler) handleSubmission(f *extractedform) {
	// //	Assuming that we are working with an array of extractedforms
	// var tstform extractedform = extractedform{
	// 	method: http.MethodPost,
	// 	action: "#",
	// 	enctype: "",
	// 	elements: []string{"a:aaa", "b:bbb", "c:ccc"},
	// 	contents: "THESE ARE THE CONTENTS",
	// }
	// fmt.Println("UNCOMMENT ME")
	// submitform("http://127.0.0.1", tstform)
	
	//submitform("http://127.0.0.1", f)
	// var targethost = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))

	// fmt.Println("Submitting:\r\n-------")
	// fmt.Println("\tSRC\t-", f.src)
	// fmt.Println("\tMET\t-", f.method)
	// fmt.Println("\tACT\t-", f.action)
	// fmt.Println("\tENC\t-", f.enctype)
	// fmt.Println("\tELM\t-", f.elements)
	// fmt.Println("\tUQL\t-", f.uqelemstring)
	// fmt.Println("\tCON\t-", f.contents)
	// fmt.Println("-------")

	// fmt.Println("CALLING submitform()")
	h.submitform(f)


}

var scounter, fcounter int
//	Submits the given form using the autogenerated values for the fields.
//	Returns the fields submitted
func (h *InjectionHandler) submitform(form *extractedform) {
	
	// fmt.Println("IN submitform()")
	// if !strings.Contains(form.action, "post-testimonial.php") {
	// 	fmt.Println("IN submitform() - RETURNING 1")

	// 	return
	// } else {
	// 	fmt.Println("IN submitform() - YES")
	// }
	
	// fmt.Println("IN submitform() - Contains")

	data := url.Values{}

	if len(form.elements) > 0 {

		// fmt.Println("IN submitform() - IF LEN > 0")

		var tmp []string
		var uq string
		for _,v := range form.elements {
			tmp = strings.SplitN(v,":",2)
			uq = u.UniqueStringAlphaNum()
			// fmt.Println("APPENDING", tmp[0] + ":" + uq)
			form.uqelemstring = append(form.uqelemstring, tmp[0] + ":" + uq)

			//	@HERE

			if tmp[0] != "submit" {
				data.Set(tmp[0],uq)
			} else {
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

	var tokens []string = u.StringCookiesToList(h.sessiontokens)
	for _,k := range tokens {
		var token []string = u.SeparateCookie(k)
		if h.debug {
			// fmt.Println("Using Token 1\t-\t",token[0],"\t-\t",token[1])
		}
		r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
	}

	//	Show Param info
	// fmt.Println("\t", urlStr)
	// for j,jj := range data {
	// 	fmt.Println("\t", j, "\t\t", jj)
	// }


	// // fmt.Println(r.Header)
	// fmt.Println("\r\n++++++++++++++++++++++++++\r\n",reflect.TypeOf(r.Header), "\r\n++++++++++++++++++++++++++\r\n")
	// fmt.Println(")))))))))", r.Method, "((((((((")
	// pp,_ := url.Parse(r.URL.String())
	// fmt.Println(")))))))))PATH\t", pp.Path, "((((((((")
	// fmt.Println(")))))))))", r.RequestURI, "((((((((")
	// fmt.Println(")))))))))", r.Proto, "((((((((")
	// // fmt.Println(")))))))))encoding\t", r.TransferEncoding, "((((((((")
	// fmt.Println(")))))))))URL\t", r.URL, "((((((((")
	// fmt.Println(")))))))))host\t", r.Host, "((((((((")

	// // fmt.Println(")))))))))", r.Header, "((((((((")
	// for kk, vv := range r.Header {
	// 	fmt.Println(")))))))))", kk, "((((((((",vv)
	// }
	// fmt.Println(")))))))))", data.Encode(), "((((((((")


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
		// fmt.Println(requestString)
		// fmt.Println(form.request)
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

}

func (h *InjectionHandler) combinedURLs() []string {

	////////////////////////////////////////////////////////////////////////////////////////
	//	Filter Gobuster output
	var gobusterOutURL string = h.outputFolder + "/links_gobuster_and_rel.txt"	//	"/root/Desktop/report/gobuster-URLs"

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

//	ParseGobuster Filters out results that are of (Status: 403)
//	Returns an array of lines.
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
	//
	//
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
			// fmt.Println("method is:\t", method)
			f.method = strings.ToUpper(method)
		}

		enctype, okenctype := form.Attr("enctype")
		if okenctype {
			// fmt.Println("enctype is:\t", enctype)
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


			// f.elements = 

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
			
			// fmt.Println("GOT EEEEEEEEEEEEEEEEEEEEEEET]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]")
			// fmt.Println(txtareaname + `:` + `BBBBB` + `&`)
			// // f.elements = append(f.elements, txtareaname + `:` + `BBBBB` + `&`)
			// fmt.Println("GOT EEEEEEEEEEEEEEEEEEEEEEET]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]")
			
			/*	//	@UNCOMMENT	

			//	APPARENTLY THIS IS NOT EXECUTED
			// os.Exit(0)

			*/

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
		// h.formsoftCheck(&f)
		if !(&f == nil) {
			forms = append(forms, f)
		}
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