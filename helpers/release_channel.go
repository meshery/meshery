package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GitTag struct holds the metadata associated with a git tag
type GitTag struct {
	Ref    string       `json:"ref,omitempty"`
	NodeID string       `json:"node_id,omitempty"`
	URL    string       `json:"url,omitempty"`
	Object GitTagObject `json:"object,omitempty"`
}

// GitTagObject holds the metadata associated with a git tag object
type GitTagObject struct {
	SHA  string `json:"sha,omitempty"`
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

// GetReleaseChannel takes in the tag name and current commit sha
// and returns the release channel
//
// To determine the channel it queries the github api to get the original commit sha
// associated with the tag, if the current commit sha != commit sha associated with the tag
// it will return "edge" or else "stable"
func GetReleaseChannel(build, commitsha string) (string, error) {
	tag, err := GetGitTag("layer5io", "meshery", build)
	if err != nil {
		return "NA", errors.Wrap(err, "couldn't get the release channel")
	}

	if tag.Object.SHA == commitsha {
		return "stable", nil
	}

	return "edge", nil
}

// GetGitTag uses github api to get the metadata associated with the
// given tag
//
// this approach is chosen over using
// 	git rev-list -1 $TAG
// under the assumption that meshery server may or may not have access to
// .git directory (if server is running within a docker container)
func GetGitTag(org, repo, tag string) (GitTag, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/tags/%s", org, repo, tag)

	resp, err := http.Get(url)
	if err != nil {
		return GitTag{}, errors.Wrap(err, "failed to query github api")
	}

	if resp.StatusCode != http.StatusOK {
		return GitTag{}, fmt.Errorf("failed to query github api, status code: %d", resp.StatusCode)
	}

	gt := GitTag{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&gt); err != nil {
		return GitTag{}, errors.Wrap(err, "failed to parse github api response, response could be malformed")
	}

	return gt, nil
}
