// Package injectiondetector is responsible for scraping the target website,
// extracting <form> tags, filter for unique forms, submitting all of them
// and identifying user controlled input which appears on the application pages.
// interface their execution with parsing.
package utilities

import (
	"fmt"
	"strings"
	// "net/http"
	// "io/ioutil"
	// "reflect"
	"bytes"
	
	"github.com/PuerkitoBio/goquery"
)

type InjectionHandler struct {
	// e     elementsHandler //	Receives Main.hCmd
}


//
func NewInjectionHandler() *InjectionHandler {

	//	Create an elementsHandler Object to be passed to the exported execHandler
	// var l elementsHandler = elementsHandler{
	// 	installationDir:      installationDir,
	// 	configFilePath:       configPath,
	// }

	//	Create execHandler
	var h InjectionHandler = InjectionHandler{}

	//fmt.Printf("Address of execHandler - %p", &h) //	Prints the address of outputFolderHandler
	return &h
}

func (h *InjectionHandler) XssURLsi() {
	
	var urls []string

	urls = append(urls, "http://192.168.1.20/vehical-details.php?vhid=2")
	urls = append(urls, "http://192.168.1.20/robots.txt")
	urls = append(urls, "http://192.168.1.20/contact-us.php")


	for _,url := range urls {
		h.xssRequestURLi(url)
	}
}

func (h *InjectionHandler) xssRequestURLi(url string) {
	var t Agent
	var or string
	// var r *http.Response

	// r = t.WrappedGet(url)			//	h.e.targetHost + ":" + strconv.Itoa(h.e.targetPort))

	// //	Extract Body
	// body, e := ioutil.ReadAll(r.Body)
	// if e != nil {
	// 	// log.Println(e)
	// 	fmt.Printf("%s", e)
	// }

	// or = string(body)

	var r string

	r = t.WrappedGet(url)
	or = r

	if strings.Contains(or, "Target Responds in HTTPS - Cannot Follow through with HTTP Methods Checking") {
		fmt.Println("-------------")
		fmt.Println("[+]\tUpgrading to HTTPS")
		// tsec := utilities.NewHTTPShandler()
		
		// fmt.Println("HTTPS test\r\n",h.RequestMethodStatus("OPTIONS", target))
		// fmt.Println("_______________________", h.Robots(target))
		// fmt.Println("_______________________", h.Head(target))
		/*
		tester := utilities.NewHTTPShandler()
		tester.TestHTTPS(h.e.targetHost)
		tester.Robots(h.e.targetHost)
		*/
		fmt.Println("-------------")
	} else {								//	@TODO	consider checking for another error
		fmt.Println("-------------")
		fmt.Println("[+]\tContinue HTTP")
		fmt.Println(url)
		// fmt.Println(or)
		//	If HTML response contains a form -> pass it to the parser
		if strings.Contains(or, "<form") {
			h.extractForms(r)
			// fmt.Println("YES")
			// r := strings.Split(or, "<form")
			// for i,j := range r {
			// 	fmt.Println(i, "\t-\t", j)
			// 	//	This WORKS
				
			// }
		} else {
			fmt.Println("NO")
		}
		fmt.Println("-------------")
	}
	// fmt.Println("-------------")
	// fmt.Println(results)
	// fmt.Println("-------------")
}

func (h *InjectionHandler) extractForms(r string) {	//(r *http.Response) {
	funky4()
}

func funky() {
    test := `<speak><p>My paragraph</p></speak>`
    doc, _ := goquery.NewDocumentFromReader(strings.NewReader(test))
    var childrenHtml []string
    doc.Find("speak").Children().Each(func(i int, s *goquery.Selection) {
        html, _ := goquery.OuterHtml(s)																	//s.Html()
        childrenHtml = append(childrenHtml, html)
    })
    if childrenHtml[0] != "<p>My paragraph</p>" {
       fmt.Println("First element html is not valid: '%s'", childrenHtml[0])
	}
	
	fmt.Println("CHILDRENHTML", childrenHtml[0])
}

