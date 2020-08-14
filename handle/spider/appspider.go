// Package spider is responsible for scraping the page
package spider

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/ren-zxcyq/nier/utilities"
)

var prespiderlinks string

var u utilities.Utils

// var crawledlinks string				//	contains both Gobuster_URLs && rel_URLs

type AppSpider struct {
	targetHost           string
	targetPort           int
	outputFolder         string

	postlinkslocation	 string
	prelinkslocation	 string
	httpprefix			 string

	debug			bool	//	@TODO	-	Consider using a similar flag to the other operations
}

func NewAppSpider(targetH string, targetP int, outFolder string) *AppSpider {

	//	Create an elementsHandler Object to be passed to the exported execHandler
	var h AppSpider = AppSpider {
		targetHost:           	targetH,
		targetPort:           	targetP,
		outputFolder:         	outFolder,
		prelinkslocation:		outFolder + "/prespiderlinks.txt",
		postlinkslocation:		outFolder + "/gospider_out",
		httpprefix:				"http://",
	}

	return &h

}

func (h *AppSpider) Prepare() {
	fmt.Println("\r\n\r\n[*]\tPreparing URL list to launch spider\r\n-------------")
	if h.debug {
		fmt.Println("********************************************")
		fmt.Println("fROM REL")
		fmt.Println(h.targetHost)
		fmt.Println(h.targetPort)
	

		fmt.Println(h.outputFolder)
		fmt.Println(h.httpprefix)
		fmt.Println("********************************************")	
	}


	var urls []string
	urls = h.gobusterandrelURLs()

	for _,j := range urls {	//	i
		// fmt.Println(i,"\t-\t",j)
		prespiderlinks += j + "\r\n"
	}

	//	Save to File -> Essentially the list used by GoSpider
	u.SaveStringToFile(h.prelinkslocation, prespiderlinks)
}

// Function Organise is responsible for parsing GoSpider Results and organising
// findings by type. ([form] [url] [javascript] [form] etc.)
func (h *AppSpider) Organize() {

	fmt.Println("\r\n\r\n[*]\tFiltering Elements Identified\r\n")
	var urls []string
	var formsexistin_URL []string
	var uploadform_URL []string
	var javascript_URLs []string
	var linkfinder_URLs []string
	var subdomain_URLs []string

	var gospiderresults []string = u.ReturnLinesFromFile(h.postlinkslocation)
	for _,line := range gospiderresults {
		if strings.HasPrefix(line,"[url]") {
			// fmt.Println("[url]")
			urls = append(urls,line)
		} else if strings.HasPrefix(line,"[form]") {
			// fmt.Println("[form]")
			formsexistin_URL = append(formsexistin_URL,line)
		} else if strings.HasPrefix(line,"[javascript]") {
			// fmt.Println("[javascript]")
			javascript_URLs = append(javascript_URLs,line)
		} else if strings.HasPrefix(line,"[upload-form]") {
			// fmt.Println("[upload-form]", line)
			uploadform_URL = append(uploadform_URL,line)
		} else if strings.HasPrefix(line,"[linkfinder]") {
			// fmt.Println("[linkfinder]", line)
			linkfinder_URLs = append(linkfinder_URLs,line)
		} else if strings.HasPrefix(line,"[subdomain]") {
			subdomain_URLs = append(subdomain_URLs,line)
		} else {
			fmt.Println("\t[ignored]\t", line)
		}
	}

	// fmt.Println("----------------------")
	// fmt.Println("before being unique")
	// fmt.Println("----------------------")

	// fmt.Println(len(urls))
	// fmt.Println(len(formsexistin_URL))
	// fmt.Println(len(uploadform_URL))
	// fmt.Println(len(linkfinder_URLs))
	// fmt.Println(len(javascript_URLs))


	urls = h.uniqueslice(urls)
	formsexistin_URL = h.uniqueslice(formsexistin_URL)
	uploadform_URL = h.uniqueslice(uploadform_URL)
	linkfinder_URLs = h.uniqueslice(linkfinder_URLs)
	javascript_URLs = h.uniqueslice(javascript_URLs)
	subdomain_URLs = h.uniqueslice(subdomain_URLs)

	// fmt.Println("----------------------")
	// fmt.Println("after uq")
	// fmt.Println("----------------------")

	// fmt.Println(len(urls))
	// fmt.Println(len(formsexistin_URL))
	// fmt.Println(len(uploadform_URL))
	// fmt.Println(len(linkfinder_URLs))
	// fmt.Println(len(javascript_URLs))


	// Convert to string
	// SaveToFile() requires a string where each line is separated by \r\n
	var urls_str string
	var formsexistin_URL_str string
	var uploadform_URL_str string
	var javascript_URLs_str string
	var linkfinder_URLs_str string
	var subdomain_URLs_str string


	fmt.Println("[*]\tURLs identified")
	for _,i := range urls {
		fmt.Println("\t",i)
		urls_str += i + "\r\n"
	}
	fmt.Println("[*]\t<forms> identified")
	for _,i := range formsexistin_URL {
		fmt.Println("\t",i)
		formsexistin_URL_str += i + "\r\n"
	}
	fmt.Println("[*]\tUpload <forms> identified")
	for _,i := range uploadform_URL {
		fmt.Println("\t",i)
		uploadform_URL_str += i + "\r\n"
	}
	fmt.Println("[*]\tlinkfinder URLs identified")
	for _,i := range linkfinder_URLs {
		fmt.Println("\t",i)
		linkfinder_URLs_str += i + "\r\n"
	}
	fmt.Println("[*]\tjavascript URLs identified")
	for _,i := range javascript_URLs {
		fmt.Println("\t",i)
		javascript_URLs_str += i + "\r\n"
	}
	fmt.Println("[*]\tSubdomains identified")
	for _,i := range subdomain_URLs {
		fmt.Println("\t",i)
		subdomain_URLs_str += i + "\r\n"
	}

	fmt.Println("\r\n\r\n[*]\tSave Element Lists into files /spider_[element].list\r\n")
	// Save to Files spider_element.list
	u.SaveStringToFile(h.outputFolder+"/spider_URLs.list",urls_str)
	u.SaveStringToFile(h.outputFolder+"/spider_Forms.list",formsexistin_URL_str)
	u.SaveStringToFile(h.outputFolder+"/spider_UploadForms.list",uploadform_URL_str)
	u.SaveStringToFile(h.outputFolder+"/spider_LinkFinder.list",linkfinder_URLs_str)
	u.SaveStringToFile(h.outputFolder+"/spider_JavascriptFiles.list",javascript_URLs_str)
	u.SaveStringToFile(h.outputFolder+"/spider_Subdomains.list",subdomain_URLs_str)

}

