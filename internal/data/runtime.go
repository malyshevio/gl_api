package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("не корректный формат runtime")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d мин.", r)

	// для строк важно упаковать данные в кавычки иначе ошибка
	quotedJSONValue := strconv.Quote(jsonValue)

	// convert string to byte array
	return []byte(quotedJSONValue), nil
}

// UnmarshalJSON implements Unmarshaler interface
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "мин." {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
