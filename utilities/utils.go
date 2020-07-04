package utilities

import (
	"bufio"
	. "fmt"
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
