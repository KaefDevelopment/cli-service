package utils

import "errors"

var (
	ErrAuthKey                  = errors.New("failed with auth key from request events")
	ErrCreateTable              = errors.New("failed with table creation")
	ErrConnectDB                = errors.New("failed to initialize database")
	ErrReadRequestDataUnmarshal = errors.New("bad request data from plugin")
	ErrTypeField                = errors.New("empty type field")
	ErrCreatedAtField           = errors.New("empty created at field")
)
