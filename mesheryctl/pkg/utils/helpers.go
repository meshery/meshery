package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	dockerComposeWebURL         = "https://api.github.com/repos/docker/compose/releases/latest"
	defaultDockerComposeVersion = "1.24.1/docker-compose"
	dockerComposeBinaryURL      = "https://github.com/docker/compose/releases/download/"
	dockerComposeBinary         = "/usr/local/bin/docker-compose"
)

var (
	// ResetFlag indicates if a reset is required
	ResetFlag bool
	// MesharyFolder is the default relative location of the meshery config
	// related configuration files.
	MesheryFolder = ".meshery"
	// DockerComposeFile is the default location within the MesheryFolder
	// where the docker compose file is located?
	DockerComposeFile = "/meshery.yaml"
	// AuthoConfigFile is the location of the auth file for performing perf testing
	AuthConfigFile = "/auth.json"
)

// SafeClose is a helper function help to close the io
func SafeClose(co io.Closer) {
	if cerr := co.Close(); cerr != nil {
		log.Error(cerr)
	}
}

// DownloadFile from url and save to configured file location
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func prereq() ([]byte, []byte) {
	ostype, err := exec.Command("uname", "-s").Output()
	if err != nil {
		log.Fatal("[ERROR] Please, install docker-compose. The error message: \n", err)
	}
	//fmt.Printf("%s\n", ostype)

	osarch, err := exec.Command("uname", "-m").Output()
	if err != nil {
		log.Fatal("[ERROR] Please, install docker-compose. The error message: \n", err)
	}
	//	fmt.Printf("%s\n", arch)
	return ostype, osarch
}

// SetFileLocation to set absolute path
func SetFileLocation() {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("[ERROR] Cannot determine location of $HOME")
	}
	MesheryFolder = path.Join(home, MesheryFolder)
	DockerComposeFile = path.Join(MesheryFolder, DockerComposeFile)
	AuthConfigFile = path.Join(MesheryFolder, AuthConfigFile)
}

//PreReqCheck prerequisites check
func PreReqCheck() {
	//Check for installed docker-compose on client system
	if err := exec.Command("docker-compose", "-v").Run(); err != nil {
		log.Info("Docker-Compose is not installed")
		//No auto installation of Docker-compose for windows
		if runtime.GOOS == "windows" {
			return
		}
		installprereq()
	}
}

func installprereq() {
	log.Info("Attempting Docker-Compose installation...")
	ostype, osarch := prereq()
	osdetails := strings.TrimRight(string(ostype), "\r\n") + "-" + strings.TrimRight(string(osarch), "\r\n")

	//checks for the latest docker-compose
	resp, err := http.Get(dockerComposeWebURL)
	dockerComposeBinaryURL := dockerComposeBinaryURL
	if err != nil {
		// download the default version as 1.24.1 if unable to fetch latest page data
		dockerComposeBinaryURL = dockerComposeBinaryURL + defaultDockerComposeVersion
	} else {
		var dat map[string]interface{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(body, &dat); err != nil {
			log.Fatal(err)
		}
		num := dat["tag_name"]
		dockerComposeBinaryURL = fmt.Sprintf(dockerComposeBinaryURL+"%v/docker-compose", num)
	}
	dockerComposeBinaryURL = dockerComposeBinaryURL + "-" + osdetails
	if err := DownloadFile(dockerComposeBinary, dockerComposeBinaryURL); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("chmod", "+x", dockerComposeBinary).Run(); err != nil {
		log.Fatal(err)
	}
	log.Info("Prerequisite Docker Compose is installed.")
}

func IsMesheryRunning() bool {
	op, err := exec.Command("docker-compose", "-f", DockerComposeFile, "ps").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(op), "meshery")
}
