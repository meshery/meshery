package pkg

import (
	"encoding/csv"
	"fmt"
)

func GetEntries(r *csv.Reader, columnNames []string) ([]map[string]string, error) {
	final := make([]map[string]string, 0)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	mapColumns := make(map[string]bool)
	for _, c := range columnNames {
		mapColumns[c] = true
	}
	var internalIndex map[int]string
	var start bool
	var startReadingFrom int
	for i, record := range records {
		start, internalIndex = startRecording(record, mapColumns)
		if start {
			startReadingFrom = i + 1
			break
		}
	}
	if internalIndex == nil {
		return nil, fmt.Errorf("could not populate entries")
	}
	for startReadingFrom < len(records) {
		record := records[startReadingFrom]
		temp := make(map[string]string)
		for i, cell := range record {
			key := internalIndex[i]
			value := cell
			temp[key] = value
		}
		final = append(final, temp)
		startReadingFrom++
	}
	return final, nil
}

func startRecording(record []string, columnNames map[string]bool) (bool, map[int]string) {
	for _, r := range record {
		if columnNames[r] { //even if one column is detected, recording the further rows will start
			m := make(map[int]string)
			for i, r0 := range record {
				m[i] = r0
			}
			return true, m
		}
	}
	return false, nil
}
