package main

import (
	"encoding/json"
	"errors"
	"syscall/js"

	"github.com/layer5io/meshkit/cytoscape"
	"gopkg.in/yaml.v3"
)

func PatternFilefromCytoJSON(this js.Value, args []js.Value) any {
	err := validateInput(args)
	if err != nil {
		return (err.Error())
	}

	cytoJSON := args[1].String()
	patternName := args[0].String()

	pattern, err := cytoscape.NewPatternFileFromCytoscapeJSJSON(patternName, []byte(cytoJSON))
	if err != nil {
		return err
	}
	byt, err := json.Marshal(pattern)
	if err != nil {
		return err
	}
	return string(byt)
}

func PolicyResultsToPattern(opaResult map[string]interface{}) (string, error) {
	marshalledResult, err := yaml.Marshal(opaResult)
	if err != nil {
		return "", err
	}
	return string(marshalledResult), nil
}

func validateInput(args []js.Value) error {
	if len(args) < 1 {
		return errors.New("invalid input provided")
	}
	for _, arg := range args {
		if arg.IsNaN() || arg.IsNull() || arg.IsUndefined() {
			return errors.New("invalid input provided")
		}
	}
	return nil
}


func main() {
	done := make(chan struct{})
	js.Global().Set("formatCytoJSONToPatternFile", js.FuncOf(PatternFilefromCytoJSON))
	
	<- done
}