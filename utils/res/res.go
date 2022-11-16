package res

type Res struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func Ok(data any) *Res {
	return &Res{
		Code: 200,
		Data: data,
	}
}

func OkMessage(data any, message string) *Res {
	return &Res{
		Code:    200,
		Data:    data,
		Message: message,
	}
}

func Err(code int, err error) *Res {
	return &Res{
		Code:    code,
		Data:    err.Error(),
		Message: err.Error(),
	}
}

func ErrMessage(code int, message string) *Res {
	return &Res{
		Code:    code,
		Data:    "",
		Message: message,
	}
}
