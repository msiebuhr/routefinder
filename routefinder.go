package routefinder

// Package routefinder implements functions for decomposing known dynamic URLs into templates.

import (
	"fmt"
	"regexp"
	"strings"
)

type Routes []struct {
	name string
	re   regexp.Regexp
}

// Constructs a Route-table from the given templates. Everything in the
// template from a : to the following / is considered a variable.
func NewRoutefinder(templates ...string) (Routes, error) {
	converterRegex, err := regexp.Compile(":[^/]+")

	if err != nil {
		return Routes{}, err
	}

	// Regex placeholders out of the template URLs
	routes := make(Routes, len(templates))

	for i, template := range templates {
		// Quote slashes &c.
		withQuotedMeta := regexp.QuoteMeta(template)

		// Switch out :[^\/]+ for capture groups
		withQuotedMeta = converterRegex.ReplaceAllStringFunc(withQuotedMeta, func(group string) string {
			return fmt.Sprintf("(?P<%s>[^/]+)", group[1:])
		})

		// Add start and end guards
		withQuotedMeta = fmt.Sprintf("^%s$", withQuotedMeta)

		re, err := regexp.Compile(withQuotedMeta)

		if err != nil {
			return Routes{}, err
		}

		routes[i].name = template
		routes[i].re = *re
	}

	return routes, nil
}

// Look up the given path in the available Routes, first-match-wins.  A match
// will return the original template string along with a map of the parsed-out
// variables. Lookup returns empty values, if no match is found.
func (r Routes) Lookup(path string) (string, map[string]string) {
	// Dump any query string
	normalizedPath := strings.SplitN(path, "?", 1)[0]

	// Check key against regex'es
	for _, template := range r {
		// Check key against regexes
		if !template.re.MatchString(normalizedPath) {
			continue
		}

		subMatchNames := template.re.SubexpNames()
		subMatchValues := template.re.FindStringSubmatch(normalizedPath)

		meta := make(map[string]string, len(subMatchNames))

		for j := 1; j < len(subMatchNames); j += 1 {
			meta[subMatchNames[j]] = subMatchValues[j]
		}

		return template.name, meta
	}

	return "", map[string]string{}
}
