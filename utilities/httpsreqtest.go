package utilities

import (
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
)

type HTTPShandler struct{}

func NewHTTPShandler() *HTTPShandler {
	var h HTTPShandler = HTTPShandler{}
	return &h
}

func (h *HTTPShandler) TestHTTPS(target string) {
	//h.actualTest()
	target = "https://" + target
	fmt.Println("HTTPS test\r\n",h.HTTPSRequestMethodStatus("OPTIONS", target))
	fmt.Println("HTTPS Response Contents:\r\n", h.HTTPSOptionsRequest(target))
}

func (h *HTTPShandler) actualTest(link string) {

    transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	client := &http.Client{Transport: transport}
	


	method := "OPTIONS"
    req, rerr := http.NewRequest(method, link, nil)	//client.Get(link)
	if rerr != nil {
        fmt.Println(rerr)
	}
	fmt.Println("Sending OPTIONS")
	response, e := client.Do(req)
	if e != nil {
		fmt.Sprint("Request Method -", e)
		os.Exit(1)
	}
	
	fmt.Println("Sent OPTIONS - ",)




    defer response.Body.Close()

    content, _ := ioutil.ReadAll(response.Body)
    s := string(content)	//strings.TrimSpace(string(content))

    fmt.Println(s)

    // out := s + " world"      
    // Not working POST...
    // resp, err := client.Post(link, "text/plain", &out)

}

func (a *HTTPShandler) HTTPSRequestMethodStatus(method, url string) string {
	// url = a.checkURL(url)
	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		// log.Println(e)
		if strings.Contains(url, "https://") {
			//log.Println("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking")
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Request Method Status - ", e)
	}
	//fmt.Println(r)
	//fmt.Println(r.Status)
	//body, e := ioutil.ReadAll(r.Body)
	//if e != nil {
	//	log.Println(e)
	//}
	// fmt.Println(body)
	return string(r.Status)
}

func (a *HTTPShandler) HTTPSOptionsRequest(url string) string {
	// url = a.checkURL(url)
	fmt.Println(url)
	fmt.Println("\r\nHTTP OPTIONS Request - Retrieve Supported HTTP Methods\r\n-------------\r\nResponse Status:")
	var r string
	r = a.RequestMethod("OPTIONS", url)
	// fmt.Println("Response to the OPTIONS HTTP Request:\r\n", r)
	return fmt.Sprintln("Response to the OPTIONS HTTP Request:\r\n", r)
}

func (a *HTTPShandler) RequestMethod(method, url string) string {
	// url = a.checkURL(url)
	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		//log.Println("RM", e)
		//	If
		// if strings.Contains(url, "https://") {
		// 	//log.Println("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking")
		// 	return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		// }
		return fmt.Sprint("Request Method -", e)

	}
	fmt.Println(r)
	fmt.Println(r.Status)
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		// log.Println(e)
		fmt.Sprintf("%s", e)
	}
	// fmt.Println(body)
	return string(body)
}
