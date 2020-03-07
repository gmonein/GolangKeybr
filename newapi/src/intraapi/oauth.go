package intraapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (api IntraAPI) getTokenUser(code string) (t *Token, err error) {
	if api.config == nil {
		api.config = loadConfig()
		if api.config == nil {
			return nil, errors.New("API config is null")
		}
	}
	tokenParam := map[string]string{
		"client_id":     api.config.UID,
		"client_secret": api.config.Secret,
		"redirect_uri":  "http://localhost:8083/oauth/marvin",
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
	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
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
