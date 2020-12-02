package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

var err error

func init() {
	kubeConfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetLogs(namespaceID, podID string, containerNames []string) (io.ReadCloser, error) {
	for _, container := range containerNames {
		req := clientset.CoreV1().Pods(namespaceID).GetLogs(
			podID,
			&apiv1.PodLogOptions{
				Follow:     true,
				Timestamps: true,
				Container:  container,
			},
		)
		readCloser, err := req.Stream(context.Background())
		if err != nil {
			log.Fatal(err.Error())
		}

		dataChannel := make(chan string, 10)

		reader := bufio.NewReader(readCloser)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				if len(line) > 0 {
					dataChannel <- line
				}
				break
			}
			if err != nil {
				break
			}
			dataChannel <- line
			fmt.Print("Pushed ", <-dataChannel)
		}
	}

	return nil, nil
}
