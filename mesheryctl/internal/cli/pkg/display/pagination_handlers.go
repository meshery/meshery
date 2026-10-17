package display

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var (
	whiteBoardPrinter    = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	nextPageKeyboardKeys = []keyboard.Key{keyboard.KeyEnter, keyboard.KeyArrowDown, keyboard.KeySpace}
	escapeKeyboardKeys   = []keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC}
)

// listPageHandler creates a pageHandler that displays each fetched page as a table.
// It tracks the current index across pages to determine when all data has been shown
func listPageHandler[T any](displayData DisplayDataAsync, processDataFunc listRowBuilder[T]) pageHandler[T] {
	startIndex := 0
	return func(data *T, currentPage int, pgSize int) (bool, error) {
		rows, totalCount := processDataFunc(data)

		// Display the total count and current page
		utils.DisplayCount(displayData.DataType, totalCount)

		if len(rows) == 0 {
			return false, nil
		}

		if displayData.DisplayCountOnly {
			return false, nil
		}

		// Display the current page number to be one-based
		_, _ = whiteBoardPrinter.Fprint(os.Stdout, "Page: ", currentPage+1)
		fmt.Println()

		// Display the data in a table
		utils.PrintToTable(displayData.Header, rows, nil)

		if displayData.IsPage {
			return false, nil
		}

		// If the URL already contains "pagesize=all", it means all data has been fetched in one go,
		// so we can break the loop without waiting for user input
		if strings.Contains(displayData.UrlPath, "pagesize=all") {
			return false, nil
		}

		if int64(startIndex+pgSize) >= totalCount {
			return false, nil
		}

		// Paging on a keypress needs somebody to press the key. Without a
		// terminal the keyboard library below fails on /dev/tty and the command
		// exits non-zero having already printed a perfectly good first page -
		// which is how `mesheryctl connection list` came to fail outright in CI
		// and in any pipeline, as soon as the account held more than one page.
		// Stop after this page instead, and say so, so the output stays usable
		// and the exit status stays honest.
		if !utils.IsInteractiveTerminal() {
			utils.Log.Infof(
				"Showing page %d only: paging through results needs an interactive terminal. Use --page and --pagesize to select a page, or --count for the total.",
				currentPage+1,
			)
			return false, nil
		}

		// Wait for user input to navigate pages
		keysEvents, err := keyboard.GetKeys(10)
		if err != nil {
			return false, err
		}

		defer func() {
			_ = keyboard.Close()
		}()

		event := <-keysEvents
		if event.Err != nil {
			utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
			return false, event.Err
		}

		if slices.Contains(escapeKeyboardKeys, event.Key) {
			return false, nil
		}

		if slices.Contains(nextPageKeyboardKeys, event.Key) {
			startIndex += pgSize
			return true, nil
		} else {
			return false, nil
		}
	}
}

