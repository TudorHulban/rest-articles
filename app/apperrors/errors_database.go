package apperrors

import (
	"database/sql"
	"errors"
)

const (
	ErrorMsgConnectionCreation = "DB connection creation: %s"
)

// ErrObjectNotFound is used to indicate that selecting an individual object
// yielded no result. Declared as type, not value, for consistency reasons.
type ErrObjectNotFound struct{}

const ErrObjectNotFoundMessage = "record not found"

func (ErrObjectNotFound) Error() string {
	return ErrObjectNotFoundMessage
}

func (ErrObjectNotFound) Unwrap() error {
	return errors.New(ErrObjectNotFoundMessage)
}

func HandleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrObjectNotFound{}
	}

	return err
}
