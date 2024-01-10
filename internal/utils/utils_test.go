package utils

import (
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
