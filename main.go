package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/733amir/doctor/linarian"
)

var (
	doctorRegex = regexp.MustCompile(`^@doctor`)

	groupingRegexes = []*regexp.Regexp{
		regexp.MustCompile(`^#\s[^#]`),
		regexp.MustCompile(`^##\s[^#]`),
		regexp.MustCompile(`^###\s[^#]`),
		regexp.MustCompile(`^####\s[^#]`),
		regexp.MustCompile(`^#####\s[^#]`),
		regexp.MustCompile(`^######\s[^#]`),
	}
)

func main() {
	i := linarian.New(bufio.NewReader(os.Stdin), 2)

	var gs group
	p := make(pointer)
	err := func() error {
		var grouping []string
		var gsum string
		for {
			line, err := i.ReadLine()
			if err != nil {
				return err
			}

			if doctorRegex.MatchString(line) {
				grouping, err = extractGroup(i)
				if err != nil {
					return err
				}
				gsum = strings.Join(grouping, "")
			} else if c, ok := p[gsum]; ok {
				c.WriteString(line)
			} else {
				p[gsum] = gs.add(grouping, line)
			}
		}
	}()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	gs.print()
}

func extractGroup(i *linarian.Linarian) ([]string, error) {
	grouping := make([]string, 0, 2)
	for {
		line, err := i.ReadLine()
		if err != nil {
			i.JumpBack(1)
			return grouping, err
		}

		if !groupingRegexes[len(grouping)].MatchString(line) {
			i.JumpBack(1)
			return grouping, nil
		}

		grouping = append(grouping, line)
	}
}

type group struct {
	childs  titleGroupMap
	content *strings.Builder
}

type titleGroupMap map[string]group

type pointer map[string]*strings.Builder

func (g *group) add(grouping []string, line string) *strings.Builder {
	if len(grouping) == 0 {
		if g.content == nil {
			g.content = new(strings.Builder)
		}
		g.content.WriteString(line)
		return g.content
	}

	if g.childs == nil {
		g.childs = make(titleGroupMap)
	}

	if child, ok := g.childs[grouping[0]]; !ok {
		c := child.add(grouping[1:], line)
		g.childs[grouping[0]] = child
		return c
	} else {
		return child.add(grouping[1:], line)
	}
}

func (g *group) print() {
	if g.content != nil {
		fmt.Print(g.content.String() + "\n")
	}

	titles := make([]string, 0, len(g.childs))
	for t := range g.childs {
		titles = append(titles, t)
	}

	sort.Strings(titles)

	for _, t := range titles {
		c := g.childs[t]
		fmt.Printf("%s\n", t)
		c.print()
	}
}
