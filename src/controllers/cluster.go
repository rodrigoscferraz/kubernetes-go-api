package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	gabs "github.com/Jeffail/gabs/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeClient struct {
	Client *kubernetes.Clientset
}

type Info struct {
	NodeName              string `json:"nodeName"`
	NodeCapacityCPU       string `json:"nodeCapacityCPU"`
	NodeCapacityMemory    string `json:"nodeCapacityMemory"`
	NodeCapacityPods      string `json:"nodeCapacityPods"`
	NodeAllocatableCPU    string `json:"nodeAllocatableCPU"`
	NodeAllocatableMemory string `json:"nodeAllocatableMemory"`
	NodeAllocatablePods   string `json:"nodeAllocatablePods"`
}

func parseCapacity(s string) string {
	return s[:len(s)-2]
}

func isNodeMaster(nodeList map[string]string) bool {
	isMaster := false
	for key, _ := range nodeList {
		if key == "node-role.kubernetes.io/master" {
			isMaster = true
		}
	}
	return isMaster
}

func (kube KubeClient) GetClusterInfo(w http.ResponseWriter, r *http.Request) {
	var listinfo []*Info

	clientset := kube.Client
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

	for i, node := range nodes.Items {

		if isNodeMaster(node.GetLabels()) {
			continue
		}

		nodeName := jsonParsed.Search("items", strconv.Itoa(i), "metadata", "name").Data()
		nodeCapacityCPU := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "cpu").Data()
		nodeCapacityMemory := fmt.Sprintf("%v", jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "memory").Data())
		nodeCapacityPods := jsonParsed.Search("items", strconv.Itoa(i), "status", "capacity", "pods").Data()
		nodeAllocatableCPU := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "cpu").Data()
		nodeAllocatableMemory := fmt.Sprintf("%v", jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "memory").Data())
		nodeAllocatablePods := jsonParsed.Search("items", strconv.Itoa(i), "status", "allocatable", "pods").Data()

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

		listinfo = append(listinfo, info)
	}

	w.Header().Add("Content-type", "application/json")
	sbyte, err := json.Marshal(listinfo)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(sbyte))

}
