package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Provider struct {
	ProviderURL  string `json:"provider_url,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
}

// AddAuthDetails Adds authentication cookies to the request
func AddAuthDetails(req *http.Request, filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		err = errors.Wrap(err, "could not read token:")
		return err
	}
	var tokenObj map[string]string
	if err := json.Unmarshal(file, &tokenObj); err != nil {
		err = errors.Wrap(err, "token file invalid:")
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
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// TODO: get this from the global config
	req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/token", bytes.NewBuffer([]byte("")))
	if err != nil {
		err = errors.Wrap(err, "error Creating the request :")
		return err
	}
	if err := AddAuthDetails(req, filepath); err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer SafeClose(resp.Body)

	if err != nil {
		err = errors.Wrap(err, "error dispatching there request :")
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "error reading body :")
		return err
	}

	if ContentTypeIsHTML(resp) {
		return errors.New("invalid body")
	}

	return os.WriteFile(filepath, data, os.ModePerm)
}

// CreateTempAuthServer creates a temporary http server
//
// It implements a custom mux and has a catch all route, the function passed as the
// parameter is binded to the catch all route
func CreateTempAuthServer(fn func(http.ResponseWriter, *http.Request)) (*http.Server, int, error) {
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
	}

	listener, err := net.Listen("tcp4", ":0")
	if err != nil {
		return nil, -1, err
	}

	mux.HandleFunc("/", fn)

	go func() {
		if err := srv.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				log.Println("error creating temporary server")
			}
		}
	}()

	return srv, listener.Addr().(*net.TCPAddr).Port, nil
}

// InitiateLogin initates the login process
func InitiateLogin(mctlCfg *config.MesheryCtlConfig) ([]byte, error) {
	// Get the providers info
	providers, err := GetProviderInfo(mctlCfg)
	if err != nil {
		return nil, err
	}

	// Let the user select a provider
	provider := selectProviderPrompt(providers)

	var token string

	log.Println("Initiating login...")

	// If the provider URL is empty then local provider
	if provider.ProviderURL == "" {
		token, err = initiateLocalProviderAuth(provider)
		if err != nil {
			return nil, err
		}
	} else {
		token, err = initiateRemoteProviderAuth(provider)
		if err != nil {
			return nil, err
		}
	}

	// Send request with the token to the meshery server
	data, err := getTokenObjFromMesheryServer(mctlCfg, provider.ProviderName, token)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetProviderInfo queries meshery API for the provider info
func GetProviderInfo(mctCfg *config.MesheryCtlConfig) (map[string]Provider, error) {
	res := map[string]Provider{}

	resp, err := http.Get(mctCfg.GetBaseMesheryURL() + "/api/providers")
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

// initiateLocalProviderAuth initiates login process for the local provider
func initiateLocalProviderAuth(provider Provider) (string, error) {
	return "", nil
}

// initiateRemoteProviderAuth intiates login process for the remote provider
func initiateRemoteProviderAuth(provider Provider) (string, error) {
	tokenChan := make(chan string, 1)

	// Create temporary server
	srv, port, err := CreateTempAuthServer(func(rw http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			fmt.Fprintf(rw, "token not found")
			return
		}

		fmt.Fprint(rw, "successfully logged in, you can close this window now")
		tokenChan <- token
	})
	if err != nil {
		return "", err
	}

	// Create provider URI
	uri, err := createProviderURI(provider, "http://localhost", port)
	if err != nil {
		return "", err
	}

	// Redirect user to the provider page
	if err := NavigateToBrowser(uri); err != nil {
		return "", err
	}

	// Pause until we get the response on the channel
	token := <-tokenChan

	// Shut down the server
	if err := srv.Shutdown(context.TODO()); err != nil {
		return token, err
	}

	return token, nil
}

func selectProviderPrompt(provs map[string]Provider) Provider {
	provArray := []Provider{}
	provNames := []string{}

	for _, prov := range provs {
		provArray = append(provArray, prov)
	}

	for _, prov := range provArray {
		provNames = append(provNames, prov.ProviderName)
	}

	prompt := promptui.Select{
		Label: "Select a Provider",
		Items: provNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return provArray[i]
	}
}

func createProviderURI(provider Provider, host string, port int) (string, error) {
	uri, err := url.Parse(provider.ProviderURL)
	if err != nil {
		return "", err
	}

	address := fmt.Sprintf("%s:%d", host, port)

	q := uri.Query()
	q.Add("source", base64.RawURLEncoding.EncodeToString([]byte(address)))
	q.Add("provider_version", "v0.3.14")

	uri.RawQuery = q.Encode()

	return uri.String(), nil
}

func getTokenObjFromMesheryServer(mctl *config.MesheryCtlConfig, provider, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, mctl.GetBaseMesheryURL()+"/api/token", nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{
		Name:     tokenName,
		Value:    token,
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name:     "meshery-provider",
		Value:    provider,
		HttpOnly: true,
	})

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
