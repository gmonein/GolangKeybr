package intraapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type apiConfig struct {
	UID         string `json:"uid"`
	Secret      string `json:"secret"`
	APIEndpoint string `json:"api_endpoint"`
}

// Token is the query interface
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Connector is the interface to query intra api
var Connector IntraAPI

// IntraAPI is the interface intra api
type IntraAPI struct {
	Token  *Token
	config *apiConfig
}

// Me return /v2/me parsed IntraUser of token
func init() {
	Connector = IntraAPI{config: loadConfig()}
}

func readConfig(path string) (*apiConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to read file", path)
		return nil, err
	}

	config := apiConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return &config, nil
}

func loadConfig() *apiConfig {
	config, err := readConfig(os.Getenv("CONFIGS_PATH") + "/api.json")
	if err != nil {
		fmt.Println("Intra connector failed to read config", err)
		return nil
	}
	return config
}

// Get with token
func (t Token) Get(endpoint string) (body []byte, err error) {
	base, err := url.Parse(Connector.config.APIEndpoint + endpoint)
	if err != nil {
		return
	}

	urlParams := url.Values{}
	urlParams.Add("access_token", t.AccessToken)
	base.RawQuery = urlParams.Encode()

	fmt.Printf("Encoded URL is %q\n", base.String())
	resp, err := http.Get(base.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	return
}

// Me call /v2/me and return an IntraUser and error
func (t *Token) Me() (user *IntraUser, err error) {
	me, err := t.Get("v2/me")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(me, &user)
	if err != nil {
		return
	}
	return
}
