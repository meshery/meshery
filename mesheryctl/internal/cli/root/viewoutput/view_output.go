package viewoutput

import (
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"io"
)

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
