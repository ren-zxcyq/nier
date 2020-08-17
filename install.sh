#!/bin/sh

#   Some or most of the tools used exist in open github repositories.
#       git clone --depth 1 https://github.com/sqlmapproject/sqlmap.git sqlmap-dev

cd && sudo apt install -y nmap whatweb wapiti wpscan python3 sqlmap golang-1.13

#   go.mod contents
#   
#   module github.com/ren-zxcyq/nier
#   
#   go 1.13
#   
#   require (
#   	github.com/PuerkitoBio/goquery 6f0f9d6b87b3ab53ab8dfc9cdda07ba5bbdd57a3
#   	github.com/dchest/uniuri 7aecb25e1fe5a22533fab90a637a8f74a9cf7340
#   	github.com/jung-kurt/gofpdf a2a0e7f8a28b2eabe1a32097f0071a0f715a8102
#   	github.com/tebeka/selenium 9a0798fcb455aca4de72bbe424f4bbb9cb021f53
#   	github.com/OJ/gobuster/v3 372f640b754c28b750ba9ae8f3298c81ed0f9759
#   )

#   tebeka selenium  v0.9.9
#   gofpdf   v1.16.2
#   goquery  v1.5.1
#   uniuri   v0.0.0-20200228104902-7aecb25e1fe5

#   Install tebeka/selenium     (Golang Selenium Headless Driver)
#go get -t -d github.com/tebeka/selenium
#sudo cd /opt
#sudo mkdir tebeka-selenium && cd tebeka-selenium/
#sudo git clone https://github.com/tebeka/selenium
#sudo git checkout 9a0798fcb455aca4de72bbe424f4bbb9cb021f53
cd ~/go/src/github.com/tebeka/selenium/vendor
go run init.go --alsologtostderr  --download_browsers --download_latest
#sudo cd ..
#   Test Selenium Installation
sudo apt-get install xvfb openjdk-11-jre
#   sudo go test         #   We just need FirefoxGeckoDriver

#   Install XSStrike
sudo cd /opt
sudo git clone https://github.com/s0md3v/XSStrike
sudo git checkout 0ecedc1bba149931e3b32e53422d5b7c089ba9dc
