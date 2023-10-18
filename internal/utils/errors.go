package utils

import "errors"

var (
	ErrAuthKey                  = errors.New("failed with auth key from request events")
	ErrConnectDB                = errors.New("failed to initialize database")
	ErrReadRequestDataUnmarshal = errors.New("bad request data from plugin")
)
