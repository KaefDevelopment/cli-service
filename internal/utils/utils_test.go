package utils

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestWriteErrorResponse1(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "standard error",
			args: args{err: errors.New("")},
		},
		{
			name: "auth error",
			args: args{err: ErrAuthKey},
		},
		{
			name: "connect error",
			args: args{err: ErrConnectDB},
		},
		{
			name: "create table error",
			args: args{err: ErrCreateTable},
		},
		{
			name: "read request error",
			args: args{err: ErrReadRequestDataUnmarshal},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteErrorResponse(tt.args.err)
		})
	}
}

// The revoked marker must be present and true so plugins can machine-detect it,
// while a plain success response must NOT carry the field (backward compat).
func TestResponseForPlugin_RevokedMarker(t *testing.T) {
	revoked, err := json.Marshal(ResponseForPlugin{Status: true, Revoked: true})
	if err != nil {
		t.Fatalf("marshal revoked: %v", err)
	}
	if got := string(revoked); got != `{"status":true,"revoked":true,"error":""}` {
		t.Fatalf("unexpected revoked json: %s", got)
	}

	success, err := json.Marshal(ResponseForPlugin{Status: true})
	if err != nil {
		t.Fatalf("marshal success: %v", err)
	}
	if got := string(success); got != `{"status":true,"error":""}` {
		t.Fatalf("success response must omit revoked field, got: %s", got)
	}
}

func TestWriteRevokedResponse(t *testing.T) {
	WriteRevokedResponse()
}

func TestWriteSuccessResponse(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{err: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteSuccessResponse()
		})
	}
}
