package handlers

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"gopkg.in/yaml.v2"
)

// To be used while adding metadata to patterns,filters and applicationss
type ContentModifier struct {
	token    string
	provider models.Provider
	prefObj  *models.Preference
	userID   string
}

func NewContentModifier(token string,
	provider models.Provider,
	prefObj *models.Preference,
	userID string) *ContentModifier {
	return &ContentModifier{
		token:    token,
		provider: provider,
		prefObj:  prefObj,
		userID:   userID,
	}
}

// TODO: Similar mechanisms for filters and applications
// Takes in response bytes, and add metadata to it based on some checks
func (mc *ContentModifier) AddMetadataForPatterns(ctx context.Context, contentBytes *[]byte) error {
	var patternsPage models.MesheryPatternPage
	err := json.Unmarshal(*contentBytes, &patternsPage)
	if err != nil {
		return err
	}
	patterns := patternsPage.Patterns
	var patternsPageMap map[string]interface{}
	patternsPageBytes, err := json.Marshal(patternsPage)
	if err != nil {
		return err
	}
	err = json.Unmarshal(patternsPageBytes, &patternsPageMap)
	if err != nil {
		return err
	}
	p := make([]map[string]interface{}, len(patterns))
	var wg sync.WaitGroup
	for i, pattern := range patterns {
		wg.Add(1)
		go func(pattern *models.MesheryPattern, i int, p *[]map[string]interface{}, token string, provider models.Provider, prefObj *models.Preference, uid string) {
			defer wg.Done()
			patterncontent := pattern.PatternFile
			temp, err := json.Marshal(pattern)
			if err != nil {
				return
			}
			err = json.Unmarshal(temp, &(*p)[i])
			if err != nil {
				return
			}
			msg, ok := mc.isPatternSupported(ctx, patterncontent)
			(*p)[i]["canSupport"] = ok
			(*p)[i]["errmsg"] = msg
		}(pattern, i, &p, mc.token, mc.provider, mc.prefObj, mc.userID)
	}
	wg.Wait()
	patternsPageMap["patterns"] = p
	*contentBytes, err = json.Marshal(patternsPageMap)
	if err != nil {
		return err
	}
	return err
}

// takes a patternfile and returns the status of its current support by using dry run
func (mc *ContentModifier) isPatternSupported(ctx context.Context, patternfile string) (msg string, ok bool) {
	var pattern map[string]interface{}
	err := yaml.Unmarshal([]byte(patternfile), &pattern)
	if err != nil {
		return err.Error(), false
	}
	patternFile, err := core.NewPatternFile([]byte(patternfile))
	if err != nil {
		return err.Error(), false
	}

	msg, err = _processPattern(
		ctx,
		mc.provider,
		patternFile,
		mc.prefObj,
		mc.userID,
		false,
		true,
		true,
		nil,
	)
	if err != nil {
		return err.Error(), false
	}

	return msg, true
}
