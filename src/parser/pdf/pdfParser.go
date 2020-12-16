package pdf

import "book_parser/src"
import "github.com/flotzilla/pdf_parser"

// wrapper around pdf_parser package
func Parse(bookFile *src.BookFile) (*src.BookInfo, error) {
	pdfInfo, err := pdf_parser.ParsePdf(bookFile.FilePath)

	if err != nil {
		return nil, err
	}

	bI := src.BookInfo{
		Title:         pdfInfo.GetTitle(),
		Author:        pdfInfo.GetAuthor(),
		Pages:         pdfInfo.PagesCount,
		ISBN:          pdfInfo.GetISBN(),
		NumberOfPages: pdfInfo.PagesCount, // TODO wtf is number of pages
		Subject:       "",
		Description:   pdfInfo.GetDescription(),
	}

	return &bI, nil
}
