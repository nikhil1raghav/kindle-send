package util

import (
	"regexp"
)

var illegalRe = regexp.MustCompile(`[\/\?<>\\:\*\|"]`)
var controlRe = regexp.MustCompile(`[\x00-\x1f\x80-\x9f]`)
var reservedRe = regexp.MustCompile(`^\.+$`)
var windowsReservedRe = regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])(\..*)?$`)
var windowsTrailingRe = regexp.MustCompile(`[\. ]+$`)
var whitespaceRe = regexp.MustCompile(`\s`)

type FileName string

func (n *FileName) sanitize(r *regexp.Regexp, replacement string) *FileName {
	var result FileName = FileName(r.ReplaceAllString(string(*n), replacement))
	return &result
}

func SanitazeFileName(fileName string, replacement string) string {
	var n FileName = FileName(fileName)
	sanitazed := n.
		sanitize(illegalRe, replacement).
		sanitize(controlRe, replacement).
		sanitize(reservedRe, replacement).
		sanitize(windowsReservedRe, replacement).
		sanitize(windowsTrailingRe, replacement)

	return string(*sanitazed)
}

func SanitazeFileNameAndReplaceWhitespaces(fileName string, replacement string) string {
	sanitazed := SanitazeFileName(fileName, replacement)
	sanitazed = whitespaceRe.ReplaceAllString(sanitazed, replacement)
	return sanitazed
}
