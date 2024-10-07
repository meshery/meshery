package registry

import (
	"google.golang.org/api/sheets/v4"
)

// GetSheetIDFromTitle returns the sheet ID for the sheet with the given title.
func GetSheetIDFromTitle(s *sheets.Spreadsheet, title string) int64 {
	for _, sheet := range s.Sheets {
		if sheet.Properties.Title == title {
			return sheet.Properties.SheetId
		}
	}
	return -1
}
