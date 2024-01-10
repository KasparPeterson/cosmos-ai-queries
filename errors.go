package cosmos_ai_queries

import "cosmossdk.io/errors"

var (
	ErrIndexTooLong     = errors.Register(ModuleName, 2, "index too long")
	ErrDuplicateAddress = errors.Register(ModuleName, 3, "duplicate address")
	ErrInvalidUser      = errors.Register(ModuleName, 4, "user address is invalid: %s")
)
