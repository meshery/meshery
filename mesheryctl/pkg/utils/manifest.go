package utils

import "encoding/json"

type ManifestList struct {
	SHA       string     `json:"sha,omitempty"`
	URL       string     `json:"url,omitempty"`
	Tree      []Manifest `json:"tree,omitempty"`
	Truncated bool       `json:"truncated,omitempty"`
}

type Manifest struct {
	Path string      `json:"path,omitempty"`
	Mode string      `json:"mode,omitempty"`
	Typ  string      `json:"type,omitempty"`
	SHA  string      `json:"sha,omitempty"`
	Size json.Number `json:"size,omitempty"`
	URL  string      `json:"url,omitempty"`
}
