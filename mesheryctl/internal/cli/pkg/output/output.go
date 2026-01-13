package output

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type OutputFormatter interface {
	Display() error
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

func NewJSONOutputFormatter[T any](data T) *JSONOutputFormatter[T] {
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

func (j *JSONOutputFormatter[T]) WithOutput(out io.Writer) *JSONOutputFormatter[T] {
	j.Out = out
	return j
}

func (j *JSONOutputFormatter[T]) WithEncoderSettings(settings JsonEncoderSettings) *JSONOutputFormatter[T] {
	j.EncoderSettings = settings
	return j
}

func (y *YAMLOutputFormatter[T]) WithOutput(out io.Writer) *YAMLOutputFormatter[T] {
	y.Out = out
	return y
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
