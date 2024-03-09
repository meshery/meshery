package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/models/oci"
	location "github.com/layer5io/meshkit/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	save string
	wd   string
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export model",
	Long:  `Export model to OCI format`,
	Args:  cobra.MinimumNArgs(1),
	Example: `
// export a model to default path
mesheryctl system model export [model-name]
mesheryctl system model export "my meshery design"

// export model to the path that given by user
mesheryctl system model export [model-name] -o [path]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		model := args[0]
		urlString := mctlCfg.GetBaseMesheryURL()

		url := fmt.Sprintf("%s/api/meshmodels/models/%s?pagesize=all", urlString, model)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		var ModelContent v1alpha1.Model

		if modelsResponse.Count == 0 {
			fmt.Println("No model(s) found for the given name ", model)
			return nil
		} else if modelsResponse.Count == 1 {
			ModelContent = modelsResponse.Models[0]
		} else {
			ModelContent = selectModelPrompt(modelsResponse.Models)
		}

		ModelName := ModelContent.Name
		ModelVersion := ModelContent.Version
		output,err := json.Marshal(ModelContent)
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		strcontent := fmt.Sprintln(output)

		tmpDir, err := oci.CreateTempOCIContentDir()
		if err != nil {
			utils.Log.Error(handlers.ErrCreateDir(err, "OCI"))
			return nil
		}
		defer os.RemoveAll(tmpDir)

		tmpDesignFile := filepath.Join(tmpDir, ModelName+".yaml")
		file, err := os.Create(tmpDesignFile)
		if err != nil {
			utils.Log.Error(handlers.ErrCreateFile(err, tmpDesignFile))
			return nil
		}
		defer file.Close()

		modelReader := strings.NewReader(strcontent)
		ymlDesign, err := io.ReadAll(modelReader)
		if err != nil {
			utils.Log.Error(handlers.ErrIOReader(err))
			return err
		}

		if _, err := file.Write(ymlDesign); err != nil {
			utils.Log.Error(handlers.ErrWritingIntoFile(err, tmpDesignFile))
			return err
		}

		ociImg, err := oci.BuildImage(tmpDir)
		if err != nil {
			utils.Log.Error(handlers.ErrBuildOCIImg(err))
			return nil
		}

		digest, err := ociImg.Digest()
		if err != nil {
			utils.Log.Error(handlers.ErrBuildOCIImg(err))
			return nil
		}

		size, err := ociImg.Size()
		if err != nil {
			utils.Log.Error(handlers.ErrBuildOCIImg(err))
			return nil
		}

		FullName := ModelName+ModelVersion
		pretifiedName := strings.ToLower(strings.Replace(FullName, " ", "", -1))
		exts := filepath.Ext(pretifiedName)
		if exts != "" {
			pretifiedName = strings.TrimSuffix(pretifiedName, exts)
		}
		tmpOCITarFilePath := filepath.Join(tmpDir, pretifiedName+".tar")
		err = oci.SaveOCIArtifact(ociImg, tmpOCITarFilePath, pretifiedName)
		if err != nil {
			utils.Log.Error(handlers.ErrSaveOCIArtifact(err))
			return nil
		}

		file, err = os.OpenFile(tmpOCITarFilePath, os.O_RDONLY, 0444)
		if err != nil {
			utils.Log.Error(handlers.ErrOpenFile(tmpOCITarFilePath))
			return nil
		}

		content, err := io.ReadAll(file)
		if err != nil {
			utils.Log.Error(handlers.ErrIOReader(err))
			return nil
		}

		err = os.WriteFile(save+"/"+pretifiedName, content, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write oci image data to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("OCI Image %s is successfully built on %s \nDigest: %v, Size: %v", pretifiedName, save, digest, size)
		return nil
	},
}

func init() {
	wd = filepath.Join(location.GetHome(), ".meshery", "content")
	ExportCmd.Flags().StringVarP(&save, "output", "o", wd, "location for exported oci file")
}
