package walker

import (
	"fmt"
	"regexp"
	"strings"
)

func search(html string, domain string) []string {
	pattern := `<a\s+href="(.*?)">(.*?)</a>`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Find matches in the HTML
	matches := regex.FindAllStringSubmatch(html, -1)

	var res []string
	for _, match := range matches {
		if strings.HasPrefix(match[1], "http") {
			res = append(res, match[1])
		} else {
			res = append(res, fmt.Sprintf("%s/%s", domain, match[1]))
		}
	}

	return res
}
