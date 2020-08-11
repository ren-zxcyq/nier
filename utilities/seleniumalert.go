package utilities

import (
	// "os"
	"fmt"
	"time"

	"io/ioutil"
	// "bufio"
	// "regexp"
	
	"net/url"
	// "resource"
	"net/http"
	"strconv"
	"strings"


	"github.com/tebeka/selenium"
	"github.com/PuerkitoBio/goquery"
)

// selenium.WebDriver config
var wd selenium.WebDriver
const (
	seleniumPath	= "/opt/tebeka-selenium/selenium/vendor/selenium-server.jar"
	geckoDriverPath = "/opt/tebeka-selenium/selenium/vendor/geckodriver"
	port			= 9999
)
var options []selenium.ServiceOption
var caps selenium.Capabilities


// App Specific
// var targetURL string = "http://192.168.1.20:80/post-testimonial.php"	//	"http://192.168.1.20"
// var deletemalicioustestimonialURL string = "http://192.168.1.20/my-testimonials.php"
// var findalertintestimonialURL string = "http://192.168.1.20"


func (h *InjectionHandler) checkForAlertInPage(url string, identifier string) (search string) {

	// setup selenium.WebDriver
	options := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),			//	X frame buffer for the browser to run in
		selenium.GeckoDriver(geckoDriverPath),	//	GeckoDriver in order to use firefox
		selenium.Output(nil),				//	Output Debug info to STDERR	//	os.Stderr
	}

	selenium.SetDebug(true)

	service, err := selenium.NewSeleniumService(seleniumPath, port, options...)
	if err != nil {
		panic(err)	//	Maybe change panic
	}
	defer service.Stop()


	//	Connect to the WebDriver instance running locally
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	wd.SetPageLoadTimeout(1*time.Second)



	// Navigate to the page
	// if err := wd.Get(targetURL); err != nil {
	// 	panic(err)
	// }

	alertPageURL := url
	fmt.Println("\r\n\r\n\r\n\r\nREQUESTING:\t",alertPageURL,"\r\n\r\n")

	flag := 2
	search:
		for {	//	 countdown:=3;countdown>0;countdown--

			if flag == 0 {
				search = "unsuccessful"
				break search
			}

			if err := wd.Get(alertPageURL); err != nil {
				fmt.Printf("wd.Get(%q) returned error: %v\r\n", alertPageURL, err)
			}

			alertText, err := wd.AlertText()
			if err != nil {
				fmt.Printf("wd.AlertText() returned error: %v\r\n", err)
			}

			fmt.Println("\r\n\r\n\r\n\r\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>",alertText,"\t",len(alertText),"\r\n\r\n\r\n\r\n")
			if len(alertText) == 0 {
				// wd.SetImplicitWaitTimeout()
				time.Sleep(2*time.Second)
				flag -= 1
				// continue	//	?set up a flag?
			} else {

				if alertText != identifier {
					if err := wd.AcceptAlert(); err != nil {
						fmt.Println("wd.AcceptAlert() returned error: %v", err)
					}
					// fmt.Println("\r\n\r\n\r\n=======================\r\n<SCRIPT>ALERT<> WITHOUT OUR FLAG\r\n=======================\r\n\r\n")
					time.Sleep(1*time.Second)
					continue
				} else {	//	implying alertText == id
					fmt.Println("\r\n\r\n\r\n=======================\r\n<SCRIPT>ALERT<> WITH OUR FLAG\r\n=======================\r\n\r\n")
					search = alertText	//	= identifier
					fmt.Println("\r\n\r\n\r\n\r\n\r\n\r\n\r\n>>>>>>>>>>>>>>>>>>>>Alert is Active - str:", alertText, "<<<<<<<<<<<<<<<<<<<<<\r\n\r\n\r\n\r\n")
					// os.Exit(123)
					break search
				}
			}
		}
	return


	/*
		if alertText == "" {
			wd.WaitWithTimeoutAndInterval(condition, 3*time.Second, 1*time.Second)
		} else if alertText != "1" {
			fmt.Printf("Expected '1' but got '%s'\r\n", alertText)
		} else {
			fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\r\nIT WORKED!\r\n^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		}
	*/
	///////////////////////////////
	// activeelement,_ := wd.ActiveElement()
	// fmt.Println(activeelement)

	// flag,_ := activeelement.IsDisplayed()
	// fmt.Println(flag)

	// contents,_ := activeelement.Text()
	// fmt.Println(contents)

	///////////////////////////////

	// if err := wd.AcceptAlert(); err != nil {
	// 	fmt.Printf("wd.AcceptAlert() returned error: %v", err)
	// }

	// os.Exit(1)
	
	
}

