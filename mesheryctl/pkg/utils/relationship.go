package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/csv"
	"github.com/meshery/schemas/models/v1alpha3"
	_rel "github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/model"
)

type RelationshipCSVHelper struct {
	SpreadsheetID        int64
	SpreadsheetURL       string
	Title                string
	CSVPath              string
	Relationships        []RelationshipCSV
	UpdatedRelationships []RelationshipCSV
}

type RelationshipCSV struct {
	RowIndex    int    `json:"-" csv:"-"`
	Model       string `json:"Model" csv:"Model"`
	Version     string `json:"Version" csv:"Version"`
	KIND        string `json:"kind" csv:"kind"`
	Type        string `json:"type" csv:"type"`
	SubType     string `json:"subType" csv:"subType"`
	Description string `json:"metadata.description" csv:"metadata.description	"`
	Styles      string `json:"metadata.styles" csv:"metadata.styles"`
	EvalPolicy  string `json:"evalPolicy" csv:"evalPolicy"`
	Selector    string `json:"selector" csv:"selector"`
	Filename    string `json:"filename" csv:"filename"`
}

func NewRelationshipCSVHelper(sheetURL, spreadsheetName string, spreadsheetID int64, localCsvPath string) (*RelationshipCSVHelper, error) {
	var csvPath string
	if localCsvPath == "" {
		sheetURL = sheetURL + "/pub?output=csv" + "&gid=" + strconv.FormatInt(spreadsheetID, 10)
		Log.Info("Downloading CSV from: ", sheetURL)
		dirPath := filepath.Join(utils.GetHome(), ".meshery", "content")
		_ = os.MkdirAll(dirPath, 0755)
		csvPath = filepath.Join(dirPath, "relationship.csv")
		err := utils.DownloadFile(csvPath, sheetURL)
		if err != nil {
			return nil, utils.ErrReadingRemoteFile(err)
		}
	} else {
		csvPath = localCsvPath
	}

	return &RelationshipCSVHelper{
		SpreadsheetID:  spreadsheetID,
		SpreadsheetURL: sheetURL,
		Title:          spreadsheetName,
		CSVPath:        csvPath,
		Relationships:  []RelationshipCSV{},
	}, nil
}
func (mrh *RelationshipCSVHelper) ParseRelationshipsSheet(modelName string) error {
	ch := make(chan RelationshipCSV, 1)
	errorChan := make(chan error, 1)
	csvReader, err := csv.NewCSVParser[RelationshipCSV](mrh.CSVPath, rowIndex, nil, func(_ []string, _ []string) bool {
		return true
	})

	if err != nil {
		return ErrFileRead(err)
	}

	rowIndex := 3 //first 2 are headers
	currentRow := rowIndex

	go func() {
		Log.Info("Parsing Relationships...")

		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			errorChan <- err
		}
	}()

	for {
		select {
		case data := <-ch:
			if modelName != "" && data.Model != modelName {
				continue
			}
			data.RowIndex = currentRow
			currentRow++
			mrh.Relationships = append(mrh.Relationships, data)
		case err := <-errorChan:
			Log.Error(err)
		case <-csvReader.Context.Done():
			return nil
		}
	}
}

func ProcessRelationships(relationshipCSVHelper *RelationshipCSVHelper, spreadsheetUpdateChan chan RelationshipCSV) {
	for _, relationship := range relationshipCSVHelper.Relationships {
		var versions []string

		if relationship.Version == "*" {
			modelBasePath := fmt.Sprintf("../server/meshmodel/%s", relationship.Model)
			dirs, err := os.ReadDir(modelBasePath)
			if err != nil {
				err = utils.ErrReadDir(err, modelBasePath)
				Log.Error(err)
				continue
			}

			for _, dir := range dirs {
				if dir.IsDir() {
					versions = append(versions, dir.Name())
				}
			}

			if len(versions) == 0 {
				continue
			}
		} else {
			versions = []string{relationship.Version}
		}

		for _, version := range versions {
			if relationship.Model == "" {
				continue
			}
			var rel _rel.RelationshipDefinition
			rel.SchemaVersion = v1alpha3.RelationshipSchemaVersion
			rel.Kind = _rel.RelationshipDefinitionKind(utils.ReplaceSpacesWithHyphenAndConvertToLowercase(relationship.KIND))
			rel.RelationshipType = utils.ReplaceSpacesWithHyphenAndConvertToLowercase(relationship.Type)
			rel.SubType = utils.ReplaceSpacesWithHyphenAndConvertToLowercase(relationship.SubType)
			rel.EvaluationQuery = &relationship.EvalPolicy

			rel.Version = "v1.0.0"
			rel.Model = model.ModelDefinition{
				Name:  relationship.Model,
				Model: model.Model{Version: version},
			}
			status := _rel.Enabled
			rel.Status = &status
			var styles _rel.RelationshipDefinitionMetadataStyles

			if relationship.Styles != "" {
				err := encoding.Unmarshal([]byte(relationship.Styles), &styles)
				if err != nil {
					Log.Error(err)
					continue
				}
				if rel.Metadata == nil {
					rel.Metadata = &_rel.Relationship_Metadata{}
				}
			}

			rel.Metadata = &_rel.Relationship_Metadata{
				Description: &relationship.Description,
				Styles:      &styles,
			}

			var selectorSet _rel.SelectorSet
			if relationship.Selector != "" {
				err := encoding.Unmarshal([]byte(relationship.Selector), &selectorSet)
				if err != nil {
					Log.Error(err)
					continue
				}
				rel.Selectors = &selectorSet
			}

			filenameGenerated := false
			if relationship.Filename == "" {
				relationship.Filename = fmt.Sprintf("%s-%s-%s-%s.json", rel.Kind, rel.RelationshipType, rel.SubType, utils.GetRandomAlphabetsOfDigit(5))
				filenameGenerated = true
			}

			fullPath, err := ConstructRelationshipPath(relationship.Model, version, rel.Version, "../server/meshmodel", relationship.Filename)
			if err != nil {
				Log.Error(err)
				continue
			}

			// Call the separate function to handle file operations
			fileUpdated, err := HandleRelationshipFile(fullPath, rel)
			if err != nil {
				Log.Error(ErrUpdateRelationshipFile(err))
				continue
			}

			if fileUpdated {
				Log.Info(fmt.Sprintf("Relationship written to %s", fullPath))
			} else {
				Log.Info(fmt.Sprintf("No changes detected in %s, skipping write.", fullPath))
			}

			if filenameGenerated {
				spreadsheetUpdateChan <- relationship
			}
		}
	}
}

func ConstructRelationshipPath(modelName, version, defVersion, outputLocation, filename string) (string, error) {
	basePath := fmt.Sprintf("%s/%s/%s/%s", outputLocation, modelName, version, defVersion)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return "", err
	}

	fullPath := filepath.Join(basePath, "relationships", filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return "", err
		}
	}
	return fullPath, nil
}

// HandleRelationshipFile checks if the file exists and compares its content.
// It writes the new data only if the file doesn't exist or the contents differ.
// Returns true if the file was written or updated, false if no changes were made.
func HandleRelationshipFile(filePath string, rel _rel.RelationshipDefinition) (bool, error) {
	newData, err := json.MarshalIndent(rel, "", "  ")
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(filePath); err == nil {
		existingData, err := os.ReadFile(filePath)
		if err != nil {
			return false, err
		}

		if bytes.Equal(existingData, newData) {
			return false, nil
		} else {
			err = os.WriteFile(filePath, newData, 0644)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	} else if os.IsNotExist(err) {
		err = os.WriteFile(filePath, newData, 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err
	}
}
