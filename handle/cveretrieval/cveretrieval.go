// Package cveretrieval retrieves CVEs related to the banners extracted by
// nmap -script=banners.
// 10 items per banner are returned and saved to a file.
package cveretrieval

import (
	"fmt"
	"strings"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/ren-zxcyq/nier/utilities"
)

type cveRetriever struct {
	outputfolder	string
}

var u utilities.Utils

func NewCVERetriever(outputfolder string) *cveRetriever {
	var h cveRetriever = cveRetriever{outputfolder: outputfolder}

	//fmt.Printf("Address of outputFolderHandler - %p", &h) //	Prints the address of outputFolderHandler
	return &h
}

func (h *cveRetriever) Retrieve() {

	banners := h.retrieveBanners()

	var cves string
	for _,j := range banners {
		cves += "\r\n\r\n\r\n[*] " + j + "\r\n"
		
		var tmp string = h.retrieveforbanner(j)
		if len(tmp) > 0 {
			cves += tmp
		} else {	//	No CVEs were found
			cves += "\r\n\tNo CVEs were discovered."
		}
	}

	// var cves map[string]string
	// for _,v := range banners {
	// 	//cves[v] = 
	// h.retrieveforbanner(v)	
	// }

	// for k,_ := range cves {
	// 	fmt.Println(k)
	// }

	var location string = h.outputfolder + "/cves.list"
	fmt.Println("\r\n[*]\tWriting CVE List Retrieved from NIST NVD:\t" + location + "\r\n")
	var u utilities.Utils
	u.SaveStringToFile(location, cves)

}

func (h *cveRetriever) retrieveforbanner(banner string) string {
	//	https://nvd.nist.gov/vuln/search/results?form_type=Basic&results_type=overview&query=apache&search_type=all
	var nvd = "https://nvd.nist.gov/vuln/search/results?"
	nvd += "form_type=Basic&results_type=overview&query=" + banner
	nvd += "&search_type=all"


	tsec := utilities.NewHTTPShandler()
	var r string
	r = "Banner:\t" + banner + "\r\n"
	r += tsec.RequestMethod("GET", nvd)

	// fmt.Println(r)
	// fmt.Println(h.filter(r))
	return h.filter(r)
}

func (h *cveRetriever) filter(r string) string {
	var retrievednum int = 13	//	We want to retrieve max 10 CVEs
	var cveresult string
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader((r)))
    doc.Find("tr").Each(func(i int, form *goquery.Selection) {
		retrievednum -= 1
		if (retrievednum > 0) {
			_, oktestid := form.Attr("data-testid")
			var cveid string
			if oktestid {
				form.Find("th").Each(func(indexth int, rownamecell *goquery.Selection) {
					cveid = rownamecell.Text()
				})

				var rowelements []string
				form.Find("td").Each(func(indexth int, infocell *goquery.Selection) {
					// ic,_ := goquery.OuterHtml(infocell)
					var rpieces []string = strings.SplitN(infocell.Text(),"Published:",2)

					if len(rpieces) > 1 {
						rowelements = append(rowelements, "\t" + strings.TrimSpace(rpieces[0]))
					} else {
						// <span id="cvss3-link">
						var cvssstrs string
						infocell.Find("span").Each(func(indexth int, spancell *goquery.Selection) {
							id,okid := spancell.Attr("id")
							if okid {
								// fmt.Println(id)
								if id == "cvss3-link" || id == "cvss2-link" {
									cvssstrs += "\tCVSS " + strings.TrimSpace(spancell.Text()) + "\t"
								}
							}
						})
						rowelements = append(rowelements, cvssstrs)
					}
				})
				cveresult += "\r\n\t" + cveid + "\r\n-\r\n"
				for _,re := range rowelements {
					cveresult += "\t" + re + "\r\n-"
				}
			}
		}
	})

		//	Filtering:
		// <tr data-testid="vuln-row-12">

		// <th nowrap="nowrap"><strong><a href="/vuln/detail/CVE-2020-11978"
		// 	 data-testid="vuln-detail-link-12">CVE-2020-11978</a></strong><br></th>
		// 	<td>
		// 	<p data-testid="vuln-summary-12">An issue was found in Apache Airflow versions 1.10.10 and below. A remote code/command injection vulnerability was discovered in one of the example DAGs shipped with Airflow which would allow any authenticated user to run arbitrary commands as the user running airflow worker/scheduler (depending on the executor in use). If you already have examples disabled by setting load_examples=False in the config then you are not vulnerable.</p> <strong>Published:</strong>
		// 	<span
		// 	 data-testid="vuln-published-on-12">July 16, 2020; 8:15:10 PM -0400</span>
		// 	 </td>
		// 	<td nowrap="nowrap">

        //     <span id="cvss3-link">
        //     <em>V3.1:</em> <a
        //      href="/vuln-metrics/cvss/v3-calculator?name=CVE-2020-11978&amp;vector=AV:N/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H&amp;version=3.1&amp;source=NIST"
        //      class="label label-danger" data-testid="vuln-cvss3-link-12">8.8 HIGH</a><br/>
        //      </span>  <span id="cvss2-link"> <em>    V2.0:</em> <a
        //      href="/vuln-metrics/cvss/v2-calculator?name=CVE-2020-11978&amp;vector=(AV:N/AC:L/Au:S/C:P/I:P/A:P)&amp;version=2.0&amp;source=NIST"
        //      class="label label-warning" data-testid="vuln-cvss2-link-12">6.5 MEDIUM</a><br/>
        //      </span>
		// </td>

		return cveresult
}

func (h *cveRetriever) retrieveBanners() []string{
	//open file
	var bannerscontents []string = u.ReturnLinesFromFile(h.outputfolder + "/nmap-banners.nmap")
	//filter out keys
	var banners []string
	for _,j := range bannerscontents {
		// fmt.Println(j)
		if strings.Contains(j,"open") {
			// row := strings.SplitN(j," ",5)
			var banner string
			var tmp string = strings.Join(strings.Fields(j), " ")
			row := strings.SplitN(tmp," ",5)
			
			if strings.Contains(row[4]," (") {	//	unauthorized)
				var tmp2 []string = strings.SplitN(row[4]," (",2)	//	unauthorized)
				banner = strings.TrimSpace(tmp2[0])

			} else {
				banner = strings.TrimSpace(row[4])
			}
			

			// fmt.Println("banner:\t", banner)
			banners = append(banners, url.QueryEscape(banner))
		}
		// banners = append(banners, j)
	}
	//append to []string
	return banners
}