package util

import (
	"regexp"
)

var domainComponentRegexp = regexp.MustCompile(`(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])`)

func expression(res ...*regexp.Regexp) *regexp.Regexp {
	var s string
	for _, re := range res {
		s += re.String()
	}

	return regexp.MustCompile(s)
}

// optional wraps the expression in a non-capturing group and makes the
// production optional.
func optional(res ...*regexp.Regexp) *regexp.Regexp {
	return regexp.MustCompile(group(expression(res...)).String() + `?`)
}

// repeated wraps the regexp in a non-capturing group to get one or more
// matches.
func repeated(res ...*regexp.Regexp) *regexp.Regexp {
	return regexp.MustCompile(group(expression(res...)).String() + `+`)
}

// group wraps the regexp in a non-capturing group.
func group(res ...*regexp.Regexp) *regexp.Regexp {
	return regexp.MustCompile(`(?:` + expression(res...).String() + `)`)
}

// capture wraps the expression in a capturing group.
func capture(res ...*regexp.Regexp) *regexp.Regexp {
	return regexp.MustCompile(`(` + expression(res...).String() + `)`)
}

// anchored anchors the regular expression by adding start and end delimiters.
func anchored(res ...*regexp.Regexp) *regexp.Regexp {
	return regexp.MustCompile(`^` + expression(res...).String() + `$`)
}

// literal compiles s into a literal regular expression, escaping any regexp
// reserved characters.
func literal(s string) *regexp.Regexp {
	re := regexp.MustCompile(regexp.QuoteMeta(s))

	if _, complete := re.LiteralPrefix(); !complete {
		panic("must be a literal")
	}

	return re
}

func CheckForRegistry(ref string) bool {
	var DomainRegexp = expression(domainComponentRegexp, repeated(literal(`.`), domainComponentRegexp), optional(literal(`:`), regexp.MustCompile(`[0-9]+`)))
	return DomainRegexp.MatchString(ref)
}
