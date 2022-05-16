package repository

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var (
	ErrorRecordNotFound           = errors.New("resource not found")
	ErrorOther                    = errors.New("an error occurred")
	ErrorUniqueConstaintViolation = errors.New("record already exists (duplicate unique key)")
)

var PGuniqueConstraintCode = "23505"

func isUniqueConstaintViolationError(err error) bool {
	if err, ok := err.(*pgconn.PgError); ok {
		return err.Code == PGuniqueConstraintCode
	}
	return false
}

func checkError(err error) error {
	if isUniqueConstaintViolationError(err) {
		return ErrorUniqueConstaintViolation
	} else if err == gorm.ErrRecordNotFound {
		return ErrorRecordNotFound
	}
	return ErrorOther
}
