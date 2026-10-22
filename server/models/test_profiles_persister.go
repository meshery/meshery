package models

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

// TestProfilesPersister assists with persisting session in store
type TestProfilesPersister struct {
	DB *database.Handler
}

type PerformanceTestConfig struct {
	ID                         core.Uuid
	PerformanceTestConfigBytes []byte
	UpdatedAt                  time.Time
}

type legacyPerformanceTestConfig struct {
	SmpVersion string                        `json:"smp_version"`
	ID         string                        `json:"id"`
	Name       string                        `json:"name"`
	Labels     map[string]string             `json:"labels"`
	Clients    []legacyPerformanceTestClient `json:"clients"`
	Duration   string                        `json:"duration"`
}

type legacyPerformanceTestClient struct {
	Internal          bool              `json:"internal"`
	LoadGenerator     string            `json:"load_generator"`
	Protocol          json.RawMessage   `json:"protocol"`
	Connections       int               `json:"connections"`
	Rps               int               `json:"rps"`
	Headers           map[string]string `json:"headers"`
	Cookies           map[string]string `json:"cookies"`
	Body              string            `json:"body"`
	ContentType       string            `json:"content_type"`
	EndpointUrls      []string          `json:"endpoint_urls"`
	SslCertificate    string            `json:"ssl_certificate"`
	AdditionalOptions string            `json:"additional_options"`
}

// decodePerformanceTestConfig accepts both the current schema-native wire shape
// and the legacy SMP JSON shape already stored in performance_test_configs. This
// prevents saved test profiles from disappearing after the SMP dependency was
// replaced by github.com/meshery/schemas performance profile models.
func decodePerformanceTestConfig(data []byte) (*perfprofile.PerformanceTestConfig, error) {
	if hasLegacyPerformanceTestConfigShape(data) {
		if config, err := decodeLegacyPerformanceTestConfig(data); err == nil {
			return config, nil
		}
	}

	testConfig := &perfprofile.PerformanceTestConfig{}
	if err := json.Unmarshal(data, testConfig); err != nil {
		if config, legacyErr := decodeLegacyPerformanceTestConfig(data); legacyErr == nil {
			return config, nil
		}
		return nil, err
	}
	return testConfig, nil
}

