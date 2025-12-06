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

func DockerComposeCommand(args ...string) *exec.Cmd {
	dockerComposeCheck.Do(func() {
		// Try v2 first
		cmd := exec.Command("docker", "compose", "version")
		if cmd.Run() == nil {
			useDockerComposeV2 = true
			return
		}

		// Check if v1 is available
		if exec.Command("docker-compose", "version").Run() == nil {
			log.Fatal("Docker Compose v1 detected. Meshery requires Docker Compose v2. Please upgrade: https://docs.docker.com/compose/install/")
		}

		// Neither v1 nor v2 available
		log.Fatal("Docker Compose not found. Meshery requires Docker Compose v2. Please install: https://docs.docker.com/compose/install/")
	})

	if useDockerComposeV2 {
		return exec.Command("docker", append([]string{"compose"}, args...)...)
	}

	// Should never reach here due to log.Fatal above
	log.Fatal("Docker Compose not found.")
	return nil
}