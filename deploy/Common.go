package deploy

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
)

func GetImages(dep v1.Deployment) string {
	images := dep.Spec.Template.Spec.Containers[0].Image
	if imgLen := len(dep.Spec.Template.Spec.Containers); imgLen > 1 {
		images += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}

	return images
}
