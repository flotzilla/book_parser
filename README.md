[![MIT License][license-shield]][license-url]

# BOOK PARSER

Console book scanner-parser for book cataloging , written in go

#### Supported formats
- [pdf](https://github.com/flotzilla/pdf_parser)
- fb2

#### Example
Build and run: `go build book_parser.go`

Executing: `./book_parser` - will scan current directory


Scan directories: 
```bash
./book_parser scan directory_name directory_name1
```

#### Configuration
Main configuration here: `conf/config.json`

Config logrus logger here: `src/logging/logging.go`

#### Links
* [fb2 description](http://www.fictionbook.org/index.php/%D0%9E%D0%BF%D0%B8%D1%81%D0%B0%D0%BD%D0%B8%D0%B5_%D1%84%D0%BE%D1%80%D0%BC%D0%B0%D1%82%D0%B0_FB2_%D0%BE%D1%82_Sclex)

### License
MIT 

[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=flat-square
[license-url]: https://github.com/flotzilla/book_parser/blob/main/LICENSE
