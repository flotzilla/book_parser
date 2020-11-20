package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ScanExt    []string
	SkippedExt []string
}

type ScanResult struct {
	BooksFoundTotalCount int
	BooksSkippedCount    int
	BooksTotalCount      int

	Books []BookFile
}

func GetConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)

	defer file.Close()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(file)
	conf := Config{}

	err = dec.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func Scan(filePath string, cnf Config) ScanResult {
	sc := ScanResult{}

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

	if err != nil {
		fmt.Println(err)
	}

	return sc
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

	return b
}
