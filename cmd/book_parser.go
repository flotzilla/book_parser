package main

import (
	"book_parser/src"
	_ "book_parser/src/logging"
	"book_parser/src/parser"
	"book_parser/src/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
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

func help(args []string) {
	logrus.Info("Available commands:")
	for _, v := range commands {
		fmt.Printf("%s - %s\n", v.Name, v.Description)
	}
}

func init() {
	cnf, err := utils.GetConfig("../conf/config.json")
	if err != nil || !cnf.HasInit {
		logrus.Error("Error during configuration init")
	}

	logrus.Debug("Configuration: \n", " * ScanExt: ", cnf.ScanExt, "\n * SkippedExt: ", cnf.SkippedExt)
}

func scan(args []string) {
	var sc src.ScanResult

	if len(args) > 0 {
		for _, el := range args {
			//TODO add gorutines
			logrus.Info("Scanning directory: ", el)
			sc = src.Scan(el, cnf)
			showScanResult(&sc)
			pr := parseScanResult(&sc)
			showParseResult(pr)
		}
	} else {
		// scan current directory
		currPath, err := os.Getwd()

		if err != nil {
			fmt.Println(err)
		}

		logrus.Info("Scanning current directory: ", currPath)
		sc = src.Scan(currPath, cnf)

		showScanResult(&sc)
		pr := parseScanResult(&sc)
		showParseResult(pr)
	}
}

func parseScanResult(result *src.ScanResult) *parser.ParseResult {
	if len(result.Books) != 0 {
		pR := parser.Parse(result)

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

func showParseResult(pr *parser.ParseResult) {
	if pr == nil {
		return
	}
	for _, el := range pr.Books {
		fmt.Println(el.BookInfo)
	}
}

func sync(args []string) {
	logrus.Info(args)
}

func handleCommand(args []string) {
	aLength := len(args)
	if aLength == 0 {
		help(args)
	} else if aLength >= 1 {
		if c, found := commands[args[0]]; found {
			c.Runner(args[1:])
		} else {
			logrus.Warn("There is no such command")
		}
	}
}

func main() {
	// TODO move logger configuration to conf/config.json file

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
			Runner:      sync,
		},
	}

	handleCommand(os.Args[1:])
}
