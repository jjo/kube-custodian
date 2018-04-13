package cleaner

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_DeleteJobs(t *testing.T) {
	obj := &batchv1.JobList{
		Items: []batchv1.Job{
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job3",
					Namespace: "ns3",
				},
				Status: batchv1.JobStatus{
					Succeeded: 0,
				},
			},
			// will be skipped from its .*-system namespace
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job4",
					Namespace: "kube-system",
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
		},
	}
	SetSystemNS("")
	// 2of4 non system Jobs Succeeded
	clientset := fake.NewSimpleClientset(obj)
	count, err := DeleteJobs(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 2)

	// no one, as the 1st two now have the required label
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteJobs(clientset, false, "", []string{"created_by"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 0)

	// only one in ns1
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteJobs(clientset, false, "ns1", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	// 3of4 Jobs Succeeded, as sysNS has been overridden
	SetSystemNS(".*sYsTEM")
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteJobs(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

}
