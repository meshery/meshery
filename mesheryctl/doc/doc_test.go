// Copyright 2023 Layer5, Inc.
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

package main

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	t.Run("Test linkHandler function", func(t *testing.T) {
		expectedLinkHandler := "/main"
		linkHandler := linkHandler("test.md")
		assert.Equal(t, expectedLinkHandler, linkHandler)
	})

	t.Run("Test GenMarkdownTreeCustom function", func(t *testing.T) {
		cmd.AddCommand(&cobra.Command{
			Use: "sub",
		})
		markDownPath := "../../docs/pages/reference/mesheryctl/"
		err := GenMarkdownTreeCustom(cmd, markDownPath, prepender, linkHandler)
		assert.NoError(t, err)
	})

	t.Run("Test GenMarkdownCustom function", func(t *testing.T) {
		cmd.Annotations = map[string]string{
			"link":    "test_link",
			"caption": "test_caption",
		}
		cmd.Example = "test_example"
		manuallyAddedContent, _ := getManuallyAddedContentMap("test.md")
		buf := &bytes.Buffer{}
		err := GenMarkdownCustom(cmd, buf, manuallyAddedContent)
		assert.NoError(t, err)
		assert.NotEmpty(t, buf.String())
		assert.Contains(t, buf.String(), "test_link")
		assert.Contains(t, buf.String(), "test_caption")
		assert.Contains(t, buf.String(), "test_example")
	})

	t.Run("Test HasSeeAlso function", func(t *testing.T) {
		hasSeeAlso := hasSeeAlso(cmd)
		assert.False(t, hasSeeAlso)
	})

	t.Run("Test GenYamlTreeCustom function", func(t *testing.T) {
		cmd.AddCommand(&cobra.Command{
			Use: "sub",
		})
		yamlPath := "../../docs/pages/reference/mesheryctl/"
		err := GenYamlTreeCustom(cmd, yamlPath, prepender, linkHandler)
		assert.NoError(t, err)
	})

	t.Run("Test getManuallyAddedContentMap function", func(t *testing.T) {
		_, err := getManuallyAddedContentMap("test.md")
		assert.NoError(t, err)
	})

	t.Run("Test GenYamlCustom function", func(t *testing.T) {
		cmd.Example = "test_example"
		buf := &bytes.Buffer{}
		err := GenYamlCustom(cmd, buf)
		assert.NoError(t, err)
		assert.NotEmpty(t, buf.String())
		assert.Contains(t, buf.String(), "test_example")
	})
}
