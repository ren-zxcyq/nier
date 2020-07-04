package handlePdf

/*
 *	Create a pdf file using gofpdf	-	"github.com/jung-kurt/gofpdf"
 *
 */

import (
	"fmt"
	"log"
	"path"

	"github.com/jung-kurt/gofpdf"
)

type pdfHandler struct {
	filename string
}

func newPdfHandler(foldername string) *pdfHandler {
	var h pdfHandler = pdfHandler{filename: path.Join(foldername, "Nier_Automata_Report.pdf")}
	//fmt.Printf("Address of pdfHandler - %p", &h) //	Prints the address of documentHandler
	//fmt.Println(foldername)
	return &h
}

func (h *pdfHandler) exCreate() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose(h.filename)

	if err != nil {
		log.Fatal(err)
	}
}

func CreatePdf(outputFolderName string) {
	fmt.Println("outputfoldername is", outputFolderName)
	pdfHandler := newPdfHandler(outputFolderName)
	pdfHandler.exCreate()
}

// ExampleFpdf_HTMLBasicNew demonstrates internal and external links with and without basic
// HTML.
// func (h *pdfHandler) pdfCreate() {
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	// First page: manual local link
// 	pdf.AddPage()
// 	pdf.SetFont("Helvetica", "", 20)
// 	_, lineHt := pdf.GetFontSize()
// 	pdf.Write(lineHt, "To find out what's new in this tutorial, click ")
// 	pdf.SetFont("", "U", 0)
// 	link := pdf.AddLink()
// 	pdf.WriteLinkID(lineHt, "here", link)
// 	pdf.SetFont("", "", 0)
// 	// Second page: image link and basic HTML with link
// 	pdf.AddPage()
// 	pdf.SetLink(link, 0, -1)
// 	pdf.Image(example.ImageFile("logo.png"), 10, 12, 30, 0, false, "", 0, "http://www.fpdf.org")
// 	pdf.SetLeftMargin(45)
// 	pdf.SetFontSize(14)
// 	_, lineHt = pdf.GetFontSize()
// 	htmlStr := `You can now easily print text mixing different styles: <b>bold</b>, ` +
// 		`<i>italic</i>, <u>underlined</u>, or <b><i><u>all at once</u></i></b>!<br><br>` +
// 		`<center>You can also center text.</center>` +
// 		`<right>Or align it to the right.</right>` +
// 		`You can also insert links on text, such as ` +
// 		`<a href="http://www.fpdf.org">www.fpdf.org</a>, or on an image: click on the logo.`
// 	html := pdf.HTMLBasicNew()
// 	html.Write(lineHt, htmlStr)
// 	fileStr := example.Filename("Fpdf_HTMLBasicNew")
// 	err := pdf.OutputFileAndClose(fileStr)
// 	example.Summary(err, fileStr)
// 	// Output:
// 	// Successfully generated pdf/Fpdf_HTMLBasicNew.pdf
// }
