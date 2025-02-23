package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KaefDevelopment/cli-service/internal/model"
)

func TestCLIService_Send_Positive(t *testing.T) {
	events := model.Events{Events: []model.Event{{
		Id:             "qwerty12345",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         model.Params{"count": "12"},
		AuthKey:        "12345",
		Send:           false,
	}}}
	hn, err := os.Hostname()
	assert.NoError(t, err)

	expectedData := fmt.Sprintf(`{"osName":"%s","deviceName":"%s","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"1","language":"golang","target":"1","branch":"master","timezone":"1","params":{"count":"12"},"pluginId":"12345"},{"createdAt":"%s","type":"REPO_INFO","timezone":"1","params":{"reposInfo":[{"project":"1","projectBaseDir":"1","repoUrls":null}]},"pluginId":"12345"}]}`, runtime.GOOS, hn, time.Now().String())

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		actualData, err := io.ReadAll(request.Body)
		assert.NoError(t, err)

		var expected, actual map[string]interface{}
		err = json.Unmarshal([]byte(expectedData), &expected)
		assert.NoError(t, err)
		err = json.Unmarshal(actualData, &actual)
		assert.NoError(t, err)
		expectedEvents := expected["events"].([]interface{})
		actualEvents := actual["events"].([]interface{})

		for i := range expectedEvents {
			expEvent := expectedEvents[i].(map[string]interface{})
			actEvent := actualEvents[i].(map[string]interface{})

			delete(expEvent, "createdAt")
			delete(actEvent, "createdAt")
			delete(expEvent, "id")
			delete(actEvent, "id")

			assert.Equal(t, expEvent, actEvent)
		}
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.sendEvents(context.Background(), events, "1.0.1")
	assert.NoError(t, actualErr)
}

func TestCLIService_Send_Error_BadStatus(t *testing.T) {
	testEvents := model.Events{Events: []model.Event{{
		Id:             "qwerty12345",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "cli-service",
		ProjectBaseDir: "C:/Users/jaros/GolandProjects/cli-service",
		Language:       "golang",
		Target:         "1",
		Branch:         "new_contract_v1",
		Timezone:       "1",
		Params:         model.Params{"count": "12"},
		AuthKey:        "12345",
		Send:           false,
	}}}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.sendEvents(context.Background(), testEvents, "1.0.1")
	assert.ErrorIs(t, actualErr, errBadStatusCode)
}

func TestCLIService_Send_Error(t *testing.T) {
	testEvents := model.Events{Events: []model.Event{{
		Id:             "qwerty12345",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "cli-service",
		ProjectBaseDir: "C:/Users/jaros/GolandProjects/cli-service",
		Language:       "golang",
		Target:         "1",
		Branch:         "new_contract_v1",
		Timezone:       "1",
		Params:         model.Params{"count": "12"},
		AuthKey:        "12345",
		Send:           false,
	}}}
	service := CLIService{}

	actualErr := service.sendEvents(context.Background(), testEvents, "1.0.1")
	assert.Error(t, actualErr)
}
