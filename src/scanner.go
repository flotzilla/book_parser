package src

import (
	cnf "book_parser/src/config"
	_ "book_parser/src/logging"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func Scan(filePath string, cnf *cnf.Config) (ScanResult, error) {
	sc := ScanResult{
		FilePath: filePath,
	}

	logrus.Trace("dive into ", filePath)
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		for _, el := range cnf.SkippedExt {
			if strings.HasSuffix(info.Name(), el) {
				sc.BooksSkippedCount++
			}
		}

		for _, el := range cnf.ScanExt {
			if strings.HasSuffix(info.Name(), el) {
				sc.BooksTotalCount++
				sc.Books = append(sc.Books, creteBookFile(info, el, path))
			}
		}
		return nil
	})

	sc.BooksFoundTotalCount = sc.BooksSkippedCount + sc.BooksTotalCount

	return sc, err
}

func creteBookFile(info os.FileInfo, ext string, path string) BookFile {
	b := BookFile{
		path,
		info.Name(),
		info.Size(),
		ext,
		info.ModTime(),
		info.Mode()&os.ModeSymlink != 0,
	}

	logrus.Trace("Created book: ", b.FilePath)
	return b
}
