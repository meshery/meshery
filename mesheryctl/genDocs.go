package main

import (
	//"bytes"
	"fmt"
	//"io"
	"log"
	"path"
	"path/filepath"
	"strings"

	//"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	//"gopkg.in/yaml.v2"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/app"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/experimental"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/mesh"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/pattern"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/perf"
	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
)

const MARKDOWN_TEMPLATE_COMMAND = `---
layout: default
title: %s
permalink: /%s
redirect_from: /%s
type: reference
display-title: false
language: en
command: %s
---
`

const MARKDOWN_TEMPLATE_SUBCOMMAND = `---
layout: default
title: %s
permalink: /%s
redirect_from: /%s
type: reference
display-title: false
language: en
command: %s
subcommand: %s
---
`

type cmdDoc struct {
	name        string
	description string
	usage       string   `yaml:",omitempty"`
	examples    []string `yaml:",omitempty"`
}

func prepender(filename string) string {
	title := filepath.Base(filename)
	base := strings.TrimSuffix(title, path.Ext(title))
	url := "reference/" + strings.ToLower(base) + "/"
	return fmt.Sprintf(MARKDOWN_TEMPLATE_COMMAND, title, url, url, "test")
}

func linkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))
	return "reference/" + strings.ToLower(base) + "/"
}

func main() {
	markDownPath := "./internal/cli/root/testDoc/"
	//yamlPath := "./internal/cli/root/testDoc/"

	fmt.Println("Scanning available commands...")
	cmd := root.TreePath()
	fmt.Println("Generating markdown docs...")

	err := doc.GenMarkdownTreeCustom(cmd, markDownPath, prepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

/*
func GenYaml(cmd *cobra.Command, w io.Writer) error {
	return GenYamlCustom(cmd, w, func (s string) string  {
		return s
	})
}

func GenYamlCustom(cmd *cobra.Command, w io.Writer, linkHandler string) error {
	yaml := cmdDoc{}
	yaml.name = cmd.CommandPath()
	yaml.description = forceMultiLine(cmd.Long)
	yaml.usage = forceString(cmd.Use)
	if len(cmd.Example) > 0 {
		yaml.examples = strings.Split(cmd.Example, "\n")
	}

	final, err := yaml.Marshal(&yamlDoc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := w.Write(final); err != nil {
		return err
	}
	return nil
}
*/
