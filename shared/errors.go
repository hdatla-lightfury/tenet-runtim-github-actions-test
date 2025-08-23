package common

import "github.com/heroiclabs/nakama-common/runtime"

const (
	OK                  = 0
	CANCELED            = 1
	UNKNOWN             = 2
	INVALID_ARGUMENT    = 3
	DEADLINE_EXCEEDED   = 4
	NOT_FOUND           = 5
	ALREADY_EXISTS      = 6
	PERMISSION_DENIED   = 7
	RESOURCE_EXHAUSTED  = 8
	FAILED_PRECONDITION = 9
	ABORTED             = 10
	OUT_OF_RANGE        = 11
	UNIMPLEMENTED       = 12
	INTERNAL            = 13
	UNAVAILABLE         = 14
	DATA_LOSS           = 15
	UNAUTHENTICATED     = 16
)

var (
	ErrBadInput           = runtime.NewError("input contained invalid data", INVALID_ARGUMENT)
	ErrInternalError      = runtime.NewError("internal server error", INTERNAL)
	ErrGuildAlreadyExists = runtime.NewError("guild name is in use", ALREADY_EXISTS)
	ErrFullGuild          = runtime.NewError("guild is full", RESOURCE_EXHAUSTED)
	ErrNotAllowed         = runtime.NewError("operation not allowed", PERMISSION_DENIED)
	ErrNoGuildFound       = runtime.NewError("guild not found", NOT_FOUND)
)
