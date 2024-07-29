package response

// Ответ HTTP-сервера
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// OK возвращает ответ для запроса с успешным выполнением
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// Error возвращает ответ для запроса с ошибкой
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
