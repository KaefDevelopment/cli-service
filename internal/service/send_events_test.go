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

	"github.com/KaefDevelopment/cli-service/internal/model"

	"github.com/stretchr/testify/assert"
)

var (
	testEvents = model.Events{Events: []model.Event{{
		Id:             "qwerty12345",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "/mnt/c/Users/jaros/GolandProjects/tts",
		Language:       "golang",
		Target:         "1",
		Branch:         "new_contract_v1",
		Timezone:       "1",
		Params:         model.Params{"count": "12"},
		AuthKey:        "12345",
		Send:           false,
	}}}
)

func TestCLIService_Send_Positive(t *testing.T) {
	hn, err := os.Hostname()
	assert.NoError(t, err)

	expectedData := fmt.Sprintf(`{"osName":"%s","deviceName":"%s","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"1","branch":"new_contract_v1","timezone":"1","params":{"count":"12"},"pluginId":"12345"},{"createdAt":"%s","type":"REPO_INFO","timezone":"1","params":{"reposInfo":[{"project":"1","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","repoUrl":"https://github.com/jaroslav1991/tts"}]},"pluginId":"12345"}]}`, runtime.GOOS, hn, time.Now().String())

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

			if expTime, ok := expEvent["createdAt"].(string); ok {
				if actTime, ok := actEvent["createdAt"].(string); ok {
					expParsed, _ := time.Parse(time.RFC3339, expTime)
					actParsed, _ := time.Parse(time.RFC3339, actTime)

					diff := expParsed.Sub(actParsed)
					assert.LessOrEqual(t, diff.Abs(), 5*time.Millisecond)
				}
			}
			delete(expEvent, "createdAt")
			delete(actEvent, "createdAt")

			assert.Equal(t, expEvent, actEvent)
		}
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.sendEvents(context.Background(), testEvents, "1.0.1")
	assert.NoError(t, actualErr)
}

func TestCLIService_Send_Error_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.sendEvents(context.Background(), testEvents, "1.0.1")
	assert.ErrorIs(t, actualErr, errBadStatusCode)
}

func TestCLIService_Send_Error(t *testing.T) {
	service := CLIService{}

	actualErr := service.sendEvents(context.Background(), testEvents, "1.0.1")
	assert.Error(t, actualErr)
}
