package utilities

import (
	"bufio"
	"encoding/base64"
	. "fmt"
	"net/url"
	"os"
	"runtime"
)

type Utils struct{}

func (h *Utils) DetectOS() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return Sprintf("Windows")
	case "darwin":
		return Sprintf("Mac OS")
	case "linux":
		return Sprintf("Linux")
	default:
		return Sprintf("%s", os)
	}
}

//	Return list by reading absPath file line-by-line
func (h *Utils) ReturnLinesFromFile(absPath string) []string {
	tfile, err := os.Open(absPath)
	if err != nil {
		Println(err)
	}
	defer tfile.Close()

	var tscanner = bufio.NewScanner(tfile)
	tlines := []string{}

	for tscanner.Scan() {
		tlines = append(tlines, tscanner.Text())
	}

	return tlines
}

func (h *Utils) PrintFileContents(file string) {

	Println("\r\nLoading Config:", file)

	//	File exists
	Println("\r\nLoading Utilities From:\r\n-------------")
	var utils []string = h.ReturnLinesFromFile(file)

	var i int = 0
	for i < len(utils) {
		Println(utils[i])
		i++
	}
	Println("-------------")
}

/*
 *	encodes t -> base64 & url encoding
 *	uses:
 *			net/url
 *			encoding/base64
 *	in:
 *		encodeParams()
 *		encodeStringBase64()
 *
 *	usage:	utilities.EncodingTest()
 */
func EncodingTest() {
	Println("Encoding Test Starting")
	t := "enc*de Me Plea$e"
	Println(t)
	Println(encodeParam(t))
	Println(encodeStringBase64(t))
}

func encodeParam(s string) string {
	return url.QueryEscape(s)
}

func encodeStringBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
