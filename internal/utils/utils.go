package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResponseForPlugin struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func WriteErrorResponse(err error) {
	response := ResponseForPlugin{
		Status: false,
		Error:  err.Error(),
	}

	if errors.Is(err, ErrAuthKey) {
		response.Error = err.Error()
	}

	if errors.Is(err, ErrConnectDB) {
		response.Error = err.Error()
	}

	if errors.Is(err, ErrReadRequestDataUnmarshal) {
		response.Error = err.Error()
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
