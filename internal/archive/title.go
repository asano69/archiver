package archive

import (
	"html"
	"regexp"
	"strings"
)

// titleTagRe matches the first <title> element in an HTML document.
// (?is) makes it case-insensitive and lets "." also match newlines,
// since a title can legitimately span multiple lines in the markup.
var titleTagRe = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

// ExtractTitle returns the page's <title> text, or fallback if none is
// found. Kept in its own file so a more elaborate strategy (e.g.
// falling back to an Open Graph tag) can be added later without
// touching the rest of the package.
func ExtractTitle(pageHTML []byte, fallback string) string {
	match := titleTagRe.FindSubmatch(pageHTML)
	if match == nil {
		return fallback
	}

	title := html.UnescapeString(string(match[1]))
	title = strings.Join(strings.Fields(title), " ")
	if title == "" {
		return fallback
	}
	return title
}
