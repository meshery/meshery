package models

import (
	"encoding/gob"

	"github.com/grafana-tools/sdk"
)

type K8SConfig struct {
	InClusterConfig   bool   `json:"inClusterConfig,omitempty"`
	K8Sfile           string `json:"k8sfile,omitempty"`
	Config            []byte `json:"config,omitempty"`
	Server            string `json:"configuredServer,omitempty"`
	ContextName       string `json:"contextName,omitempty"`
	ClusterConfigured bool   `json:"clusterConfigured,omitempty"`
	// ConfiguredServer  string `json:"configuredServer,omitempty"`
}

// just used to send contexts to the UI
type K8SContext struct {
	ContextName string `json:"context-name"`
	ClusterName string `json:"cluster-name"`
	// ContextDisplayName string `json:"context-display-name"`
	IsCurrentContext bool `json:"current-context"`
}

type Grafana struct {
	GrafanaURL    string `json:"grafanaURL,omitempty"`
	GrafanaAPIKey string `json:"grafanaAPIKey,omitempty"`
	// GrafanaBoardSearch string          `json:"grafanaBoardSearch,omitempty"`
	GrafanaBoards []*SelectedGrafanaConfig `json:"selectedBoardsConfigs,omitempty"`
}

type SelectedGrafanaConfig struct {
	GrafanaBoard         *GrafanaBoard `json:"board,omitempty"`
	GrafanaPanels        []*sdk.Panel  `json:"panels,omitempty"`
	SelectedTemplateVars []string      `json:"templateVars,omitempty"`
}

type Prometheus struct {
	PrometheusURL                   string                   `json:"prometheusURL,omitempty"`
	SelectedPrometheusBoardsConfigs []*SelectedGrafanaConfig `json:"selectedPrometheusBoardsConfigs,omitempty"`
}

type Session struct {
	// User         *User       `json:"user,omitempty"`
	K8SConfig    *K8SConfig  `json:"k8sConfig,omitempty"`
	MeshAdapters []*Adapter  `json:"meshAdapters,omitempty"`
	Grafana      *Grafana    `json:"grafana,omitempty"`
	Prometheus   *Prometheus `json:"prometheus,omitempty"`
}

func init() {
	gob.Register(&Session{})
	gob.Register(map[string]interface{}{})
}

type SessionPersister interface {
	Read(userId string) (*Session, error)
	Write(userId string, data *Session) error
	Delete(userId string) error

	// Lock(userId string)
	// Unlock(userId string)
	Close()
}
