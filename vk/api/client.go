package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	tokenURL = "https://oauth.vk.com/token"
	apiURL   = "https://api.vk.com/method/%s"
	apiVersion = "5.103"
)

const (
	DeviceIPhone = iota
	DeviceWPhone
	DeviceAndroid
)

type ratelimiter struct {
	requestsCount   int
	lastRequestTime time.Time
}

type VKClient struct {
	Token  Token
	Client *http.Client
	rl     *ratelimiter
	cb     *callbackHandler
}

type VKGroupBot struct {
	VKClient
	Group
	cb *botsCallBackHandler
}

type TokenOptions struct {
	ServiceToken    bool
	ValidateOnStart bool
}

func newVKClientBlank() *VKClient {
	return &VKClient{
		Client: &http.Client{},
		rl:     &ratelimiter{},
		cb: &callbackHandler{
			events: make(map[string]func(*LongPollMessage)),
		},
	}
}

func NewVKClient(device int, user string, password string) (*VKClient, error) {
	vkclient := newVKClientBlank()

	token, err := vkclient.auth(device, user, password)
	if err != nil {
		return nil, err
	}

	vkclient.Token = token

	return vkclient, nil
}

func NewVKClientWithToken(token string, options *TokenOptions) (*VKClient, error) {
	vkclient := newVKClientBlank()
	vkclient.Token.AccessToken = token

	if options != nil {
		if options.ValidateOnStart {
			uid, err := vkclient.requestSelfID()
			if err != nil {
				return nil, err
			}
			vkclient.Token.UID = uid

			if !options.ServiceToken {
				if err := vkclient.updateSelfUser(); err != nil {
					return nil, err
				}
			}
		}
	}

	return vkclient, nil
}

func NewVKGroupBot(token string, options *TokenOptions) (*VKGroupBot, error) {
	vkclient, err := NewVKClientWithToken(token, options)
	if err != nil {
		return nil, err
	}

	resp, err := vkclient.MakeRequest("groups.getById", nil)
	if err != nil {
		return nil, err
	}
	
	var res []*Group
	json.Unmarshal(resp.Response, &res)
	return &VKGroupBot{
		VKClient: *vkclient,
		Group: *res[0],
		cb: &botsCallBackHandler{
			events: make(map[string]func(*BotsLongPollObject)),
		},

	}, nil
}

func (client *VKClient) updateSelfUser() error {
	me, err := client.UsersGet([]int{client.Token.UID})
	if err != nil {
		return err
	}

	client.Token.FirstName = me[0].FirstName
	client.Token.LastName = me[0].LastName
	client.Token.PicSmall = me[0].Photo
	client.Token.PicMedium = me[0].PhotoMedium
	client.Token.PicBig = me[0].PhotoBig

	return nil
}

func (s *ratelimiter) Wait() {
	if s.requestsCount == 3 {
		secs := time.Since(s.lastRequestTime).Seconds()
		ms := int((1 - secs) * 1000)
		if ms > 0 {
			duration := time.Duration(ms * int(time.Millisecond))
			//fmt.Println("attempted to make more than 3 requests per second, "+
			//"sleeping for", ms, "ms")
			time.Sleep(duration)
		}

		s.requestsCount = 0
	}
}

func (s *ratelimiter) Update() {
	s.requestsCount++
	s.lastRequestTime = time.Now()
}

func (client *VKClient) auth(device int, user string, password string) (Token, error) {
	client.rl.Wait()
	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return Token{}, err
	}
	client.rl.Update()

	clientID := ""
	clientSecret := ""

	switch device {
	case DeviceIPhone:
		clientID = "3140623"
		clientSecret = "VeWdmVclDCtn6ihuP1nt"
	case DeviceWPhone:
		clientID = "3697615"
		clientSecret = "AlVXZFMUqyrnABp8ncuU"
	case DeviceAndroid:
		clientID = "2274003"
		clientSecret = "hHbZxrka2uZ6jB1inYsH"
	default:
		clientID = "3140623"
		clientSecret = "VeWdmVclDCtn6ihuP1nt"
	}

	q := req.URL.Query()
	q.Add("grant_type", "password")
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("username", user)
	q.Add("password", password)
	q.Add("v", "5.40")
	req.URL.RawQuery = q.Encode()

	client.rl.Wait()
	resp, err := client.Client.Do(req)
	if err != nil {
		return Token{}, err
	}
	client.rl.Update()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	var token Token
	json.Unmarshal(body, &token)

	if token.Error != "" {
		return Token{}, errors.New(token.Error + ": " + token.ErrorDescription)
	}

	return token, nil
}

func (client *VKClient) requestSelfID() (uid int, err error) {
	resp, err := client.MakeRequest("users.get", url.Values{})
	if err != nil {
		return 0, err
	}

	rawdata, err := resp.Response.MarshalJSON()
	if err != nil {
		return 0, err
	}

	data := make([]struct {
		ID int `json:"id"`
	}, 1)

	if err := json.Unmarshal(rawdata, &data); err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, nil
	}

	return data[0].ID, nil
}

func (client *VKClient) MakeRequest(method string, params url.Values) (APIResponse, error) {
	client.rl.Wait()

	endpoint := fmt.Sprintf(apiURL, method)
	if params == nil {
		params = url.Values{}
	}

	params.Set("access_token", client.Token.AccessToken)
	params.Set("v", apiVersion)

	resp, err := client.Client.PostForm(endpoint, params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	client.rl.Update()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var apiresp APIResponse
	json.Unmarshal(body, &apiresp)

	if apiresp.ResponseError.ErrorCode != 0 {
		return apiresp, errors.New("Error code: " + strconv.Itoa(apiresp.ResponseError.ErrorCode) + ", " + apiresp.ResponseError.ErrorMsg)
	}
	return apiresp, nil
}

func (client *VKClient) Execute(code string) (map[string]interface{}, error) {
	token := client.Token.AccessToken
	params := url.Values{}
	params.Set("code", code)
	params.Set("access_token", token)
	params.Set("v", apiVersion)

	resp, err := client.MakeRequest("execute", params)
	if err != nil {
		return nil, errors.New("client 278 " + err.Error())
	}

	response := make(map[string]interface{})
	if err := json.Unmarshal(resp.Response, &response); err != nil {
		return response, errors.New("client 283" + err.Error())
	}
	return response, nil
}