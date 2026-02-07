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
func CreateResponseSuccess[T any](data T) Response[T] {
	return Response[T]{
		Status:  200,
		Message: "success",
		Record:  data,
	}
}
