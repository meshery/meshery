package main

import (
"encoding/json"
"fmt"
)

type Registry struct {
Name string `json:"name"`
Type string `json:"type"`
}

func main() {
fmt.Println("Meshery Backend Learning - JSON Structs")

reg := Registry{
Name: "Docker Hub",
Type: "OCI",
}

jsonData, _ := json.Marshal(reg)
fmt.Printf("JSON Output: %s\n", string(jsonData))

var decoded Registry
json.Unmarshal(jsonData, &decoded)
fmt.Printf("Decoded Name: %s\n", decoded.Name)
}
