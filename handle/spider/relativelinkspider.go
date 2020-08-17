// Package relativelinkspider is responsible for scraping relative links
// from the target website, extracting href="" attributes.
package spider

import (
	"fmt"
	"strings"
	// "path"
	// "net/http"
	"io/ioutil"
	// "reflect"
	// "bytes"
	// "os"

	// "net/url"
	// "resource"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/ren-zxcyq/nier/utilities"
)

var extlinks []string
var links string				//	contains both Gobuster_URLs && rel_URLs

var t utilities.Agent
var target string

type RelativeLinkSpider struct {
	target			string
	targetport		int

	sessiontokens	string
	outputFolder	string
	httpprefix		string

	debug			bool	//	@TODO	-	Consider using a similar flag to the other operations
}

type extractedlink struct {
	contents		string
}

func NewRelativeLinkSpider(target string, targetport int, outputFolder string, stokens string) *RelativeLinkSpider {

	//	Create InjectionHandler
	var h RelativeLinkSpider = RelativeLinkSpider{target: target, targetport: targetport, sessiontokens: stokens, outputFolder: outputFolder, httpprefix: "http://"}

	//fmt.Printf("Address of InjectionHandler - %p", &h) //	Prints the address of the Handler
	return &h
}

func (h *RelativeLinkSpider) ReqURLs() {

	fmt.Println("\r\n\r\n[*]\tCrawling URLs provided to discover href attributes\r\n-------------")
	if h.debug {
		fmt.Println("********************************************")
		fmt.Println("fROM REL")
		fmt.Println(h.target)
		fmt.Println(h.targetport)
	
		fmt.Println("Tokens:")
		fmt.Println(h.sessiontokens)
		fmt.Println(h.outputFolder)
		fmt.Println(h.httpprefix)
		fmt.Println("********************************************")	
	}


	var urls []string
	urls = h.gobusterURLs()
	if h.debug {
		// fmt.Println("--------------------------------")
		// fmt.Println("URLs:")	
	}

	for _,url := range urls {
		// fmt.Println(url)
		h.requestURLi(url)
		links += url + "\r\n"
	}

	if h.debug {
		// fmt.Println("--------------------------------")
		// fmt.Println("Links Identified:")
	}

	for i,_ := range extlinks {
		if h.debug {
			// fmt.Println(extlinks[i])
			// h.getURL(extlinks[i])
		}
		links += extlinks[i] + "\r\n"
	}
	if h.debug {
		// fmt.Println("--------------------------------")
	
		// fmt.Println("--------------------------------")
		// fmt.Println("Handle Submission:")
		// // var target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
		// for i,_ := range extlinks {
		// 	// h.injRequestURLi(url)
		// 	//	parse target & src
		// 	// h.handleSubmission(form)
		// 	h.handleSubmission(&extlinks[i])
		// }
	}

	// fmt.Println("\r\n--------------------------------\r\nMerged Links:\r\n")
	var location string = h.outputFolder + "/links_gobuster_and_rel.txt"
	fmt.Println("\r\n[*]\tWriting Merged Link Lists to file:\t" + location + "\r\n")
	fmt.Println(links)
	//	Save to file
	// var u utilities.Utils
	u.SaveStringToFile(location, links)
	// fmt.Println("--------------------------------")

}


