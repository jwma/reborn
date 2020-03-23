package reborn

import "fmt"

type UnsupportedValueTypeError struct {
	Value interface{}
}

func (u UnsupportedValueTypeError) Error() string {
	return fmt.Sprintf("unsupported value type, %T", u.Value)
}
