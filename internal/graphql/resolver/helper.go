package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// getResource can get any namespaced kubernetes resource based on the query
// given to the function
//
// Example query: If label version: v2 and name: app is a filter
// criterion then query should look like ["metadata.labels.version = v2", "metadata.name = app"]
//
// If no query is provided then all of the matching resources will be returned
//
// If no namespace is given then resource across the cluster will be listed
func getResource(
	ctx context.Context,
	client dynamic.Interface,
	group string,
	version string,
	resource string,
	namespace string,
	queries []string,
) ([]unstructured.Unstructured, error) {
	res, err := client.Resource(schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}).Namespace(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := []unstructured.Unstructured{}

	for _, i := range res.Items {
		mp := flattenJSON(i.Object)

		if isValidResource(queries, mp) {
			items = append(items, i)
		}
	}

	return items, nil
}

func watchResource(
	ctx context.Context,
	client dynamic.Interface,
	group string,
	version string,
	resource string,
	namespace string,
	queries []string,
) (<-chan watch.Event, error) {
	mc := make(chan watch.Event, 2)

	wi, err := client.Resource(schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}).Namespace(namespace).Watch(ctx, v1.ListOptions{})
	if err != nil {
		return mc, err
	}

	go func() {
		ch := wi.ResultChan()

		for {
			select {
			case ev := <-ch:
				// Serialize to map
				byt, _ := json.Marshal(ev.Object)

				mp := map[string]interface{}{}
				_ = json.Unmarshal(byt, &mp)

				if isValidResource(queries, flattenJSON(mp)) {
					mc <- ev
				}
			}
		}
	}()

	return mc, nil
}

func isValidResource(queries []string, obj map[string]interface{}) bool {
	for _, query := range queries {
		parsedQuery := strings.Split(query, " = ")

		match, ok := obj[parsedQuery[0]]
		if !ok || (parsedQuery[1] != "*" && match != parsedQuery[1]) {
			return false
		}
	}

	return true
}

func flattenJSON(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})

	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := flattenJSON(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		case []interface{}:
			newChild := map[string]interface{}{}
			for i, v := range child {
				newChild[fmt.Sprint(i)] = v
			}
			nm := flattenJSON(newChild)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = v
		}
	}

	return o
}
