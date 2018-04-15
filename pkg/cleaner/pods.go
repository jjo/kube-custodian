package cleaner

import (
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
)

// DeletePods is main entry point from cmd/delete.go
func (c *Common) DeletePods() (int, error) {
	return c.DeletePodsCond(c.Namespace,
		func(pod *corev1.Pod) bool {
			if c.skipFromMeta(&pod.ObjectMeta) {
				return false
			}
			return true
		})
}

// DeletePodsCond is passed a generic closure to select Pods to delete
func (c *Common) DeletePodsCond(namespace string, filterIn func(*corev1.Pod) bool) (int, error) {

	count := 0
	pods, err := c.clientset.Core().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List pods: %v", err)
		return count, err
	}

	podsArray := []corev1.Pod{}

	for _, pod := range pods.Items {
		log.Debugf("Pod %s.%s ...", pod.Namespace, pod.Name)
		if !filterIn(&pod) {
			continue
		}

		log.Debugf("Pod %q marked for deletion", pod.Name)
		podsArray = append(podsArray, pod)
	}

	for _, pod := range podsArray {
		log.Debugf("Pod %q about to be deleted", pod.Name)

		log.Infof("%s  Deleting Pod %s.%s ...", c.dryRunStr, pod.Namespace, pod.Name)
		if !c.DryRun {
			if err := c.clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete Pod: %v", err)
				continue
			}
			count++
		}
	}
	return count, nil
}
