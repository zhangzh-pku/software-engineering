package account

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func GetToken() (token string, err error) {

	url := "http://8.130.74.159:21080/api/user/login"
	method := "POST"

	payload := strings.NewReader(`{
    "username": "rgtest",
    "password": "12345678"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response LoginResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	return response.Token, nil
}
