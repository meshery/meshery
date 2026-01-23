package display

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type OutputFormatter[T any] interface {
	Display() error
	WithOutput(io.Writer) OutputFormatter[T]
}

type OutputFormatterFactory[T any] struct{}

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

type YAMLOutputFormatter[T any] struct {
	Data T
	Out  io.Writer
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

func NewYAMLOutputFormatter[T any](data T) OutputFormatter[T] {
	return &YAMLOutputFormatter[T]{
		Data: data,
		Out:  nil,
	}
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
