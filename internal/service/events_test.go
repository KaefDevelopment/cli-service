package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/utils"
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
	repo.EXPECT().Create(events).Return(nil)

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(events)
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
	repo.EXPECT().Create(events).Return(utils.ErrReadRequestDataUnmarshal)

	service := NewCLIService(repo, txp, "", "")

	err := service.CreateEvents(events)
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
	repo.EXPECT().Create(events).Return(nil)
	repo.EXPECT().MarkSent().Return(nil)
	repo.EXPECT().GetMarked().Return(events, nil)

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(events)
	assert.NoError(t, err)

	resEvents, err := service.lockEvents(repo)
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
	repo.EXPECT().Create(events).Return(nil)
	repo.EXPECT().MarkSent().Return(nil)
	repo.EXPECT().GetMarked().Return(events, errors.New("get marked error"))

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(events)
	assert.NoError(t, err)

	_, err = service.lockEvents(repo)
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
	repo.EXPECT().Create(events).Return(nil)
	repo.EXPECT().MarkSent().Return(errors.New("mark sent error"))

	service := NewCLIService(repo, txp, "", "12345")

	err := service.CreateEvents(events)
	assert.NoError(t, err)

	_, err = service.lockEvents(repo)
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
	repo.EXPECT().Delete(events).Return(nil)

	service := NewCLIService(repo, txp, "", "12345")

	err := service.unlockEvents(repo, events)
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
	repo.EXPECT().Delete(events).Return(errors.New("delete error"))

	service := NewCLIService(repo, txp, "", "12345")

	err := service.unlockEvents(repo, events)
	assert.Error(t, err)
}
