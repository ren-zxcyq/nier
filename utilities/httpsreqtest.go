// Package utilites holds package agnostic functions.
// httpsmethods.go is the equivalent to httpmethods.go for
// https requests.
package utilities

import (
	"fmt"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	// "os"
)

type HTTPShandler struct{
	client		*http.Client
}

func NewHTTPShandler() *HTTPShandler {
	var h HTTPShandler = HTTPShandler{}

	transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	h.client = &http.Client{Transport: transport}

	return &h
}

// func (h *HTTPShandler) TestHTTPS(target string) {
// 	//h.actualTest()
// 	// target = "https://" + target	//@TODO

// 	fmt.Println("HTTPS Response Contents:\r\n", h.OptionsRequest(target))
// 	fmt.Println("HTTPS test\r\n",h.RequestMethodStatus("OPTIONS", target))
// 	// fmt.Println("_______________________", h.Robots(target))
// 	// fmt.Println("_______________________", h.Head(target))
// 	fmt.Println("_______________________")
// 	h.OptionsVerify(target)
// 	//os.Exit(0)
// }

func (h *HTTPShandler) Robots(url string) string {
	url = h.Urlprefixhttps(url)
	fmt.Println("\r\nRetrieving", url, "/robots.txt\r\n-------------\r\nResponse Status:")

	//	Body
	// fmt.Println("Contents of robots.txt\r\n-------------\r\n")
	return h.RequestMethod("GET", url + "/robots.txt")
}

func (h *HTTPShandler) Head(url string) string {
	return h.RequestMethod("HEAD", url)
}

func (h *HTTPShandler) RequestMethodStatus(method, url string) string {
	url = h.Urlprefixhttps(url)
	rq, e := http.NewRequest(method, url, nil)
	r, e := h.client.Do(rq)
	if e != nil {
		// log.Println(e)
		return fmt.Sprint("Request Method Status - ", e)
	}
	return string(r.Status)
}

func (h *HTTPShandler) OptionsRequest(url string) string {
	url = h.Urlprefixhttps(url)
	// // fmt.Println(url)
	// fmt.Println("\r\nHTTP OPTIONS Request - Retrieve Supported HTTP Methods\r\n-------------")	//	\r\nResponse Status:")
	var r string
	r = h.RequestMethod("OPTIONS", url)
	// fmt.Println("Response to the OPTIONS HTTPS Request:\r\n", r)
	return fmt.Sprintln("Response to the OPTIONS HTTPS Request:\r\n", r)
}

func (h *HTTPShandler) RequestMethod(method, url string) string {
	url = h.Urlprefixhttps(url)
	rq, e := http.NewRequest(method, url, nil)
	//var client http.Client
	r, e := h.client.Do(rq)
	if e != nil {
		return fmt.Sprint("Request Method -", e)
	}
	// fmt.Println(r)
	// fmt.Println(r.Status)
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		// log.Println(e)
		return fmt.Sprintf("%s", e)
	}
	// fmt.Println(body)
	return string(body)
}


func (h *HTTPShandler) OptionsVerify(url string) []string {
	url = h.Urlprefixhttps(url)
	var rs string
	var res []string
	options := make([]string, 7)
	options = []string{"CONNECT", "GET", "HEAD", "PATCH", "POST", "PUT", "TRACE"}

	// fmt.Println("\r\nSupported Options:\r\n-------------")
	for i := 0; i < len(options); i++ {
		rs = h.RequestMethodStatus(options[i], url)
		// fmt.Println(options[i], "-", res)	//	return and parse on caller
		res = append(res, options[i] + "-" + rs)
	}
	return res
}

// Get URL, check for http:// or https:// prefix
// add https:// if not present.
// Exit program if url contains https:// until implemented
func (h *HTTPShandler) Urlprefixhttps(url string) string {
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = "https://" + url
	} else if strings.Contains(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
		url = "https://" + url
	} else if strings.Contains(url, "https://") {
		//return the url as is
	}
	return url
}