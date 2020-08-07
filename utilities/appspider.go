// Package spider is responsible for scraping the page
package utilities

import (
	"fmt"
	"strings"
	"strconv"

	// "github.com/PuerkitoBio/goquery"


)

var prespiderlinks string
// var crawledlinks string				//	contains both Gobuster_URLs && rel_URLs

type AppSpider struct {
	targetHost           string
	targetPort           int
	outputFolder         string

	httpprefix			 string

	debug			bool	//	@TODO	-	Consider using a similar flag to the other operations
}

func NewAppSpider(targetH string, targetP int, outFolder string) *AppSpider {

	//	Create an elementsHandler Object to be passed to the exported execHandler
	var h AppSpider = AppSpider {
		targetHost:           	targetH,
		targetPort:           	targetP,
		outputFolder:         	outFolder,
		httpprefix:				"http://",
	}

	return &h

}

func (h *AppSpider) Crawl() {

	// // " -Pn -p- -vv -sTV -T5 --script=banner -oA " + filepath.ToSlash(path.Join(h.e.outputFolder, "/nmap-banners")) + " " + h.e.targetHost)
	// ./gospider `-S /root/Desktop/report/urls.txt --depth 0 --no-redirect -t 50 -c 3 -o out --cookie "PHPSESSID:cp2d3bu9sllivtmhcuc6iq0895" --blacklist "log"`

}

func (h *AppSpider) Prepare() {
	fmt.Println("\r\n\r\n[*]\tSpider Launched towards the Application\r\n-------------")
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

	for i,j := range urls {
		fmt.Println(i,"\t-\t",j)
		prespiderlinks += j + "\r\n"
	}

	var location string = h.outputFolder + "/prespiderlinks.txt"
	// fmt.Println("\r\n[*]\tWriting Merged Link Lists to file:\t" + location + "\r\n")
	// fmt.Println(crawledlinks)
	//	Save to file
	// var u utilities.Utils
	u.SaveStringToFile(location, prespiderlinks)
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

//	parseRelandGobusterURLs Filters out results that are of (Status: 403)
//	Returns an array of lines.
func (h *AppSpider) parseRelandGobusterURLs(cmdout string) []string {
	
	var extract []string
	strCont, err := u.StringToLines(cmdout)
	if err != nil {
		fmt.Println("Failed while separating lines in formatted tool output")
	}
	for _,v := range strCont {

		if strings.Contains(v,"log") {	//	Filter
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