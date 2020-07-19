package utilities

import (
	"fmt"
	"io/ioutil"
	// "log"
	"net/http"
	"strings"
)

type Agent struct{}

func (a *Agent) Robots(url string) string {
	url = a.checkURL(url)
	fmt.Println("\r\nRetrieving", url, "/robots.txt\r\n-------------\r\nResponse Status:")

	r, e := http.Get(url + "/robots.txt")
	if e != nil {
		// fmt.Printf("Error Encountered")
		// log.Println(e)
		if strings.Contains(url, "https://") {
			//log.Println("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking")
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Error encountered while requesting /robots.txt", e)
	}

	//	r.Status contains the status
	fmt.Println(r.Status)

	//	Extract Body
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		// log.Println(e)
		return fmt.Sprintf("%s", e)
	}

	//	Body
	fmt.Println("Contents of robots.txt\r\n-------------\r\n", string(body))
	return string(body)
}

func (a *Agent) Head(url string) string {
	url = a.checkURL(url)
	r, e := http.Head(url)

	if e != nil {
		// fmt.Printf("Error Encountered")
		// log.Println(e)
		if strings.Contains(url, "https://") {
			//log.Println("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking")
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Error encountered while requesting /robots.txt", e)

	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		//log.Println(e)
		return fmt.Sprintf("%s", e)
	}
	// fmt.Println(string(body))

	return fmt.Sprint(string(body))
}

func (a *Agent) RequestMethod(method, url string) string {
	url = a.checkURL(url)
	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		//log.Println("RM", e)
		//	If
		if strings.Contains(url, "https://") {
			//log.Println("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking")
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Request Method -", e)

	}
	fmt.Println(r)
	fmt.Println(r.Status)
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		// log.Println(e)
		return fmt.Sprintf("%s", e)
	}
	// fmt.Println(body)
	return string(body)
}

func (a *Agent) RequestMethodStatus(method, url string) string {
	url = a.checkURL(url)
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

func (a *Agent) OptionsRequest(url string) string {
	// url = a.checkURL(url)
	fmt.Println(url)
	fmt.Println("\r\nHTTP OPTIONS Request - Retrieve Supported HTTP Methods\r\n-------------\r\nResponse Status:")
	var r string
	r = a.RequestMethod("OPTIONS", url)
	// fmt.Println("Response to the OPTIONS HTTP Request:\r\n", r)
	return fmt.Sprintln("Response to the OPTIONS HTTP Request:\r\n", r)
}

func (a *Agent) OptionsVerify(url string) {
	// url = a.checkURL(url)
	options := make([]string, 7)                                                  //	9
	options = []string{"CONNECT", "GET", "HEAD", "PATCH", "POST", "PUT", "TRACE"} //	"DELETE", "OPTIONS",

	fmt.Println("\r\nSupported Options:\r\n-------------")
	var res string
	// @TODO	request each and every one.
	for i := 0; i < len(options); i++ {
		res = a.RequestMethodStatus(options[i], url)
		fmt.Println(options[i], "-", res)
	}

}

// Get URL, check for http:// or https:// prefix
// add http:// if not present.
// Exit program if url contains https:// until implemented
func (a *Agent) checkURL(url string) string {
	if !strings.Contains(url, "http://") || !strings.Contains(url, "https://") {
		url = "http://" + url
	}
	return url
}