// function removeelements connects to 192.168.1.20/admin/, logs in using admin:april
// goes to the managetestimonials page and deactivates every single testimonial
func (h *InjectionHandler) removeelements() {
	// GET WITH COOKIES urltoget
	rurl,errurl := url.ParseRequestURI("http://192.168.1.20/admin/")	//h.httpprefix +  + urltoget)	//	h.httpprefix + form.action)	//	targethost)
	if errurl != nil {
		fmt.Printf("\r\nparsing error - %s", errurl, "\r\n")
	}

	urlStr := rurl.String()	
	client := &http.Client{}

	data := url.Values{}
	data.Set("username","admin")
	data.Set("password","april")
	data.Set("login","")

	// fmt.Println(urlStr)
	// fmt.Println(data)
	// fmt.Println(strconv.Itoa(len(data.Encode())))

	r,e := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))	//	URL-encoded payload
	if e != nil {
		fmt.Printf("\r\nERR 1 - %s",e,"\r\n")
	}

	r.Header.Set("Referer", "http://192.168.1.20/admin/")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("\r\nERR 2 - %s",err,"\r\n")
	}

	//	Extract Body
	body, _ := ioutil.ReadAll(resp.Body)

	// Check if we are logged in
	// fmt.Println(string(body))
	if !strings.Contains(string(body),`document.location = 'change-password.php'`) {
		fmt.Println("NO")
		
		
	} else {

		fmt.Println("YES")

		// Extract cookies
		var cookie []*http.Cookie
		cookie = resp.Cookies()
		for i,j := range cookie {
			fmt.Println(i,"-",j)
		}

		// Use Cookies to get testimonials page
		// "http://192.168.1.20/admin/testimonials.php"
		req_testimonials,_ := http.NewRequest(http.MethodGet,"http://192.168.1.20/admin/testimonials.php", nil)

		req_testimonials.Header.Set("Referer", "http://192.168.1.20/admin/change-password.php")
		req_testimonials.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		req_testimonials.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
		for i,_ := range cookie {
			// fmt.Println(i,"-",j)
			// fmt.Println(cookie[i])
			req_testimonials.AddCookie(cookie[i])
		}
		
		
		res_testimonials,_ := client.Do(req_testimonials)
		if err != nil {
			fmt.Println("Failed while getting testimonials")
		}

		body_testimonials,_ := ioutil.ReadAll(res_testimonials.Body)

		// fmt.Println(string(body_testimonials))

		if !strings.Contains(string(body_testimonials),`User Testimonials`) {
			fmt.Println("testimonials NO")
		} else {
			fmt.Println("testimonials YES")

			var rows [][]string
			var row []string

			//filter
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body_testimonials)))
			if err != nil {
				fmt.Println("Testimonials table: No url found")
				fmt.Println(err)
			}

			// Find each table
			doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
				tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
					
					rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
						row = append(row, tablecell.Text())
					})
					

					rowhtml.Find("a").Each(func(indexa int, tablelink *goquery.Selection) {
						link, oklink := tablelink.Attr("href")
						if oklink {
							// fmt.Println("lllllllllllll",link)
							row = append(row, link)
						}
					})

					rows = append(rows, row)
					row = nil
				})
			})

			for _,rv := range rows {	//	rk
				// fmt.Println("row", rk,":","\t",rv)
				lrv := len(rv)
				if !(lrv == 0) {

					fmt.Println("IS IS NOT ZERO")

					// for _,rsv := range rv {
					// 	fmt.Println("---",rsv)
					// }


					// fmt.Println("key to submit and deactivate", rv[0])
					status := strings.TrimSpace(rv[lrv-2])
					// fmt.Println("status of row:", status)
					if status == "Inactive" {
						fmt.Println("====It's not")
					} else if status == "Active" {
						fmt.Println("====IT IS")


						//	Deactivate row request						"http://192.168.1.20/admin/testimonials.php"		127.0.0.1:9988
						deactivate_testimonial,_ := http.NewRequest(http.MethodGet,"http://192.168.1.20/admin/"+rv[lrv-1], nil)				///admin/testimonials.php
						deactivate_testimonial.Header.Set("Referer", "http://192.168.1.20/admin/testimonials.php")

						for i,_ := range cookie {
							// fmt.Println(cookie[i])
							deactivate_testimonial.AddCookie(cookie[i])
						}

						// fmt.Println("SENDING")
						reqdeactivate_testimonials,_ := client.Do(deactivate_testimonial)
						if err != nil {
							fmt.Println("Failed while getting testimonials")
						}

						// fmt.Println("SENT")
						bodydeactivate_testimonials,_ := ioutil.ReadAll(reqdeactivate_testimonials.Body)
						fmt.Println(reqdeactivate_testimonials.Status,"\t-\t",len(bodydeactivate_testimonials))
					}
				} else {
					fmt.Println("IT IS ZERO")
				}
			}

			// get again?
				
		}
	}
}