package grouper

import (
	"fmt"
	"io"
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

type Group struct {
	childs  titleToGroup
	content *content
}

func (g *Group) add(grouping []string, line string) *content {
	if len(grouping) == 0 {
		if g.content == nil {
			g.content = new(content)
		}
		g.content.add(line)
		return g.content
	}

	if g.childs == nil {
		g.childs = make(titleToGroup)
	}

	if child, ok := g.childs[grouping[0]]; !ok {
		c := child.add(grouping[1:], line)
		g.childs[grouping[0]] = child
		return c
	} else {
		return child.add(grouping[1:], line)
	}
}

func (g *Group) Build() string {
	b := strings.Builder{}
	g.build(&b)
	return b.String()
}

func (g *Group) build(b *strings.Builder) {
	if g.content != nil {
		fmt.Fprint(b, g.content.String()+"\n")
	}

	titles := make([]string, 0, len(g.childs))
	for t := range g.childs {
		titles = append(titles, t)
	}

	sort.Strings(titles)

	for _, t := range titles {
		c := g.childs[t]
		fmt.Fprintf(b, "%s\n\n", t)
		c.build(b)
	}
}

type contentPointer map[string]*content

type content struct {
	strings.Builder
}

func (c *content) add(d string) {
	if c == nil {
		panic("content is nil")
	}

	c.WriteString(d)
}

type titleToGroup map[string]Group

func Parse(i *linarian.Linarian) (string, error) {
	var gs Group
	p := make(contentPointer)

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
		return "", err
	}

	return gs.Build(), nil
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

		grouping = append(grouping, strings.TrimSpace(line))
	}
}
