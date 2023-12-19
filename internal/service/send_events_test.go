package service

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaroslav1991/cli-service/internal/model"
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

	requestData := fmt.Sprintf(`{"osName":"%s","deviceName":"%s","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"1","branch":"new_contract_v1","timezone":"1","params":{"count":"12"},"pluginId":"12345"}]}`, runtime.GOOS, hn)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		assert.NoError(t, err)
		assert.Equal(t, requestData, string(body))
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.sendEvents(testEvents, "1.0.1")
	assert.NoError(t, actualErr)
}

func TestCLIService_Send_Error_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.sendEvents(testEvents, "1.0.1")
	assert.ErrorIs(t, actualErr, errBadStatusCode)
}

func TestCLIService_Send_Error(t *testing.T) {
	service := CLIService{}

	actualErr := service.sendEvents(testEvents, "1.0.1")
	assert.Error(t, actualErr)
}

func TestCLIService_Send_Empty(t *testing.T) {
	actualData := model.Events{Events: nil}

	service := CLIService{}

	actualErr := service.sendEvents(actualData, "1.0.1")
	assert.NoError(t, actualErr)
}
