package model

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
	meshkitutils "github.com/layer5io/meshkit/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importModelCmd = &cobra.Command{
	Use:   "import",
	Short: "import models from mesheryctl command",
	Long:  "import model by specifying the directory, file. Use 'import model [filepath]' or 'import model  [directory]'.",
	Example: `
	import model  /path/to/[file.yaml|file.json]
	import model  /path/to/models
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisites
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model import [ file | filePath ]\nRun 'mesheryctl model import --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("[ file | filepath ] isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		fileInfo, err := os.Stat(path)
		if err != nil {
			utils.Log.Error(meshkitutils.ErrReadDir(err, path))
			return err
		}
		if fileInfo.IsDir() {
			tarData, err := compressDirectory(path)
			if err != nil {
				return err
			}
			fileName := filepath.Base(path) + ".tar.gz"
			err = sendToAPI(tarData, fileName, "dir")
			if err != nil {
				utils.Log.Error(err)
				return err
			}
		} else {
			if meshkitutils.IsYaml(path) {
				fileData, err := os.ReadFile(path)
				if err != nil {
					utils.Log.Error(meshkitutils.ErrReadFile(err, path))
					return nil
				}
				err = sendToAPI(fileData, path, "file")
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			} else if meshkitutils.IsTarGz(path) || meshkitutils.IsZip(path) {
				fileData, err := os.ReadFile(path)
				if err != nil {
					err = meshkitutils.ErrReadFile(err, path)
					utils.Log.Error(err)
					return nil
				}
				err = sendToAPI(fileData, path, "dir")
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			} else {
				err = utils.ErrInvalidFile(errors.New("invalid file format"))
				utils.Log.Error(err)
				return nil
			}
		}
		return nil
	},
}

func compressDirectory(dirpath string) ([]byte, error) {
	tw := meshkitutils.NewTarWriter()
	defer tw.Close()

	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return meshkitutils.ErrFileWalkDir(err, path)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return handlers.ErrOpenFile(path)
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return meshkitutils.ErrReadFile(err, path)
		}

		relPath, err := filepath.Rel(filepath.Dir(dirpath), path)
		if err != nil {
			return meshkitutils.ErrRelPath(err, path)
		}

		if err := tw.Compress(relPath, fileData); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err = io.Copy(gzipWriter, tw.Buffer)
	if err != nil {
		return nil, meshkitutils.ErrCopyFile(err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, meshkitutils.ErrCloseFile(err)
	}

	return buf.Bytes(), nil
}

func sendToAPI(data []byte, name string, dataType string) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}

	baseURL := mctlCfg.GetBaseMesheryURL()
	url := baseURL + "/api/meshmodels/register"
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	var formFile io.Writer
	if dataType == "file" {
		formFile, _ = writer.CreateFormFile("file", filepath.Base(name))
	} else {
		formFile, _ = writer.CreateFormField("dir")
	}

	_, err = formFile.Write(data)
	if err != nil {
		err = meshkitutils.ErrWriteFile(err, name)
		return err
	}

	_ = writer.Close()

	req, err := utils.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := utils.MakeRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = models.ErrDoRequest(err, resp.Request.Method, url)
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = models.ErrDataRead(err, "response body")
		return err
	}

	bodyString := string(bodyBytes)
	utils.Log.Info(bodyString)
	utils.Log.Info("Models imported successfully")
	return nil
}
