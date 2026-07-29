package dto

// Error codes returned in Response.Record (see CreateResponseErrorCode) let
// a caller tell apart two kinds of failure:
//
//   - Pusat  (1xxxx): a definitive answer from this server's own business
//     logic — e.g. the token was checked and rejected. The caller should
//     trust this result as-is.
//   - Local  (9xxxx): a technical problem in this service itself — e.g. a
//     dependency is down or an unexpected error occurred. The caller
//     should treat the result as unknown/retryable, not as a rejection.
//
// Add new codes here, in the range matching their source, as new failure
// cases are introduced.
const (
	// Pusat: memberinfo
	ErrCodeInvalidTokenPusat  = 10001
	ErrCodeInvalidAgentPusat  = 10004
	ErrCodeInvalidMarketPusat = 10005

	// Pusat: checkout
	ErrCodeInsufficientBalancePusat = 10002
	ErrCodePasaranOfflinePusat      = 10003

	// Local
	ErrCodeInternalLocal = 90001
)

// MsgInternalError is the only text a client ever sees for a Local
// (ErrCodeInternalLocal) failure — raw Go/SQL error strings must never reach
// the player (they can leak schema/query details, and are meaningless to a
// player anyway). The real err.Error() still has to go somewhere useful:
// callers log it via connection.Log.Error(...) (console + Telegram alert)
// right before returning this message.
const MsgInternalError = "Aplikasi bermasalah, Please Contact Admin"
