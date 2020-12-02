package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetLogsFromPod(kubeURL string, namespace string, podName string) string {
	//Example "127.0.0.1:8001/api/v1/namespaces/default/pods/linux/log"
	url := fmt.Sprintf("http://%s/api/v1/namespaces/%s/pods/%s/log", kubeURL, namespace, podName)
	fmt.Println(url)

	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	return string(body)
}
