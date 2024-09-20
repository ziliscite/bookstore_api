package tools

import (
	"regexp"
	"strings"
)

// Slugify function to convert a string into a URL-friendly slug
func Slugify(s string) string {
	// Convert the string to lowercase
	s = strings.ToLower(s)

	// Replace spaces and underscores with a hyphen
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove all non-alphanumeric characters except for hyphens
	re := regexp.MustCompile(`[^\w-]+`)
	s = re.ReplaceAllString(s, "")

	// Replace multiple hyphens with a single hyphen
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")

	// Trim leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
