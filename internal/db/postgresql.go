package db

import (
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"time"
)

var ErrIntOverflow = errors.New("integer overflow")

type pgIntToPgInt64Converter interface {
	Int64Value() (pgtype.Int8, error)
}

func newTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: !t.IsZero()}
}

func newTimestampTZ(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: !t.IsZero()}
}

func pgIntToUInt64(n pgIntToPgInt64Converter) uint64 {
	v, _ := n.Int64Value()
	return uint64(v.Int64)
}

func uInt64ToPgInt8(n uint64) (pgtype.Int8, error) {
	pgInt := pgtype.Int8{}
	if n > math.MaxInt64 {
		return pgInt, ErrIntOverflow
	}
	if err := (&pgInt).Scan(int64(n)); err != nil {
		return pgInt, err
	}
	return pgInt, nil
}