// promptPageHandler creates a pageHandler that presents each fetched page as a selection prompt.
func promptPageHandler[T any, R any](displayData DisplayDataAsync, processData promptLabelBuilder[R], extractItem itemExtractor[T, R], selectedItem *R) pageHandler[T] {
	return func(data *T, currentPage int, pgSize int) (bool, error) {
		rows, totalCount := extractItem(data)

		if len(rows) == 0 {
			if currentPage == 0 {
				msg := "no results found"
				if displayData.ErrNotFoundMsg != "" {
					msg = displayData.ErrNotFoundMsg
				}
				return false, utils.ErrNotFound(errors.New(msg))

			}
			return false, nil
		}

		// Auto-select the sole match. The test is on totalCount, not on the size
		// of this page: a page can hold one row while further pages hold more
		// (the caller chooses the page size, and --pagesize 1 makes every page
		// look like a single match), and selecting then would act on the first
		// of several candidates as though the search had been unambiguous.
		if totalCount == 1 && len(rows) == 1 && currentPage == 0 {
			*selectedItem = rows[0]
			return false, nil
		}

		// When a search term was provided and all matching candidates are available
		// on this page, auto-select if exactly one candidate has an exact name match
		// (case-insensitively). This allows commands like `model view <name>` to succeed
		// without failing as ambiguous in non-interactive environments when a query
		// returns multiple substring/related matches but only one exact name match.
		if searchTerm := strings.TrimSpace(displayData.SearchTerm); searchTerm != "" && currentPage == 0 && int64(len(rows)) == totalCount {
			var exactMatches []R
			for _, row := range rows {
				for _, name := range getCandidateNames(displayData, row) {
					if strings.EqualFold(strings.TrimSpace(name), searchTerm) {
						exactMatches = append(exactMatches, row)
						break
					}
				}
			}
			if len(exactMatches) == 1 {
				*selectedItem = exactMatches[0]
				return false, nil
			}
		}

		// More than one match needs a choice, and a choice needs somebody to
		// make it. Unlike paging there is no sensible non-interactive fallback -
		// picking for the user would act on a resource they did not name - so
		// this is a genuine failure. Say what happened and what to do about it,
		// rather than letting promptui fail on /dev/tty with an error that
		// describes neither.
		if !utils.IsInteractiveTerminal() {
			return false, ErrAmbiguousSelection(totalCount)
		}

		picked, itemSelected, err := SelectFromPagedResults(rows, processData, pgSize, currentPage, totalCount)
		if err != nil {
			return false, err
		}
		if itemSelected {
			*selectedItem = picked
			return false, nil
		}
		return true, nil
	}
}

// getCandidateNames retrieves candidate resource names from item R,
// preferring a caller-supplied ExtractNames function when provided.
func getCandidateNames(displayData DisplayDataAsync, item any) []string {
	if displayData.ExtractNames != nil {
		if names := displayData.ExtractNames(item); len(names) > 0 {
			return names
		}
	}
	return extractCandidateNames(item)
}

// extractCandidateNames extracts potential names from an item via common interfaces
// or reflection across known name fields (Name, DisplayName, Kind, Component.Kind).
func extractCandidateNames(item any) []string {
	if item == nil {
		return nil
	}

	type itemNamer interface {
		ItemNames() []string
	}
	if in, ok := item.(itemNamer); ok {
		return in.ItemNames()
	}

	type singleNamer interface {
		ItemName() string
	}
	if sn, ok := item.(singleNamer); ok {
		return []string{sn.ItemName()}
	}

	type getNameInterface interface {
		GetName() string
	}
	if gn, ok := item.(getNameInterface); ok {
		return []string{gn.GetName()}
	}

	type nameInterface interface {
		Name() string
	}
	if n, ok := item.(nameInterface); ok {
		return []string{n.Name()}
	}

	val := reflect.ValueOf(item)
	for val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	var names []string

	// Check field "Name"
	if f := val.FieldByName("Name"); f.IsValid() && f.Kind() == reflect.String {
		if s := f.String(); s != "" {
			names = append(names, s)
		}
	}

	// Check field "DisplayName"
	if f := val.FieldByName("DisplayName"); f.IsValid() {
		if f.Kind() == reflect.String {
			if s := f.String(); s != "" {
				names = append(names, s)
			}
		} else if f.CanInterface() {
			if stringer, ok := f.Interface().(fmt.Stringer); ok {
				if s := stringer.String(); s != "" {
					names = append(names, s)
				}
			}
		}
	}

	// Check field "Kind"
	if f := val.FieldByName("Kind"); f.IsValid() && f.Kind() == reflect.String {
		if s := f.String(); s != "" {
			names = append(names, s)
		}
	}

	// Check nested "Component.Kind"
	if comp := val.FieldByName("Component"); comp.IsValid() {
		for comp.Kind() == reflect.Pointer || comp.Kind() == reflect.Interface {
			if comp.IsNil() {
				break
			}
			comp = comp.Elem()
		}
		if comp.Kind() == reflect.Struct {
			if kind := comp.FieldByName("Kind"); kind.IsValid() && kind.Kind() == reflect.String {
				if s := kind.String(); s != "" {
					names = append(names, s)
				}
			}
		}
	}

	return names
}
