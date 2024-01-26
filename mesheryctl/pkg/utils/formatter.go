package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/layer5io/meshkit/logger"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TerminalFormatter is exported
type TerminalFormatter struct{}

// Format defined the format of output for Logrus logs
// Format is exported
func (f *TerminalFormatter) Format(entry *log.Entry) ([]byte, error) {
	return append([]byte(entry.Message), '\n'), nil
}

// Call this function to setup logrus
func SetupLogrusFormatter() {
	//log formatter for improved UX
	log.SetFormatter(new(TerminalFormatter))
}

// Initialize Meshkit Logger instance
func SetupMeshkitLogger(debugLevel bool, output io.Writer) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if debugLevel {
		logLevel = int(log.DebugLevel)
	}

	logger, err := logger.New("mesheryctl", logger.Options{
		Format:   logger.TerminalLogFormat,
		LogLevel: logLevel,
		Output:   output,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	Log = logger
}

type Paginator struct {
	data                 []string
	lastVisibleItemIndex int
	whiteBgPrinter       *color.Color
	userMsg              string
}

func NewPaginator(header []string, rows [][]string) *Paginator {
	var buf bytes.Buffer
	renderedTableString := writeTableToBuffer(&buf, header, rows)
	lastVisibleItemIndex := 26
	data := strings.Split(renderedTableString, "\n")
	userMsg := "Press Enter to load more. Press ESC or Ctrl+C to stop."
	extraSpace := len(data[0]) - len(userMsg)
	if extraSpace > 0 {
		userMsg += strings.Repeat(" ", extraSpace+2)
	}
	whiteBgPrinter := color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	return &Paginator{
		data:                 data,
		lastVisibleItemIndex: lastVisibleItemIndex,
		whiteBgPrinter:       whiteBgPrinter,
		userMsg:              userMsg,
	}
}

func (p *Paginator) Render() {
	outputTable := strings.Join(p.data[:p.lastVisibleItemIndex], "\n")
	fmt.Println(outputTable)
	p.whiteBgPrinter.Print(p.userMsg)
}

func (p *Paginator) AddLine() bool {
	ClearLine()
	fmt.Println(p.data[p.lastVisibleItemIndex])
	p.lastVisibleItemIndex++
	p.whiteBgPrinter.Print(p.userMsg)
	return p.lastVisibleItemIndex >= len(p.data)-1
}

// PrintTableToBuffer writes the given data into a table format and return formatted string
func writeTableToBuffer(buf *bytes.Buffer, header []string, data [][]string) string {
	table := tablewriter.NewWriter(buf)
	table.SetHeader(header)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
	return buf.String()
}
