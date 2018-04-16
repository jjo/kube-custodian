package cleaner

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deployUpdater struct {
	deploy appsv1.Deployment
}

func (u *deployUpdater) Update(c *Common) error {
	_, err := c.clientset.AppsV1().Deployments(u.deploy.Namespace).Update(&u.deploy)
	return err
}

func (u *deployUpdater) Delete(c *Common) error {
	return c.clientset.AppsV1().Deployments(u.deploy.Namespace).Delete(u.deploy.Name, &metav1.DeleteOptions{})
}

func (u *deployUpdater) Meta() *metav1.ObjectMeta {
	return &u.deploy.ObjectMeta
}

// DeleteDeployments ...
func (c *Common) DeleteDeployments() (int, error) {
	count := 0
	deploys, err := c.clientset.AppsV1().Deployments(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List deploys: %v", err)
		return count, err
	}

	for _, deploy := range deploys.Items {
		log.Debugf("Deploy %s.%s ...", deploy.Namespace, deploy.Name)
		if c.skipFromMeta(&deploy.ObjectMeta) {
			continue
		}

		log.Debugf("Deploy %s.%s about to be touched ...", deploy.Namespace, deploy.Name)
		count += c.updateState(&deployUpdater{deploy: deploy})
	}
	return count, nil
}
