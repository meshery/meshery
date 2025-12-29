// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	dockerCmd "github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

// ComposeClient is a wrapper around the docker compose library
type ComposeClient struct {
	service api.Compose
	cli     *dockerCmd.DockerCli
}

// NewComposeClient creates a new ComposeClient instance
func NewComposeClient() (*ComposeClient, error) {
	// Create Docker CLI with combined streams going to discard
	dockerCli, err := dockerCmd.NewDockerCli(
		dockerCmd.WithCombinedStreams(io.Discard),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker cli: %w", err)
	}

	// Initialize the Docker CLI with default options
	clientOpts := cliflags.NewClientOptions()
	err = dockerCli.Initialize(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize docker cli: %w", err)
	}

	// Create compose service
	service := compose.NewComposeService(dockerCli)

	return &ComposeClient{
		service: service,
		cli:     dockerCli,
	}, nil
}

// LoadProject loads a compose project from the given file path
func LoadProject(ctx context.Context, composefile string) (*types.Project, error) {
	projectOptions, err := cli.NewProjectOptions(
		[]string{composefile},
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithName("meshery"), // Explicitly set project name to match docker-compose behavior
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project options: %w", err)
	}

	project, err := cli.ProjectFromOptions(ctx, projectOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	for name, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  project.WorkingDir,
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False",
		}
		project.Services[name] = s
	}

	return project, nil
}

// Up starts the containers defined in the compose file (equivalent to docker-compose up -d)
func (c *ComposeClient) Up(ctx context.Context, composefile string) error {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return err
	}

	upOptions := api.UpOptions{
		Create: api.CreateOptions{
			RemoveOrphans: false,
			Recreate:      api.RecreateDiverged, // Recreate containers if config diverged
		},
		Start: api.StartOptions{
			Project: project, // Pass the project to avoid rebuilding from containers
			Attach:  nil,
		},
	}

	return c.service.Up(ctx, project, upOptions)
}

// Stop stops containers (equivalent to docker-compose stop)
func (c *ComposeClient) Stop(ctx context.Context, composefile string) error {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return err
	}

	stopOptions := api.StopOptions{}

	return c.service.Stop(ctx, project.Name, stopOptions)
}

// Down stops and removes containers (equivalent to docker-compose rm -f)
func (c *ComposeClient) Down(ctx context.Context, composefile string) error {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return err
	}

	downOptions := api.DownOptions{
		RemoveOrphans: false,
	}

	return c.service.Down(ctx, project.Name, downOptions)
}

// Remove removes stopped containers (equivalent to docker-compose rm -f)
func (c *ComposeClient) Remove(ctx context.Context, composefile string) error {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return err
	}

	removeOptions := api.RemoveOptions{
		Force: true,
	}

	return c.service.Remove(ctx, project.Name, removeOptions)
}

// Ps returns status of containers (equivalent to docker-compose ps)
func (c *ComposeClient) Ps(ctx context.Context, composefile string) ([]api.ContainerSummary, error) {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return nil, err
	}

	psOptions := api.PsOptions{
		All: true,
	}

	return c.service.Ps(ctx, project.Name, psOptions)
}

// Pull pulls images for services (equivalent to docker-compose pull)
func (c *ComposeClient) Pull(ctx context.Context, composefile string) error {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return err
	}

	pullOptions := api.PullOptions{}

	return c.service.Pull(ctx, project, pullOptions)
}

// LogConsumer implements the api.LogConsumer interface
type LogConsumer struct {
	Out io.Writer
	mu  sync.Mutex
}

// Log implements api.LogConsumer
func (l *LogConsumer) Log(containerName, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.Out, "%s | %s\n", containerName, message)
}

// Err implements api.LogConsumer
func (l *LogConsumer) Err(containerName, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.Out, "%s | %s\n", containerName, message)
}

// Status implements api.LogConsumer
func (l *LogConsumer) Status(containerName, msg string) {
	// Skip status messages
}

// Register implements api.LogConsumer
func (l *LogConsumer) Register(containerName string) {
	// Skip registration
}

// Logs streams logs from containers (equivalent to docker-compose logs -f)
func (c *ComposeClient) Logs(ctx context.Context, composefile string, follow bool, out io.Writer) error {
	project, err := LoadProject(ctx, composefile)
	if err != nil {
		return err
	}

	logConsumer := &LogConsumer{Out: out}

	logOptions := api.LogOptions{
		Follow: follow,
	}

	return c.service.Logs(ctx, project.Name, logConsumer, logOptions)
}

// IsMesheryRunning checks if meshery containers are running using the compose library
func (c *ComposeClient) IsMesheryRunning(ctx context.Context, composefile string) (bool, error) {
	containers, err := c.Ps(ctx, composefile)
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		if container.Name == "meshery" || container.Service == "meshery" {
			return container.State == "running", nil
		}
	}

	return false, nil
}

// GetPsOutput returns a formatted string of container status (similar to docker-compose ps output)
func (c *ComposeClient) GetPsOutput(ctx context.Context, composefile string) (string, error) {
	containers, err := c.Ps(ctx, composefile)
	if err != nil {
		return "", err
	}

	if len(containers) == 0 {
		return "", nil
	}

	output := "NAME\tIMAGE\tCOMMAND\tSERVICE\tCREATED\tSTATUS\tPORTS\n"
	for _, container := range containers {
		ports := ""
		for _, p := range container.Publishers {
			if ports != "" {
				ports += ", "
			}
			ports += fmt.Sprintf("%s:%d->%d/%s", p.URL, p.PublishedPort, p.TargetPort, p.Protocol)
		}
		// Convert Unix timestamp to time.Time for formatting
		createdTime := time.Unix(container.Created, 0)
		output += fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			container.Name,
			container.Image,
			container.Command,
			container.Service,
			createdTime.Format(time.RFC3339),
			container.State,
			ports,
		)
	}

	return output, nil
}

// ContainsMesheryContainer checks if the containers include any meshery-related container
func ContainsMesheryContainer(containers []api.ContainerSummary) bool {
	for _, container := range containers {
		if strings.Contains(container.Name, "meshery") || strings.Contains(container.Service, "meshery") {
			return true
		}
	}
	return false
}
