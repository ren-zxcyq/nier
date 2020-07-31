// Package injectiondetector is responsible for scraping the target website,
// extracting <form> tags, filter for unique forms, submitting all of them
// and identifying user controlled input which appears on the application pages.
// interface their execution with parsing.
package utilities

import (
	"fmt"
	"strings"
	// "path"
	// "net/http"
	// "io/ioutil"
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
var extforms []extractedform
var target string

type InjectionHandler struct {
	target			string
	targetport		int
	// targethost		string	//	A concatenation of proto://target:targetport
	sessiontokens	string
	httpprefix		string
}


//
func NewInjectionHandler(target string, targetport int, stokens string) *InjectionHandler {

	//	Create InjectionHandler
	var h InjectionHandler = InjectionHandler{target: target, targetport: targetport, sessiontokens: stokens, httpprefix: "http://"}

	//fmt.Printf("Address of InjectionHandler - %p", &h) //	Prints the address of the Handler
	return &h
}

func (h *InjectionHandler) InjFormCheck() {
	
	var urls []string

	// urls = append(urls, "http://192.168.1.20/vehical-details.php?vhid=2")
	// urls = append(urls, "http://192.168.1.20/robots.txt")
	// urls = append(urls, "http://192.168.1.20/contact-us.php")
	urls = h.gobusterURLs()


	for _,url := range urls {
		h.injRequestURLi(url)
	}

	
	fmt.Println("--------------------------------")
	fmt.Println("Forms Identified:")
	for _,j := range extforms {
		fmt.Println(j.contents)
	}
	// fmt.Println(extforms)
	fmt.Println("--------------------------------")

	fmt.Println("--------------------------------")
	fmt.Println("Handle Submission:")
	// var target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
	for i,_ := range extforms {
		// h.injRequestURLi(url)
		//	parse target & src
		// h.handleSubmission(form)
		h.handleSubmission(&extforms[i])
	}

	fmt.Println("--------------------------------")
	fmt.Println("Checking application for strings submitted")
	h.checkforuqstrings(urls)//"http://" + h.target + h.targetport + urls)
	fmt.Println("--------------------------------")

}

func (h *InjectionHandler) checkforuqstrings(urls []string) {
	var uqstrings []string = h.getUQstrings()
	// var urls []string
	// for _,form := range extforms {
	// 	// fmt.Println(form.action)

	// 	fmt.Println("Checking URL\t-\t", form.action,"\t-\t", len(form.action))
	// 	for _,form := range extforms {
	// 		if h.urlunique(form.action) || len(urls) == 0 {
	// 			urls = append(urls, form.action)
	// 		}
	// 	}
	// }
	// fmt.Println("=============Actual URLs Requested\t-\t",urls)
	for _,v := range urls {
		// fmt.Println("=============Actual URLs Requested\t-\t",v)
		r := t.WrappedGet(target + v)
		fmt.Println("Target:\t",target + v)
		fmt.Println("Resp Len:\t",len(r))

		for _,i := range uqstrings {
			if strings.Contains(r, i) {
				fmt.Println("== Checking for:",i,":\t-\tFOUND injection")
			} else {
				fmt.Println("== Checking for:",i,":\t-\tNope")
			}
		}
	}
}

func (h *InjectionHandler) getUQstrings() []string {
	var list []string
	for _,f := range extforms {
		for _,k := range f.uqelemstring {
			var t []string = strings.Split(k,":")
			// fmt.Println("-",t[1],"-")
			list = append(list,t[1])
		}
	}
	// for _,i := range list {
	// 	fmt.Println("list\t-\t", i)
	// }
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
	// var or string
	// var extforms []extractedform
	var tmpforms []extractedform

	// r = t.WrappedGet(url)			//	h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))

	var r string

	//	Check & add if not present - http://
	target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
	fmt.Println("==============\t\t",target)
	fmt.Println("==============\t\t",url)

	r = t.WrappedGet(target)	//	h.target + url)
	// or = r

	if strings.Contains(r, "Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking") {	//	or
		// fmt.Println("-------------")
		fmt.Println("[+]\tUpgrading to HTTPS\t",url)
		h.httpprefix = "https://"
		// tsec := utilities.NewHTTPShandler()
		//	Swap http:// prefix to https://
		
		// fmt.Println("HTTPS test\r\n",h.RequestMethodStatus("OPTIONS", target))
		// fmt.Println("_______________________", h.Robots(target))
		// fmt.Println("_______________________", h.Head(target))
		/*
		tester := utilities.NewHTTPShandler()
		tester.TestHTTPS(h.e.targetHost)
		tester.Robots(h.e.targetHost)
		*/
		// fmt.Println("-------------")
	} else {								//	@TODO	consider checking for another error
		// fmt.Println("-------------")
		fmt.Println("[+]\tContinue HTTP\t",url)
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
			
			


			//	CHECK IF THEY ARE UNIQUE

			// fmt.Println("YES")
			// r := strings.Split(or, "<form")
			// for i,j := range r {
			// 	fmt.Println(i, "\t-\t", j)
			// 	//	This WORKS
				
			// }



		} else {
			// fmt.Println("NO")
			//	A <form> is not contained in the page
		}
		// fmt.Println("-------------")
	}
	// fmt.Println("-------------")
	// fmt.Println(results)
	// fmt.Println("-------------")

	

	/////////////////////////////////////////////////////////////

	/////////////////////////////////////////////////////////////

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

func uniquenesscheck(tocheck []extractedform) {					//, list
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
	h.submitform(f)


	//	Currently results in the following
		// root@kali:~# nc -lvnp 9999
		// listening on [any] 9999 ...
		// connect to [127.0.0.1] from (UNKNOWN) [127.0.0.1] 55944
		// POST / HTTP/1.1
		// Host: 127.0.0.1:9999
		// User-Agent: Go-http-client/1.1
		// Content-Length: 17
		// Accept-Encoding: gzip

		// a=aaa&b=bbb&c=ccc
}

type extractedform struct {	//	@TODO	add src
	src				string
	method			string
	action			string
	enctype			string
	
	elements		[]string
	uqelemstring	[]string

	contents		string
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
		// fmt.Println("\tEntire Form:\r\n", formhtml)
		f.contents = formhtml
		

		//	Get <form> Attributes
		action, okaction := form.Attr("action")
		if okaction {	//	basically set form.action to be used by submitform()
			// fmt.Println("action is\t", action)
			if strings.HasPrefix(action, "http://") || strings.HasPrefix(action, "https://") {
				f.action = action
			} else {
				//	If action="#" || action="/test.py"	=>	concatenate with source
				//	If action=""						=>	just use source
				f.action = f.src + action
			}
			//	If action="http://"					=>	use action	(but probably just drop)
		} else {
			f.action = f.src
		}

		method, okmethod := form.Attr("method")
		if okmethod {
			// fmt.Println("method is:\t", method)
			f.method = method
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
	// fmt.Println("=====================\r\n-------------\r\n",len(forms),"\r\n-------------\r\n=====================\r\n")
	return forms
	// htmlResult, _ := doc.Html()
	// fmt.Println(htmlResult)
	//------------------------------------------------------------------------



	// inputsSelector := new(goquery.Selection)
	// inputsSelector = visitNodes(inputsSelector, s, "input")
	// //n := inputsSelector
	// nT := inputsSelector.Text()
	// nt,_ := inputsSelector.Attr("type")


	// buttonsSelector := new(goquery.Selection)
	// buttonsSelector = visitNodes(buttonsSelector, s, "button")
	// // fmt.Println(buttonsSelector.Size())
	// b := buttonsSelector
	// bT := buttonsSelector.Text()
	// bt,_ := b.Attr("type")
	// bn,_ := b.Attr("name")

	// button := s.Find("button")
	// bText := s.Find("button").Text()
	// btype,_ := button.Attr("type")
	// bname,_ := button.Attr("name")
}

func (h *InjectionHandler) formsoftCheck(f *extractedform) *extractedform {
	var keywords []string = []string{"delete", "remove", "edit", "log"}
	for _,word := range keywords {
		if strings.Contains(f.contents, word) {
			// fmt.Println("\r\n[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[\r\niT does contain one of the keywords\r\n[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[\r\n",f.contents)
			// *f = nil
			*f = extractedform{}
			break
		} else {
			// fmt.Println("\r\n[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[\r\niT does not contain any of the keywords\r\n[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[")
		}
	}
		return f
}

//	Submits the given form using the autogenerated values for the fields.
//	Returns the fields submitted
func (h *InjectionHandler) submitform(form *extractedform) {
	
	//	for each element with a name value - create a unique string
	// vvv := u.UniqueStringAlphaNum()
	// fmt.Println(vvv)

	// fmt.Println("tar\t", targethost)
	// fmt.Println("form.action\t",form.action)
	
	//	targethost -> https://127.0.0.1:9999/resource/

	// // apiUrl := "https://api.com"
	// apiUrl := "http://127.0.0.1:9999"
	// // resource := "/user/"


    data := url.Values{}
	/////////////////////////////////////////////////
	elementsnumber := len(form.elements)
	// fmt.Println("submitting elements:", elementsnumber)

	//	Generate a unique string for each submission point
	//	Append the element:uniquestring concatenation to form.uqelemstring\
	//	data.Set(element, uniquestring)
	if elementsnumber > 0 {
		var tmp []string
		for _,v := range form.elements {
			tmp = strings.SplitN(v,":",2)
			uq := u.UniqueStringAlphaNum()
			// fmt.Println("APPENDING", tmp[0] + ":" + uq)
			form.uqelemstring = append(form.uqelemstring, tmp[0] + ":" + uq)
			data.Set(tmp[0],uq)
		}
	} else {
		//	<form> contains no elements
		//	Notify & move on to the next
		return
	}

	fmt.Println("/////////////////////////////////////////////////\r\nSUBMITTING TO:\t", form.action)	//	h.httpprefix + form.action)
	fmt.Println("\r\n",form.contents)
	fmt.Println("\r\n",form.elements)


	rurl,_ := url.ParseRequestURI(form.action)	//	h.httpprefix + form.action)	//	targethost)
	urlStr := rurl.String()
	
	client := &http.Client{}

	r,_ := http.NewRequest(form.method, urlStr, strings.NewReader(data.Encode()))	//	URL-encoded payload
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	for _,k := range u.StringCookiesToList(h.sessiontokens) {
		var token []string = u.SeparateCookie(k)
		fmt.Println("Using Token\t-\t",token[0],"\t-\t",token[1])
		r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
	}

	resp, _ := client.Do(r)
	fmt.Println("/////////////////////////////////////////////////\r\nRESPONSE:\t",resp.Status)

	
	// //	Populate data{}
	// 									// var submitteddata []string
	// if elementsnumber > 0 {
	// 	var tmp []string
	// 	for _,v := range form.elements {
	// 		// submitteddata = append(submitteddata, "AAAAA")
	// 		tmp = strings.SplitN(v,":",2)
	// 		data.Set(tmp[0],tmp[1])
	// 		fmt.Println(tmp[0],":",tmp[1])
	// 	}
	// } else {
	// 	fmt.Println("HELLO")
	// }									
	// /////////////////////////////////////////////////

	// // data.Set("name", "foo")
    // // data.Set("surname", "bar")

    // rurl, _ := url.ParseRequestURI(apiUrl)
    // // rurl.Path = resource
    // urlStr := rurl.String() // "https://api.com/user/"

	// client := &http.Client{}

	// //	http.MethodPost
    // r, _ := http.NewRequest(form.method, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
    // // r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    // // r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    // r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    // resp, _ := client.Do(r)
    // fmt.Println(resp.Status)
}

func (h *InjectionHandler) gobusterURLs() []string {

	////////////////////////////////////////////////////////////////////////////////////////
	//	Filter Gobuster output
	var gobusterOutURL string = "/root/Desktop/report/gobuster-URLs"

	var gobuster string = u.ReturnFileContentsStr(gobusterOutURL)
	res := parseGobuster(gobuster)
	// fmt.Println("******\r\n",gobuster)
	fmt.Println("******\r\n",res, "\r\n******")

	var r []string
	for _,k := range res {
		// fmt.Println(i,"\t---\t","http://192.168.1.20" + k)
		// r = append(r, "http://192.168.1.20" + k)
		r = append(r, k)
	}
	return r

	////////////////////////////////////////////////////////////////////////////////////////
	// fmt.Println("res = ", res)
	// strCont, err := u.StringToLines(res)
	// if err != nil {
	// 	log.Println("Failed while separating lines in formatted tool output")
	// }
	// pdf = h.singlelinetable(pdf, strCont)

	// pdf = h.singlelinetable(pdf, res)
	// return pdf
}

//	ParseGobuster Filters out results that are of (Status: 403)
//	Returns an array of lines.
func parseGobuster(cmdout string) []string {
	
	var extract []string
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		fmt.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {	//	k
		// fmt.Println(k,"-",v)
		if strings.Contains(v,"Status: 403") {
			continue
		} else {
			// extract += string(v)
			if !(len(v) == 0) {
				var tmp []string = strings.SplitN(v," (Status",2)
				extract = append(extract,strings.Trim(tmp[0], " "))
			}
		}
	}
	return u.Uniquestrslice(extract)
}