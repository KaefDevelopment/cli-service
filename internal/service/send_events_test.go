package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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
	requestData := `{"events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"1","branch":"new_contract_v1","timezone":"1","params":{"count":"12"}}]}`
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		assert.NoError(t, err)
		assert.Equal(t, requestData, string(body))
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.Send(testEvents)
	assert.NoError(t, actualErr)
}

func TestCLIService_Send_Error_500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.Send(testEvents)
	assert.ErrorIs(t, actualErr, errInternalServer)
}

func TestCLIService_Send_Error(t *testing.T) {
	service := CLIService{}

	actualErr := service.Send(testEvents)
	assert.Error(t, actualErr)
}

func TestCLIService_Send_Empty(t *testing.T) {
	actualData := model.Events{Events: nil}

	service := CLIService{}

	actualErr := service.Send(actualData)
	assert.NoError(t, actualErr)
}
