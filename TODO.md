TODO
=====

#### implement formats
* .mobi format
* .epub format

#### Features
* extract image
  - determine the first page
  - extract all
* fix for pdf multiple xref table
  - normal xref table / trailer [done]
  - without trailer 
  - multiple xref/table [done]
  - fix trailer invalid [done]
  - fix decrypted xref
  - test on this pdf
   * os-dev
   * kodak tri-x 400
   * kodak tmax 3200
  
  - parse root element [done]
  - parse metadata from root [done]
  - parse info from trailer [done]
* clean-up
* write tests


###### fix for updating to latest pdf_parser
```bash
$ export GO111MODULE=on
$ export GOPROXY=direct
$ export GOSUMDB=off
$ go get -v -u github.com/flotzilla/pdf_parser@version
```