func (h *RelativeLinkSpider) requestURLi(url string) {

	var r string

	//	Check & add if not present - http://
	target = t.Urlprefixhttp(h.target + `:` + strconv.Itoa(h.targetport))
	if h.debug {
		fmt.Println("--------------\t\t",target,"\t\t",url)
	}
	if strings.Contains(url,"log") {
		fmt.Println("[*]\t[Skipping URL]",url,"[Reason]: \"log\" is contained in the URL")
		return
	}

	// r = t.WrappedGet(target)	//	h.target + url
	r = h.getURL(&target)

	if h.debug {
		fmt.Println("9((((((((((((((((((0))))))))))))))))))))))))9")
	}
	if strings.Contains(r, "Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking") {	//	or
		// fmt.Println("-------------")
		//	Upgrading to
		fmt.Println("[*]\tHTTPS\t",target,"\t",url)
		h.httpprefix = "https://"
	} else {								//	@TODO	consider checking for another error
		// fmt.Println("-------------")
		//	Continue HTTP
		fmt.Println("[*]\tHTTP\t",target,"\t",url)
		h.httpprefix = "http://"

		//	?filter?


		doc, _ := goquery.NewDocumentFromReader(strings.NewReader((r)))
		doc.Find("*[href]").Each(func(i int, href *goquery.Selection) {
			
			// t,_ := goquery.OuterHtml(href)
			// href.Attr("href")
			// fmt.Println("+++++++=\t",t)

			link, oklink := href.Attr("href")
			if oklink {
				// fmt.Println("=======>\t",link)

				if !strings.HasPrefix(link,"/") && !strings.HasPrefix(link,"#") {	//	&& !strings.HasPrefix(link,"#") 
					//	@TODO	Consider symbol checking
					link = "/" + link
					// extlinks = append(extlinks, "/" + link)
				}

				if !h.isstrinlinks(link) && h.isvalidrellink(link) && !strings.HasPrefix(link,"#") {
					extlinks = append(extlinks, link)
				}
			}
		
		})
	}

}

func (h *RelativeLinkSpider) isstrinlinks(str string) bool {

	var r bool
	for _,v := range extlinks {

		if str == v {	//	v.contents
			r = true
		}
	}

	return r
}

//	Checks if the given string is a common image file or
//	if it contains a colon, common in "http://..." links.
//	Essentially, we want to limit links to local html or link containing files.
func (h *RelativeLinkSpider) isvalidrellink(str string) bool {
	var r bool = true
	var filter []string = []string{":",".png",".jpg", ".jpeg", ".gif", ".ico"}
	for _,v := range filter {

		// fmt.Println("STR IS", str, "\t-\t", v)
		if strings.Contains(str, v) {	//	v.contents
			r = false
		}
		// fmt.Println("r = ",r)
	}
	return r
}

// Submits the given form using the autogenerated values for the fields.
// Returns the fields submitted
func (h *RelativeLinkSpider) getURL(url *string) string {
	
	/////////////////////////////////////////////////
	if h.debug {
		fmt.Println("/////\r\nREQUESTING:\t",*url)
		// rurl,_ := url.ParseRequestURI(url)	//	h.httpprefix + form.action)	//	targethost)
		// urlStr := rurl.String()
	}
	
	client := &http.Client{}

	r,_ := http.NewRequest(http.MethodGet, *url, nil)	//	URL-encoded payload		//	urlStr
	// r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if len(h.sessiontokens) > 0 {
		for _,k := range u.StringCookiesToList(h.sessiontokens) {
			var token []string = u.SeparateCookie(k)
			if h.debug {
				fmt.Println("Using Token\t-\t",token[0],"\t-\t",token[1])
			}
			r.AddCookie(&http.Cookie{Name: token[0], Value: token[1]})
		}
	}

	resp, _ := client.Do(r)
	if h.debug {
		fmt.Println("/////\r\nRESPONSE:\t",resp.Status)
	}
	
	//	Extract Body
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Printf("%s", e)
	}

	//	Body
	return string(body)
}

func (h *RelativeLinkSpider) gobusterURLs() []string {

	////////////////////////////////////////////////////////////////////////////////////////
	//	Filter Gobuster output
	var gobusterOutURL string = h.outputFolder + "/gobuster-URLs"

	var gobuster string = u.ReturnFileContentsStr(gobusterOutURL)
	res := h.parseGobuster(gobuster)
	if h.debug {
		fmt.Println("******\r\n",res, "\r\n******")
	}

	var r []string
	for _,k := range res {
		r = append(r, k)
	}
	return r
}

// ParseGobuster Filters out results that are of (Status: 403)
// Returns an array of lines.
func (h *RelativeLinkSpider) parseGobuster(cmdout string) []string {
	
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
			if !(len(v) == 0) {
				var tmp []string = strings.SplitN(v," (Status",2)
				extract = append(extract,strings.Trim(tmp[0], " "))
			}
		}
	}
	return u.Uniquestrslice(extract)
}