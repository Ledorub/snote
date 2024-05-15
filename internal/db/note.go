package db

import (
	"context"
	"fmt"
	"github.com/ledorub/snote-api/internal"
	"log"
)

type NoteRepository struct {
	logger  *log.Logger
	queries *Queries
}

func NewNoteRepository(logger *log.Logger, queries *Queries) *NoteRepository {
	return &NoteRepository{logger: logger, queries: queries}
}

func (r *NoteRepository) Create(ctx context.Context, note *internal.NoteModel) (*internal.NoteModel, error) {
	createdNote, err := r.queries.CreateNote(ctx, CreateNoteParams{
		Content:           *note.Content,
		CreatedAt:         newTimestampTZ(note.CreatedAt),
		ExpiresAt:         newTimestamp(note.ExpiresAt),
		ExpiresAtTimezone: note.ExpiresAtTimeZone,
		KeyHash:           note.KeyHash,
	})
	if err != nil {
		return &internal.NoteModel{}, fmt.Errorf("creation failed: %w", err)
	}

	note.ID = pgIntToUInt64(createdNote.ID)
	note.Content = &createdNote.Content
	note.CreatedAt = createdNote.CreatedAt.Time
	note.ExpiresAt = createdNote.ExpiresAt.Time
	note.ExpiresAtTimeZone = createdNote.ExpiresAtTimezone
	note.KeyHash = createdNote.KeyHash
	return note, nil
}

func (r *NoteRepository) Get(ctx context.Context, id uint64) (*internal.NoteModel, error) {
	pgId, err := uInt64ToPgInt8(id)
	if err != nil {
		return nil, err
	}
	note, err := r.queries.GetNote(ctx, pgId)
	if err != nil {
		return nil, fmt.Errorf("retrieving failed: %w", err)
	}
	return &internal.NoteModel{
		ID:                pgIntToUInt64(note.ID),
		Content:           &note.Content,
		CreatedAt:         note.CreatedAt.Time,
		ExpiresAt:         note.ExpiresAt.Time,
		ExpiresAtTimeZone: note.ExpiresAtTimezone,
		KeyHash:           note.KeyHash,
	}, nil
}

func (r *NoteRepository) Delete(ctx context.Context, id uint64) error {
	pgInt, err := uInt64ToPgInt8(id)
	if err != nil {
		return err
	}
	if err := r.queries.DeleteNote(ctx, pgInt); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}
	return nil
}
