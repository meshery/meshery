package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

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
	// MesheryFolder is the default relative location of the meshery config
	// related configuration files.
	MesheryFolder = ".meshery"
	// DockerComposeFile is the default location within the MesheryFolder
	// where the docker compose file is located?
	DockerComposeFile = "/meshery.yaml"
	// AuthConfigFile is the location of the auth file for performing perf testing
	AuthConfigFile = "/auth.json"
)

const tokenName = "token"
const providerName = "meshery-provider"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset generates a random string with a given length
func StringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	// + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

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

// IsMesheryRunning checks if the meshery server containers are up and running
func IsMesheryRunning() bool {
	op, err := exec.Command("docker-compose", "-f", DockerComposeFile, "ps").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(op), "meshery")
}

// AddAuthDetails Adds authentication cookies to the request
func AddAuthDetails(req *http.Request, filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Errorf("File read failed : %v", err.Error())
		return err
	}
	var tokenObj map[string]string
	if err := json.Unmarshal(file, &tokenObj); err != nil {
		log.Errorf("Token file invalid : %v", err.Error())
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:     tokenName,
		Value:    tokenObj[tokenName],
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name:     providerName,
		Value:    tokenObj[providerName],
		HttpOnly: true,
	})
	return nil
}

// UpdateAuthDetails checks gets the token (old/refreshed) from meshery server and writes it back to the config file
func UpdateAuthDetails(filepath string) error {
	// TODO: get this from the global config
	req, err := http.NewRequest("GET", "http://localhost:9081/api/user", bytes.NewBuffer([]byte("")))
	if err != nil {
		return err
	}
	if err := AddAuthDetails(req, filepath); err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer SafeClose(resp.Body)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, data, os.ModePerm)
}

// UploadFileWithParams returns a request configured to upload files with other values
func UploadFileWithParams(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if err = file.Close(); err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileContents)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}
