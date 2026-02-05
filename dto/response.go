package dto

type Response[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Record  T      `json:"record"`
}

func CreateResponseError(message string) Response[string] {
	return Response[string]{
		Code:    "99",
		Message: message,
		Record:  "",
	}
}
func CreateResponseErrorData(message string, data map[string]string) Response[map[string]string] {
	return Response[map[string]string]{
		Code:    "99",
		Message: message,
		Record:  data,
	}
}
func CreateResponseSuccess[T any](data T) Response[T] {
	return Response[T]{
		Code:    "00",
		Message: "success",
		Record:  data,
	}
}
