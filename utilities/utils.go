package utilities

import (
	"bufio"
	"encoding/base64"
	. "fmt"
	"go/build"
	"net/url"
	"os"
	"runtime"

	"github.com/dchest/uniuri"
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

func (h *Utils) GetGOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	// fmt.Println(gopath)
	return gopath
}

func (h *Utils) GetGOROOT() string {
	//fmt.Println(build.Default.GOROOT)
	// fmt.Println(runtime.GOROOT())
	return runtime.GOROOT()
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

	// Println("\r\nLoading Config:", file)

	//	File exists
	Println("\r\nLoading Utilities From\t-\t", file, "\r\n-------------")
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
func (h *Utils) EncodingTest() {
	Println("\r\nEncoding Test Starting:\r\n-------------")
	t := "enc*de Me Plea$e"
	Println(t)
	Println(h.encodeParam(t))
	Println(h.encodeStringBase64(t))
	Println(h.UniqueString())
}

func (h *Utils) encodeParam(s string) string {
	return url.QueryEscape(s)
}

func (h *Utils) encodeStringBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

/*
 *	Return a Random String using - github.com/uniuri (cryptographically secure string)
 */
func (h *Utils) UniqueString() string {
	//	Default uniuri.StdChars contains only alphanum
	uniuri.StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+=-`~,<.>/?;:'\"")
	s := uniuri.New() //	default: 16 letters
	// s := uniuri.NewLen(32)	//	set our own
	return s
}
