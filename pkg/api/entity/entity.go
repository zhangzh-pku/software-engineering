package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zhangzh-pku/software-engineering/pkg/api/account"
)

type Pin struct {
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Link struct {
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type LinkID struct {
	ID string `json:"_id"`
}

func (p *Pin) Create() (string, error) {
	url := "http://8.130.74.159:21080/api/pin/create"
	method := "POST"

	pinJson, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	payload := bytes.NewBuffer(pinJson)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	token, err := account.GetToken()
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Printf("PinID: %s\n", string(body))
	return string(body), nil
}

func (l *Link) Create() (string, error) {
	url := "http://8.130.74.159:21080/api/link/create"
	method := "POST"

	linkJson, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	payload := bytes.NewBuffer(linkJson)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	token, err := account.GetToken()
	fmt.Println(token)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("body:")
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var link LinkID
	err = json.Unmarshal(body, &link)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Printf("LinkID: %s\n", link.ID)
	return link.ID, nil
}
