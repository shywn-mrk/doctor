package markdown

import "github.com/gomarkdown/markdown"

func GenerateHTML(m string) (string, error) {
	h := markdown.ToHTML([]byte(m), nil, nil)

	return string(h), nil
}
