package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/zhangzh-pku/software-engineering/pkg/api/account"
)

const (
	DockerPinType        = "docker_image"
	ScriptPinType        = "run_script"
	DissertationPinType  = "dissertation_doi"
	DatasetPinType       = "dataset"
	ReproductionPinType  = "reproduction"
	DockerKey            = "image"
	ScriptKey            = "script"
	DissertationKey      = "disseration"
	DataPathKey          = "data_path"
	ReproductionLinkType = "reproduction_link"
	OtherLinkType        = "ref_link"
)

// Manager manages Pins and Links
type DoManager struct {
	Pins  map[string]*Pin
	Links map[string]*Link
	mu    sync.Mutex
}

type LinkResp struct {
	Link
	Id string `json:"_id"`
}
type PinResp struct {
	Pin
	Id    string `json:"_id"`
	Owner string `json:"owner"`
}

// Use once.Do to ensure singleton
var (
	managerInstance *DoManager
	once            sync.Once
)

// GetManagerInstance returns the singleton Manager
func GetManagerInstance() *DoManager {
	once.Do(func() {
		managerInstance = &DoManager{
			Pins:  make(map[string]*Pin),
			Links: make(map[string]*Link),
		}
	})
	return managerInstance
}

// AddPin adds a Pin to the Manager
func (m *DoManager) AddPin(pin *Pin) (string, error) {
	pinID, err := pin.Create()
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Pins[pinID] = pin
	return pinID, nil
}

// AddLink adds a Link to the Manager
func (m *DoManager) AddLink(link *Link) (string, error) {
	linkID, err := link.Create()
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Links[linkID] = link
	return linkID, nil
}

func (m *DoManager) GetLink(id string) (*Link, error) {

	if link, ok := m.Links[id]; ok {
		return link, nil
	} else {
		url := "http://8.130.74.159:21080/api/link/" + id
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
		token, err := account.GetToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		link := &Link{
			Type:     data["type"].(string),
			Metadata: data["metadata"].(map[string]interface{}),
			From:     data["from"].(string),
			To:       data["to"].(string),
		}
		m.Links[id] = link
		return link, nil
	}
}

func (m *DoManager) GetPin(id string) (*Pin, error) {
	if pin, ok := m.Pins[id]; ok {
		return pin, nil
	} else {
		url := "http://8.130.74.159:21080/api/pin/" + strings.ReplaceAll(id, "/", "%2F")
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
		token, err := account.GetToken()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		fmt.Println(data)
		pin := &Pin{
			Type:     data["type"].(string),
			Metadata: data["metadata"].(map[string]interface{}),
		}
		m.Pins[id] = pin
		return pin, nil
	}
}

func (m *DoManager) getAllFromLinks(pinId string) ([]*Link, error) {
	baseUrl := "http://8.130.74.159:21080/api/link/list"
	params := url.Values{}
	params.Add("from", pinId)

	// Create a new HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", baseUrl, params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	token, err := account.GetToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Send the request and read the response body
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response body into a slice of links
	var _links []*LinkResp
	err = json.Unmarshal(body, &_links)
	if err != nil {
		return nil, err
	}
	// fuck golang!!!!!
	var links []*Link
	for i := range _links {
		links = append(links, &_links[i].Link)
	}
	return links, nil
}

func (m *DoManager) getAllToLinks(pinId string) ([]*Link, error) {
	baseUrl := "http://8.130.74.159:21080/api/link/list"
	params := url.Values{}
	params.Add("to", pinId)

	// Create a new HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", baseUrl, params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	token, err := account.GetToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Send the request and read the response body
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response body into a slice of links
	var _links []*LinkResp
	err = json.Unmarshal(body, &_links)
	if err != nil {
		return nil, err
	}
	// fuck golang!!!!!
	var links []*Link
	for i := range _links {
		links = append(links, &_links[i].Link)
	}
	return links, nil
}

