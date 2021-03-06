package main

import (
	"flag"
	"fmt"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var kubecfg string

func main() {

	var config *rest.Config

	flag.StringVar(&kubecfg, "kubeconfig", "", "Path to kubeconfig")
	flag.Parse()
	fmt.Println(kubecfg)

	if kubecfg == "" {
		config, _ = rest.InClusterConfig()
	} else {
		config, _ = clientcmd.BuildConfigFromFlags("", kubecfg)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	bc := BlinktController{
		client:  client,
		blinkts: make(map[string]string),
	}

	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&core_v1.Pod{},
		0, //Skip resync
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(new interface{}) { bc.PodAdded(new) },
			UpdateFunc: func(old, new interface{}) { bc.PodDeleted(old); bc.PodAdded(new) },
			DeleteFunc: func(new interface{}) { bc.PodDeleted(new) },
		},
	)
	bc.controller = controller

	go bc.Run()

	controller.Run(wait.NeverStop)
}
