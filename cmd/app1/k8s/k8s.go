package k8s

import (
	"bysj0.3/cmd/app1/meta"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNodeStatus 获取 Kubernetes 集群中所有节点的实时状态
func GetNodeStatus() ([]meta.NodeResourceUsage, error) {
	//获取k8s客户端
	clientset, _ := NewKubernetesClient()
	//获取所有节点列表
	//lientset是使用Kubernetes客户端的配置创建的客户端集合，.CoreV1()表示访问Core API，.Nodes()表示访问Nodes资源，.List()表示获取所有节点的列表。
	nodes, _ := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	//获取Metrics客户端
	metricsClient, _ := NewMetricsClient()
	//遍历节点
	var nodeResourceUsage []meta.NodeResourceUsage
	//nodes.Items是一个节点对象的数组，其中每个节点对象都包含了节点的详细信息
	for _, node := range nodes.Items {
		nodeName := node.GetName()
		//获取当前实时使用的状态
		metrics, err := metricsClient.NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("Error getting metrics for node %s: %v", nodeName, err)
		}
		//当前使用状态
		statuscpuUsages := metrics.Usage.Cpu().MilliValue()
		statusmemoryUsae := metrics.Usage.Memory().Value()
		//全部资源
		allcpuCapacity := node.Status.Capacity.Cpu().MilliValue()
		allmemoryCapacity := node.Status.Capacity.Memory().Value()

		//计算百分比
		cpuUsage := float64(statuscpuUsages) / float64(allcpuCapacity) * 100
		memoryUsage := float64(statusmemoryUsae) / float64(allmemoryCapacity) * 100
		nodeResourceUsage = append(nodeResourceUsage, meta.NodeResourceUsage{Name: nodeName, CpuUsage: cpuUsage, MemoryUsage: memoryUsage})
	}
	return nodeResourceUsage, nil
}

// 获取一下总的节点资源状态
func GeNodeAll() {
	//获取一下总的
}
