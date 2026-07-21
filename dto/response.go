package dto

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Record  T      `json:"record"`
}

func CreateResponseError(status int, message string) Response[string] {
	return Response[string]{
		Status:  status,
		Message: message,
		Record:  "",
	}
}
func CreateResponseErrorData(status int, message string, data map[string]string) Response[map[string]string] {
	return Response[map[string]string]{
		Status:  status,
		Message: message,
		Record:  data,
	}
}

type ErrorCode struct {
	Code int `json:"code"`
}

// CreateResponseErrorCode attaches a numeric error code to the response so
// other servers checking this endpoint can identify that the error came
// from this (pusat) server, independent of the HTTP status.
func CreateResponseErrorCode(status int, code int, message string) Response[ErrorCode] {
	return Response[ErrorCode]{
		Status:  status,
		Message: message,
		Record:  ErrorCode{Code: code},
	}
}
func CreateResponseSuccess[T any](data T) Response[T] {
	return Response[T]{
		Status:  200,
		Message: "success",
		Record:  data,
	}
}
