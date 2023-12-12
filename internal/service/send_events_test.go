package service

import (
	"github.com/golang/mock/gomock"
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
	requestData := `{"osName":"windows","deviceName":"LAPTOP-ELAQQAHD","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"1","branch":"new_contract_v1","timezone":"1","params":{"count":"12"},"pluginId":"12345"}]}`
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		assert.NoError(t, err)
		assert.Equal(t, requestData, string(body))
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.Send(testEvents, "1.0.1", "12345")
	assert.NoError(t, actualErr)
}

func TestCLIService_Send_Error_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	service := CLIService{httpAddr: server.URL}

	actualErr := service.Send(testEvents, "1.0.1", "12345")
	assert.ErrorIs(t, actualErr, errBadStatusCode)
}

func TestCLIService_Send_Error(t *testing.T) {
	service := CLIService{}

	actualErr := service.Send(testEvents, "1.0.1", "12345")
	assert.Error(t, actualErr)
}

func TestCLIService_Send_Empty(t *testing.T) {
	actualData := model.Events{Events: nil}

	service := CLIService{}

	actualErr := service.Send(actualData, "1.0.1", "12345")
	assert.NoError(t, actualErr)
}

func TestCLIService_Send(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
	}))

	defer server.Close()

	events := model.Events{Events: []model.Event{{
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
		Send:           true,
	}, {
		Id:             "qwerty123456",
		CreatedAt:      "2",
		Type:           "2",
		Project:        "2",
		ProjectBaseDir: "/mnt/c/Users/jaros/GolandProjects/tts",
		Language:       "golang",
		Target:         "2",
		Branch:         "new_contract_v1",
		Timezone:       "2",
		Params:         model.Params{"count": "13"},
		AuthKey:        "12345",
		Send:           true,
	}}}

	expectedEvents := model.Events{Events: []model.Event{{
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
		Send:           true,
	}, {
		Id:             "qwerty123456",
		CreatedAt:      "2",
		Type:           "2",
		Project:        "2",
		ProjectBaseDir: "/mnt/c/Users/jaros/GolandProjects/tts",
		Language:       "golang",
		Target:         "2",
		Branch:         "new_contract_v1",
		Timezone:       "2",
		Params:         model.Params{"count": "13"},
		AuthKey:        "12345",
		Send:           true,
	}}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Create(events).Return(nil)
	repo.EXPECT().Get().Return(events, nil)
	repo.EXPECT().Update().Return(nil)

	service := CLIService{repo: repo, httpAddr: server.URL, authKey: "12345"}

	err := service.CreateEvents(events)
	assert.NoError(t, err)

	actualEvents, err := service.GetEvents()
	assert.NoError(t, err)

	err = service.UpdateEvents()
	assert.NoError(t, err)

	actualErr := service.Send(events, "1.0.1", "12345")
	assert.ErrorIs(t, actualErr, errBadStatusCode)
	assert.Equal(t, expectedEvents, actualEvents)
}
