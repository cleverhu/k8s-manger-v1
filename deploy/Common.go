package deploy

import (
	"context"
	"fmt"
	"k8s-manger-v1/lib"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func GetImages(dep appsv1.Deployment) string {
	images := dep.Spec.Template.Spec.Containers[0].Image
	if imgLen := len(dep.Spec.Template.Spec.Containers); imgLen > 1 {
		images += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return images
}

func GetLabels(m map[string]string) string {
	labels := ""
	for k, v := range m {
		if labels != "" {
			labels += ","
		}
		labels += fmt.Sprintf("%s=%s", k, v)
	}

	return labels
}

func GetPodLabelsByDeployment(deploy *appsv1.Deployment) string {
	selector, _ := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	rs, _ := lib.K8sClient.AppsV1().ReplicaSets(deploy.Namespace).List(context.Background(), metav1.ListOptions{LabelSelector: selector.String()})
	podLabel := ""
	for _, item := range rs.Items {
		if item.Annotations["deployment.kubernetes.io/revision"] != deploy.Annotations["deployment.kubernetes.io/revision"] {
			continue
		}
		for _, v := range item.OwnerReferences {
			if v.Name == deploy.Name {
				s, _ := metav1.LabelSelectorAsSelector(item.Spec.Selector)
				podLabel = s.String()
			}
		}
	}

	return podLabel
}

func GetImagesByPod(containers []corev1.Container) string {
	images := containers[0].Image
	if imgLen := len(containers); imgLen > 1 {
		images += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return images
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
