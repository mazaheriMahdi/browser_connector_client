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
		Id:  uuid.MustParse((*responseBody)["sessionId"]),
		Url: c.url,
	}, nil
}
func NewBrowserConnector(url string) BrowserConnector {
	return &Connector{
		url: url,
	}
}
