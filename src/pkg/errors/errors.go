package errors

import "errors"

var (
	AuthorizationError = errors.New("Authorization error")
	IsNotExist         = errors.New("Object is not exist")
	AlreadyExist       = errors.New("Object already exist")
	IncorrectRequest   = errors.New("Request is incorrect")
	EncryptionError    = errors.New("Encryption error")
)

var Is = errors.Is