// Consider removing pre-allocation of cap([]string)
func (h *AppSpider) uniqueslice(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}	//	Either this or prepopulate the slice length. Prepopulation is more effective in large slices (Below)
	// list := make([]string,len(strSlice))		//	If I pre-populate the slice length, returned sliceSSSS is longer.
	for _,entry := range strSlice {
		if _,found := keys[entry]; !found {
			keys[entry] = true
			list = append(list,entry)
		}
	}
	return list
}

func (h *AppSpider) isstrinlinks(str string) bool {

	var r bool
	for _,v := range extlinks {

		if str == v {	//	v.contents
			r = true
		}
	}

	return r
}

// //	Checks if the given string is a common image file or
// //	if it contains a colon, common in "http://..." links.
// //	Essentially, we want to limit links to local html or link containing files.
// func (h *AppSpider) isvalidrellink(str string) bool {
// 	var r bool = true
// 	var filter []string = []string{":",".png",".jpg", ".jpeg", ".gif", ".ico"}
// 	for _,v := range filter {

// 		// fmt.Println("STR IS", str, "\t-\t", v)
// 		if strings.Contains(str, v) {	//	v.contents
// 			r = false
// 		}
// 		// fmt.Println("r = ",r)
// 	}
// 	return r
// }

// func (h *AppSpider) uniquenesscheck(tocheck []string) {					//	tocheck []extractedlink

// 	if len(extlinks) > 0 {

// 		for _,v := range tocheck {
// 			if !isstrinforms(v) {	//	v.contents
// 				//	Appending
// 				extlinks = append(extlinks, v)
// 			} else {
// 				//	Not Appending
// 			}
// 		}

// 	} else {
// 		extlinks = append(extlinks, tocheck[0])
// 		h.uniquenesscheck(tocheck[1:])
// 	}
// }


func (h *AppSpider) gobusterandrelURLs() []string {

	//	Filter output
	var linksretrievedURL string = "/root/Desktop/report/links_gobuster_and_rel.txt"

	var linksretrieved string = u.ReturnFileContentsStr(linksretrievedURL)
	res := h.parseRelandGobusterURLs(linksretrieved)

	var r []string
	for _,k := range res {
		r = append(r, k)
	}
	return r
}

// parseRelandGobusterURLs Filters out results that are of (Status: 403)
// Returns an array of lines.
func (h *AppSpider) parseRelandGobusterURLs(cmdout string) []string {
	
	var extract []string
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		fmt.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {

		if strings.Contains(v,"log") {	//	@FILTER
			continue
		} else {
			// extract += string(v)
			if !(len(v) == 0) {
				// // var tmp []string = strings.SplitN(v," (Status",2)
				// extract = append(extract,strings.Trim(tmp[0], " "))
				extract = append(extract, h.httpprefix + h.targetHost + ":" + strconv.Itoa(h.targetPort) + v)
			}
		}
	}
	return u.Uniquestrslice(extract)
}