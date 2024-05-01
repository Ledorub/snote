package internal

import (
	"errors"
	"github.com/ledorub/snote-api/internal/validator"
	"time"
)

type Note struct {
	ID                string
	Content           *string
	CreatedAt         time.Time
	ExpiresIn         time.Duration
	ExpiresAt         time.Time
	ExpiresAtTimeZone *time.Location
	KeyHash           []byte
}

func (n *Note) CheckErrors() error {
	v := validator.Validator{}

	v.Check(len(n.ID) == 10, "id should consist of 10 letters and/or digits")
	v.Check(validator.ValidateASCIIAlphaNumeric(n.ID), "id should consist of latin letters and/or digits")
	v.Check(len(*n.Content) != 0, "content should be provided")
	v.Check(len(*n.Content) <= 1_048_576, "content should not exceed 1 MB")
	v.Check(
		validator.ValidateTimeInRange(n.CreatedAt, time.Now().Add(-1*time.Minute), time.Now()),
		"time of the creation should be in range (now - 1 min, now]",
	)
	v.Check(len(n.KeyHash) == 64, "key hash must be exactly 64 bytes long")

	isExpiresInSet := n.ExpiresIn != 0
	isExpiresAtSet := !n.ExpiresAt.IsZero() && n.ExpiresAtTimeZone != nil
	hasConflict := isExpiresInSet && isExpiresAtSet || !(isExpiresInSet || isExpiresAtSet)
	v.Check(hasConflict, "either expiration date and time zone or expiration timeout should be provided")
	if n.ExpiresIn != 0 {
		year := 24 * 60 * 365 * time.Minute
		v.Check(
			n.ExpiresIn >= 10*time.Minute && n.ExpiresIn <= 1*year,
			"expiration timeout should be in range [10 min, 365 days]",
		)
	} else {
		localCreatedAt := n.CreatedAt.In(n.ExpiresAtTimeZone)
		expiresAtLowerBound := time.Date(
			localCreatedAt.Year(), localCreatedAt.Month(), localCreatedAt.Day(),
			localCreatedAt.Hour(), localCreatedAt.Minute()+9, localCreatedAt.Second(), localCreatedAt.Nanosecond(),
			localCreatedAt.Location(),
		)
		expiresAtUpperBound := time.Date(
			localCreatedAt.Year()+1, localCreatedAt.Month(), localCreatedAt.Day(),
			localCreatedAt.Hour(), localCreatedAt.Minute(), localCreatedAt.Second(), localCreatedAt.Nanosecond(),
			localCreatedAt.Location(),
		)
		v.Check(
			validator.ValidateTimeInRange(n.ExpiresAt, expiresAtLowerBound, expiresAtUpperBound),
			"expiration date should be in range (local time + 9 min, local time + 1 year]",
		)
	}

	var validationErrors []error
	for _, err := range v.GetErrors() {
		validationErrors = append(validationErrors, err)
	}
	return errors.Join(validationErrors...)
}

func NewNote(
	id string,
	content *string,
	expiresIn time.Duration,
	expiresAt time.Time,
	expiresAtTimeZone string,
	keyHash []byte,
) (*Note, error) {
	tz, err := time.LoadLocation(expiresAtTimeZone)
	if err != nil {
		return &Note{}, err
	}

	note := &Note{
		ID:                id,
		Content:           content,
		CreatedAt:         time.Now(),
		ExpiresIn:         expiresIn,
		ExpiresAt:         expiresAt,
		ExpiresAtTimeZone: tz,
		KeyHash:           keyHash,
	}
	return note, nil
}

type NoteModel struct {
	ID                uint64
	Content           *string
	CreatedAt         time.Time
	ExpiresAt         time.Time
	ExpiresAtTimeZone string
	KeyHash           []byte
}
