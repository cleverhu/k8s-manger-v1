package lib

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

var K8sClient *kubernetes.Clientset

func init() {
	K8sClient = initClient()
}

func initClient() *kubernetes.Clientset {
	config := &rest.Config{Host: "http://47.101.175.193:8009"}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientSet
}
