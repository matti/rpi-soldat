package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// BetterNode ...
type BetterNode struct {
	Name string
}

func getBetterNodes() []BetterNode {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Panicln("clientcmd", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("clientset", err)
	}

	labelSelector := fmt.Sprintf("kubernetes.io/arch=%s", "arm64")
	nodeList, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Panicln("node list", err)
	}

	var nodes []BetterNode
	for _, node := range nodeList.Items {
		nodes = append(nodes, BetterNode{
			Name: node.Name,
		})
	}
	return nodes
}

func getTempString() string {
	bytes, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		bytes = []byte("-1000\n")
	}

	s := strings.TrimSuffix(string(bytes), "\n")
	i, _ := strconv.Atoi(s)
	f := float64(i) / float64(1000)
	return fmt.Sprintf("%.2f", f)
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./views/*")
	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"hostname":    os.Getenv("HOSTNAME"),
			"temperature": getTempString(),
			"nodes":       getBetterNodes(),
		})
	})

	log.Println("listen :8080")
	r.Run(":8080")
}
