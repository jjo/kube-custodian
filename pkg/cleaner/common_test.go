package cleaner

import (
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_SysNS(t *testing.T) {
	SetSkipNSRe("")
	assertEqual(t, skipNamespace("kube-system"), true)
	assertEqual(t, skipNamespace("kube-foo"), true)
	assertEqual(t, skipNamespace("metallb-system"), true)
	assertEqual(t, skipNamespace("monitoring"), true)
	assertEqual(t, skipNamespace("bob"), false)
	SetSkipNSRe("xyz")
	assertEqual(t, skipNamespace("kube-system"), false)
}
