package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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

func MigrateToNewConfigPath(newConfigPath, oldConfigPath string) error {
	if err := os.Mkdir(newConfigPath, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		log.Println(err)
		return fmt.Errorf("create new config path: %w", err)
	}

	if err := MakeHiddenConfigFolder(newConfigPath); err != nil {
		return fmt.Errorf("make hidden new config folder: %w", err)
	}

	if err := os.RemoveAll(oldConfigPath); err != nil && !errors.Is(err, os.ErrExist) {
		log.Println(err)
		return fmt.Errorf("delete old config path: %w", err)
	}

	return nil
}
