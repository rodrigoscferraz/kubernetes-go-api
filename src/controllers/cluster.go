package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orion-api/src/kube"
	"strconv"

	gabs "github.com/Jeffail/gabs/v2"

	//corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/tools/clientcmd"
)

func GetClusterInfo(w http.ResponseWriter, r *http.Request) {

	type Info struct {
		NodeName              string `json:"nodeName"`
		NodeCapacityCPU       string `json:"nodeCapacityCPU"`
		NodeCapacityMemory    string `json:"nodeCapacityMemory"`
		NodeCapacityPods      string `json:"nodeCapacityPods"`
		NodeAllocatableCPU    string `json:"nodeAllocatableCPU"`
		NodeAllocatableMemory string `json:"nodeAllocatableMemory"`
		NodeAllocatablePods   string `json:"nodeAllocatablePods"`
	}

	_, clientset := kube.Kubeconf()
	// mc, err := metrics.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	data, err := json.Marshal(nodes)
	if err != nil {
		panic(err.Error())
	}

	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		panic(err)
	}
	availableNodes := 4
	i := 0
	resp := gabs.New()
	for i < availableNodes {

		nodeName := jsonParsed.Search("items", strconv.Itoa(i), "metadata", "name").Data()
		nodeCapacityCPU := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "cpu").Data()
		nodeCapacityMemory := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "memory").Data()
		nodeCapacityPods := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "pods").Data()
		nodeAllocatableCPU := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "cpu").Data()
		nodeAllocatableMemory := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "memory").Data()
		nodeAllocatablePods := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "pods").Data()

		info := &Info{
			NodeName:              fmt.Sprintf("%v", nodeName),
			NodeCapacityCPU:       fmt.Sprintf("%v", nodeCapacityCPU),
			NodeCapacityMemory:    fmt.Sprintf("%v", nodeCapacityMemory),
			NodeCapacityPods:      fmt.Sprintf("%v", nodeCapacityPods),
			NodeAllocatableCPU:    fmt.Sprintf("%v", nodeAllocatableCPU),
			NodeAllocatableMemory: fmt.Sprintf("%v", nodeAllocatableMemory),
			NodeAllocatablePods:   fmt.Sprintf("%v", nodeAllocatablePods),
		}

		resp.ArrayAppend(info)
		//data, _ := json.Marshal(info)

		i++

	}
	w.Write(resp.Bytes())

}
