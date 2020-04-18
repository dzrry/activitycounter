package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	apiMethodURL = "https://api.vk.com/method/"
	tokenURL     = "https://oauth.vk.com/token/"
	apiVersion   = "5.103"
)

const (
	DeviceAndroid = iota
	DeviceApple
)

type VK struct {
	Token  Token
	Client *http.Client
}

func NewVK(device int, login, password string) (*VK, error) {
	client := newBlankVK()

	token, err := client.auth(device, login, password)
	if err != nil {
		return nil, errors.New("client 33 " + err.Error())
	}

	client.Token = token
	return client, nil
}

func newBlankVK() *VK {
	return &VK{
		Client: &http.Client{},
	}
}

func (client *VK) auth(device int, username, password string) (Token, error) {
	var clientID, clientSecret string
	switch device {
	case DeviceAndroid:
		clientID = "2274003"
		clientSecret = "hHbZxrka2uZ6jB1inYsH"
	case DeviceApple:
		clientID = "3140623"
		clientSecret = "VeWdmVclDCtn6ihuP1nt"
	// Windows client
	default:
		clientID = "3697615"
		clientSecret = "AlVXZFMUqyrnABp8ncuU"
	}

	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return Token{}, err
	}

	q := req.URL.Query()
	q.Add("grand_type", "password")
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("username", username)
	q.Add("password", password)
	q.Add("v", apiVersion)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Client.Do(req)
	if err != nil {
		return Token{}, errors.New("client 77 " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Token{}, errors.New("client 83 " + err.Error())
	}

	var token Token
	json.Unmarshal(body, &token)

	if token.Error != "" {
		return token, errors.New(token.Error + ": " + token.ErrorDescription)
	}

	return token, nil
}

func (client *VK) MakeRequest(method string, params url.Values) (Response, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("access_token", client.Token.AccessToken)
	params.Set("v", apiVersion)

	u := apiMethodURL + method

	resp, err := client.Client.PostForm(u, params)
	if err != nil {
		return Response{}, errors.New("client 107 " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, errors.New("client 113 " + err.Error())
	}

	var response Response
	json.Unmarshal(body, &response)

	if response.ResponseError.ErrorCode != 0 {
		return response, errors.New(
			"Error code: " +
				strconv.Itoa(response.ResponseError.ErrorCode) + ", " +
				response.ResponseError.ErrorMsg)
	}

	return response, nil
}