func (m *DoManager) GetAllLinks(pinId string) ([]*Link, error) {
	links1, err := m.getAllFromLinks(pinId)
	if err != nil {
		return nil, err
	}
	links2, err := m.getAllToLinks(pinId)
	if err != nil {
		return nil, err
	}
	return append(links1, links2...), nil
}

func (m *DoManager) GenerateReproduction(dockerImage string, runScript string, dissertationDOI string, path string) error {
	var DockerImagePin, RunScriptPin, DissertationDOIPin, DatasetPin Pin
	DockerImagePin.Type = DockerPinType
	RunScriptPin.Type = ScriptPinType
	DissertationDOIPin.Type = DissertationPinType
	DatasetPin.Type = DatasetPinType
	DissertationDOIPin.Metadata = make(map[string]interface{})
	RunScriptPin.Metadata = make(map[string]interface{})
	DockerImagePin.Metadata = make(map[string]interface{})
	DatasetPin.Metadata = make(map[string]interface{})
	DatasetPin.Metadata[DataPathKey] = path
	DockerImagePin.Metadata[DockerKey] = dockerImage
	RunScriptPin.Metadata[ScriptKey] = runScript
	DissertationDOIPin.Metadata[DissertationKey] = dissertationDOI
	// 五个pin 对应论文 复现 数据 docker镜像 脚本
	dockerPinID, err := m.AddPin(&DockerImagePin)
	if err != nil {
		return err
	}
	runScriptID, err := m.AddPin(&RunScriptPin)
	if err != nil {
		return err
	}
	dissertationDOIID, err := m.AddPin(&DissertationDOIPin)
	if err != nil {
		return err
	}
	datasetID, err := m.AddPin(&DatasetPin)
	if err != nil {
		return err
	}
	var reproductionPin Pin
	reproductionPin.Metadata = make(map[string]interface{})
	reproductionPin.Type = ReproductionPinType
	reproductionID, err := m.AddPin(&reproductionPin)
	if err != nil {
		return err
	}
	link1 := &Link{
		From: dissertationDOIID,
		To:   reproductionID,
		Type: ReproductionLinkType,
	}

	link := &Link{
		From: reproductionID,
		To:   dissertationDOIID,
		Type: OtherLinkType,
	}
	link.Metadata = make(map[string]interface{})

	link1.Metadata = make(map[string]interface{})
	link2 := &Link{
		From: reproductionID,
		To:   runScriptID,
		Type: OtherLinkType,
	}
	link2.Metadata = make(map[string]interface{})
	link3 := &Link{
		From: reproductionID,
		To:   dockerPinID,
		Type: OtherLinkType,
	}
	link3.Metadata = make(map[string]interface{})
	link4 := &Link{
		From: reproductionID,
		To:   datasetID,
		Type: OtherLinkType,
	}
	link4.Metadata = make(map[string]interface{})
	_, err = m.AddLink(link)
	if err != nil {
		return err
	}
	_, err = m.AddLink(link1)
	if err != nil {
		return err
	}
	_, err = m.AddLink(link2)
	if err != nil {
		return err
	}
	_, err = m.AddLink(link3)
	if err != nil {
		return err
	}
	_, err = m.AddLink(link4)
	if err != nil {
		return err
	}
	return nil
}

func (m *DoManager) GetAllReproductions() ([]string, error) {
	url := "http://8.130.74.159:21080/api/pin/list"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	token, err := account.GetToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var _pins []PinResp
	err = json.Unmarshal(body, &_pins)
	if err != nil {
		return nil, err
	}
	var pins []string
	for i := range _pins {
		if _pins[i].Type == ReproductionPinType {
			pins = append(pins, _pins[i].Id)
		}
	}
	return pins, nil
}

func (m *DoManager) DeleteAllReproductions() error {
	l, _ := m.GetAllReproductions()
	for i := range l {
		id := l[i]
		url := "http://8.130.74.159:21080/api/pin/" + strings.ReplaceAll(id, "/", "%2F")
		method := "DELETE"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}
		token, err := account.GetToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		fmt.Println("delete body:" + string(body))
	}
	return nil
}
