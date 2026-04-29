package utils

import (
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type TableHeader string

func (h TableHeader) Format() string {
	c := color.New(color.Bold).SprintFunc()
	return c(tw.Title(string(h)))
}

func generateTableOptions() []tablewriter.Option {
	options := []tablewriter.Option{
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.Separators{
					ShowHeader:     tw.Off,
					ShowFooter:     tw.Off,
					BetweenRows:    tw.Off,
					BetweenColumns: tw.Off},
				Lines: tw.Lines{
					ShowTop:        tw.Off,
					ShowBottom:     tw.Off,
					ShowHeaderLine: tw.Off,
					ShowFooterLine: tw.Off,
				},
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			MaxWidth: 255,
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{
					Global: tw.AlignLeft,
				},
				Formatting: tw.CellFormatting{
					AutoWrap:   tw.WrapNone,
					MergeMode:  tw.MergeNone,
					AutoFormat: tw.Off, // switch off AutoFormat
				},
				Padding: tw.CellPadding{
					Global: tw.Padding{
						Left:      tw.Empty,
						Right:     "  ",
						Top:       tw.Empty,
						Bottom:    tw.Empty,
						Overwrite: true,
					},
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap:   tw.WrapNone,
					MergeMode:  tw.MergeNone,
					AutoFormat: tw.Off,
				},
				Alignment: tw.CellAlignment{
					Global: tw.AlignLeft,
				},
				Padding: tw.CellPadding{
					Global: tw.Padding{
						Left:      tw.Empty,
						Right:     "  ",
						Top:       tw.Empty,
						Bottom:    tw.Empty,
						Overwrite: true,
					},
				},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	}
	return options
}

func renderTable(table *tablewriter.Table, data [][]string, header, footer []string) {
	tableHeader := make([]any, len(header))
	for i, h := range header {
		tableHeader[i] = TableHeader(h)
	}
	table.Header(tableHeader...)
	err := table.Bulk(data)
	if err != nil {
		Log.Error(ErrTableRender(err))

	}
	if footer != nil {
		table.Footer(footer)
	}
	err = table.Render()
	if err != nil {
		Log.Error(ErrTableRender(err))
	}
}

// PrintToTable prints the given data into a table format
func PrintToTable(header []string, data [][]string, footer []string) {
	options := generateTableOptions()

	table := tablewriter.NewTable(os.Stdout,
		options...,
	)

	renderTable(table, data, header, footer)
}
