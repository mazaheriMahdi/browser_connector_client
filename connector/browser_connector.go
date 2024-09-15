package connector

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type Connector struct {
	sessionId uuid.UUID
	url       string
}

type BrowserConnector interface {
	CreateSession() (Session, error)
}

func (c Connector) CreateSession() (Session, error) {

	marshal, _ := json.Marshal(map[string]string{})
	response, err := http.Post(c.url+"Session", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return BrowserSession{}, err
	}
	defer response.Body.Close()

	responseBody := &map[string]string{}
	err = json.NewDecoder(response.Body).Decode(responseBody)
	if err != nil {
		return BrowserSession{}, err
	}
	return BrowserSession{
		id:  uuid.MustParse((*responseBody)["sessionId"]),
		url: c.url,
	}, nil
}
func NewBrowserConnector(url string) BrowserConnector {
	return &Connector{
		url: url,
	}
}

type BrowserSession struct {
	id  uuid.UUID
	url string
}

type Session interface {
	GetPageContent() (string, error)
	Goto(url string) error
	ImplicitWait(seconds int32) error
	DeleteSession() error
}

func (c BrowserSession) Goto(url string) error {
	marshal, _ := json.Marshal(map[string]any{
		"url": url,
	})
	response, err := http.Post(c.url+"Session/"+c.id.String()+"/Goto", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody := &map[string]string{}
	err = json.NewDecoder(response.Body).Decode(responseBody)
	if err != nil {
		return err
	}
	return nil
}

func (c BrowserSession) GetPageContent() (string, error) {
	response, err := http.Get(c.url + "Session/" + c.id.String() + "/Content")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody := &map[string]string{}
	err = json.NewDecoder(response.Body).Decode(responseBody)
	if err != nil {
		return "", err
	}
	return (*responseBody)["content"], nil
}

func (c BrowserSession) ImplicitWait(seconds int32) error {
	marshal, _ := json.Marshal(map[string]any{
		"seconds": seconds,
	})
	response, err := http.Post(c.url+"Session/"+c.id.String()+"/ImplicitWait", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c BrowserSession) DeleteSession() error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", c.url+"Session/"+c.id.String(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return nil
}
