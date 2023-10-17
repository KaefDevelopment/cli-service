package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestCLIService_ReadRequestData_Positive(t *testing.T) {
	requestData := `{
				"events":
					[
						{
							"id":"qwerty12345",
							"createdAt":"1",
							"type":"1",
							"project":"1",
							"projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts",
							"language":"golang",
							"target":"1",
							"branch":"",
							"timezone":"1",
							"params":{"count":"12"}
						}
					]
				}`

	service := CLIService{authKey: "12345"}

	actualData, actualErr := service.ReadRequestData(requestData)
	assert.NoError(t, actualErr)

	assert.Equal(t, model.Events{Events: []model.Event{{
		Id:             "qwerty12345",
		CreatedAt:      "1",
		Type:           "1",
		Project:        "1",
		ProjectBaseDir: "/mnt/c/Users/jaros/GolandProjects/tts",
		Language:       "golang",
		Target:         "1",
		Branch:         "",
		Timezone:       "1",
		Params:         model.Params{"count": "12"},
		AuthKey:        "12345",
		Send:           false,
	}}}, actualData)
}

func TestCLIService_ReadRequestData_UnmarshalError(t *testing.T) {
	requestData := `{"bad data request"}`

	service := CLIService{authKey: "12345"}

	_, actualErr := service.ReadRequestData(requestData)
	assert.Error(t, actualErr)

}

func TestCLIService_ReadRequestData_AuthKeyError(t *testing.T) {
	requestData := `{
				"events":
					[
						{
							"id":"qwerty12345",
							"createdAt":"1",
							"type":"1",
							"project":"1",
							"projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts",
							"language":"golang",
							"target":"1",
							"branch":"",
							"timezone":"1",
							"params":{"count":"12"}
						}
					]
				}`

	service := CLIService{}

	_, actualErr := service.ReadRequestData(requestData)
	assert.Error(t, actualErr)

}
