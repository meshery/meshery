package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Provider struct {
	ProviderURL  string `json:"provider_url,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
}

// AddAuthDetails Adds authentication cookies to the request
func AddAuthDetails(req *http.Request, filepath string) error {
	file, err := ioutil.ReadFile(filepath)
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "error reading body :")
		return err
	}

	if ContentTypeIsHTML(resp) {
		return errors.New("invalid body")
	}

	return ioutil.WriteFile(filepath, data, os.ModePerm)
}

func CreateTempAuthServer(fn func(http.ResponseWriter, *http.Request)) (*http.Server, int, error) {
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
	}

	listener, err := net.Listen("tcp4", ":0")
	if err != nil {
		return nil, -1, err
	}

	mux.HandleFunc("/api/user/token", fn)

	if err := srv.Serve(listener); err != nil {
		return nil, -1, err
	}

	return srv, listener.Addr().(*net.TCPAddr).Port, nil
}

func InitiateLogin() {
	// Get mesheryctl config
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
	}

	// Get the providers info
	providers, err := GetProviderInfo(mctlCfg)
	if err != nil {

	}

	// Let the user select a provider
	provider := selectProviderPrompt(providers)

	var token string

	// If the provider URL is empty then local provider
	if provider.ProviderURL == "" {
		token, err = initiateLocalProviderAuth(provider)
		if err != nil {

		}
	} else {
		token, err = initiateRemoteProviderAuth(provider)
		if err != nil {

		}
	}

	// Send request with the token to the meshery server
	data, err := getTokenObjFromMesheryServer(mctlCfg, token)
	if err != nil {

	}

	// Save the token json file to the filesystem
	fmt.Println(string(data))
}

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

func initiateLocalProviderAuth(provider Provider) (string, error) {
	return "", nil
}

func initiateRemoteProviderAuth(provider Provider) (string, error) {
	stopChan := make(chan struct{}, 1)

	// Create temporary server
	srv, port, err := CreateTempAuthServer(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.String())
		stopChan <- struct{}{}
	})
	if err != nil {
		return "", err
	}

	// Create provider URI
	uri, err := createProviderURI(provider, "localhost", port)
	if err != nil {
		return "", err
	}

	// Redirect user to the provider page
	if err := NavigateToBrowser(uri); err != nil {
		return "", err
	}

	// Pause until we get the response on the channel
	<-stopChan

	// Shut down the server
	srv.Shutdown(context.TODO())

	return "", nil
}

func selectProviderPrompt(provs map[string]Provider) Provider {
	provArray := []Provider{}

	for _, prov := range provs {
		provArray = append(provArray, prov)
	}

	// Print providers
	for i, prov := range provArray {
		fmt.Printf("%d. %s\n", i+1, prov.ProviderName)
	}

	// Scan user response
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()

		resp, err := strconv.Atoi(scanner.Text())
		if resp > 0 && resp < len(provArray) || err != nil {
			fmt.Printf("invalid input: enter a number between 1 - %d\n", len(provArray))
			continue
		}

		return provArray[resp]
	}

	return Provider{}
}

func createProviderURI(provider Provider, host string, port int) (string, error) {
	uri, err := url.Parse(provider.ProviderURL)
	if err != nil {
		return "", err
	}

	address := fmt.Sprintf("%s:%d", host, port)

	uri.Query().Add("source", base64.RawURLEncoding.EncodeToString([]byte(address)))
	uri.Query().Add("provider_version", "v0.3.14")

	return uri.String(), nil
}

func getTokenObjFromMesheryServer(mctl *config.MesheryCtlConfig, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, mctl.GetBaseMesheryURL()+"/api/gettoken", nil)
	if err != nil {
		return nil, err
	}

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
