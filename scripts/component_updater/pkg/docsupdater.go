package pkg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const XMLTAG = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE svg>"

func WriteToFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(content)
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

// UpdateSVGString updates the width and height attributes of an SVG file and returns the modified SVG as a string.
func UpdateSVGString(svgStr string, width, height int) (string, error) {
	// Create a reader for the SVG string.
	r := strings.NewReader(svgStr)

	// Create a decoder for the SVG string.
	d := xml.NewDecoder(r)

	// Create a buffer for the modified SVG string.
	var b bytes.Buffer

	// Create an encoder for the buffer.
	e := xml.NewEncoder(&b)

	// Iterate through the tokens in the SVG string.
	for {
		// Read the next token.
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// If the token is an element name, check if it is an "svg" element.
		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "svg" {
				// Set the width and height attributes to the desired values.
				updatedH := false
				updatedW := false
				xmlnsindex := -1
				for i, a := range se.Attr {
					if a.Name.Local == "width" {
						se.Attr[i].Value = strconv.Itoa(width)
						updatedW = true
					} else if a.Name.Local == "height" {
						se.Attr[i].Value = strconv.Itoa(height)
						updatedH = true
					}
					if a.Name.Local == "xmlns" {
						xmlnsindex = i
					}
				}
				if !updatedH {
					se.Attr = append(se.Attr, xml.Attr{
						Name: xml.Name{
							Local: "height",
						},
						Value: strconv.Itoa(height),
					})
				}
				if !updatedW {
					se.Attr = append(se.Attr, xml.Attr{
						Name: xml.Name{
							Local: "width",
						},
						Value: strconv.Itoa(width),
					})
				}
				if xmlnsindex > -1 {
					se.Attr = append(se.Attr[0:xmlnsindex], se.Attr[xmlnsindex+1:]...)
				}
			} else {
				for i, a := range se.Attr {
					xmlnsindex := -1
					nahbro := 0
					if a.Name.Local == "xmlns" {
						fmt.Println("found at ", i)
						fmt.Println(a.Name)
						xmlnsindex = i
						nahbro++
					}
					if xmlnsindex > -1 {
						se.Attr = append(se.Attr[0:xmlnsindex], se.Attr[xmlnsindex+1:]...)
					}
				}
			}
			t = se
		}
		// Write the modified token to the buffer.
		if err := e.EncodeToken(t); err != nil {
			return "", err
		}
	}

	// Flush the encoder's buffer to the buffer.
	if err := e.Flush(); err != nil {
		return "", err
	}
	var svg string
	if b.String() != "" {
		svg = XMLTAG + b.String()
	}
	return svg, nil
}
