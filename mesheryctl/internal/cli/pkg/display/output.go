package display

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"gopkg.in/yaml.v3"
)

type OutputFormatter[T any] interface {
	Display() error
	WithOutput(io.Writer) OutputFormatter[T]
}

type OutputFormatterSaver[T any] interface {
	Save() error
	WithFilePath(string) OutputFormatterSaver[T]
}

type OutputFormatterFactory[T any] struct{}

type OutputFormatterSaverFactory[T any] struct {
	OutputFormatter OutputFormatter[T]
}

func (o *OutputFormatterFactory[T]) New(format string, data T) (OutputFormatter[T], error) {
	switch format {
	case "json":
		return NewJSONOutputFormatter(data), nil
	case "yaml":
		return NewYAMLOutputFormatter(data), nil
	default:
		return nil, ErrUnsupportedFormat(format)
	}
}

func (o *OutputFormatterSaverFactory[T]) New(format string, outputFormatter OutputFormatter[T]) (OutputFormatterSaver[T], error) {
	switch format {
	case "json":
		jsonFormatter, ok := outputFormatter.(*JSONOutputFormatter[T])
		if !ok {
			return nil, ErrUnsupportedFormat(format)
		}
		return NewJSONOutputFormatterSaver(*jsonFormatter), nil
	case "yaml":
		yamlFormatter, ok := outputFormatter.(*YAMLOutputFormatter[T])
		if !ok {
			return nil, ErrUnsupportedFormat(format)
		}
		return NewYAMLOutputFormatterSaver(*yamlFormatter), nil
	default:
		return nil, ErrUnsupportedFormat(format)
	}
}

type JsonEncoderSettings struct {
	SetEscapeHTML bool
	IndentPrefix  string
	IndentValue   string
}

type JSONOutputFormatter[T any] struct {
	Data            T
	EncoderSettings JsonEncoderSettings
	Out             io.Writer
}

type JSONOutputFormatterSaver[T any] struct {
	OutputFormatter JSONOutputFormatter[T]
	FilePath        string
}

type YAMLOutputFormatter[T any] struct {
	Data T
	Out  io.Writer
}

type YAMLOutputFormatterSaver[T any] struct {
	OutputFormatter YAMLOutputFormatter[T]
	FilePath        string
}

func NewJSONOutputFormatter[T any](data T) OutputFormatter[T] {
	return &JSONOutputFormatter[T]{
		Data: data,
		EncoderSettings: JsonEncoderSettings{
			SetEscapeHTML: false,
			IndentPrefix:  "",
			IndentValue:   "  ",
		},
		Out: nil,
	}
}

func (j *JSONOutputFormatterSaver[T]) WithFilePath(filePath string) OutputFormatterSaver[T] {
	j.FilePath = filePath
	return j
}

func NewJSONOutputFormatterSaver[T any](outputFormatter JSONOutputFormatter[T]) OutputFormatterSaver[T] {
	return &JSONOutputFormatterSaver[T]{
		OutputFormatter: outputFormatter,
		FilePath:        "",
	}
}

func (j *JSONOutputFormatter[T]) WithOutput(out io.Writer) OutputFormatter[T] {
	j.Out = out
	return j
}

func (j *JSONOutputFormatter[T]) WithEncoderSettings(settings JsonEncoderSettings) OutputFormatter[T] {
	j.EncoderSettings = settings
	return j
}

func (j *JSONOutputFormatter[T]) Display() error {
	// Default to stdout if no output is provided
	if j.Out == nil {
		j.Out = os.Stdout
	}
	encoder := json.NewEncoder(j.Out)

	// Configure the JSON encoder settings.
	encoder.SetEscapeHTML(j.EncoderSettings.SetEscapeHTML)
	encoder.SetIndent(j.EncoderSettings.IndentPrefix, j.EncoderSettings.IndentValue)

	if err := encoder.Encode(j.Data); err != nil {
		return ErrEncodingData(err, "json")
	}

	return nil
}

func (j *JSONOutputFormatterSaver[T]) Save() error {
	if j.FilePath == "" {
		return ErrOutputFileNotSpecified()
	}

	var output []byte
	var err error

	fmt.Println()
	if output, err = json.MarshalIndent(j.OutputFormatter.Data, j.OutputFormatter.EncoderSettings.IndentPrefix,
		j.OutputFormatter.EncoderSettings.IndentValue); err != nil {
		return utils.ErrMarshalIndent(errors.Wrap(err, "failed to format output in JSON"))
	}

	err = writeFile(j.FilePath, output)
	if err != nil {
		return utils.ErrCreateFile(j.FilePath, errors.Wrap(err, "failed to save output as JSON file"))
	}

	utils.Log.Info("Connection saved to file: ", j.FilePath)
	return nil
}

func NewYAMLOutputFormatter[T any](data T) OutputFormatter[T] {
	return &YAMLOutputFormatter[T]{
		Data: data,
		Out:  nil,
	}
}

func NewYAMLOutputFormatterSaver[T any](outputFormatter YAMLOutputFormatter[T]) OutputFormatterSaver[T] {
	return &YAMLOutputFormatterSaver[T]{
		OutputFormatter: outputFormatter,
		FilePath:        "",
	}
}

func (y *YAMLOutputFormatterSaver[T]) WithFilePath(filePath string) OutputFormatterSaver[T] {
	y.FilePath = filePath
	return y
}

func (y *YAMLOutputFormatterSaver[T]) Save() error {
	if y.FilePath == "" {
		return ErrOutputFileNotSpecified()
	}

	var output []byte
	var err error

	fmt.Println()
	if output, err = yaml.Marshal(y.OutputFormatter.Data); err != nil {
		return utils.ErrMarshal(errors.Wrap(err, "failed to format output in yaml"))
	}

	err = writeFile(y.FilePath, output)
	if err != nil {
		return utils.ErrCreateFile(y.FilePath, errors.Wrap(err, "failed to save output as yaml file"))
	}

	utils.Log.Info("Connection saved to file: ", y.FilePath)
	return nil
}

func writeFile(filePath string, data []byte) error {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return utils.ErrCreateFile(filePath, errors.Wrap(err, "failed to write data to file"))
	}

	return nil
}

func (y *YAMLOutputFormatter[T]) WithOutput(out io.Writer) OutputFormatter[T] {
	y.Out = out
	return y
}

func (y *YAMLOutputFormatter[T]) Display() error {
	// Default to stdout if no output is provided
	if y.Out == nil {
		y.Out = os.Stdout
	}

	encoder := yaml.NewEncoder(y.Out)

	if err := encoder.Encode(y.Data); err != nil {
		return ErrEncodingData(err, "yaml")
	}

	return nil
}