// hasLegacyPerformanceTestConfigShape cheaply detects persisted SMP JSON before
// unmarshalling into the new schema type. encoding/json does not fail on unknown
// snake_case fields, so without this check old configs would decode with empty
// client fields instead of being normalized.
func hasLegacyPerformanceTestConfigShape(data []byte) bool {
	var raw struct {
		SmpVersion json.RawMessage              `json:"smp_version"`
		Clients    []map[string]json.RawMessage `json:"clients"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return false
	}
	if len(raw.SmpVersion) > 0 {
		return true
	}
	legacyClientKeys := []string{
		"load_generator",
		"content_type",
		"endpoint_urls",
		"ssl_certificate",
		"additional_options",
	}
	for _, client := range raw.Clients {
		for _, key := range legacyClientKeys {
			if _, ok := client[key]; ok {
				return true
			}
		}
	}
	return false
}

// decodeLegacyPerformanceTestConfig maps old SMP field names to the canonical
// schema-native names used by perfprofile.PerformanceTestConfig.
func decodeLegacyPerformanceTestConfig(data []byte) (*perfprofile.PerformanceTestConfig, error) {
	var legacy legacyPerformanceTestConfig
	if err := json.Unmarshal(data, &legacy); err != nil {
		return nil, err
	}
	config := &perfprofile.PerformanceTestConfig{
		SmpVersion: legacy.SmpVersion,
		ID:         legacy.ID,
		Name:       legacy.Name,
		Labels:     legacy.Labels,
		Duration:   legacy.Duration,
		Clients:    make([]perfprofile.PerformanceTestClient, 0, len(legacy.Clients)),
	}
	for _, client := range legacy.Clients {
		config.Clients = append(config.Clients, perfprofile.PerformanceTestClient{
			Internal:          client.Internal,
			LoadGenerator:     client.LoadGenerator,
			Protocol:          decodeLegacyProtocol(client.Protocol),
			Connections:       client.Connections,
			Rps:               client.Rps,
			Headers:           client.Headers,
			Cookies:           client.Cookies,
			Body:              client.Body,
			ContentType:       client.ContentType,
			EndpointUrls:      client.EndpointUrls,
			SslCertificate:    client.SslCertificate,
			AdditionalOptions: client.AdditionalOptions,
		})
	}
	return config, nil
}

// decodeLegacyProtocol converts old SMP protocol encodings to the string values
// expected by the schema model. Stored configs may contain protobuf enum numbers
// from encoding/json or enum names from externally-authored files.
func decodeLegacyProtocol(raw json.RawMessage) string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return ""
	}
	var protocolNumber int
	if err := json.Unmarshal(raw, &protocolNumber); err == nil {
		switch protocolNumber {
		case 1:
			return "http"
		case 2:
			return "tcp"
		case 3:
			return "udp"
		case 4:
			return "grpc"
		default:
			return ""
		}
	}
	var protocol string
	if err := json.Unmarshal(raw, &protocol); err != nil {
		return ""
	}
	protocol = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(protocol)), "protocol_")
	if protocolNumber, err := strconv.Atoi(protocol); err == nil {
		return decodeLegacyProtocol(json.RawMessage(strconv.Itoa(protocolNumber)))
	}
	if protocol == "invalid" {
		return ""
	}
	return protocol
}

// UserTestProfiles - represents a page of user test configs
type UserTestProfiles struct {
	Page        uint64                               `json:"page"`
	PageSize    uint64                               `json:"pageSize"`
	TotalCount  int                                  `json:"totalCount"`
	TestConfigs []*perfprofile.PerformanceTestConfig `json:"testConfigs"`
}

// GetTestConfigs - gets result for the page and pageSize
func (s *TestProfilesPersister) GetTestConfigs(page, pageSize uint64) ([]byte, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}
	order := defaultOrderUpdatedAtDesc
	query := s.DB.Order(order)
	total := int64(0)
	s.DB.Model(&PerformanceTestConfig{}).Count(&total)
	testConfigs := []*perfprofile.PerformanceTestConfig{}
	var p []*PerformanceTestConfig
	Paginate(uint(page), uint(pageSize))(query).Find(&p)
	for _, config := range p {
		testConfig, err := decodePerformanceTestConfig(config.PerformanceTestConfigBytes)
		if err != nil {
			return nil, err
		}
		testConfigs = append(testConfigs, testConfig)
	}
	bd, err := json.Marshal(&UserTestProfiles{
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  int(total),
		TestConfigs: testConfigs,
	})
	if err != nil {
		return nil, ErrMarshal(err, "result data")
	}

	return bd, nil
}

// GetTestConfig - gets result for a specific key
func (s *TestProfilesPersister) GetTestConfig(key core.Uuid) (*perfprofile.PerformanceTestConfig, error) {
	if s.DB == nil {
		return nil, ErrDBConnection
	}
	var u PerformanceTestConfig
	err := s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).First(&u).Error
	if err != nil {
		return nil, err
	}
	testConfig, err := decodePerformanceTestConfig(u.PerformanceTestConfigBytes)
	if err != nil {
		return nil, err
	}
	return testConfig, nil
}

// DeleteTestConfig - delete result for a specific key
func (s *TestProfilesPersister) DeleteTestConfig(key core.Uuid) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	return s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).Delete(&PerformanceTestConfig{}).Error
}

// WriteTestConfig persists the result
func (s *TestProfilesPersister) WriteTestConfig(key core.Uuid, result []byte) error {
	if s.DB == nil {
		return ErrDBConnection
	}

	if result == nil {
		return ErrResultData()
	}
	var p PerformanceTestConfig
	if err := s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).First(&p).Error; err == nil {
		return s.UpdateTestConfig(key, p)
	}
	p.ID = key
	p.PerformanceTestConfigBytes = result
	return s.DB.Model(&PerformanceTestConfig{}).Create(&p).Error
}

func (s *TestProfilesPersister) UpdateTestConfig(key core.Uuid, p PerformanceTestConfig) error {
	return s.DB.Model(&PerformanceTestConfig{}).Where("id = ?", key).UpdateColumns(p).Error
}