func funky2() {
	test := `<form action="#" method="get" id="header-search-form"><input type="text" placeholder="Search..." class="form-control"><button type="submit"><i class="fa fa-search" aria-hidden="true"></i></button></form>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(test))
	var childrenHtml []string
	doc.Find("form").Children().Each(func(i int, s *goquery.Selection) {
		html, _ := goquery.OuterHtml(s)
		childrenHtml = append(childrenHtml, html)
	})
	// if childrenHtml[0] != `<input type="text" placeholder="Search..." class="form-control"/>` {
	// 	fmt.Println("Not Valid: ", childrenHtml[0])
	// }
	fmt.Println(len(childrenHtml))
	fmt.Println("CHILDRENHTML", childrenHtml[1])

}

func funky3() {
	test := `<form  method="post">
	<div class="form-group">
	  <label class="control-label">Full Name <span>*</span></label>
	  <input type="text" name="fullname" class="form-control white_bg" id="fullname" required>
	</div>
	<div class="form-group">
	  <label class="control-label">Email Address <span>*</span></label>
	  <input type="email" name="email" class="form-control white_bg" id="emailaddress" required>
	</div>
	<div class="form-group">
	  <label class="control-label">Phone Number <span>*</span></label>
	  <input type="text" name="contactno" class="form-control white_bg" id="phonenumber" required>
	</div>
	<div class="form-group">
	  <label class="control-label">Message <span>*</span></label>
	  <textarea class="form-control white_bg" name="message" rows="4" required></textarea>
	</div>
	<div class="form-group">
	  <button class="btn" type="submit" name="send" type="submit">Send Message <span class="angle_arrow"><i class="fa fa-angle-right" aria-hidden="true"></i></span></button>
	</div>
  </form>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(test))

	var childrenHtml []string
	var inputHtml []string
	var buttonHtml []string

	doc.Find("form").Children().Each(func(i int, s *goquery.Selection) {
		html, _ := goquery.OuterHtml(s)
		childrenHtml = append(childrenHtml, html)
	})

	fmt.Println("CHILDRENHTML", childrenHtml[1])
	fmt.Println(len(childrenHtml))


	for _,j := range childrenHtml {
		// fmt.Println("START", j)
		j, _ := goquery.NewDocumentFromReader(strings.NewReader(j))
		j.Find("input").Siblings().Each(func(i int, ss *goquery.Selection) {
			//fmt.Println("11111111")
			htmll, _ := ss.Html()	//goquery.OuterHtml(ss)
			fmt.Println("ALLOCATING", htmll)
			//fmt.Println("44444444")
			inputHtml = append(inputHtml, htmll)
			//fmt.Println("55555555")

		})
		
		//fmt.Println("77777777")
		//fmt.Println("SSSSSSSSSSSSSS", j.Find("input").Children().Size())

		if j.Find("input").Children().Size() > 0 {
			fmt.Println("HAS CHILDREN")
			j.Find("input").Children().Each(func(i int, sss *goquery.Selection) {
				//fmt.Println("11111111")
	
				htmll1, _ := sss.Html()		//goquery.OuterHtml(sss)
				inputHtml = append(inputHtml, htmll1)
			})
		}

		if len(inputHtml) > 0 {
			// fmt.Println("77777777")
			// fmt.Println("INPUT-html", inputHtml[0])
			// fmt.Println("88888888")	
		}

		j.Find("button").Siblings().Each(func(i int, ssss *goquery.Selection) {
			//fmt.Println("222222222")

			htmlll, _ := ssss.Html()		//goquery.OuterHtml(ssss)
			buttonHtml = append(buttonHtml, htmlll)
		})
		if j.Find("button").Children().Size() > 0 {
			fmt.Println("HAS CHILDREN")

			j.Find("button").Children().Each(func(i int, sssss *goquery.Selection) {
				//fmt.Println("222222222")
				htmll2, _ := goquery.OuterHtml(sssss)
				buttonHtml = append(buttonHtml, htmll2)
			})	
		}
		if len(buttonHtml) > 0 {
			// fmt.Println("77777777")
			// fmt.Println("button-html", buttonHtml[0])
			// fmt.Println("88888888")
		}
	}

	for nia, niania := range inputHtml {
		fmt.Println(nia, "-", niania)
	}
	
	for ain, ainain := range buttonHtml {
		fmt.Println(ain, "-", ainain)
	}
	
	
}


