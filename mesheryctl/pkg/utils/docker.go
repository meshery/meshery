package utils

import (
	"os/exec"
	"sync"
	log "github.com/sirupsen/logrus"
)

var (
	dockerComposeCheck sync.Once
	useDockerComposeV2 bool
)

func DockerComposeCmd(args ...string) *exec.Cmd {
	// Try v2 first
	dockerComposeCheck.Do(func() {
		cmd := exec.Command("docker", "compose", "version")
		useDockerComposeV2 = cmd.Run() == nil
	})

	if useDockerComposeV2 {
		return exec.Command("docker", append([]string{"compose"}, args...)...)
	}

	// Check if v1 is available
	if exec.Command("docker-compose", "version").Run() == nil {
		log.Fatal("Docker Compose v1 detected. Meshery requires Docker Compose v2. Please upgrade: https://docs.docker.com/compose/install/")
	}

	// Neither v1 nor v2 available
	return exec.Command("docker-compose", args...)
}