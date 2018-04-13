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
			appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubernetes-dashboard",
					Namespace: "kube-system",
				},
			},
			appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prometheus",
					Namespace: "monitoring",
				},
			},
		},
	}
	t.Logf("Should delete all deploys except those in kube-system and monitoring NS")
	SetSkipMeta("", []string{"xxx"})
	clientset := fake.NewSimpleClientset(obj)
	count, err := DeleteDeployments(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	t.Logf("Should delete only deploys in ns1")
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteDeployments(clientset, false, "ns1")
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	t.Logf("Should delete only one deploy, as the other two candidates have the 'created_by' label")
	SetSkipMeta("", nil)
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteDeployments(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// all, as sysNS has been overridden
	t.Logf("Should delete all deploys, as namespaceRE and skipLabels don't match any")
	SetSkipMeta(".*sYsTEM", []string{"xxx"})
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteDeployments(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 5)

}
