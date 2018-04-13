package cleaner

import (
	"regexp"
)

// SystemNS has default "system" namespaces regexp
const (
	SkipNSRe = "kube-.*|.*(-system|monitoring|logging|ingress)"
)

var skipNSRegexp *regexp.Regexp

func init() {
	SetSkipNSRe("")
}

func skipNamespace(namespace string) bool {
	return skipNSRegexp.MatchString(namespace)
}

// SetSkipNSRe is used from cmd/delete.go flags
func SetSkipNSRe(namespaceRe string) {
	if namespaceRe != "" {
		skipNSRegexp = regexp.MustCompile(namespaceRe)
	} else {
		skipNSRegexp = regexp.MustCompile(SkipNSRe)
	}
}
