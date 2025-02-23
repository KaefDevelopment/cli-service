package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/utils"
)

func TestCLIService_CreateEvents_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Create(gomock.Any(), events).Return(nil)

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(context.Background(), events)
	assert.NoError(t, err)

}

func TestCLIService_CreateEvents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Create(gomock.Any(), events).Return(utils.ErrReadRequestDataUnmarshal)

	service := NewCLIService(repo, txp, "", "")

	err := service.CreateEvents(context.Background(), events)
	assert.Error(t, err)

}

func TestCLIService_lockEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Create(gomock.Any(), events).Return(nil)
	repo.EXPECT().MarkSent(gomock.Any()).Return(nil)
	repo.EXPECT().GetMarked(gomock.Any()).Return(events, nil)

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(context.Background(), events)
	assert.NoError(t, err)

	resEvents, err := service.lockEvents(context.Background(), repo)
	assert.NoError(t, err)
	assert.Equal(t, events, resEvents)

}

func TestCLIService_lockEvents_ErrorGetMarked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Create(gomock.Any(), events).Return(nil)
	repo.EXPECT().MarkSent(gomock.Any()).Return(nil)
	repo.EXPECT().GetMarked(gomock.Any()).Return(events, errors.New("get marked error"))

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(context.Background(), events)
	assert.NoError(t, err)

	_, err = service.lockEvents(context.Background(), repo)
	assert.Error(t, err)
}

func TestCLIService_lockEvents_ErrorMarkSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Create(gomock.Any(), events).Return(nil)
	repo.EXPECT().MarkSent(gomock.Any()).Return(errors.New("mark sent error"))

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(context.Background(), events)
	assert.NoError(t, err)

	_, err = service.lockEvents(context.Background(), repo)
	assert.Error(t, err)
}

func TestCLIService_unlockEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Delete(gomock.Any(), events).Return(nil)

	service := NewCLIService(repo, txp, "", "12345")

	err := service.unlockEvents(context.Background(), repo, events)
	assert.NoError(t, err)
}

func TestCLIService_unlockEvents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)
	repo.EXPECT().Delete(gomock.Any(), events).Return(errors.New("delete error"))

	service := NewCLIService(repo, txp, "", "12345")

	err := service.unlockEvents(context.Background(), repo, events)
	assert.Error(t, err)
}

func TestCLIService_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	service := NewCLIService(repo, txp, "", "")

	txp.EXPECT().Transaction(gomock.Any()).Return(nil)

	err := service.Send(context.Background(), "")
	assert.NoError(t, err)
}

func TestCLIService_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	service := NewCLIService(repo, txp, "", "")

	txp.EXPECT().Transaction(gomock.Any()).Return(errors.New("error with transaction"))

	err := service.Send(context.Background(), "")
	assert.Error(t, err)
}

func TestCLIService_sendWithLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	events := model.Events{Events: []model.Event{{
		Id:             "qwerty12345",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "master",
		Timezone:       "MSK",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	hn, err := os.Hostname()

	expectedData := fmt.Sprintf(`{"osName":"%s","deviceName":"%s","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"1","language":"golang","target":"1","branch":"master","timezone":"MSK","pluginId":"12345"},{"createdAt":"%s","type":"REPO_INFO","timezone":"MSK","params":{"reposInfo":[{"project":"1","projectBaseDir":"1","repoUrls":null}]},"pluginId":"12345"}]}`, runtime.GOOS, hn, time.Now().String())
	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	repo.EXPECT().MarkSent(gomock.Any()).Return(nil)
	repo.EXPECT().GetMarked(gomock.Any()).Return(events, nil)

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

	service := NewCLIService(repo, txp, server.URL, "")

	err = service.sendEvents(context.Background(), events, "1.0.1")
	assert.NoError(t, err)

	repo.EXPECT().Delete(gomock.Any(), events).Return(nil)

	err = service.sendWithLock(context.Background(), repo, "1.0.1")
	assert.NoError(t, err)
}

func TestCLIService_sendWithLock_ErrorUnlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	hn, err := os.Hostname()

	expectedData := fmt.Sprintf(`{"osName":"%s","deviceName":"%s","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"1","language":"golang","target":"1","branch":"master","timezone":"1","pluginId":"12345"},{"createdAt":"%s","type":"REPO_INFO","timezone":"1","params":{"reposInfo":[{"project":"1","projectBaseDir":"1","repoUrls":null}]},"pluginId":"12345"}]}`, runtime.GOOS, hn, time.Now().String())

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	repo.EXPECT().MarkSent(gomock.Any()).Return(nil)
	repo.EXPECT().GetMarked(gomock.Any()).Return(events, nil)

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

	service := NewCLIService(repo, txp, server.URL, "")

	err = service.sendEvents(context.Background(), events, "1.0.1")
	assert.NoError(t, err)

	repo.EXPECT().Delete(gomock.Any(), events).Return(errors.New("error with unlock"))

	err = service.sendWithLock(context.Background(), repo, "1.0.1")
	assert.Error(t, err)
}

func TestCLIService_sendWithLock_ErrorSendEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	hn, err := os.Hostname()

	requestData := fmt.Sprintf(`{"osName":"%s","deviceName":"%s","cliVersion":"1.0.1","events":[{"id":"qwerty12345","createdAt":"1","type":"1","project":"1","projectBaseDir":"1","language":"golang","target":"1","branch":"master","timezone":"1","pluginId":"12345"}]}`, runtime.GOOS, hn)

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	repo.EXPECT().MarkSent(gomock.Any()).Return(nil)
	repo.EXPECT().GetMarked(gomock.Any()).Return(events, nil)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		assert.NoError(t, err)
		assert.Equal(t, requestData, string(body))
	}))

	defer server.Close()

	service := NewCLIService(repo, txp, "", "")

	err = service.sendEvents(context.Background(), events, "1.0.1")
	assert.Error(t, err)

	err = service.sendWithLock(context.Background(), repo, "1.0.1")
	assert.Error(t, err)
}

func TestCLIService_sendWithLock_ErrorLockSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	repo.EXPECT().MarkSent(gomock.Any()).Return(errors.New("error with mark sent"))

	service := NewCLIService(repo, txp, "", "")

	err := service.sendWithLock(context.Background(), repo, "1.0.1")
	assert.Error(t, err)
}

func TestCLIService_sendWithLock_ErrorLockGetMarked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	txp := NewMockTxProvider(ctrl)

	repo.EXPECT().MarkSent(gomock.Any()).Return(nil)
	repo.EXPECT().GetMarked(gomock.Any()).Return(model.Events{}, errors.New("error with get marked"))

	service := NewCLIService(repo, txp, "", "")

	err := service.sendWithLock(context.Background(), repo, "1.0.1")
	assert.Error(t, err)
}
