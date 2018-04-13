package cleaner

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeleteDeployments(t *testing.T) {
	obj := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp3",
					Namespace: "ns3",
				},
			},
			// will be skipped from its .*-system namespace
			appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp4",
					Namespace: "kube-system",
				},
			},
		},
	}
	SetSystemNS("")
	// All deploys except kube-system's
	clientset := fake.NewSimpleClientset(obj)
	count, err := DeleteDeployments(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	// only one, as the 1st two have the required label
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteDeployments(clientset, false, "", []string{"created_by"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// only 1of4 in ns1
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteDeployments(clientset, false, "ns1", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// 4of4 deleted, as sysNS has been overridden
	SetSystemNS(".*sYsTEM")
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteDeployments(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 4)

}
