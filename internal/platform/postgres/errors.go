package postgres

import "github.com/lib/pq"

// errDuplicatedKey is the PostgreSQL code for duplicated key errors.
const errDuplicatedKey = pq.ErrorCode("23505")

func IsAlreadyExistsError(err error) bool {
	pgerr, ok := err.(*pq.Error)
	return ok && pgerr.Code == errDuplicatedKey
}
