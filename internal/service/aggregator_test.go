package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KaefDevelopment/cli-service/internal/model"
)

func TestCLIService_Aggregate_BranchNotFound(t *testing.T) {
	service := CLIService{authKey: "12345"}

	events := model.Events{Events: []model.Event{{
		Id:             "qwerty123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	expected := model.Events{Events: []model.Event{{
		Id:             "qwerty123",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "1",
		Language:       "golang",
		Target:         "1",
		Branch:         "",
		Timezone:       "1",
		Params:         nil,
		AuthKey:        "12345",
		Send:           false,
	}}}

	service.Aggregate(events)
	assert.Equal(t, expected, events)
}

func TestCLIService_Aggregate_IDEmpty(t *testing.T) {
	service := CLIService{authKey: "12345"}

	getIDFn = func() string {
		return "qwerty123"
	}

	events := model.Events{Events: []model.Event{{
		Id:             "",
		ProjectBaseDir: "test",
	}}}

	expected := model.Events{Events: []model.Event{{
		Id:             getIDFn(),
		ProjectBaseDir: "test",
	}}}

	service.Aggregate(events)
	assert.Equal(t, expected, events)
}

func TestCLIService_Aggregate_FindBranchInGit(t *testing.T) {
	service := CLIService{authKey: "12345"}

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		ProjectBaseDir: "C:/Users/jaros/GolandProjects/tts",
		Branch:         "",
	}}}

	getBranchFn = func(_ string) string {
		return "new_contract_v1"
	}

	expected := model.Events{Events: []model.Event{{
		Id:             "123",
		ProjectBaseDir: "C:/Users/jaros/GolandProjects/tts",
		Branch:         "new_contract_v1",
	}}}

	service.Aggregate(events)
	assert.Equal(t, expected, events)
}

func TestCLIService_Aggregate_BranchNotEmpty(t *testing.T) {
	service := CLIService{authKey: "12345"}

	events := model.Events{Events: []model.Event{{
		Id:     "123",
		Branch: "master",
	}}}

	expected := model.Events{Events: []model.Event{{
		Id:     "123",
		Branch: "master",
	}}}

	service.Aggregate(events)
	assert.Equal(t, expected, events)
}

func TestCLIService_Aggregate_ProjectBaseDir_Empty(t *testing.T) {
	service := CLIService{authKey: "12345"}

	getBranchFn, getIDFn = func(_ string) string { return "" }, func() string { return "" }

	events := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "123456789",
		Type:           "test",
		ProjectBaseDir: "",
	}}}

	expected := model.Events{Events: []model.Event{{
		Id:             "123",
		CreatedAt:      "123456789",
		Type:           "test",
		ProjectBaseDir: "",
	}}}

	service.Aggregate(events)
	assert.Equal(t, expected, events)
}
