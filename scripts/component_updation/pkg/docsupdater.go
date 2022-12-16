package pkg

import (
	"os"
	"strings"
)

const template string = `---
title: <Project Name>
subtitle: <Page Subtitle>
integrationIcon: ../../../assets/images/service-mesh-icons/aws-app-mesh.svg
darkModeIntegrationIcon: ../../../assets/images/service-mesh-icons/aws-app-mesh_white.svg
docURL: <Docs URL>
category: <Category>
subcategory: <Sub-Category>
featureList: [<Feature 1>,<Feature 2>,<Feature 3>]
workingSlides: [
    ../_images/meshmap-visualizer.png,
    ../_images/meshmap-designer.png]
howItWorks: <howItWorks>
howItWorksDetails: howItWorksDetails
published: <Publish>
---
<p>
   <About Project>
</p>
<p>
   <Standard Blurb>
</p>`

// func createEmptyMarkdown(path string) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = file.Write([]byte(template))
// 	return err
// }

type TemplateAttributes struct {
	Title                   string
	Subtitle                string
	DocURL                  string
	Category                string
	Subcategory             string
	FeatureList             string
	HowItWorks              string
	HowItWorksDetails       string
	AboutProject            string
	StandardBlurb           string
	WorkingSlides           string
	Published               string
	IntegrationIcon         string
	DarkModeIntegrationIcon string
	FullPage                string
}

func (t TemplateAttributes) CreateMarkDown() string {
	// markdown := "---\n"
	// markdown += "title: " + t.Title + "\n"
	// markdown += "subtitle: " + t.Subtitle + "\n"
	// markdown += "integrationIcon: " + t.IntegrationIcon + "\n"
	// markdown += "darkModeIntegrationIcon: " + t.DarkModeIntegrationIcon + "\n"
	// markdown += "docURL: " + t.DocURL + "\n"
	// markdown += "category: " + t.Category + "\n"
	// markdown += "subcategory: " + t.Subcategory + "\n"
	// markdown += "featureList: " + t.FeatureList + "\n"
	// markdown += "workingSlides: " + t.WorkingSlides + "\n"
	// markdown += "howItWorks: " + t.HowItWorks + "\n"
	// markdown += "howItWorksDetails: " + t.HowItWorksDetails + "\n"
	// markdown += "published: " + t.Published + "\n"
	// markdown += "---\n"
	// markdown += t.AboutProject + "\n"
	// markdown += t.StandardBlurb
	markdown := t.FullPage
	markdown = strings.ReplaceAll(markdown, "published: TRUE", "published: true")
	markdown = strings.ReplaceAll(markdown, "published: FALSE", "published: false")
	return markdown
}
func WriteMarkDown(path string, md string) error {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(md)
	if err != nil {
		panic(err)
	}
	// Close the file to save the changes.
	err = file.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
func WriteSVG(path string, svg string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.WriteString(svg)
	if err != nil {
		return err
	}
	// Close the file to save the changes.
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

// func (t *templateAttributes) fillAttributes(path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
// 	scanner := bufio.NewScanner(file)
// 	scanner.Split(bufio.ScanLines)
// 	aboutProjectRead := false
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.HasPrefix(line, "title:") {
// 			t.title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
// 		} else if strings.HasPrefix(line, "subtitle:") {
// 			t.subtitle = strings.TrimSpace(strings.TrimPrefix(line, "subtitle:"))
// 		} else if strings.HasPrefix(line, "docURL:") {
// 			t.docURL = strings.TrimSpace(strings.TrimPrefix(line, "docURL:"))
// 		} else if strings.HasPrefix(line, "category:") {
// 			t.category = strings.TrimSpace(strings.TrimPrefix(line, "category:"))
// 		} else if strings.HasPrefix(line, "subcategory:") {
// 			t.subcategory = strings.TrimSpace(strings.TrimPrefix(line, "subcategory:"))
// 		} else if strings.HasPrefix(line, "subtitle:") {
// 			t.howItWorks = strings.TrimSpace(strings.TrimPrefix(line, "howItWorks:"))
// 		} else if strings.HasPrefix(line, "howItWorksDetails:") {
// 			t.howItWorksDetails = strings.TrimSpace(strings.TrimPrefix(line, "howItWorksDetails:"))
// 		} else if strings.HasPrefix(line, "workingSlides:") {
// 			t.workingSlides = strings.TrimSpace(strings.TrimPrefix(line, "workingSlides:"))
// 			if !strings.HasSuffix(t.workingSlides, "]") {
// 				for scanner.Scan() {
// 					subline := scanner.Text()
// 					t.workingSlides += "\n" + subline
// 					if strings.HasSuffix(subline, "]") {
// 						break
// 					}
// 				}
// 			}
// 		} else if strings.HasPrefix(line, "integrationIcon") {
// 			t.integrationIcon = strings.TrimSpace(strings.TrimPrefix(line, "integrationIcon:"))
// 		} else if strings.HasPrefix(line, "darkModeIntegrationIcon") {
// 			t.darkModeIntegrationIcon = strings.TrimSpace(strings.TrimPrefix(line, "darkModeIntegrationIcon:"))
// 		} else if strings.HasPrefix(line, "published") {
// 			t.published = strings.TrimSpace(strings.TrimPrefix(line, "published:"))
// 		} else if strings.HasPrefix(line, "<p>") {
// 			if !aboutProjectRead {
// 				t.aboutProject = line
// 				if !strings.HasSuffix(t.aboutProject, "</p>") {
// 					for scanner.Scan() {
// 						subline := scanner.Text()
// 						t.aboutProject += "\n" + subline
// 						if strings.HasSuffix(subline, "</p>") {
// 							break
// 						}
// 					}
// 				}
// 			} else {
// 				t.standardBlurb = line
// 				if !strings.HasSuffix(t.aboutProject, "</p>") {
// 					for scanner.Scan() {
// 						subline := scanner.Text()
// 						t.standardBlurb += "\n" + subline
// 						if strings.HasSuffix(subline, "</p>") {
// 							break
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }
