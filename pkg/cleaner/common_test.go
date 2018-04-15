package cleaner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_SkipMeta(t *testing.T) {
	var c *Common
	c = &Common{
		SkipNamespaceRE: CommonDefaults.SkipNamespaceRE,
		SkipLabels:      CommonDefaults.SkipLabels,
	}
	c.Init(nil)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-system"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-foo"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "foo-system"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "monitoring"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "bob"}), false)
	c = &Common{
		SkipNamespaceRE: "xyz",
		SkipLabels:      CommonDefaults.SkipLabels,
	}
	c.Init(nil)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-system"}), false)
}
