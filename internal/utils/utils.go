package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type ResponseForPlugin struct {
	Status bool `json:"status"`
	// Revoked is a positive, machine-readable marker telling the forking plugin
	// that this plugin binding was destroyed on the backend (account deleted or
	// plugin unlinked). On revoked:true the plugin MUST drop to a logged-out
	// state and regenerate its uid. Emitted only on a clean REVOKED signal
	// (HTTP 410 + body code PLUGIN_REVOKED); omitted otherwise for backward compat.
	Revoked bool   `json:"revoked,omitempty"`
	Error   string `json:"error"`
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

// WriteRevokedResponse emits an unambiguous revoked marker on stdout so the
// forking plugin logs out and regenerates its uid. Status stays true because
// the cli itself did its job (events were persisted / wiped as required);
// the actionable signal for the plugin is revoked:true.
func WriteRevokedResponse() {
	WriteResponse(ResponseForPlugin{
		Status:  true,
		Revoked: true,
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
