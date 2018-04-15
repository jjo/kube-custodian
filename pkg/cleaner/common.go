package cleaner

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/jjo/kube-custodian/pkg/utils"
)

type Common struct {
	DryRun          bool
	Namespace       string
	SkipNamespaceRE string
	SkipLabels      []string

	clientset kubernetes.Interface

	skipNamespaceRegexp *regexp.Regexp
	timeStamp           int64
	dryRunStr           string
}

const (
	kubeCustodianAnnotationTime = "kube-custodian.bitnami.com/expiration-time"
)

var CommonDefaults = &Common{
	SkipNamespaceRE: "kube-.*|.*(-system|monitoring|logging|ingress)",
	SkipLabels:      []string{"created_by"},
}

func (c *Common) Init(clientset kubernetes.Interface) {
	c.skipNamespaceRegexp = regexp.MustCompile(c.SkipNamespaceRE)
	c.timeStamp = time.Now().Unix()
	c.dryRunStr = map[bool]string{true: "[dry-run]", false: ""}[c.DryRun]
	c.clientset = clientset
}

func (c Common) Run() {
	c.DeleteDeployments()
	c.DeleteStatefulSets()
	c.DeleteJobs()
	c.DeletePods()
}

func (c *Common) skipFromMeta(meta *metav1.ObjectMeta) bool {
	skipIt := false
	switch {
	case c.skipNamespaceRegexp.MatchString(meta.Namespace):
		log.Debugf("%s.%s skipped from meta.Namespace", meta.Name, meta.Namespace)
		skipIt = true
	case utils.LabelsSubSet(meta.Labels, c.SkipLabels):
		log.Debugf("%s.%s skipped from meta.Labels", meta.Name, meta.Labels)
		skipIt = true
	}
	return skipIt
}

func (c *Common) updateState(Update func() error, Delete func() error, objMeta *metav1.ObjectMeta) int {

	fqName := fmt.Sprintf("%s.%s", objMeta.Name, objMeta.Namespace)
	cnt := 0
	annotations := objMeta.GetAnnotations()
	if valueStr, found := annotations[kubeCustodianAnnotationTime]; found {
		log.Debugf("%s already has annotation %s: %s",
			fqName, kubeCustodianAnnotationTime, valueStr)
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			log.Errorf("%s: failed to convert %s to integer", fqName, valueStr)
		}
		expiredSecs := c.timeStamp - (int64(value) + 60)
		if expiredSecs > 0 {
			log.Debugf("%s%s TTL expired %d seconds ago, deleting",
				c.dryRunStr, fqName, expiredSecs)
			/*
			   if !dryRun {
			       if err := Delete(); err != nil {
			           log.Errorf("failed to delete %s with error: %v", fqName, err)
			       } else {
			           cnt++
			       }
			   }
			*/
		}
	} else {
		timeStampStr := fmt.Sprintf("%d", c.timeStamp)
		log.Debugf("[%s]%s updating annotation %s: %s",
			c.dryRunStr, fqName, kubeCustodianAnnotationTime, timeStampStr)
		if !c.DryRun {
			metav1.SetMetaDataAnnotation(objMeta,
				kubeCustodianAnnotationTime, timeStampStr)
			if err := Update(); err != nil {
				log.Errorf("failed to update %s with error: %v", fqName, err)
			} else {
				cnt++
			}
		}
	}
	return cnt
}
