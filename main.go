package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	doctorRegex = regexp.MustCompile(`^/\*+?\s+?\@doctor`)

	headerTemplate = "$header"
	headerRegexes  = []*regexp.Regexp{
		regexp.MustCompile(`^\*+?\s+?(?P<header>\#\s+?.*)$`),
		regexp.MustCompile(`^\*+?\s+?(?P<header>\##\s+?.*)$`),
		regexp.MustCompile(`^\*+?\s+?(?P<header>\###\s+?.*)$`),
		regexp.MustCompile(`^\*+?\s+?(?P<header>\####\s+?.*)$`),
		regexp.MustCompile(`^\*+?\s+?(?P<header>\#####\s+?.*)$`),
		regexp.MustCompile(`^\*+?\s+?(?P<header>\######\s+?.*)$`),
	}

	leadingRegex  = regexp.MustCompile(`^\s*\**\s`)
	trailingRegex = regexp.MustCompile(`\*/`)
)

func main() {
	i := NewLastLines(bufio.NewReader(os.Stdin), 2)

	for ii := 0; ii < 100; ii++ {
		// fmt.Printf("=> i: %#v\n", i)
		line, err := i.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		line = strings.Trim(line, " \t\n\r")
		// fmt.Printf("=> %v\n", line)

		if doctorRegex.MatchString(line) {
			// Extracting group information
			headers, err := extractGroup(i)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", headers)
		} else {
			fmt.Printf("%s\n", removeTrailing(removeLeading(line)))
		}

	}
}

func extractGroup(i *LastLines) ([]string, error) {
	headers := []string{}
	for len(headers) < 6 {
		// fmt.Printf("==> i: %#v\n", i)
		line, err := i.ReadLine()
		if err != nil {
			// fmt.Printf("==> err: %v\n", err)
			i.JumpBack(1)
			return headers, err
		}
		line = strings.Trim(line, " \t\n\r")

		p := headerRegexes[len(headers)]
		if !p.MatchString(line) {
			// fmt.Printf("==> no: %v\n", line)
			// fmt.Printf("==> i: %#v\n", i)
			i.JumpBack(1)
			// fmt.Printf("==> i: %#v\n", i)
			return headers, nil
		}

		for _, submatches := range p.FindAllStringSubmatchIndex(line, -1) {
			headers = append(headers, string(p.ExpandString(nil, headerTemplate, line, submatches)))
		}
		// fmt.Printf("==> h: %v\n", headers)
	}

	return headers, nil
}

func removeLeading(line string) string {
	return leadingRegex.ReplaceAllLiteralString(line, "")
}

func removeTrailing(line string) string {
	return trailingRegex.ReplaceAllLiteralString(line, "")
}
