package intraapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type apiConfig struct {
	UID         string `json:"uid"`
	Secret      string `json:"secret"`
	APIEndpoint string `json:"api_endpoint"`
}

// IntraAPI is the interface intra api
type IntraAPI struct {
	config *apiConfig
}

// Token is the query interface
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Toto         string `json:"toto"`
}

// func oauthHandler(w http.ResponseWriter, r *http.Request) {
// 	token := r.URL.Query()["code"]
// 	fmt.Println(token)
// 	http.Get("https://api.intra.42.fr/v2/users/me")
// }

// Connector is the interface to query intra api
var Connector IntraAPI

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

func init() {
	config, err := readConfig("api.json")
	if err != nil {
		fmt.Println("Intra connector failed to read config", err)
		return
	}
	Connector = IntraAPI{config: config}
}

func (api IntraAPI) getTokenUser(code string) (t *Token, err error) {
	tokenParam := map[string]string{
		"client_id":     api.config.UID,
		"client_secret": api.config.Secret,
		"redirect_uri":  "http://localhost:8082/oauth",
		"grant_type":    "authorization_code",
		"scope":         "public",
		"code":          code}
	data, err := json.Marshal(tokenParam)
	if err != nil {
		return
	}

	resp, err := http.Post(api.config.APIEndpoint+"/oauth/token", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &t)
	return
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

// Me return /v2/me parsed IntraUser of token
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

// GetUserFromCode send /v2/me from oauth authorization_code
func GetUserFromCode(code string) (user *IntraUser, token *Token, err error) {
	token, err = Connector.getTokenUser(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	user, err = token.Me()
	return
}
