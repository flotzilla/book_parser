package main

import (
	"book_parser/src"
	_ "book_parser/src/logging"
	"book_parser/src/parser"
	"book_parser/src/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type CommandRunner func([]string)

type Command struct {
	Name        string
	Description string
	Runner      CommandRunner
}

var (
	commands map[string]Command
	cnf      src.Config
)

func init() {
	conf, err := utils.GetConfig("../conf/config.json")
	if err != nil {
		logrus.Error(err)
	}

	cnf = *conf

	logrus.Debug("Configuration:",
		"\n * ScanExt: ", cnf.ScanExt,
		"\n * SkippedExt: ", cnf.SkippedExt,
		"\n * With Covers: ", cnf.WithCoverImages)
}

func main() {
	commands = map[string]Command{
		"help": Command{
			Name:        "help",
			Description: "Show help info",
			Runner:      help,
		},
		"scan": Command{
			Name:        "scan",
			Description: "Will scan folder for books in case of any path arguments or read config file",
			Runner:      scan,
		},
		"sync": Command{
			Name:        "sync",
			Description: "Sync parsed data to central database",
			Runner:      syncCommand,
		},
	}

	handleCommand(os.Args[1:])
}

func help(args []string) {
	fmt.Println("Available commands:")
	for _, v := range commands {
		fmt.Printf(" %s - %s\n", v.Name, v.Description)
	}
}

func scan(args []string) {
	if len(args) > 0 {
		var wg sync.WaitGroup
		logrus.Trace(args)
		for _, el := range args {
			wg.Add(1)

			el := el
			go func(path string) {
				scanAndParse(el)
				defer wg.Done()
			}(el)
		}

		wg.Wait()
		logrus.Debug("All scans finished")
	} else {
		currPath, err := os.Getwd()

		if err != nil {
			logrus.Error(err)
		}
		scanAndParse(currPath)
	}
}

func scanAndParse(currPath string) {
	logrus.Info("Scanning directory: ", currPath)
	sc, err := src.Scan(currPath, &cnf)

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	showScanResult(&sc)
	pr := parseScanResult(&sc)
	showParseResult(pr)
}

func parseScanResult(result *src.ScanResult) *src.ParseResult {
	if len(result.Books) != 0 {
		pR := parser.Parse(result, &cnf)

		if len(pR.Errors) != 0 {
			logrus.Error("Some errors found")
			for _, err := range pR.Errors {
				logrus.Error(err)
			}
		}

		return pR
	}

	return nil
}

func showScanResult(sc *src.ScanResult) {
	logrus.Info("Found files: ", sc.BooksFoundTotalCount, ", skipped: ", sc.BooksSkippedCount,
		", total: ", sc.BooksTotalCount)
}

func showParseResult(pr *src.ParseResult) {
	if pr == nil {
		return
	}
	for _, el := range pr.Books {
		logrus.Info(el.BookInfo.Title)
		logrus.Info(el.BookInfo.Authors)
	}
}

func syncCommand(args []string) {
	logrus.Info("Showing args", args)
	// TODO send results to server
}

func handleCommand(args []string) {
	aLength := len(args)
	if aLength == 0 {
		help(args)
	} else if aLength >= 1 {
		if c, found := commands[args[0]]; found {
			c.Runner(args[1:])
		} else {
			fmt.Println("There is no such command")
		}
	}
}
