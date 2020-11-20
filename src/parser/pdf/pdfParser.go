package pdf

import "book_parser/src"
import "github.com/flotzilla/pdf_parser"

func Parse(bookFile *src.BookFile) *src.BookInfo {
	pdfInfo := pdf_parser.ParsePdf(bookFile.FilePath)

	bI := src.BookInfo{
		Title:         pdfInfo.GetTitle(),
		Author:        pdfInfo.GetAuthor(),
		Pages:         pdfInfo.PagesCount,
		ISBN:          pdfInfo.GetISBN(),
		NumberOfPages: pdfInfo.PagesCount, // TODO wtf is number of pages
		Subject:       "",
		Description:   pdfInfo.GetDescription(),
	}

	return &bI
}
