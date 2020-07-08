package handleFolder

import (
	. "fmt"
	"os"
)

type outputFolderHandler struct {
	folderName string
}

func newFolderHandler(fName string) *outputFolderHandler {
	var h outputFolderHandler = outputFolderHandler{folderName: fName}

	//Printf("Address of outputFolderHandler - %p", &h) //	Prints the address of outputFolderHandler
	return &h
}

/*
 *	Creates a folder if it does not exist. Caller is responsible for handling that case.
 *	p passed in here is actually the merged cwd & user supplied [-o] param.
 */
func (h *outputFolderHandler) makeFolder() bool {
	Println("\nIn create")
	var p string = h.folderName
	if wh, err := os.Stat(p); os.IsNotExist(err) {

		Printf("\nCreating: %s\n", p)
		Printf("\n%s\n", wh)

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
	//Println("Creating Folder", outputFolder)
	var res bool
	if res = fHandler.makeFolder(); res == false {
		//	Folder NOT generated. Exit
		Println("\n-------------\nFolder already exists. Exiting.")
		os.Exit(1)
	}
}