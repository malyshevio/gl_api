package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d мин.", r)

	// для строк важно упаковать данные в кавычки иначе ошибка
	quotedJSONValue := strconv.Quote(jsonValue)

	// convert string to byte array
	return []byte(quotedJSONValue), nil
}
