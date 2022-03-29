package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orion-api/src/kube"
	"strconv"

	gabs "github.com/Jeffail/gabs/v2"

	//corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/tools/clientcmd"
)

type Info struct {
	NodeName              string `json:"nodeName"`
	NodeCapacityCPU       string `json:"nodeCapacityCPU"`
	NodeCapacityMemory    string `json:"nodeCapacityMemory"`
	NodeCapacityPods      string `json:"nodeCapacityPods"`
	NodeAllocatableCPU    string `json:"nodeAllocatableCPU"`
	NodeAllocatableMemory string `json:"nodeAllocatableMemory"`
	NodeAllocatablePods   string `json:"nodeAllocatablePods"`
}

var listinfo []*Info

func parseCapacity(s string) string {
	return s[:len(s)-2]
}

func GetClusterInfo(w http.ResponseWriter, r *http.Request) {

	_, clientset := kube.Kubeconf()
	// mc, err := metrics.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		// LabelSelector: "node-role.kubernetes.io/master=true",
	})
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

	// availableNodes := len(nodes.Items)
	// i := 0
	// resp := gabs.New()

	

	for i, node := range nodes.Items {

		nodeName := jsonParsed.Search("items", strconv.Itoa(i), "metadata", "name").Data()
		nodeCapacityCPU := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "cpu").Data()
		nodeCapacityMemory := fmt.Sprintf("%v", jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "memory").Data())
		nodeCapacityPods := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "pods").Data()
		nodeAllocatableCPU := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "cpu").Data()
		nodeAllocatableMemory := fmt.Sprintf("%v", jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "memory").Data())
		nodeAllocatablePods := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "pods").Data()

		// sb, _ := json.Marshal(nodeRole)
		// fmt.Println(string(sb))
		// node-role.kubernetes.io/master
		isMaster := false

		for key, _ := range node.GetLabels() {
			if key == "node-role.kubernetes.io/master" {
				isMaster = true
			}
		}

		if isMaster {
			continue
		}

		nodeCapacityMemory = parseCapacity(nodeCapacityMemory)
		nodeAllocatableMemory = parseCapacity(nodeAllocatableMemory)

		info := &Info{
			NodeName:              fmt.Sprintf("%v", nodeName),
			NodeCapacityCPU:       fmt.Sprintf("%v", nodeCapacityCPU),
			NodeCapacityMemory:    nodeCapacityMemory,
			NodeCapacityPods:      fmt.Sprintf("%v", nodeCapacityPods),
			NodeAllocatableCPU:    fmt.Sprintf("%v", nodeAllocatableCPU),
			NodeAllocatableMemory: nodeAllocatableMemory,
			NodeAllocatablePods:   fmt.Sprintf("%v", nodeAllocatablePods),
		}

		// resp.ArrayAppend(info)
		listinfo = append(listinfo, info)
		//data, _ := json.Marshal(info)
	}

	w.Header().Add("Content-type", "application/json")
	sbyte, err := json.Marshal(listinfo)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(sbyte))

}
