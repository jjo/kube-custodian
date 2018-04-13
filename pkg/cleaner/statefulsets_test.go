package cleaner

import (
	"testing"

	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeleteStatefulSets(t *testing.T) {
	obj := &v1beta1.StatefulSetList{
		Items: []v1beta1.StatefulSet{
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts3",
					Namespace: "ns3",
				},
			},
			// sysNS will be skipped:
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fatDb",
					Namespace: "kube-system",
				},
			},
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prometheus",
					Namespace: "monitoring",
				},
			},
		},
	}
	SetSystemNS("")
	// All deploys except kube-system's
	clientset := fake.NewSimpleClientset(obj)
	count, err := DeleteStatefulSets(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	// only one, as the 1st two have the required label
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteStatefulSets(clientset, false, "", []string{"created_by"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// only one in ns1
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteStatefulSets(clientset, false, "ns1", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// all, as sysNS has been overridden
	SetSystemNS(".*sYsTEM")
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteStatefulSets(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 5)

}
