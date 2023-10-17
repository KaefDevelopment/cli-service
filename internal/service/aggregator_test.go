package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// C:\Users\jaros\GolandProjects\tts

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
		Id: "",
	}}}

	expected := model.Events{Events: []model.Event{{
		Id: getIDFn(),
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
