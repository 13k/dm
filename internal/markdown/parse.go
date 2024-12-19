package markdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

//nolint:ireturn // must return a Node
func parseBuffer(buf []byte) ast.Node {
	parser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	).Parser()

	reader := text.NewReader(buf)

	return parser.Parse(reader)
}

// ParseList parses markdown source from the given buffer, finds the first top-level list and
// returns the list items source. Returns nil if no top-list list was found.
func ParseList(src []byte) []string {
	var (
		root ast.Node
		list ast.Node
	)

	root = parseBuffer(src)

	for c := root.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == ast.KindList {
			list = c

			break
		}
	}

	if list == nil {
		return nil
	}

	var items []string

	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		if item.Kind() != ast.KindListItem {
			continue
		}

		var itemSrc []byte

		for itemChild := item.FirstChild(); itemChild != nil; itemChild = itemChild.NextSibling() {
			switch itemChild.Kind() {
			case ast.KindTextBlock, ast.KindHeading:
			default:
				continue
			}

			for i := range itemChild.Lines().Len() {
				line := itemChild.Lines().At(i)
				itemSrc = append(itemSrc, line.Value(src)...)
			}
		}

		items = append(items, string(itemSrc))
	}

	return items
}
