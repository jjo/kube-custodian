package cleaner

import (
	"regexp"
)

// SystemRE set from flags at cmd/delete.go
var SystemRE *regexp.Regexp

func isSystemNS(namespace string) bool {
	return SystemRE.MatchString(namespace)
}
