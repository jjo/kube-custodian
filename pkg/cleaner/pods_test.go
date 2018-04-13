package cleaner

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeletePodsCond(t *testing.T) {
	obj := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod3",
					Namespace: "ns3",
				},
			},
			// sysNS will be skipped:
			corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubernetes-dashboard-deadbeefed-quack",
					Namespace: "kube-system",
				},
			},
			corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prometheus",
					Namespace: "monitoring",
				},
			},
		},
	}
	// The only condition *always* checked is skipping "system" Pods
	SetSystemNS("")
	clientset := fake.NewSimpleClientset(obj)
	count, err := DeletePodsCond(clientset, false, "", func(pod *corev1.Pod) bool {
		return true
	})
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	// Select a single pod
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeletePodsCond(clientset, false, "", func(pod *corev1.Pod) bool {
		return pod.Labels["created_by"] == "foo"
	})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// Change system NS to whatever, we should get them all
	SetSystemNS("sYsTEM")
	clientset = fake.NewSimpleClientset(obj)
	// The only condition always present is skipping "system" Pods
	count, err = DeletePodsCond(clientset, false, "", func(pod *corev1.Pod) bool {
		return true
	})
	assertEqual(t, err, nil)
	assertEqual(t, count, 5)

}
