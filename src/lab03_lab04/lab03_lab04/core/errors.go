package core

import "fmt"

// UnknownCodeError - ошибка, когда по коду ничего не нашлось (неизвестная локация, предмет и т.д.)
type UnknownCodeError struct {
	What string
	Code string
}

func (e *UnknownCodeError) Error() string {
	return fmt.Sprintf("%s с кодом '%s' не существует", e.What, e.Code)
}
