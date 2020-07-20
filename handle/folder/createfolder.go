// Package handlefolder checks if the requested output folder exists and creates it if not.
// If the folder exists, the package terminates execution.
package handlefolder

import (
	"fmt"
	"os"
)

type outputFolderHandler struct {
	folderName string
}

func newFolderHandler(fName string) *outputFolderHandler {
	var h outputFolderHandler = outputFolderHandler{folderName: fName}

	//fmt.Printf("Address of outputFolderHandler - %p", &h) //	Prints the address of outputFolderHandler
	return &h
}

/*
 *	Creates a folder if it does not exist. Caller is responsible for handling that case.
 *	p passed in here is actually the merged cwd & user supplied [-o] param.
 */
func (h *outputFolderHandler) makeFolder() bool {
	var p string = h.folderName
	if n, err := os.Stat(fmt.Sprintf(p)); os.IsNotExist(err) {

		fmt.Printf("\nCreating: %s\n", p)
		fmt.Printf("\n%s\n", n)

		//os.Mkdir(p, 0777)	//	Create a single folder
		os.MkdirAll(p, 0777) //	Create Folder & any parents
		return true
	} else {
		//	Folder does not exist.
		return false
	}
}

func CreateFolder(outputFolder string) {
	fHandler := newFolderHandler(outputFolder)
	//fmt.Println("Creating Folder", outputFolder)
	var res bool
	if res = fHandler.makeFolder(); res == false {
		//	Folder NOT generated. Exit
		fmt.Println("\n-------------\nFolder already exists. Exiting.")
		os.Exit(1)
	}
}
