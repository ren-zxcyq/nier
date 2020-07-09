package handlePdf

/*
 *	Create a pdf file using gofpdf	-	"github.com/jung-kurt/gofpdf"
 *
 */

import (
	"fmt"
	"log"
	"time"
	"path"
	"path/filepath"
	"github.com/jung-kurt/gofpdf"
	//"github.com/jung-kurt/gofpdf/internal/example"
)

const fontname = "Courier"	//"Times", "Arial", "Helvetica"

type pdfHandler struct {
	installationDir string
	filename string
	foldername string
}

func newPdfHandler(installDir, foldername string) *pdfHandler {
	var h pdfHandler = pdfHandler{installationDir: installDir, foldername: foldername, filename: path.Join(foldername, "Nier_Automata_Report.pdf")}
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

func CreatePdf(installDir, outputFolderName string) {
	fmt.Println("outputfoldername is", outputFolderName)
	pdfHandler := newPdfHandler(installDir, outputFolderName)
	// pdfHandler.exCreate()
	err := pdfHandler.pdfCreate()
	if err != nil {
		panic(err)
	}
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

/*	pdfCreate()
//	This should be okay so far but let's replace it with the website version.
func (h *pdfHandler) imageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "image", fileStr)
}


func (h *pdfHandler) pdfCreate() error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Set Header
	pdf.SetTopMargin(30)
	pdf.SetHeaderFuncMode(func() {
		pdf.Image(h.imageFile("avatar.png"), 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(30, 10, "Title", "1", 0, "C", false, 0, "")
		pdf.Ln(20)
	}, true)
	/*
	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, "Nier - Report", "0", 0, "CM", false, 0, "")

	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	pdf.ImageOptions(
		"avatar.jpg",
		20, 20,
		140, 100,	//
		false,
		gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
		0,
		"",
	)

	pdf.CellFormat(190, 7, "- by zxcyq", "0", 0, "CM", false, 0, "")
	
	return pdf.OutputFileAndClose(h.filename)
}
*/

/*
 *	Creates a document -> Sets Header & Creates a table
 *			uses h.newReport() to do so
 */
func (h *pdfHandler) pdfCreate() error {

	//	Create a new PDF doc & write title & current date
	pdf := h.newReport()


	tableCols := []string{"No", "Tool", "Description"}
	// var data []
	tableCont := [][]string {
		{"1", "ping", "ping -c 1"},
		{"2", "nmap", "version scan"},
		{"3", "nikto", "vuln testing"},
	}

	//	Create Table Header & Fill
	pdf = h.header(pdf, tableCols)
	pdf = h.table(pdf, tableCont)

	//	Add Logo
	pdf = h.image(pdf)

	if pdf.Err() {
		log.Fatalf("Failed while creating the PDF Report: %s\n", pdf.Error())
		return pdf.Error()
	}

	//	Save
	err := h.savePDF(pdf)
	if err != nil {
		log.Fatalf("Cannot save PDF: %s\n", err)
		return err
	}
	return nil
}

//	This creates the Document template
func (h *pdfHandler) newReport() *gofpdf.Fpdf {
	//	New() creates

	//	"mm" unit for expressing lengths & sizes
	//	"Letter" -> Paper format
	//	"L" or "P" -> Orientation
	//	path to a font directory
	pdf := gofpdf.New("P", "mm", "Letter", "")

	//
	pdf.AddPage()

	//	set font
	//		"B" -> bold
	//		int -> size
	pdf.SetFont(fontname, "B", 28)


	//	Write a text cell of length 40 & height 10.
	//		-	no starting coords
	//			Cell() moves the current pos to the end of the cell
	//					next Cell() will continue after
	pdf.Cell(40, 10, "Nier - Report")

	//	Ln()	->	Moves the current pos to a new line
	//				(optional height param)
	pdf.Ln(12)

	pdf.SetFont(fontname, "", 20)
	pdf.Cell(40,10, time.Now().Format("Mon Sep 9, 2020"))
	pdf.Ln(12)
	//	Note on Cell() & Ln()
	//	
	//	Cell() -> No Coordinates
	//		document keeps them internally
	//		advances to the right by the length of the cell being written
	//	Ln() -> moves current position back to the left border & down
	//			by the provided value

	return pdf
}

func (h *pdfHandler) header (pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont(fontname, "B", 16)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		//	pdf.CellFormat() -> format the new Cell -> +border +background_fill
		pdf.CellFormat(40, 7, str, "1", 0, "", true, 0, "")
	}

	//	Pass	-1 ->	Ln()	i.e. use the height of the last printed Cell as the line height
	pdf.Ln(-1)
	return pdf
}

func (h *pdfHandler) table(pdf *gofpdf.Fpdf, tbl [][]string) *gofpdf.Fpdf {
	
	//	Font & Fill Color
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(255,255,255)

	//	Allign columns according to their contents
	align := []string{"L", "C", "L", "R", "R", "R"}	//"No.","Tool"
	for _, line := range tbl {
		for i, str := range line {
			//	CellFormat() -> Create a visible border around the cell
			//	alignStr param is used to align the cell content either Left or Right
			pdf.CellFormat(40, 7, str, "1", 0, align[i], false, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}

// ImageFile returns a qualified filename in which the path to the image
// directory is prepended to the specified filename.
func (h *pdfHandler) imageFile(fileStr string) string {
	return filepath.Join(h.installationDir, "image", fileStr)
}

func (h *pdfHandler) image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	//	ImageOptions() takes a file path, x,y,width & height params
	//						&	ImageOptions struct to specify a couple of options
	// pdf.ImageOptions(
	// 	h.imageFile("avatar.jpg"),
	// 	255, 10,
	// 	25, 25,
	// 	false,
	// 	gofpdf.ImageOptions{ImageType:"JPG", ReadDpi: true},
	// 	0,
	// 	"",
	// )
	pdf.ImageOptions(
		"image/avatar.jpg",
		// 20, 20,
		// 140, 100,	//
		25, 70,
		140, 100,
		false,
		gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
		0,
		"",
	)
	if pdf.Err() {
		fmt.Printf("Failed while adding image to the PDF Report: %s\n", pdf.Error())
		
	}

	return pdf
}

func (h *pdfHandler) savePDF(pdf *gofpdf.Fpdf) error {
	return pdf.OutputFileAndClose(h.filename)
}