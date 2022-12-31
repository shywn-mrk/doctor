package grouper

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode"

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
	content content
}

type titleToGroup map[string]*Group

func (g *Group) add(grouping []string, sortKey, line string) *strings.Builder {
	if len(grouping) == 0 {
		if g.content == nil {
			g.content = make(content)
		}
		return g.content.add(sortKey, line)
	}

	if g.childs == nil {
		g.childs = make(titleToGroup)
	}

	child, ok := g.childs[grouping[0]]
	if !ok || child == nil {
		child = new(Group)
	}
	c := child.add(grouping[1:], sortKey, line)
	g.childs[grouping[0]] = child
	return c

	// if child, ok := g.childs[grouping[0]]; !ok || child == nil {
	// 	child = new(Group)
	// 	c := child.add(grouping[1:], sortKey, line)
	// 	g.childs[grouping[0]] = child
	// 	return c
	// } else {
	// 	return child.add(grouping[1:], sortKey, line)
	// }
}

func (g *Group) Build() string {
	b := strings.Builder{}
	g.build(&b)
	return b.String()
}

func (g *Group) build(b *strings.Builder) {
	if g.content != nil {
		sortedContents := sortByKeys(g.content)

		for _, c := range sortedContents {
			fmt.Fprintf(b, c.value.String()+"\n")
		}
	}

	sortedChilds := sortByKeys(g.childs)

	for _, c := range sortedChilds {
		fmt.Fprintf(b, "%s\n\n", c.key)
		c.value.build(b)
	}
}

type contentPointer map[string]*strings.Builder

type content map[string]*strings.Builder

func (c content) add(sortKey, d string) *strings.Builder {
	if c == nil {
		panic("content is nil")
	}

	b, ok := c[sortKey]
	if !ok {
		b = new(strings.Builder)
	}

	b.WriteString(d)
	c[sortKey] = b

	return b
}

func Parse(i *linarian.Linarian) (string, error) {
	var gs Group
	p := make(contentPointer)

	err := func() error {
		var grouping []string
		var gsum, sortKey string
		var ok bool
		for {
			line, err := i.ReadLine()
			if err != nil {
				return err
			}

			if doctorRegex.MatchString(line) {
				options := extractOptions(line)

				sortKey, ok = options["sort"]
				if !ok {
					sortKey = ""
				}

				grouping, err = extractGroup(i)
				if err != nil {
					return err
				}
				gsum = strings.Join(grouping, "") + sortKey
			} else if c, ok := p[gsum]; ok {
				c.WriteString(line)
			} else {
				p[gsum] = gs.add(grouping, sortKey, line)
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

type options map[string]string

func extractOptions(data string) options {
	inQuote := false
	f := func(c rune) bool {
		switch {
		case unicode.In(c, unicode.Quotation_Mark):
			inQuote = !inQuote
		case unicode.In(c, unicode.White_Space):
			return !inQuote
		}
		return false
	}
	items := strings.FieldsFunc(data, f)

	m := make(options)
	for _, item := range items {
		x := strings.Split(item, "=")
		if len(x) < 2 {
			continue
		}

		m[x[0]] = strings.Join(x[1:], "")
	}

	return m
}

type keyValue[V any] struct {
	key   string
	value V
}

func sortByKeys[V any](data map[string]V) []keyValue[V] {
	result := make([]keyValue[V], 0, len(data))
	for k, v := range data {
		result = append(result, keyValue[V]{
			key:   k,
			value: v,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].key, result[j].key) != 1
	})

	return result
}
