package test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mazaheriMahdi/browser_connector_client/connector"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBrowserSession_SHOULD_GoTo_Where_UrlIsValid(t *testing.T) {
	sessionID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Session/"+sessionID.String()+"/Goto", r.URL.Path)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		var body map[string]string
		json.Unmarshal(bodyBytes, &body)
		assert.Equal(t, "https://example.com", body["Url"])
		content := map[string]string{}
		json.NewEncoder(w).Encode(content)
	}))
	defer server.Close()

	browserSession := connector.BrowserSession{
		Id:  sessionID,
		Url: server.URL,
	}

	err := browserSession.Goto("https://example.com")
	assert.NoError(t, err)
}

func TestBrowserSession_SHOULD_ReturnContent_Where_GetPageContentIsCalled(t *testing.T) {
	sessionID := uuid.New()
	expectedContent := "Page Content"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Session/"+sessionID.String()+"/Content", r.URL.Path)
		content := map[string]string{"content": expectedContent}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(content)
	}))
	defer server.Close()

	browserSession := connector.BrowserSession{

		Id:  sessionID,
		Url: server.URL,
	}

	content, err := browserSession.GetPageContent()
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestBrowserSession_SHOULD_ImplicitWait_Where_WaitTimeIsValid(t *testing.T) {
	sessionID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Session/"+sessionID.String()+"/ImplicitWait", r.URL.Path)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		var body map[string]int32
		json.Unmarshal(bodyBytes, &body)
		assert.Equal(t, int32(10), body["seconds"])
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	browserSession := connector.BrowserSession{
		Id:  sessionID,
		Url: server.URL,
	}

	err := browserSession.ImplicitWait(10)
	assert.NoError(t, err)
}

func TestBrowserSession_SHOULD_DeleteSession_Where_SessionIsActive(t *testing.T) {
	sessionID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Session/"+sessionID.String(), r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	browserSession := connector.BrowserSession{
		Id:  sessionID,
		Url: server.URL,
	}

	err := browserSession.DeleteSession()
	assert.NoError(t, err)
}

func TestBrowserSession_SHOULD_Scroll_Where_XAndYAreValid(t *testing.T) {
	sessionID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Session/"+sessionID.String()+"/Scroll", r.URL.Path)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		var body map[string]int64
		json.Unmarshal(bodyBytes, &body)
		assert.Equal(t, int64(100), body["x"])
		assert.Equal(t, int64(200), body["y"])
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	browserSession := connector.BrowserSession{
		Id:  sessionID,
		Url: server.URL,
	}

	err := browserSession.Scroll(100, 200)
	assert.NoError(t, err)
}

func TestBrowserSession_SHOULD_ReturnScreenshot_Where_ScreenshotIsCalled(t *testing.T) {
	sessionID := uuid.New()
	expectedImage := []byte{0xFF, 0xD8, 0xFF, 0xE0}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/Session/"+sessionID.String()+"/Screenshot", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write(expectedImage)
	}))
	defer server.Close()

	browserSession := connector.BrowserSession{
		Id:  sessionID,
		Url: server.URL,
	}

	screenshot, err := browserSession.Screenshot()

	assert.NoError(t, err)
	assert.Equal(t, expectedImage, screenshot)
}
