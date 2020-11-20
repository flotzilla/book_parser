package main

import (
	"book_parser/src"
	"book_parser/src/parser"
	"fmt"
	"os"
)

type CommandRunner func([]string)

type Command struct {
	Name        string
	Description string
	Runner      CommandRunner
}

var commands map[string]Command

func help(args []string) {
	fmt.Println("Available commands:")
	for _, v := range commands {
		fmt.Printf("%s - %s\n", v.Name, v.Description)
	}
}

func scan(args []string) {

	// TODO fix this
	// import config from user folder, win/lin/mac
	config, err := src.GetConfig("/home/bbyte/GolandProjects/book_parser/conf/scanner.json")
	//fmt.Println(runtime.GOOS)
	//fmt.Println(config)

	if err != nil {
		fmt.Println("Error during configuration init")
	}

	var sc src.ScanResult

	if len(args) > 0 {
		for _, el := range args {
			fmt.Printf("Scanning directory: %s\n", el)
			sc = src.Scan(el, *config)
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

		fmt.Printf("Scanning current directory: %s\n", currPath)
		sc = src.Scan(currPath, *config)

		showScanResult(&sc)
		pr := parseScanResult(&sc)
		showParseResult(pr)
	}
}

func parseScanResult(result *src.ScanResult) *parser.ParseResult {
	if len(result.Books) != 0 {
		pR := parser.Parse(result)

		if len(pR.Errors) != 0 {
			for _, err := range pR.Errors {
				fmt.Println(err)
			}
		}

		return pR
	}

	return nil
}

func showScanResult(sc *src.ScanResult) {
	fmt.Printf("Found files: %d, skipped: %d, total: %d\n",
		sc.BooksFoundTotalCount, sc.BooksSkippedCount, sc.BooksTotalCount)

	//for _, el := range sc.Books {
	//	fmt.Println("=========")
	//	fmt.Println(el)
	//}
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
	fmt.Println(args)
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
			Runner:      sync,
		},
	}

	handleCommand(os.Args[1:])
}