func funky4() {
	test := `<form  method="post">
	<div class="form-group">
	  <label class="control-label">Full Name <span>*</span></label>
	  <input type="text" name="fullname" class="form-control white_bg" id="fullname" required>
	</div>
	<div class="form-group">
	  <label class="control-label">Email Address <span>*</span></label>
	  <input type="email" name="email" class="form-control white_bg" id="emailaddress" required>
	</div>
	<div class="form-group">
	  <label class="control-label">Phone Number <span>*</span></label>
	  <input type="text" name="contactno" class="form-control white_bg" id="phonenumber" required>
	</div>
	<div class="form-group">
	  <label class="control-label">Message <span>*</span></label>
	  <textarea class="form-control white_bg" name="message" rows="4" required></textarea>
	</div>
	<div class="form-group">
	  <button class="btn" type="submit" name="send" type="submit">Send Message <span class="angle_arrow"><i class="fa fa-angle-right" aria-hidden="true"></i></span></button>
	</div>
  </form>
  <form action="#" method="get" id="header-search-form">
  <input type="text" placeholder="Search..." class="form-control">
  <button type="submit"><i class="fa fa-search" aria-hidden="true"></i></button>
</form> 
  `

	// doc, _ := goquery.NewDocumentFromReader(strings.NewReader(test))

	// // var childrenHtml []string
	// var formStr string
	// // var inputHtml []string
	// // var buttonHtml []string

	// doc.Find("form").Children().Each(func(i int, s *goquery.Selection) {
	// 	html, _ := goquery.OuterHtml(s)
	// 	// childrenHtml = append(childrenHtml, html)
	// 	formStr += html
	// })

	// // fmt.Println("CHILDRENHTML", childrenHtml[1])
	// // fmt.Println(len(childrenHtml))

	// fmt.Println("CHILDRENHTML", formStr)
	// // fmt.Println(len(formStr))


	// form, _ := goquery.NewDocumentFromReader(strings.NewReader(formStr))
	
    doc, _ := goquery.NewDocumentFromReader(strings.NewReader((test)))
    doc.Find("form").Each(func(i int, form *goquery.Selection) {

		fmt.Println("FOUND FORM")
        form.Find("button").Each(func(j int, b *goquery.Selection) {
			fmt.Println("FOUND BUTTON")
			btype, oktype := b.Attr("type")
			if oktype {
				if strings.Contains(btype, "submit") {
					fmt.Println("Has TYPE - submit",)
				} else {
					fmt.Println("Has TYPE -", btype)
				}
			}
			bname, okname := b.Attr("name")
			if okname {
				fmt.Println("Has NAME -", bname)
			} 

			/*if b.Attr("type")  {
                if p := h1.Next(); p != nil {
                    if ps := p.Children().First(); ps != nil && ps.HasClass("text") {
                        ps.ReplaceWithHtml(
                            fmt.Sprintf("<span class=\"text\">%s%s</span>)", s.Text(), ps.Text()))
                        h1.Remove()
                    }
                }
			}
			*/
		})

		form.Find("input").Each(func(j int, in *goquery.Selection) {
			fmt.Println("FOUND INPUT")
			intype, oktype := in.Attr("type")
			if oktype {
				if strings.Contains(intype, "text") {
					fmt.Println("Has TYPE - text")
				} else {
					fmt.Println("Has TYPE -", intype)
				}
			}
			inname, okname := in.Attr("name")
			if okname {
				fmt.Println("Has NAME -", inname)
			}
		})

    })
	// htmlResult, _ := doc.Html()
	// fmt.Println(htmlResult)



	// inputsSelector := new(goquery.Selection)
	// inputsSelector = visitNodes(inputsSelector, s, "input")
	// //n := inputsSelector
	// nT := inputsSelector.Text()
	// nt,_ := inputsSelector.Attr("type")


	// buttonsSelector := new(goquery.Selection)
	// buttonsSelector = visitNodes(buttonsSelector, s, "button")
	// // fmt.Println(buttonsSelector.Size())
	// b := buttonsSelector
	// bT := buttonsSelector.Text()
	// bt,_ := b.Attr("type")
	// bn,_ := b.Attr("name")

	// button := s.Find("button")
	// bText := s.Find("button").Text()
	// btype,_ := button.Attr("type")
	// bname,_ := button.Attr("name")
}


