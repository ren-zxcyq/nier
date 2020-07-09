package utilities

import (
	. "fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Agent struct{}

func (a *Agent) Robots(url string) {

	Println("\r\nRetrieving", url, "/robots.txt\r\n-------------\r\nResponse Status:")

	r, e := http.Get(url + "/robots.txt")
	if e != nil {
		Printf("Error Encountered")
		log.Panicln(e)
	}

	//	r.Status contains the status
	Println(r.Status)

	//	Extract Body
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Panicln(e)
	}

	//	Body
	Println("Contents of robots.txt\r\n-------------\r\n", string(body))
}

func (a *Agent) Head(url string) {
	r, e := http.Head(url)

	if e != nil {
		Printf("Error Encountered")
		log.Panicln(e)
	}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Panicln(e)
	}
	Println(string(body))
}

func (a *Agent) RequestMethod(method, url string) string {
	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		log.Panicln(e)
	}
	Println(r)
	Println(r.Status)
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Panicln(e)
	}
	// Println(body)
	return string(body)
}

func (a *Agent) RequestMethodStatus(method, url string) string {

	rq, e := http.NewRequest(method, url, nil)
	var client http.Client
	r, e := client.Do(rq)
	if e != nil {
		log.Panicln(e)
	}
	//Println(r)
	//Println(r.Status)
	//body, e := ioutil.ReadAll(r.Body)
	//if e != nil {
	//	log.Panicln(e)
	//}
	// Println(body)
	return string(r.Status)
}

func (a *Agent) OptionsRequest(url string) {
	Println("\r\nHTTP OPTIONS Request - Retrieve Supported HTTP Methods\r\n-------------\r\nResponse Status:")
	var r string
	r = a.RequestMethod("OPTIONS", url)
	Println("Response to the OPTIONS HTTP Request:\r\n", r)
}

func (a *Agent) OptionsVerify(url string) {
	options := make([]string, 7)                                                  //	9
	options = []string{"CONNECT", "GET", "HEAD", "PATCH", "POST", "PUT", "TRACE"} //	"DELETE", "OPTIONS",

	Println("\r\nSupported Options:\r\n-------------")
	var res string
	// @TODO	request each and every one.
	for i := 0; i < len(options); i++ {
		res = a.RequestMethodStatus(options[i], url)
		Println(options[i], "-", res)
	}
}
