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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/spf13/cobra"
)

// TestDoc is the main function for the TestDoc test
func TestDoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Doc tests")
}

var _ = Describe("Tests for Doc", func() {

	// Create a new cobra command before each test
	var cmd *cobra.Command
	BeforeEach(func() {
		cmd = &cobra.Command{
			Use: "test",
		}
	})

	// TestLinkHandler is the link handler for the markdown file
	const TestLinkHandler = "/main"
	Context("Test linkHandler function", func() {
		It("should return the correct linkHandler", func() {
			// call linkHandler
			linkHandler := linkHandler("test.md")
			// check if linkHandler is correct
			Expect(linkHandler).To(Equal(TestLinkHandler))
		})
	})

	// TestGenMarkdownTreeCustom is the test for GenMarkdownTreeCustom function
	Context("Test GenMarkdownTreeCustom function", func() {
		It("should return nil", func() {
			// add a subcommand
			cmd.AddCommand(&cobra.Command{
				Use: "sub",
			})
			// path for docs
			markDownPath := "../../docs/pages/reference/mesheryctl/"

			// call GenMarkdownTreeCustom
			err := GenMarkdownTreeCustom(cmd, markDownPath, prepender, linkHandler)
			// check if err is nil
			Expect(err).To(BeNil())
		})
	})

	// TestGenMarkdownCustom is the test for GenMarkdownCustom function
	Context("Test GenMarkdownCustom function", func() {
		It("should contain specific sub-strings", func() {
			// add annotations "link" and "caption" to the command
			cmd.Annotations = map[string]string{
				"link":    "test_link",
				"caption": "test_caption",
			}

			// add Example for cmd for test
			cmd.Example = "test_example"

			// io.Writer
			buf := &bytes.Buffer{}
			// call GenMarkdownCustom
			err := GenMarkdownCustom(cmd, buf)
			// check if err is nil
			Expect(err).To(BeNil())
			// check if buf is not empty
			Expect(buf.String()).NotTo(BeEmpty())
			// check if buf contains the correct string
			Expect(buf.String()).To(ContainSubstring("test_link"))
			Expect(buf.String()).To(ContainSubstring("test_caption"))
			Expect(buf.String()).To(ContainSubstring("test_example"))
		})
	})

	// TestHasSeeAlso is the test for HasSeeAlso function
	Context("Test HasSeeAlso function", func() {
		It("should return false", func() {
			// call HasSeeAlso
			hasSeeAlso := hasSeeAlso(cmd)
			// check if hasSeeAlso is false
			Expect(hasSeeAlso).To(BeFalse())
		})
	})

	// TestGenYamlTreeCustom is the test for GenYamlTreeCustom function
	Context("Test GenYamlTreeCustom function", func() {
		It("should return nil", func() {
			// add a subcommand
			cmd.AddCommand(&cobra.Command{
				Use: "sub",
			})
			// path for docs
			yamlPath := "../../docs/pages/reference/mesheryctl/"
			// call GenYamlTreeCustom
			err := GenYamlTreeCustom(cmd, yamlPath, prepender, linkHandler)
			// check if err is nil
			Expect(err).To(BeNil())
		})
	})

	//TestGenYamlCustom is the test for GenYamlCustom function
	Context("Test GenYamlCustom function", func() {
		It("should return specific sub-strings", func() {
			// add Example for cmd for test
			cmd.Example = "test_example"
			// io.Writer
			buf := &bytes.Buffer{}
			// call GenYamlCustom
			err := GenYamlCustom(cmd, buf)
			// check if err is nil
			Expect(err).To(BeNil())
			// check if buf is not empty
			Expect(buf.String()).NotTo(BeEmpty())
			// check if buf contains the correct string
			Expect(buf.String()).To(ContainSubstring("test_example"))
		})
	})
})
