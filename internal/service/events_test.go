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

func TestCLIService_GetKeys_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().GetAuthKeys().Return([]string{"12345"}, nil)

	service := NewCLIService(repo, "", "12345", true)
	actual, err := service.GetKeys()
	assert.NoError(t, err)
	assert.Equal(t, []string{"12345"}, actual)
}

func TestCLIService_GetKeys_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().GetAuthKeys().Return(nil, utils.ErrAuthKey)

	service := NewCLIService(repo, "", "", true)
	_, err := service.GetKeys()
	assert.ErrorIs(t, err, utils.ErrAuthKey)
}

func TestCLIService_GetEvents_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventsByAuthKey := model.EventsByAuthKey{Events: []model.Events{{Events: []model.Event{{
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
		Send:           true,
	}}}}}

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Get([]string{"12345"}).Return(eventsByAuthKey, nil)

	service := NewCLIService(repo, "", "12345", true)
	actualEvents, err := service.GetEvents([]string{"12345"})
	assert.NoError(t, err)

	assert.Equal(t, model.EventsByAuthKey{Events: []model.Events{{Events: []model.Event{{
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
		Send:           true,
	}}}}}, actualEvents)
}

func TestCLIService_GetEvents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Get([]string{}).Return(model.EventsByAuthKey{}, errors.New("some error"))

	service := NewCLIService(repo, "", "", true)
	_, err := service.GetEvents([]string{})
	assert.Error(t, err)
}

func TestCLIService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := model.EventsByAuthKey{Events: []model.Events{{Events: []model.Event{{
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
		Send:           true,
	}}}}}

	repo := NewMockRepository(ctrl)
	repo.EXPECT().Drop(events).Return(nil)

	service := NewCLIService(repo, "", "12345", true)
	err := service.Delete(events)
	assert.NoError(t, err)
}