func brokenTestWithVisitNodes() {
	

	var reee string = `<!DOCTYPE HTML>
	<html lang="en">
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<meta name="keywords" content="">
	<meta name="description" content="">
	<title>CarForYou - Responsive Car Dealer HTML5 Template</title>
	<!--Bootstrap -->
	<link rel="stylesheet" href="assets/css/bootstrap.min.css" type="text/css">
	<!--Custome Style -->
	<link rel="stylesheet" href="assets/css/style.css" type="text/css">
	<!--OWL Carousel slider-->
	<link rel="stylesheet" href="assets/css/owl.carousel.css" type="text/css">
	<link rel="stylesheet" href="assets/css/owl.transitions.css" type="text/css">
	<!--slick-slider -->
	<link href="assets/css/slick.css" rel="stylesheet">
	<!--bootstrap-slider -->
	<link href="assets/css/bootstrap-slider.min.css" rel="stylesheet">
	<!--FontAwesome Font Style -->
	<link href="assets/css/font-awesome.min.css" rel="stylesheet">
	
	<!-- SWITCHER -->
			<link rel="stylesheet" id="switcher-css" type="text/css" href="assets/switcher/css/switcher.css" media="all" />
			<link rel="alternate stylesheet" type="text/css" href="assets/switcher/css/red.css" title="red" media="all" data-default-color="true" />
			<link rel="alternate stylesheet" type="text/css" href="assets/switcher/css/orange.css" title="orange" media="all" />
			<link rel="alternate stylesheet" type="text/css" href="assets/switcher/css/blue.css" title="blue" media="all" />
			<link rel="alternate stylesheet" type="text/css" href="assets/switcher/css/pink.css" title="pink" media="all" />
			<link rel="alternate stylesheet" type="text/css" href="assets/switcher/css/green.css" title="green" media="all" />
			<link rel="alternate stylesheet" type="text/css" href="assets/switcher/css/purple.css" title="purple" media="all" />
			
	<!-- Fav and touch icons -->
	<link rel="apple-touch-icon-precomposed" sizes="144x144" href="assets/images/favicon-icon/apple-touch-icon-144-precomposed.png">
	<link rel="apple-touch-icon-precomposed" sizes="114x114" href="assets/images/favicon-icon/apple-touch-icon-114-precomposed.html">
	<link rel="apple-touch-icon-precomposed" sizes="72x72" href="assets/images/favicon-icon/apple-touch-icon-72-precomposed.png">
	<link rel="apple-touch-icon-precomposed" href="assets/images/favicon-icon/apple-touch-icon-57-precomposed.png">
	<link rel="shortcut icon" href="assets/images/favicon-icon/favicon.png">
	<link href="https://fonts.googleapis.com/css?family=Lato:300,400,700,900" rel="stylesheet">
	 <style>
		.errorWrap {
		padding: 10px;
		margin: 0 0 20px 0;
		background: #fff;
		border-left: 4px solid #dd3d36;
		-webkit-box-shadow: 0 1px 1px 0 rgba(0,0,0,.1);
		box-shadow: 0 1px 1px 0 rgba(0,0,0,.1);
	}
	.succWrap{
		padding: 10px;
		margin: 0 0 20px 0;
		background: #fff;
		border-left: 4px solid #5cb85c;
		-webkit-box-shadow: 0 1px 1px 0 rgba(0,0,0,.1);
		box-shadow: 0 1px 1px 0 rgba(0,0,0,.1);
	}
		</style>
	</head>
	<body>
	
	<<!-- Start Switcher -->
	<div class="switcher-wrapper">	
		<div class="demo_changer">
			<div class="demo-icon customBgColor"><i class="fa fa-cog fa-spin fa-2x"></i></div>
			<div class="form_holder">
				<div class="row">
					<div class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
						<div class="predefined_styles">
							<div class="skin-theme-switcher">
								<h4>Color</h4>
								<a href="#" data-switchcolor="red" class="styleswitch" style="background-color:#de302f;"> </a>
								<a href="#" data-switchcolor="orange" class="styleswitch" style="background-color:#f76d2b;"> </a>
								<a href="#" data-switchcolor="blue" class="styleswitch" style="background-color:#228dcb;"> </a>
								<a href="#" data-switchcolor="pink" class="styleswitch" style="background-color:#FF2761;"> </a>
								<a href="#" data-switchcolor="green" class="styleswitch" style="background-color:#2dcc70;"> </a>
								<a href="#" data-switchcolor="purple" class="styleswitch" style="background-color:#6054c2;"> </a>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div><!-- /Switcher -->  
			
	<!--Header-->
	
	<header>
	  <div class="default-header">
		<div class="container">
		  <div class="row">
			<div class="col-sm-3 col-md-2">
			  <div class="logo"> <a href="index.php"><img src="assets/images/logo.png" alt="image"/></a> </div>
			</div>
			<div class="col-sm-9 col-md-10">
			  <div class="header_info">
				<div class="header_widgets">
				  <div class="circle_icon"> <i class="fa fa-envelope" aria-hidden="true"></i> </div>
				  <p class="uppercase_text">For Support Mail us : </p>
				  <a href="mailto:info@example.com">rickastley@astleycars.com</a> </div>
				<div class="header_widgets">
				  <div class="circle_icon"> <i class="fa fa-phone" aria-hidden="true"></i> </div>
				  <p class="uppercase_text">Service Helpline Call Us: </p>
				  <a href="tel:999898989898">+999898989898</a> </div>
				<div class="social-follow">
				  <ul>
					<li><a href="#"><i class="fa fa-facebook-square" aria-hidden="true"></i></a></li>
					<li><a href="#"><i class="fa fa-twitter-square" aria-hidden="true"></i></a></li>
					<li><a href="#"><i class="fa fa-linkedin-square" aria-hidden="true"></i></a></li>
					<li><a href="#"><i class="fa fa-google-plus-square" aria-hidden="true"></i></a></li>
					<li><a href="#"><i class="fa fa-instagram" aria-hidden="true"></i></a></li>
				  </ul>
				</div>
	   Welcome To Astley Car rental portal          </div>
			</div>
		  </div>
		</div>
	  </div>
	  
	  <!-- Navigation -->
	  <nav id="navigation_bar" class="navbar navbar-default">
		<div class="container">
		  <div class="navbar-header">
			<button id="menu_slide" data-target="#navigation" aria-expanded="false" data-toggle="collapse" class="navbar-toggle collapsed" type="button"> <span class="sr-only">Toggle navigation</span> <span class="icon-bar"></span> <span class="icon-bar"></span> <span class="icon-bar"></span> </button>
		  </div>
		  <div class="header_wrap">
			<div class="user_login">
			  <ul>
				<li class="dropdown"> <a href="#" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false"><i class="fa fa-user-circle" aria-hidden="true"></i> 
	Steve Brown<i class="fa fa-angle-down" aria-hidden="true"></i></a>
				  <ul class="dropdown-menu">
						   <li><a href="profile.php">Profile Settings</a></li>
				  <li><a href="updatepassword.php">Update Password</a></li>
				<li><a href="my-booking.php">My Booking</a></li>
				<li><a href="post-testimonial.php">Post a Testimonial</a></li>
			  <li><a href="my-testimonials.php">My Testimonial</a></li>
				<li><a href="logout.php">Sign Out</a></li>
						  </ul>
				</li>
			  </ul>
			</div>
			<div class="header_search">
			  <div id="search_toggle"><i class="fa fa-search" aria-hidden="true"></i></div>
			  <form action="#" method="get" id="header-search-form">
				<input type="text" placeholder="Search..." class="form-control">
				<button type="submit"><i class="fa fa-search" aria-hidden="true"></i></button>
			  </form>
			</div>
		  </div>
		  <div class="collapse navbar-collapse" id="navigation">
			<ul class="nav navbar-nav">
			  <li><a href="index.php">Home</a>    </li>
				   
			  <li><a href="page.php?type=aboutus.php">About Us</a></li>
			  <li><a href="car-listing.php">Car Listing</a>
			  <li><a href="page.php?type=faqs.php">FAQs</a></li>
			  <li><a href="contact-us.php">Contact Us</a></li>
	
			</ul>
		  </div>
		</div>
	  </nav>
	  <!-- Navigation end --> 
	  
	</header><!-- /Header --> 
	
	<!--Page Header-->
	<section class="page-header contactus_page">
	  <div class="container">
		<div class="page-header_wrap">
		  <div class="page-heading">
			<h1>Contact Us</h1>
		  </div>
		  <ul class="coustom-breadcrumb">
			<li><a href="#">Home</a></li>
			<li>Contact Us</li>
		  </ul>
		</div>
	  </div>
	  <!-- Dark Overlay-->
	  <div class="dark-overlay"></div>
	</section>
	<!-- /Page Header--> 
	
	<!--Contact-us-->
	<section class="contact_us section-padding">
	  <div class="container">
		<div  class="row">
		  <div class="col-md-6">
			<h3>Get in touch using the form below</h3>
					  <div class="contact_form gray-bg">
			  <form  method="post">
				<div class="form-group">
				  <label class="control-label">Full Name <span>*</span></label>
				  <input type="text" name="fullname" class="form-control white_bg" id="fullname" required>
				</div>
				<div class="form-group">
				  <label class="control-label">Email Address <span>*</span></label>
				  <input type="email" name="email" class="form-control white_bg" id="emailaddress" required>
				</div>
				<div class="form-group">
				  <label class="control-label">Phone Number <span>*</span></label>
				  <input type="text" name="contactno" class="form-control white_bg" id="phonenumber" required>
				</div>
				<div class="form-group">
				  <label class="control-label">Message <span>*</span></label>
				  <textarea class="form-control white_bg" name="message" rows="4" required></textarea>
				</div>
				<div class="form-group">
				  <button class="btn" type="submit" name="send" type="submit">Send Message <span class="angle_arrow"><i class="fa fa-angle-right" aria-hidden="true"></i></span></button>
				</div>
			  </form>
			</div>
		  </div>
		  <div class="col-md-6">
			<h3>Contact Info</h3>
			<div class="contact_detail">
							<ul>
				<li>
				  <div class="icon_wrap"><i class="fa fa-map-marker" aria-hidden="true"></i></div>
				  <div class="contact_info_m">Test Demo test demo																									</div>
				</li>
				<li>
				  <div class="icon_wrap"><i class="fa fa-phone" aria-hidden="true"></i></div>
				  <div class="contact_info_m"><a href="tel:61-1234-567-90">test@test.com</a></div>
				</li>
				<li>
				  <div class="icon_wrap"><i class="fa fa-envelope-o" aria-hidden="true"></i></div>
				  <div class="contact_info_m"><a href="mailto:contact@exampleurl.com">8585233222</a></div>
				</li>
			  </ul>
					</div>
		  </div>
		</div>
	  </div>
	</section>
	<!-- /Contact-us--> 
	
	
	<!--Footer -->
	
	<footer>
	  <div class="footer-top">
		<div class="container">
		  <div class="row">
		  
			<div class="col-md-6">
			  <h6>About Us</h6>
			  <ul>
	
			
			  <li><a href="page.php?type=aboutus.php">About Us</a></li>
				<li><a href="page.php?type=faqs/[h[">FAQs</a></li>
				<li><a href="page.php?type=privacy/php">Privacy</a></li>
			  <li><a href="page.php?type=terms.php">Terms of use</a></li>
	 
			 </ul>
			</div>
	  
			<div class="col-md-3 col-sm-6">
			  <h6>Subscribe Newsletter</h6>
			  <div class="newsletter-form">
				<form method="post">
				  <div class="form-group">
					<input type="email" name="subscriberemail" class="form-control newsletter-input" required placeholder="Enter Email Address" />
				  </div>
				  <button type="submit" name="emailsubscibe" class="btn btn-block">Subscribe <span class="angle_arrow"><i class="fa fa-angle-right" aria-hidden="true"></i></span></button>
				</form>
				<p class="subscribed-text">*We send great deals and latest auto news to our subscribed users very week.</p>
			  </div>
			</div>
		  </div>
		</div>
	  </div>
	  <div class="footer-bottom">
		<div class="container">
		  <div class="row">
			<div class="col-md-6 col-md-push-6 text-right">
			  <div class="footer_widget">
				<p>Connect with Us:</p>
				<ul>
				  <li><a href="#"><i class="fa fa-facebook-square" aria-hidden="true"></i></a></li>
				  <li><a href="#"><i class="fa fa-twitter-square" aria-hidden="true"></i></a></li>
				  <li><a href="#"><i class="fa fa-linkedin-square" aria-hidden="true"></i></a></li>
				  <li><a href="#"><i class="fa fa-google-plus-square" aria-hidden="true"></i></a></li>
				  <li><a href="#"><i class="fa fa-instagram" aria-hidden="true"></i></a></li>
				</ul>
			  </div>
			</div>
			<div class="col-md-6 col-md-pull-6">
			  <p class="copy-right">Copyright &copy; 2017 Astley car rental. All Rights Reserved</p>
			</div>
		  </div>
		</div>
	  </div>
	</footer><!-- /Footer--> 
	
	<!--Back to top-->
	<div id="back-top" class="back-top"> <a href="#top"><i class="fa fa-angle-up" aria-hidden="true"></i> </a> </div>
	<!--/Back to top--> 
	
	<!--Login-Form -->
	
	<div class="modal fade" id="loginform">
	  <div class="modal-dialog" role="document">
		<div class="modal-content">
		  <div class="modal-header">
			<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
			<h3 class="modal-title">Login</h3>
		  </div>
		  <div class="modal-body">
			<div class="row">
			  <div class="login_wrap">
				<div class="col-md-12 col-sm-6">
				  <form method="post">
					<div class="form-group">
					  <input type="text" class="form-control" name="email" placeholder="Email address*">
					</div>
					<div class="form-group">
					  <input type="password" class="form-control" name="password" placeholder="Password*">
					</div>
					<div class="form-group checkbox">
					  <input type="checkbox" id="remember">
				   
					</div>
					<div class="form-group">
					  <input type="submit" name="login" value="Login" class="btn btn-block">
					</div>
				  </form>
				</div>
			   
			  </div>
			</div>
		  </div>
		  <div class="modal-footer text-center">
			<p>Don't have an account? <a href="#signupform" data-toggle="modal" data-dismiss="modal">Signup Here</a></p>
			<p><a href="#forgotpassword" data-toggle="modal" data-dismiss="modal">Forgot Password ?</a></p>
		  </div>
		</div>
	  </div>
	</div><!--/Login-Form --> 
	
	<!--Register-Form -->
	
	
	<script>
	function checkAvailability() {
	$("#loaderIcon").show();
	jQuery.ajax({
	url: "check_availability.php",
	data:'emailid='+$("#emailid").val(),
	type: "POST",
	success:function(data){
	$("#user-availability-status").html(data);
	$("#loaderIcon").hide();
	},
	error:function (){}
	});
	}
	</script>
	<script type="text/javascript">
	function valid()
	{
	if(document.signup.password.value!= document.signup.confirmpassword.value)
	{
	alert("Password and Confirm Password Field do not match  !!");
	document.signup.confirmpassword.focus();
	return false;
	}
	return true;
	}
	</script>
	<div class="modal fade" id="signupform">
	  <div class="modal-dialog" role="document">
		<div class="modal-content">
		  <div class="modal-header">
			<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
			<h3 class="modal-title">Sign Up</h3>
		  </div>
		  <div class="modal-body">
			<div class="row">
			  <div class="signup_wrap">
				<div class="col-md-12 col-sm-6">
				  <form  method="post" name="signup" onSubmit="return valid();">
					<div class="form-group">
					  <input type="text" class="form-control" name="fullname" placeholder="Full Name" required="required">
					</div>
						  <div class="form-group">
					  <input type="text" class="form-control" name="mobileno" placeholder="Mobile Number" maxlength="10" required="required">
					</div>
					<div class="form-group">
					  <input type="email" class="form-control" name="emailid" id="emailid" onBlur="checkAvailability()" placeholder="Email Address" required="required">
					   <span id="user-availability-status" style="font-size:12px;"></span> 
					</div>
					<div class="form-group">
					  <input type="password" class="form-control" name="password" placeholder="Password" required="required">
					</div>
					<div class="form-group">
					  <input type="password" class="form-control" name="confirmpassword" placeholder="Confirm Password" required="required">
					</div>
					<div class="form-group checkbox">
					  <input type="checkbox" id="terms_agree" required="required" checked="">
					  <label for="terms_agree">I Agree with <a href="#">Terms and Conditions</a></label>
					</div>
					<div class="form-group">
					  <input type="submit" value="Sign Up" name="signup" id="submit" class="btn btn-block">
					</div>
				  </form>
				</div>
				
			  </div>
			</div>
		  </div>
		  <div class="modal-footer text-center">
			<p>Already got an account? <a href="#loginform" data-toggle="modal" data-dismiss="modal">Login Here</a></p>
		  </div>
		</div>
	  </div>
	</div>
	<!--/Register-Form --> 
	
	<!--Forgot-password-Form -->
	<!--/Forgot-password-Form --> 
	
	<!-- Scripts --> 
	<script src="assets/js/jquery.min.js"></script>
	<script src="assets/js/bootstrap.min.js"></script> 
	<script src="assets/js/interface.js"></script> 
	<!--Switcher-->
	<script src="assets/switcher/js/switcher.js"></script>
	<!--bootstrap-slider-JS--> 
	<script src="assets/js/bootstrap-slider.min.js"></script> 
	<!--Slider-JS--> 
	<script src="assets/js/slick.min.js"></script> 
	<script src="assets/js/owl.carousel.min.js"></script>
	
	</body>
	
	`
	// body, e := ioutil.ReadAll(r.Body)
	// if e != nil {
	// 	// log.Println(e)
	// 	fmt.Printf("%s", e)
	// }

	// doc, err := goquery.NewDocumentFromResponse(r)	//res.Body)
	// if err != nil {
	// 	// log.Fatal(err)
	// 	fmt.Println("ERROR WHILE PARSING RES.BODY")
	// }

	// // res := doc.Find("form").Contents().Text()
	// // fmt.Println("RES", res)
	// // // doc.Find("form").Each(func(i int, s *goquery.Selection) {
	// // // 	method, _ := s.Attr("method")
	// // // 	fmt.Println(method, s.Text())
	// // // 	fmt.Printf("Description field: %s\n", method)
	// // // })
	// doc.Find("form").Each(func(index int, item *goquery.Selection) {
    //     title := item.Text()
    //     // linkTag := item.Find("a")
    //     // link, _ := linkTag.Attr("href")
    //     fmt.Printf("Post #%d: %s\n", index, title)
	// })
	// fmt.Println(reflect.TypeOf(doc))
	// fmt.Println(body)
    // txt,e := doc.Find("form").Html()
	// if e!=nil {
	// 	fmt.Println("EEEEEEEEEEEEE")
	// }
	// fmt.Println(txt)
	// doc.Find("body > form").Each(func(i int, s *goquery.Selection) {
	// })

	

	reeeee := bytes.NewReader([]byte(reee))

	doc, err := goquery.NewDocumentFromReader(reeeee)
    if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Println("doc is ",doc)

//		<button class="btn" type="submit" name="send" type="submit">Send Message <span class="angle_arrow"><i class="fa fa-angle-right" aria-hidden="true"></i></span></button>
/*
<form  method="post">
<div class="form-group">
  <label class="control-label">Full Name <span>*</span></label>
  <input type="text" name="fullname" class="form-control white_bg" id="fullname" required>
</div>
<div class="form-group">
  <label class="control-label">Email Address <span>*</span></label>
  <input type="email" name="email" class="form-control white_bg" id="emailaddress" required>
</div>
<div class="form-group">
  <label class="control-label">Phone Number <span>*</span></label>
  <input type="text" name="contactno" class="form-control white_bg" id="phonenumber" required>
</div>
<div class="form-group">
  <label class="control-label">Message <span>*</span></label>
  <textarea class="form-control white_bg" name="message" rows="4" required></textarea>
</div>
<div class="form-group">
  <button class="btn" type="submit" name="send" type="submit">Send Message <span class="angle_arrow"><i class="fa fa-angle-right" aria-hidden="true"></i></span></button>
</div>
</form>
*/

		/*
			  <form action="#" method="get" id="header-search-form">
				<input type="text" placeholder="Search..." class="form-control">
				<button type="submit"><i class="fa fa-search" aria-hidden="true"></i></button>
			  </form>

		*/


	doc.Find("body form").Each(func(i int, s *goquery.Selection) {
		// fmt.Println("YOOOOOOOOOOOOOOO")
		// input := s.Find("input").Text()

		// <input type="text" name="fullname" class="form-control white_bg" id="fullname" required>


		inputsSelector := new(goquery.Selection)
		inputsSelector = visitNodes(inputsSelector, s, "input")
		//n := inputsSelector
		nT := inputsSelector.Text()
		nt,_ := inputsSelector.Attr("type")


		buttonsSelector := new(goquery.Selection)
		buttonsSelector = visitNodes(buttonsSelector, s, "button")
		// fmt.Println(buttonsSelector.Size())
		b := buttonsSelector
		bT := buttonsSelector.Text()
		bt,_ := b.Attr("type")
		bn,_ := b.Attr("name")

		button := s.Find("button")
		bText := s.Find("button").Text()
		btype,_ := button.Attr("type")
		bname,_ := button.Attr("name")
		
		fmt.Printf("Form > Input: %d: %s - %s\n", i, nT, nt)	//input, bText)
		fmt.Printf("Form > Button %d: %s - %s - %s\n", i, bText, btype, bname)
		fmt.Printf("Form > ButtoN %d: %s - %s - %s\n", i, bT, bt, bn)
	})

	// query, e := goquery.NewDocumentFromResponse(r)		//Reader(strings.NewReader(someHtml))
	// if e!=nil {
	// 	fmt.Println("EEEEEEEEEEEEE")
	// }
	// sel:= query.Find("body form")
	// fmt.Println(sel.Text())

	// query.Find("body form").Each(func(i int, s *goquery.Selection) {
	// 	s.ChildrenFiltered("ul")
	// })

	// doc.Find("a").Siblings().Each(func(i int, s *goquery.Selection) {
	// 	fmt.Printf("%d, Sibling text: %s\n", i, s.Text())
	// })

	// Children().Each().Contents().Each(func(i int, s *goquery.Selection) {		//Each(func(i int, s *goquery.Selection) 
	// 	html, _ := s.Html()

	// 	fmt.Println("-------------\r\n",html,"\r\n-------------\r\n")
	// })
	// // Find form input elements
	// doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
	// 	// For each item found, get the band and title
	// 	band := s.Find("a").Text()
	// 	title := s.Find("i").Text()
	// 	fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// 	})

	//	Validated
	// test()

}

