package pdf_parser

import "book_parser/src"
import "github.com/flotzilla/pdf_parser"

type PdfParser struct{}

// wrapper around pdf_parser package
func (parser *PdfParser) Parse(bookFile *src.BookFile, withCover bool) (*src.BookInfo, error) {
	pdfInfo, err := pdf_parser.ParsePdf(bookFile.FilePath)

	if err != nil {
		return nil, err
	}

	// TODO fix this
	bI := src.BookInfo{
		Title:         pdfInfo.GetTitle(),
		Author:        pdfInfo.GetAuthor(),
		ISBN:          pdfInfo.GetISBN(),
		NumberOfPages: pdfInfo.PagesCount,
		Description:   pdfInfo.GetDescription(),
	}

	// TODO finish this
	//if withCover {
	//bI.CoverPage = pdf_parser.
	//}

	return &bI, nil
}
