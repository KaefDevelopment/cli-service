package utils

import (
	"encoding/json"
	"fmt"
)

type ResponseForPlugin struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func WriteErrorResponse(err error) {
	WriteResponse(ResponseForPlugin{
		Status: false,
		Error:  err.Error(),
	})
}

func WriteSuccessResponse() {
	WriteResponse(ResponseForPlugin{
		Status: true,
	})
}

func WriteResponse(response ResponseForPlugin) {
	dataBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("marshal response failed:", err)
		return
	}

	fmt.Println(string(dataBytes))
}
