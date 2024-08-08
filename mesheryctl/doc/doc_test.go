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

package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPrepender(t *testing.T) {
	t.Run("Simple file", func(t *testing.T) {
		expected := `---
layout: default
title: test
permalink: reference/test/main
redirect_from: reference/test/main/
type: reference
display-title: "false"
language: en
command: test
subcommand: nil
---

`
		assert.Equal(t, expected, prepender("test.md"))
	})

	t.Run("File with subcommands", func(t *testing.T) {
		expected := `---
layout: default
title: test-sub
permalink: reference/test/sub
redirect_from: reference/test/sub/
type: reference
display-title: "false"
language: en
command: sub
subcommand: nil
---

`
		assert.Equal(t, expected, prepender("test-sub.md"))

		expected = `---
layout: default
title: test-sub-sub
permalink: reference/test/sub/sub
redirect_from: reference/test/sub/sub/
type: reference
display-title: "false"
language: en
command: sub
subcommand: sub
---

`
		assert.Equal(t, expected, prepender("test-sub-sub.md"))

		expected = `---
layout: default
title: test-sub-sub-sub
permalink: reference/test/sub/sub/sub
redirect_from: reference/test/sub/sub/sub/
type: reference
display-title: "false"
language: en
command: sub
subcommand: sub
---

`
		assert.Equal(t, expected, prepender("test-sub-sub-sub.md"))
	})
}

func TestDoc(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	t.Run("Test linkHandler function", func(t *testing.T) {
		assert.Equal(t, "/main", linkHandler("test.md"))
		assert.Equal(t, "sub", linkHandler("test-sub-sub.md"))
		assert.Equal(t, "sub/sub", linkHandler("test-sub-sub-sub.md"))
		assert.Equal(t, "sub", linkHandler("test-sub-sub-sub-sub.md"))
	})

	t.Run("Test GenMarkdownTreeCustom function", func(t *testing.T) {
		cmd.AddCommand(&cobra.Command{
			Use: "sub",
		})
		markDownPath := "../../docs/pages/reference/mesheryctl/"
		err := GenMarkdownTreeCustom(cmd, markDownPath, prepender, linkHandler)
		assert.NoError(t, err)
	})

	t.Run("Test HasSeeAlso function", func(t *testing.T) {

		assert.False(t, hasSeeAlso(cmd))

		parentCmd := &cobra.Command{Use: "parent"}
		childCmd := &cobra.Command{Use: "child"}
		parentCmd.AddCommand(childCmd)
		assert.True(t, hasSeeAlso(childCmd))
	})

	t.Run("Test GenMarkdownCustom with initial setup", func(t *testing.T) {
		cmd.Annotations = map[string]string{
			"link":    "test_link",
			"caption": "test_caption",
		}
		cmd.Example = "test_example"
		cmd.Long = "test_long"
		manuallyAddedContent, _ := getManuallyAddedContentMap("test.md")
		buf := &bytes.Buffer{}
		err := GenMarkdownCustom(cmd, buf, manuallyAddedContent)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "test_link")
		assert.Contains(t, output, "test_caption")
		assert.Contains(t, output, "test_example")
	})

	t.Run("Test GenMarkdownCustom with parent command and manually added content", func(t *testing.T) {

		parentCmd := &cobra.Command{Use: "parent"}
		parentCmd.AddCommand(cmd)
		cmd.Long = "Find test_long"
		cmd.Run = func(cmd *cobra.Command, args []string) {}
		cmd.Example = "// test_example"
		file, _ := os.CreateTemp("", "test.md")
		defer os.Remove(file.Name())
		_, err := file.WriteString("{% include example.md %}")
		assert.NoError(t, err)
		file.Close()
		manuallyAddedContent, _ := getManuallyAddedContentMap(file.Name())
		buf := &bytes.Buffer{}
		err = GenMarkdownCustom(cmd, buf, manuallyAddedContent)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "test_link")
		assert.Contains(t, output, "codeblock-pre")
		assert.Contains(t, output, "test_caption")
		assert.Contains(t, output, "test_example")
		assert.Contains(t, output, "See Also")
		assert.Contains(t, output, "preserving-manually-added-documentation")
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

		file, err := os.CreateTemp("", "test.md")
		assert.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString("{% include example.md %}")
		assert.NoError(t, err)
		file.Close()

		contentMap, err := getManuallyAddedContentMap(file.Name())
		assert.NoError(t, err)
		assert.Contains(t, contentMap, 0)
		assert.Equal(t, "example.md", contentMap[0])
	})

	t.Run("Test GenYamlCustom function", func(t *testing.T) {
		cmd.Example = "test_example"
		buf := &bytes.Buffer{}
		err := GenYamlCustom(cmd, buf)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "test_example")
	})

	t.Run("Test printOptions function", func(t *testing.T) {
		cmdWithFlags := &cobra.Command{Use: "testWithFlags"}
		cmdWithFlags.Flags().String("flag1", "default1", "description1")
		cmdWithFlags.Flags().String("flag2", "default2", "description2")

		buf := &bytes.Buffer{}
		err := printOptions(buf, cmdWithFlags)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "## Options")
		assert.Contains(t, output, "--flag1")
		assert.Contains(t, output, "--flag2")

		parentCmd := &cobra.Command{Use: "parent"}
		parentCmd.PersistentFlags().String("parentFlag", "defaultParent", "parent description")
		parentCmd.AddCommand(cmdWithFlags)

		buf = &bytes.Buffer{}
		err = printOptions(buf, cmdWithFlags)
		assert.NoError(t, err)
		output = buf.String()
		assert.Contains(t, output, "## Options inherited from parent commands")
		assert.Contains(t, output, "--parentFlag")
	})
}
