package service

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/ledorub/snote-api/internal"
	"github.com/mr-tron/base58"
	"log"
	"time"
)

type noteRepository interface {
	Create(ctx context.Context, note *internal.NoteModel) (*internal.NoteModel, error)
	Delete(ctx context.Context, id uint64) error
}

type idEncDec interface {
	Encode(intID uint64) string
	Decode(strID string) (uint64, error)
}

type NoteService struct {
	logger   *log.Logger
	repo     noteRepository
	idEncDec idEncDec
}

func New(logger *log.Logger, repo noteRepository) *NoteService {
	return &NoteService{logger: logger, repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, note *internal.Note) (*internal.Note, error) {
	if err := note.CheckErrors(); err != nil {
		return &internal.Note{}, fmt.Errorf("note creation failed: %w", err)
	}

	expiresAt, tz := calcExpirationDate(note.ExpiresAt, note.ExpiresAtTimeZone, note.ExpiresIn)
	newNote := &internal.NoteModel{
		Content:           note.Content,
		CreatedAt:         note.CreatedAt,
		ExpiresAt:         expiresAt,
		ExpiresAtTimeZone: tz.String(),
		KeyHash:           note.KeyHash,
	}
	createdNote, err := s.repo.Create(ctx, newNote)
	if err != nil {
		return &internal.Note{}, fmt.Errorf("note creation failed: %w", err)
	}

	tz, err = stringToTimeZone(createdNote.ExpiresAtTimeZone)
	if err != nil {
		return &internal.Note{}, fmt.Errorf("note creation failed: %w", err)
	}

	encodedID := s.idEncDec.Encode(createdNote.ID)
	note.ID = encodedID
	note.Content = createdNote.Content
	note.CreatedAt = createdNote.CreatedAt
	note.ExpiresIn = 0
	note.ExpiresAt = createdNote.ExpiresAt
	note.ExpiresAtTimeZone = tz
	note.KeyHash = createdNote.KeyHash
	return note, nil
}

func calcExpirationDate(expiresAt time.Time, tz *time.Location, expiresIn time.Duration) (time.Time, *time.Location) {
	if expiresIn != 0 {
		exp := time.Now().UTC().Add(expiresIn)
		return exp, time.UTC
	}
	return expiresAt.UTC(), tz
}

func stringToTimeZone(tzID string) (*time.Location, error) {
	tz, err := time.LoadLocation(tzID)
	if err != nil {
		return tz, fmt.Errorf("unable to convert tz ID to time zone: %w", err)
	}
	return tz, nil
}

type B58IDEncDec struct{}

func (ed *B58IDEncDec) Encode(id uint64) string {
	bin := make([]byte, 8)
	binary.BigEndian.PutUint64(bin, id)
	return base58.Encode(bin)
}

func (ed *B58IDEncDec) Decode(str string) (uint64, error) {
	decoded, err := base58.Decode(str)
	if err != nil {
		return 0, fmt.Errorf("id decoding error: %w", err)
	}
	num := binary.BigEndian.Uint64(decoded)
	return num, nil
}
