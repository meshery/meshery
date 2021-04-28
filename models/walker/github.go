package walker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type GithubContentAPI struct {
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	SHA         string `json:"sha,omitempty"`
	Size        int64  `json:"size,omitempty"`
	URL         string `json:"url,omitempty"`
	HTMLURL     string `json:"html_url,omitempty"`
	GitURL      string `json:"git_url,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
	Type        string `json:"type,omitempty"`
	Content     string `json:"content,omitempty"`
	Encoding    string `json:"encoding,omitempty"`
}

type GithubDirectoryContentAPI []GithubContentAPI

type GithubFileInterceptor func(GithubContentAPI) error
type GithubDirInterceptor func(GithubDirectoryContentAPI) error

type Github struct {
	owner           string
	repo            string
	branch          string
	root            string
	recurse         bool
	fileInterceptor GithubFileInterceptor
	dirInterceptor  GithubDirInterceptor
}

func NewGithub() *Github {
	return &Github{
		branch: "main",
	}
}

func (g *Github) Owner(owner string) *Github {
	g.owner = owner
	return g
}

func (g *Github) Repo(repo string) *Github {
	g.repo = repo
	return g
}

func (g *Github) Branch(branch string) *Github {
	g.branch = branch
	return g
}

func (g *Github) Root(root string) *Github {
	g.root = root

	if strings.HasSuffix(root, "/**") {
		g.recurse = true
		g.root = strings.TrimSuffix(root, "/**")
	}

	return g
}

func (g *Github) RegisterFileInterceptor(i GithubFileInterceptor) *Github {
	g.fileInterceptor = i
	return g
}

func (g *Github) RegisterDirInterceptor(i GithubDirInterceptor) *Github {
	g.dirInterceptor = i
	return g
}

func (g *Github) Walk() error {
	// Check if a file is requested
	isFile := true
	if g.root == "" || filepath.Ext(g.root) == "" {
		isFile = false
	}

	return g.walker(g.root, isFile)
}

func (g *Github) walker(path string, isFile bool) error {
	githubAPIURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		g.owner,
		g.repo,
		path,
		g.branch,
	)

	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden {
			respJSON := map[string]interface{}{}

			if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
				return fmt.Errorf("GithubAPI responded with: forbidden")
			}

			message, ok := respJSON["message"].(string)
			if !ok {
				return fmt.Errorf("GithubAPI responded with: forbidden")
			}

			return fmt.Errorf(message)
		}
		return fmt.Errorf("file not found")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Error("failed to close response body", err)
		}
	}()

	if isFile {
		respBody := GithubContentAPI{}

		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			logrus.Error("[GithubWalker]: failed to decode API response", err)
			return err
		}

		if g.fileInterceptor != nil {
			return g.fileInterceptor(respBody)
		}

		return nil
	}

	respBody := GithubDirectoryContentAPI{}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		logrus.Error("[GithubWalker]: failed to decode API response", err)
		return err
	}

	var wg sync.WaitGroup
	for _, r := range respBody {
		nextPath := r.Path
		typ := r.Type
		isFile := typ == "file"

		wg.Add(1)
		go func() {
			if g.recurse || isFile {
				if err := g.walker(nextPath, isFile); err != nil {
					logrus.Error("[GithubWalker]: error occured while processing github node ", err)
				}
			}

			if g.dirInterceptor != nil {
				if err := g.dirInterceptor(respBody); err != nil {
					logrus.Error("[GithubWalker]: error occured while executing directory interceptor function ", err)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}
