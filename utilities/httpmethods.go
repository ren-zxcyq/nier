// Package utilites holds package agnostic functions.
// httpmethods.go defines an Agent struct, which can be used to perform
// http requests.
//
// Example calls can be seen below:
// host := "http://192.168.1.20"
// port := 80
// // var tar string = host + ":" + string(port)
// var tar string = host + ":" + strconv.Itoa(port)
// var u utilities.Utils
// u.EncodingTest()
// fmt.Println("-------------")
// var a utilities.Agent
// //a.Robots("http://www.google.com")
// //a.Head("http://192.168.1.20")
// //a.OptionsRequest("http://192.168.1.20")
// //a.OptionsVerify("http://192.168.1.20")
// a.OptionsVerify(tar)
// fmt.Println("-------------")
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
	url = a.Urlprefixhttp(url)
	fmt.Println("\r\nRetrieving", url, "/robots.txt\r\n-------------\r\nResponse Status:")

	r, e := http.Get(url + "/robots.txt")
	if e != nil {
		if strings.Contains(url, "https://") {
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Error encountered while requesting /robots.txt", e)
	}

	//	r.Status contains the status
	fmt.Println(r.Status)

	//	Extract Body
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return fmt.Sprintf("%s", e)
	}

	//	Body
	// fmt.Println("Contents of robots.txt\r\n-------------\r\n", string(body))
	return string(body)
}

func (a *Agent) Head(url string) string {
	url = a.Urlprefixhttp(url)
	r, e := http.Head(url)

	if e != nil {
		if strings.Contains(url, "https://") {
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Error encountered while requesting HEAD, ",url, " - ", e)

	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return fmt.Sprintf("%s", e)
	}

	return fmt.Sprint(string(body))
}

func (a *Agent) RequestMethod(method, url string) string {
	url = a.Urlprefixhttp(url)
	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		if strings.Contains(url, "https://") {
			return fmt.Sprint("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		return fmt.Sprint("Request Method -", e)

	}
	// fmt.Println(r)					//	Maybe Add them in one var
	// fmt.Println(r.Status)
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return fmt.Sprintf("%s", e)
	}
	return string(body)
}

func (a *Agent) RequestMethodStatus(method, url string) string {
	url = a.Urlprefixhttp(url)
	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		if strings.Contains(url, "https://") {
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
	url = a.Urlprefixhttp(url)
	// // fmt.Println(url)
	// fmt.Println("\r\nHTTP OPTIONS Request - Retrieve Supported HTTP Methods\r\n-------------")	//	\r\nResponse Status:")
	var r string
	r = a.RequestMethod("OPTIONS", url)
	// fmt.Println("Response to the OPTIONS HTTP Request:\r\n", r)
	return fmt.Sprintln("Response to the OPTIONS HTTP Request:\r\n", r)
}

func (a *Agent) OptionsVerify(url string) []string {
	url = a.Urlprefixhttp(url)
	var rs string
	var res []string
	options := make([]string, 7)                                                  //	9
	options = []string{"CONNECT", "GET", "HEAD", "PATCH", "POST", "PUT", "TRACE"} //	"DELETE", "OPTIONS",

	// fmt.Println("\r\nSupported Options:\r\n-------------")
	for i := 0; i < len(options); i++ {
		rs = a.RequestMethodStatus(options[i], url)
		// fmt.Println(options[i], "-", res)	//	return and parse on caller
		res = append(res, options[i] + "-" + rs)
	}
	return res
}

// Get URL, check for http:// or https:// prefix
// add http:// if not present.
// Exit program if url contains https:// until implemented
func (a *Agent) Urlprefixhttp(url string) string {
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "http://" + url
	} else if !strings.Contains(url, "http://") && strings.Contains(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
		url = "http://" + url
	} else if strings.Contains(url, "http://") {
		//	Leave unchanged
	}
	return url
}

func (a *Agent) WrappedGet(url string) string {
	url = a.Urlprefixhttp(url)
	// fmt.Println("\r\nRetrieving", url, "\r\n-------------\r\n")	//	Response Status:")

	r, e := http.Get(url)
	if e != nil {
		if strings.Contains(url, "https://") {
			fmt.Println("Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking - Error here is:\r\n", e)
		}
		fmt.Println("Error encountered while requesting", url, e)
	}

	//	r.Status contains the status
	fmt.Println("[+]\t\t",r.Status)

	//	Extract Body
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Printf("%s", e)
	}

	//	Body
	return string(body)
}