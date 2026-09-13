// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package printer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Printer struct {
	format string
	out    io.Writer
}

// AddFormatFlag helper-func that helps to add flag in command
func AddFormatFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVarP(target, "output-format", "o", "", "(optional) format to display in [json|yaml|table|string]")
}

// New returns a new Printer. It returns an error only if out is nil.
func New(format string, out io.Writer) (*Printer, error) {
	if out == nil {
		return nil, errors.New("writer is not provided")
	}

	return &Printer{format: strings.ToLower(strings.TrimSpace(format)), out: out}, nil
}

// Print writes data to p.out according to the configured format.
// For "json"/"yaml" it marshals data directly; for "", "table", "string"
// it delegates rendering to humanRender.
func (p *Printer) Print(data any, humanRender func(w io.Writer) error) error {
	switch p.format {
	case "json":
		enc := json.NewEncoder(p.out)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case "yaml":
		enc := yaml.NewEncoder(p.out)
		enc.SetIndent(2)
		defer func() {
			_ = enc.Close()
		}()
		return enc.Encode(data)
	case "", "table", "string":
		if humanRender != nil {
			return humanRender(p.out)
		}
		return nil
	default:
		return utils.ErrFlagsInvalid(fmt.Errorf("invalid value for --output-format '%s': valid values are json, yaml, table, string", p.format))
	}
}
