package utils

type ManifestList struct {
	SHA       string     `json:"sha"`
	URL       string     `json:"url"`
	Tree      []Manifest `json:"tree"`
	Truncated bool       `json:"truncated"`
}

type Manifest struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Typ  string `json:"type"`
	SHA  string `json:"sha"`
	Size string `json:"size"`
	URL  string `json:"url"`
}
