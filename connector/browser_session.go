package connector

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type BrowserSession struct {
	Id  uuid.UUID
	Url string
}

type Session interface {
	GetPageContent() (string, error)
	Goto(url string) error
	ImplicitWait(seconds int32) error
	DeleteSession() error
	Scroll(x int64, y int64) error
	Clean() error
	Screenshot() ([]byte, error)
	Click(selector string) error
	Wait(selector string) error
}

func (c BrowserSession) Goto(url string, pageHeight int64, pageWeight int64) error {
	marshal, _ := json.Marshal(
		map[string]any{
			"url":        url,
			"pageHeight": pageHeight,
			"pageWeight": pageWeight,
		},
	)
	response, err := http.Post(c.Url+"/Session/"+c.Id.String()+"/Goto", "application/json", bytes.NewBuffer(marshal))
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
	response, err := http.Get(c.Url + "/Session/" + c.Id.String() + "/Content")
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
	marshal, _ := json.Marshal(
		map[string]any{
			"seconds": seconds,
		},
	)
	response, err := http.Post(
		c.Url+"/Session/"+c.Id.String()+"/ImplicitWait", "application/json", bytes.NewBuffer(marshal),
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c BrowserSession) DeleteSession() error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", c.Url+"/Session/"+c.Id.String(), nil)
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

func (c BrowserSession) Scroll(x int64, y int64) error {
	marshal, _ := json.Marshal(
		map[string]any{
			"x": x,
			"y": y,
		},
	)
	response, err := http.Post(c.Url+"/Session/"+c.Id.String()+"/Scroll", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c BrowserSession) Clean() error {
	marshal, _ := json.Marshal(map[string]any{})
	response, err := http.Post(c.Url+"/Session/"+c.Id.String()+"/Clean", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func (c BrowserSession) Screenshot() ([]byte, error) {
	response, err := http.Get(c.Url + "/Session/" + c.Id.String() + "/Screenshot")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	imageBytes, err := io.ReadAll(response.Body)
	return imageBytes, err
}

func (c BrowserSession) Click(selector string) error {
	marshal, _ := json.Marshal(
		map[string]any{
			"selector": selector,
		},
	)
	response, err := http.Post(c.Url+"/Session/"+c.Id.String()+"/Click", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return err
}

func (c BrowserSession) Wait(selector string) error {
	marshal, _ := json.Marshal(
		map[string]any{
			"selector": selector,
		},
	)
	response, err := http.Post(c.Url+"/Session/"+c.Id.String()+"/Wait", "application/json", bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return err
}
