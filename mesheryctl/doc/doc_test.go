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
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
		filename := "docs/content/en/reference/mesheryctl/_index.md"
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
		filename := "docs/content/en/reference/mesheryctl/adapter/_index.md"
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
		filename := "docs/content/en/reference/mesheryctl/adapter/deploy.md"
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
		filename := "docs/content/en/reference/mesheryctl/exp/relationship/generate.md"
		assert.Equal(t, expected, prepender(filename))
	})
}

func TestDoc(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}

	t.Run("Test linkHandler function (directory structure)", func(t *testing.T) {
		assert.Equal(t,
			"/reference/mesheryctl",
			linkHandler("docs/content/en/reference/mesheryctl/_index.md"),
		)

		assert.Equal(t,
			"/reference/mesheryctl/adapter",
			linkHandler("docs/content/en/reference/mesheryctl/adapter/_index.md"),
		)

		assert.Equal(t,
			"/reference/mesheryctl/adapter/deploy",
			linkHandler("docs/content/en/reference/mesheryctl/adapter/deploy.md"),
		)

		assert.Equal(t,
			"/reference/mesheryctl/exp/relationship/generate",
			linkHandler("docs/content/en/reference/mesheryctl/exp/relationship/generate.md"),
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

func TestCommandDataYAMLTags(t *testing.T) {
	data := commandDataFile{
		Order: []string{"global"},
		Commands: map[string]commandData{
			"global": {
				Key:         "global",
				Name:        "mesheryctl",
				Description: "Mesheryctl root command",
				Usage:       "mesheryctl [flags]",
				URL:         "/reference/mesheryctl/",
				Flags: []flagData{{
					Key:         "config",
					Name:        "--config, -c",
					Description: "path to config file",
				}},
				Subcommands: []commandData{{
					Key:         "version",
					Name:        "version",
					Description: "Print version",
					Usage:       "mesheryctl version",
					URL:         "/reference/mesheryctl/version/",
				}},
			},
		},
	}

	content, err := yaml.Marshal(data)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "order:")
	assert.Contains(t, output, "commands:")
	assert.Contains(t, output, "key: global")
	assert.Contains(t, output, "name: mesheryctl")
	assert.Contains(t, output, "description: Mesheryctl root command")
	assert.Contains(t, output, "usage: mesheryctl [flags]")
	assert.Contains(t, output, "url: /reference/mesheryctl/")
	assert.Contains(t, output, "flags:")
	assert.Contains(t, output, "subcommands:")
}

func TestGenerateCommandData(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "mesheryctl",
		Short: "Mesheryctl root command",
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file")

	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "Manage Meshery lifecycle",
	}
	systemCmd.Flags().Bool("dry-run", false, "preview the operation")

	contextCmd := &cobra.Command{
		Use:   "context",
		Short: "Manage contexts",
	}
	viewCmd := &cobra.Command{
		Use:   "view [context-name]",
		Short: "View a context",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	viewCmd.Flags().StringP("output", "o", "", "output format")
	contextCmd.AddCommand(viewCmd)
	systemCmd.AddCommand(contextCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	expCmd := &cobra.Command{
		Use:   "exp",
		Short: "Preview experimental commands",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	hiddenCmd := &cobra.Command{
		Use:    "hidden",
		Short:  "Hidden command",
		Hidden: true,
		Run:    func(cmd *cobra.Command, args []string) {},
	}
	helpTopicCmd := &cobra.Command{
		Use:   "topic",
		Short: "Additional help topic",
	}

	rootCmd.AddCommand(versionCmd, hiddenCmd, helpTopicCmd, systemCmd, expCmd)

	outputPath := filepath.Join(t.TempDir(), "cmds.yml")
	require.NoError(t, GenerateCommandData(rootCmd, outputPath))

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Code generated by mesheryctl/doc/doc.go; DO NOT EDIT.")

	var generated commandDataFile
	require.NoError(t, yaml.Unmarshal(content, &generated))

	assert.Equal(t, []string{"global", "system", "exp"}, generated.Order)

	globalCommand := generated.Commands["global"]
	assert.Equal(t, "mesheryctl", globalCommand.Name)
	assert.Equal(t, "/reference/mesheryctl/", globalCommand.URL)
	require.Len(t, globalCommand.Flags, 2)
	assert.Equal(t, "config", globalCommand.Flags[0].Key)
	assert.Equal(t, "--config, -c", globalCommand.Flags[0].Name)
	assert.Equal(t, "help", globalCommand.Flags[1].Key)
	assert.Equal(t, "--help, -h", globalCommand.Flags[1].Name)
	require.Len(t, globalCommand.Subcommands, 1)
	assert.Equal(t, "version", globalCommand.Subcommands[0].Name)

	systemCommand := generated.Commands["system"]
	assert.Equal(t, "system", systemCommand.Name)
	assert.Equal(t, "Manage Meshery lifecycle", systemCommand.Description)
	assert.Equal(t, "mesheryctl system [flags]", systemCommand.Usage)
	assert.Equal(t, "/reference/mesheryctl/system/", systemCommand.URL)
	require.Len(t, systemCommand.Flags, 1)
	assert.Equal(t, "dry-run", systemCommand.Flags[0].Key)
	require.Len(t, systemCommand.Subcommands, 1)
	assert.Equal(t, "context", systemCommand.Subcommands[0].Name)
	assert.Equal(t, "/reference/mesheryctl/system/context/", systemCommand.Subcommands[0].URL)
	require.Len(t, systemCommand.Subcommands[0].Subcommands, 1)
	assert.Equal(t, "view", systemCommand.Subcommands[0].Subcommands[0].Name)
	assert.Equal(t, "/reference/mesheryctl/system/context/view/", systemCommand.Subcommands[0].Subcommands[0].URL)
	require.Len(t, systemCommand.Subcommands[0].Subcommands[0].Flags, 1)
	assert.Equal(t, "--output, -o", systemCommand.Subcommands[0].Subcommands[0].Flags[0].Name)

	expCommand := generated.Commands["exp"]
	assert.Equal(t, "exp", expCommand.Name)
	assert.Equal(t, "Preview experimental commands", expCommand.Description)
	assert.Equal(t, "/reference/mesheryctl/exp/", expCommand.URL)
}
