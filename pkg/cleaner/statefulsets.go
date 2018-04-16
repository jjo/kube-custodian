package cleaner

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type stsUpdater struct {
	sts  v1beta1.StatefulSet
	meta metav1.ObjectMeta
}

func (u *stsUpdater) Update(c *Common) error {
	_, err := c.clientset.AppsV1beta1().StatefulSets(u.sts.Namespace).Update(&u.sts)
	return err
}

func (u *stsUpdater) Delete(c *Common) error {
	return c.clientset.AppsV1beta1().StatefulSets(u.sts.Namespace).Delete(u.sts.Name, &metav1.DeleteOptions{})
}

func (u *stsUpdater) Meta() *metav1.ObjectMeta {
	return &u.sts.ObjectMeta
}

// DeleteStatefulSets ...
func (c *Common) DeleteStatefulSets() (int, error) {

	count := 0
	stss, err := c.clientset.AppsV1beta1().StatefulSets(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List StatefulSets: %v", err)
		return count, err
	}

	for _, sts := range stss.Items {
		log.Debugf("StatefulSet %s.%s ...", sts.Namespace, sts.Name)
		if c.skipFromMeta(&sts.ObjectMeta) {
			continue
		}

		log.Debugf("StatefulSet %s.%s about to be touched ...", sts.Namespace, sts.Name)

		count += c.updateState(&stsUpdater{sts: sts})
	}
	return count, nil
}
