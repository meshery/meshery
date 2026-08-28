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
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPrepender(t *testing.T) {

	t.Run("Root mesheryctl index", func(t *testing.T) {
		expected := `---
title: mesheryctl
display_title: false
command: mesheryctl
subcommand: nil
---

`
		filename := "docs/content/en/reference/references/mesheryctl/_index.md"
		assert.Equal(t, expected, prepender(filename))
	})

	t.Run("First-level command (_index.md)", func(t *testing.T) {
		expected := `---
title: mesheryctl-adapter
display_title: false
command: adapter
subcommand: nil
---

`
		filename := "docs/content/en/reference/references/mesheryctl/adapter/_index.md"
		assert.Equal(t, expected, prepender(filename))
	})

	t.Run("Leaf command (single subcommand)", func(t *testing.T) {
		expected := `---
title: mesheryctl-adapter-deploy
display_title: false
command: adapter
subcommand: deploy
---

`
		filename := "docs/content/en/reference/references/mesheryctl/adapter/deploy.md"
		assert.Equal(t, expected, prepender(filename))
	})

	t.Run("Nested command (two levels deep)", func(t *testing.T) {
		expected := `---
title: mesheryctl-exp-relationship-generate
display_title: false
command: exp
subcommand: relationship
---

`
		filename := "docs/content/en/reference/references/mesheryctl/exp/relationship/generate.md"
		assert.Equal(t, expected, prepender(filename))
	})
}

func TestDoc(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	t.Run("Test linkHandler function (directory structure)", func(t *testing.T) {
		assert.Equal(t,
			"/reference/references/mesheryctl",
			linkHandler("docs/content/en/reference/references/mesheryctl/_index.md"),
		)

		assert.Equal(t,
			"/reference/references/mesheryctl/adapter",
			linkHandler("docs/content/en/reference/references/mesheryctl/adapter/_index.md"),
		)

		assert.Equal(t,
			"/reference/references/mesheryctl/adapter/deploy",
			linkHandler("docs/content/en/reference/references/mesheryctl/adapter/deploy.md"),
		)

		assert.Equal(t,
			"/reference/references/mesheryctl/exp/relationship/generate",
			linkHandler("docs/content/en/reference/references/mesheryctl/exp/relationship/generate.md"),
		)
	})

	t.Run("Test GenMarkdownTreeCustom function", func(t *testing.T) {
		cmd.AddCommand(&cobra.Command{
			Use: "sub",
		})
		markDownPath := t.TempDir()
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
		defer func() { _ = os.Remove(file.Name()) }()
		_, err := file.WriteString("{{< example >}}")
		assert.NoError(t, err)
		_ = file.Close()
		manuallyAddedContent, _ := getManuallyAddedContentMap(file.Name())
		buf := &bytes.Buffer{}
		err = GenMarkdownCustom(cmd, buf, manuallyAddedContent)
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "test_link")
		assert.Contains(t, output, "codeblock-pre")
		// The copy-to-clipboard button in docs/static/js/main.js binds to '.clipboardjs'.
		assert.Contains(t, output, "<code class='clipboardjs'>")
		assert.Contains(t, output, "test_caption")
		assert.Contains(t, output, "test_example")
		assert.Contains(t, output, "See Also")
		assert.Contains(t, output, "preserving-manually-added-documentation")
	})
	t.Run("Test getManuallyAddedContentMap function", func(t *testing.T) {
		_, err := getManuallyAddedContentMap("test.md")
		assert.NoError(t, err)

		file, err := os.CreateTemp("", "test.md")
		assert.NoError(t, err)
		defer func() { _ = os.Remove(file.Name()) }()

		_, err = file.WriteString("{{< example >}}")
		assert.NoError(t, err)
		_ = file.Close()

		contentMap, err := getManuallyAddedContentMap(file.Name())
		assert.NoError(t, err)
		assert.Contains(t, contentMap, 0)
		assert.Equal(t, "example", contentMap[0])
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

func TestWriteCommandBlock(t *testing.T) {
	t.Run("emits the clipboardjs hook the docs copy button binds to", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writeCommandBlock(buf, "mesheryctl system start")

		expected := "<pre class='codeblock-pre'>\n" +
			"<div class='codeblock'>\n" +
			"<code class='clipboardjs'>mesheryctl system start</code>\n" +
			"</div>\n</pre> \n\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("escapes markup so placeholders survive the HTML parser", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writeCommandBlock(buf, "mesheryctl relationship search [--kind <kind>] & more")

		out := buf.String()
		// <kind> would otherwise be parsed as a tag and vanish from the rendered page.
		assert.Contains(t, out, "[--kind &lt;kind&gt;] &amp; more")
		assert.NotContains(t, out, "<kind>")
	})

	t.Run("does not double-escape an ampersand it just inserted", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writeCommandBlock(buf, "a < b")

		assert.Contains(t, buf.String(), "a &lt; b")
		assert.NotContains(t, buf.String(), "&amp;lt;")
	})
}

func TestGenMarkdownCustomExamples(t *testing.T) {
	t.Run("every example command block is copyable", func(t *testing.T) {
		cmd := &cobra.Command{Use: "start", Short: "short"}
		cmd.Run = func(cmd *cobra.Command, args []string) {}
		cmd.Example = "// Start Meshery\nmesheryctl system start\n// Start in debug mode\nmesheryctl system start --debug"

		buf := &bytes.Buffer{}
		err := GenMarkdownCustom(cmd, buf, map[int]string{})
		assert.NoError(t, err)
		output := buf.String()

		// Two examples plus the synopsis usage line are all copyable.
		assert.Equal(t, 3, strings.Count(output, "<code class='clipboardjs'>"))
		assert.Contains(t, output, "<code class='clipboardjs'>mesheryctl system start</code>")
		assert.Contains(t, output, "<code class='clipboardjs'>mesheryctl system start --debug</code>")
		// Description lines stay prose, not code blocks.
		assert.Contains(t, output, "Start Meshery\n")
	})

	t.Run("whitespace-only example lines produce no empty block", func(t *testing.T) {
		cmd := &cobra.Command{Use: "start", Short: "short"}
		cmd.Example = "mesheryctl system start\n\t\t\n   \n\t"

		buf := &bytes.Buffer{}
		err := GenMarkdownCustom(cmd, buf, map[int]string{})
		assert.NoError(t, err)

		// Only the single real example is emitted; the blank lines are dropped.
		assert.Equal(t, 1, strings.Count(buf.String(), "<code class='clipboardjs'>"))
	})
}
