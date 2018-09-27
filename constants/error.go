package constants

import "errors"

// Errors that migrant will return
var (
	// ErrServiceInvalid is returned when this service is invalid.
	ErrServiceInvalid = errors.New("service is invalid")

	// ErrMigrationMissing is returned when some migrations are missing.
	ErrMigrationMissing = errors.New("migrations are missing")
	// ErrMigrationMismatch is returned when migrations are mismatch.
	ErrMigrationMismatch = errors.New("migrations are mismatch")
	// ErrMigrationNotSupported is returned when migration not supported.
	ErrMigrationNotSupported = errors.New("migrations are not supported")

	// ErrRedisActionNotSupported is returned when redis action not supported.
	ErrRedisActionNotSupported = errors.New("redis action not supported")
)
