package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type BlinktController struct {
	client     kubernetes.Interface
	controller cache.Controller
	blinkts    map[string]string
}

func (c *BlinktController) Run() {
	ticker := time.NewTicker(700 * time.Millisecond)
	for t := range ticker.C {
		fmt.Println("Tick at", t)
		for k, v := range c.blinkts {
			fmt.Printf("key[%s] value[%s]\n", k, v)
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(generate())
			http.Post("http://"+v+":5000/set_color", "application/json; charset=utf-8", b)
		}
	}
}

func generate() LedColor {
	g := rand.Int31n(255)
	b := rand.Int31n(g - 100)
	r := rand.Int31n(g - 50)
	l := rand.Int31n(8)
	return LedColor{Red: r, Blue: b, Green: g, Led: l}
}

func (c *BlinktController) PodDeleted(obj interface{}) {
	pod := obj.(*core_v1.Pod)

	if pod.Labels["name"] == "http-blinkt" {
		fmt.Println("is http blinkt " + pod.Status.PodIP)
		delete(c.blinkts, pod.Status.HostIP)
	}
}

func (c *BlinktController) PodAdded(obj interface{}) {
	pod := obj.(*core_v1.Pod)

	if pod.Labels["name"] == "http-blinkt" {
		fmt.Println("is http blinkt " + pod.Status.PodIP)
		c.blinkts[pod.Status.HostIP] = pod.Status.PodIP
	}
}
