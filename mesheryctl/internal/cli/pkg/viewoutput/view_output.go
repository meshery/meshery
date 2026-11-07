// Copyright 2024 Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package viewoutput

import (
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"io"
)

// HandleViewOutput formats and displays the provided data in the specified output format.
// If a save path is provided, it also saves the formatted output to a file.
// It returns an error if formatting, displaying, or saving fails.
func HandleViewOutput[T any](
	outputFormat string,
	data T,
	output io.Writer,
	savePath string,
) error {
	outputFormatterFactory := display.OutputFormatterFactory[T]{}

	outputFormatter, err := outputFormatterFactory.New(outputFormat, data)
	if err != nil {
		return err
	}

	outputFormatter = outputFormatter.WithOutput(output)

	if err := outputFormatter.Display(); err != nil {
		return err
	}

	//if the savePath is empty only the output is displayed, and isnt written to a file.4
	if savePath == "" {
		return nil
	}

	outputFormatterSaverFactory := display.OutputFormatterSaverFactory[T]{}

	outputFormatterSaver, err := outputFormatterSaverFactory.New(
		outputFormat,
		outputFormatter,
	)
	if err != nil {
		return err
	}

	outputFormatterSaver = outputFormatterSaver.WithFilePath(savePath)

	return outputFormatterSaver.Save()
}
