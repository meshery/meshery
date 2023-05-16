package main

import (
	"fmt"
	"syscall/js"
)

func extractResourceFromSchema(schema string, extractResource []string) (string, error) {

	return "", nil
}

func extractResourceFromSchemaWrapper() js.Func {
	// this function intends to validate the arguments sent from the javascript,
	// like number of arguments, type-checking of arguments, marshalling, unmarshalling
	// etc, if required
	jsFuncValidator := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			result := map[string]string{
				"error": "invalid number of arguments passed to wasm wrapper function",
			}
			return result
		}

		inputSchema := args[0]
		fmt.Println("[Wasm]", inputSchema)
		fmt.Println("[Wasm]", args, "wait \n")
		return inputSchema
	})

	return jsFuncValidator
}

func main() {
	js.Global().Set("extractResourceFromSchemaWasm", extractResourceFromSchemaWrapper())
	<-make(chan bool) // this channel to keep running the wasm server to server this function
}
