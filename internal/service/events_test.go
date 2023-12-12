package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
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
	repo.EXPECT().Create(events).Return(nil)

	service := NewCLIService(repo, "", "12345", true)

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
	repo.EXPECT().Create(events).Return(utils.ErrReadRequestDataUnmarshal)

	service := NewCLIService(repo, "", "", true)

	err := service.CreateEvents(events)
	assert.Error(t, err)

}

func TestCLIService_UpdateEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Update().Return(nil)

	service := NewCLIService(repo, "", "12345", true)
	err := service.UpdateEvents()
	assert.NoError(t, err)

}

func TestCLIService_GetEvents_Positive(t *testing.T) {
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
	repo.EXPECT().Get().Return(events, nil)

	service := NewCLIService(repo, "", "12345", true)
	actualEvents, err := service.GetEvents()
	assert.NoError(t, err)

	assert.Equal(t, model.Events{Events: []model.Event{{
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
	}}}, actualEvents)
}

func TestCLIService_GetEvents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Get().Return(model.Events{}, errors.New("some error"))

	service := NewCLIService(repo, "", "", true)
	_, err := service.GetEvents()
	assert.Error(t, err)
}

func TestCLIService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Drop().Return(nil)

	service := NewCLIService(repo, "", "12345", true)
	err := service.Delete()
	assert.NoError(t, err)
}