func visitNodes(dst, src *goquery.Selection, tagname string) *goquery.Selection {
	// var j int = 0
	// fmt.Println("-----",tagname,"-----")
	fmt.Println("--------------------------------------------------------------PRE")
	src.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == tagname {
			dst = dst.AddSelection(s)
		} else {
			// fmt.Println(tagname, "====", dst.Children().Size())
			// dst = dst.Contents().Children()
			// dst = visitNodes(dst, s, tagname)
			// fmt.Println("ITERATE NOW", dst.Children())
			// src.Contents().Each(func(i int),  *goquery.Selection)

			
			// fmt.Println(goquery.NodeName(s))
			s.Contents().Siblings().Each(func(ii int, ss *goquery.Selection) {
				
				// fmt.Println(goquery.NodeName(ss))
				if goquery.NodeName(ss) == tagname {
				// 	j += 1
					fmt.Println("YES")
					dst = dst.AddSelection(ss)

					} else {
					dst = visitNodes(dst, ss, tagname)
				}
			})
			
		}
	})
	// fmt.Println(j)
	fmt.Println("--------------------------------------------------------------AFT")

	return dst
}

func test() {
	var s = `<html><body>
<form name="query" action="http://www.example.net/action.php" method="post">
    <textarea type="text" name="nameiknow">The text I want</textarea>
    <div id="button">
        <input type="submit" value="Submit" />
    </div>
</form>
</body></html>`

	r := bytes.NewReader([]byte(s))
    doc, _ := goquery.NewDocumentFromReader(r)
    // text := doc.Find("textarea").Text()
    text,_ := doc.Find("form").Html()
	fmt.Println(text)
}

func onlinetest() {
	doc, err := goquery.NewDocument("http://metalsucks.net")
    if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Println("doc is ",doc)
    doc.Find("body div ul li").Each(func(i int, s *goquery.Selection) {
		fmt.Println("YOOOOOOOOOOOOOOO")
		band := s.Find("a").Text()
        title := s.Find("href").Text()
        fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}