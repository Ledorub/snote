package db

import (
	"context"
	"fmt"
	"github.com/ledorub/snote-api/internal"
	"log"
	"time"
)

type TaskRepository struct {
	logger  *log.Logger
	queries *Queries
}

func NewTaskRepository(logger *log.Logger, queries *Queries) *TaskRepository {
	return &TaskRepository{logger: logger, queries: queries}
}

func (r *TaskRepository) Create(
	ctx context.Context,
	content string,
	createdAt time.Time,
	expiresAt time.Time,
	expiresAtTimeZone *time.Location,
	keyHash []byte,
) (internal.Note, error) {
	note, err := r.queries.CreateNote(ctx, CreateNoteParams{
		Content:           content,
		CreatedAt:         newTimestampTZ(createdAt),
		ExpiresAt:         newTimestamp(expiresAt),
		ExpiresAtTimezone: expiresAtTimeZone.String(),
		KeyHash:           keyHash,
	})
	if err != nil {
		return internal.Note{}, fmt.Errorf("insertion failed: %w", err)
	}
	tz, err := time.LoadLocation(note.ExpiresAtTimezone)
	if err != nil {
		return internal.Note{}, fmt.Errorf("received invalid TZ from the DB: %w", err)
	}

	return internal.Note{
		ID:                pgIntToUInt64(note.ID),
		Content:           note.Content,
		CreatedAt:         note.CreatedAt.Time,
		ExpiresAt:         note.ExpiresAt.Time,
		ExpiresAtTimeZone: tz,
		KeyHash:           keyHash,
	}, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uint64) error {
	pgInt, err := uInt64ToPgInt8(id)
	if err != nil {
		return err
	}
	if err := r.queries.DeleteNote(ctx, pgInt); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}
	return nil
}
