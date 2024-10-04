package utils

import (
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
	SpreadsheetID  int64
	SpreadsheetURL string
	Title          string
	CSVPath        string
	Relationships  []RelationshipCSV
}
type RelationshipCSV struct {
	Model       string `json:"Model" csv:"Model"`
	Version     string `json:"Version" csv:"Version"`
	KIND        string `json:"kind" csv:"kind"`
	Type        string `json:"type" csv:"type"`
	SubType     string `json:"subType" csv:"subType"`
	Description string `json:"metadata.description" csv:"metadata.description	"`
	Styles      string `json:"metadata.styles" csv:"metadata.styles"`
	EvalPolicy  string `json:"evalPolicy" csv:"evalPolicy"`
	Selector    string `json:"selector" csv:"selector"`
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
			mrh.Relationships = append(mrh.Relationships, data)
		case err := <-errorChan:
			Log.Error(err)

		case <-csvReader.Context.Done():
			return nil
		}
	}
}
func ProcessRelationships(relationshipCSVHelper *RelationshipCSVHelper) {
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
			rel.Kind = _rel.RelationshipDefinitionKind(utils.ReplaceSpacesAndConvertToLowercase(relationship.KIND))
			rel.RelationshipType = utils.ReplaceSpacesAndConvertToLowercase(relationship.Type)
			rel.SubType = utils.ReplaceSpacesAndConvertToLowercase(relationship.SubType)
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

			fullPath, err := ConstructRelationshipPath(relationship.Model, version, rel.Version, "../server/meshmodel", rel)
			if err != nil {
				Log.Error(err)
				continue
			}

			err = utils.WriteJSONToFile(fullPath, rel)
			if err != nil {
				Log.Error(err)
				continue
			}

			Log.Info(fmt.Sprintf("Relationship written to %s", fullPath))
		}
	}
}

func ConstructRelationshipPath(modelName, version, defVersion, outputLocation string, rel _rel.RelationshipDefinition) (string, error) {
	// Construct the base path using model and version
	basePath := fmt.Sprintf("%s/%s/%s/v1.0.0", outputLocation, modelName, version)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return "", err
	}

	var fromKind, toKind string

	if rel.Selectors != nil && len(*rel.Selectors) > 0 {
		firstSelector := (*rel.Selectors)[0]
		if len(firstSelector.Allow.From) > 0 && len(firstSelector.Allow.To) > 0 {
			if firstSelector.Allow.From[0].Kind != nil {
				fromKind = *firstSelector.Allow.From[0].Kind
			}
			if firstSelector.Allow.To[0].Kind != nil {
				toKind = *firstSelector.Allow.To[0].Kind
			}
		}
	}
	// Construct the filename
	//filename:- kind-type-subtype-from-{selector from kind}-to-{selector to kind}.json
	filename := fmt.Sprintf("%s-%s-%s-from-%s-to-%s.json",
		rel.Kind, rel.RelationshipType, rel.SubType,
		utils.ReplaceSpacesAndConvertToLowercase(fromKind),
		utils.ReplaceSpacesAndConvertToLowercase(toKind),
	)

	relationshipsPath := filepath.Join(basePath, "relationships")

	if _, err := os.Stat(relationshipsPath); os.IsNotExist(err) {
		err = os.MkdirAll(relationshipsPath, os.ModePerm)
		if err != nil {
			err = utils.ErrCreateDir(err, relationshipsPath)
			return "", err
		}
	}
	fullPath := filepath.Join(relationshipsPath, filename)
	return fullPath, nil
}
