package pattern

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/oci"
	location "github.com/layer5io/meshkit/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	save string
	wd string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export pattern",
	Long:  `Export pattern to OCI format`,
	Args:  cobra.MinimumNArgs(1),
	Example: `
// export a pattern to default path
mesheryctl pattern export [pattern-name | ID]
mesheryctl pattern export "my meshery design"

// export pattern to the path that given by user
mesheryctl pattern export [pattern-name | ID] -o [path]
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
		pattern := ""
		isID := false

		if len(args) > 0 {
			pattern, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "pattern")
			if err != nil {
				utils.Log.Error(ErrPatternInvalidNameOrID(err))
				return nil
			}
		}

		urlString := mctlCfg.GetBaseMesheryURL()
		if isID {
			// if pattern is a valid uuid, then directly fetch the pattern
			urlString += "/api/pattern/" + pattern
		} else {
			// else search pattern by name
			urlString += "/api/pattern?search=" + url.QueryEscape(pattern)
		}

		req, err := utils.NewRequest("GET", urlString, nil)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		var PatternContent,PatternName string
		if isID {
			var patterns models.MesheryPattern
			err = json.Unmarshal(body, &patterns)
			if err != nil {
				utils.Log.Error(utils.ErrUnmarshal(err))
				return nil
			}
			PatternContent = patterns.PatternFile
			PatternName = patterns.Name
		} else {
			var patterns models.PatternsAPIResponse
			err = json.Unmarshal(body, &patterns)
			if err != nil {
				utils.Log.Error(utils.ErrUnmarshal(err))
				return nil
			}
			for _, v := range patterns.Patterns {
				PatternContent = v.PatternFile
				PatternName = v.Name
			}
		}

		tmpDir, err := oci.CreateTempOCIContentDir()
		if err != nil {
			utils.Log.Error(handlers.ErrCreateDir(err, "OCI"))
			return nil
		}
		defer os.RemoveAll(tmpDir)

		tmpDesignFile := filepath.Join(tmpDir, PatternName+".yaml")
		file, err := os.Create(tmpDesignFile)
		if err != nil {
			utils.Log.Error(handlers.ErrCreateFile(err, tmpDesignFile))
			return nil
		}
		defer file.Close()

		patternReader := strings.NewReader(PatternContent)
		ymlDesign, err := io.ReadAll(patternReader)
		if err != nil {
			utils.Log.Error(handlers.ErrIOReader(err))
			return nil
		}

		if _, err := file.Write(ymlDesign); err != nil {
			utils.Log.Error(handlers.ErrWritingIntoFile(err, tmpDesignFile))
			return nil
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

		pretifiedName := strings.ToLower(strings.Replace(PatternName, " ", "", -1))
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
		fmt.Printf("OCI Image %s is successfully built on %s \nDigest: %v, Size: %v",pretifiedName, save,digest, size)
		return nil
	},
}

func init(){
	wd = filepath.Join(location.GetHome(), ".meshery", "content")
	exportCmd.Flags().StringVarP(&save,"output","o",wd,"location for exported oci file")
}
