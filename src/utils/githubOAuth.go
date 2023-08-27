package utils

import (
	"api/src/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type GitHubOauthToken struct {
	AccessToken string
}

type GitHubUserResult struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
	Email string `json:"email"`
}

func GetGitHubOauthToken(code string, config config.Config) (*GitHubOauthToken, error) {
	const rootURl = "https://github.com/login/oauth/access_token"

	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", config.GitHubClientID)
	values.Add("client_secret", config.GitHubClientSecret)
	log.Printf("Values: %#v\n", values)

	query := values.Encode()

	queryString := fmt.Sprintf("%s?%s", rootURl, bytes.NewBufferString(query))
	req, err := http.NewRequest("POST", queryString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsedQuery, err := url.ParseQuery(string(resBody))
	if err != nil {
		return nil, err
	}

	tokenBody := &GitHubOauthToken{
		AccessToken: parsedQuery["access_token"][0],
	}

	return tokenBody, nil
}

func GetGitHubUser(access_token string) (*GitHubUserResult, error) {
	rootUrl := "https://api.github.com/user"

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	userBody := new(GitHubUserResult)
	if err := json.Unmarshal(resBody, userBody); err != nil {
		return nil, err
	}

	return userBody, nil
}
