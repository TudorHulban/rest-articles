package db

import (
	"database/sql"
	"errors"
)

// ErrObjectNotFound is used to indicate that selecting an individual object
// yielded no result. Declared as type, not value, for consistency reasons.
type ErrObjectNotFound struct{}

const errObjectNotFoundMessage = "object not found"

func (ErrObjectNotFound) Error() string {
	return errObjectNotFoundMessage
}

func (ErrObjectNotFound) Unwrap() error {
	return errors.New(errObjectNotFoundMessage)
}

func HandleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrObjectNotFound{}
	}

	return err
}
