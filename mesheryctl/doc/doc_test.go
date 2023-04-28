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
			// create a new cobra command
			cmd := &cobra.Command{
				Use: "test",
			}

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

})
