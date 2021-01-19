package pdf_parser

import "book_parser/src"
import "github.com/flotzilla/pdf_parser"

type PdfParser struct{}

// wrapper around pdf_parser package
func (parser *PdfParser) Parse(bookFile *src.BookFile, withCover bool) (*src.BookInfo, error) {
	b, err := pdf_parser.ParsePdf(bookFile.FilePath)

	if err != nil {
		return nil, err
	}

	pub := src.PublisherInfo{
		BookName:  b.GetTitle(),
		Publisher: b.GetPublisherInfo(),
		ISBN:      b.GetISBN(),
	}

	var authors []string
	authors = append(authors, b.GetAuthor())

	// TODO fix this
	bI := src.BookInfo{
		Title:         b.GetTitle(),
		Author:        b.GetAuthor(),
		Authors:       authors,
		ISBN:          b.GetISBN(),
		NumberOfPages: b.PagesCount,
		Description:   b.GetDescription(),
		Language:      b.GetLanguage(),
		Date:          b.GetDate(),
		PublisherInfo: pub,
	}

	// TODO finish this
	//if withCover {
	//bI.CoverPage = pdf_parser.
	//}

	return &bI, nil
}
