package user

const (
	// status
	statusError = "error"
	statusOK    = "ok"

	// error codes
	codeSuccess = iota // 0
	codeNotFound
	codeInvalidRequest // Invalid request
	codeFailure        // Other failures
)
