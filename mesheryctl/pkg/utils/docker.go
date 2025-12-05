package utils

import (
	"os/exec"
	"sync"
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
	//else v1
    return exec.Command("docker-compose", args...)
}