package reborn

import "fmt"

type UnsupportedValueTypeError struct {
	Value interface{}
}

func (u UnsupportedValueTypeError) Error() string {
	return fmt.Sprintf("unsupported value type, %T", u.Value)
}

type LoadFromDBError struct {
	err error
}

func (l LoadFromDBError) Error() string {
	return fmt.Sprintf("failed to load config from db, error: %s", l.err.Error())
}

type SyncDefaultsToDBError struct {
	err error
}

func (s SyncDefaultsToDBError) Error() string {
	return fmt.Sprintf("failed to sync defaults to db, error: %s", s.err.Error())
}
