package api_test

import (
    "github.com/ipedrazas/botd/api"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

var (
    server   *httptest.Server
    reader   io.Reader
    appsUrl  string
)

func init() {
    server = httptest.NewServer(api.Handlers())
    appsUrl = fmt.Sprintf("%s/apps", server.URL)
}

func TestCreateApp(t *testing.T) {
    appJson := `{"name": "dsp-test", "type": "platform"}`

    reader = strings.NewReader(appJson)
    request, err := http.NewRequest("POST", appsUrl, reader)
    res, err := http.DefaultClient.Do(request)

    if err != nil {
        t.Error(err)
    }

    if res.StatusCode != 201 {
        t.Errorf("Success expected: %d", res.StatusCode)
    }
}

func TestUniqueAppname(t *testing.T) {
    appJson := `{"name": "dsp-test", "type": "platform"}`

    reader = strings.NewReader(appJson)
    request, err := http.NewRequest("POST", appsUrl, reader)
    res, err := http.DefaultClient.Do(request)

    if err != nil {
        t.Error(err)
    }

    if res.StatusCode != 400 {
        t.Error("Bad Request expected: %d", res.StatusCode)
    }
}

func TestListApps(t *testing.T) {
    reader = strings.NewReader("")
    request, err := http.NewRequest("GET", appsUrl, reader)
    res, err := http.DefaultClient.Do(request)

    if err != nil {
        t.Error(err)
    }

    if res.StatusCode != 200 {
        t.Errorf("Success expected: %d", res.StatusCode)
    }
}
