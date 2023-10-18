package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResponseForPlugin struct {
	Status   string `json:"status"`
	Response any    `json:"response"`
}

func WriteErrorResponse(err error) {
	response := ResponseForPlugin{
		Status:   "Error",
		Response: nil,
	}

	if errors.Is(err, ErrAuthKey) {
		response.Response = err.Error()
	}

	if errors.Is(err, ErrConnectDB) {
		response.Response = err.Error()
	}

	if errors.Is(err, ErrReadRequestDataUnmarshal) {
		response.Response = err.Error()
	}

	WriteSuccessResponse(response)
}

func WriteSuccessResponse(data ResponseForPlugin) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal response failed:", err)
		return
	}

	fmt.Println(string(dataBytes))
}
