package service

import (
	"context"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ledorub/snote-api/internal"
	"github.com/ledorub/snote-api/internal/validator"
	"github.com/mr-tron/base58"
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrDoesNotExist = errors.New("does not exist")

type noteRepository interface {
	Create(ctx context.Context, note *internal.NoteModel) (*internal.NoteModel, error)
	Get(ctx context.Context, id uint64) (*internal.NoteModel, error)
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

func New(logger *log.Logger, repo noteRepository, idEncoderDecoder idEncDec) *NoteService {
	return &NoteService{logger: logger, repo: repo, idEncDec: idEncoderDecoder}
}

func (s *NoteService) CreateNote(ctx context.Context, note *internal.Note) (*internal.Note, error) {
	if err := note.CheckErrors(); err != nil {
		return &internal.Note{}, fmt.Errorf("note creation failed: %w", err)
	}

	expiresAt, tz := calcExpirationDate(note.ExpiresAt, note.ExpiresAtTimeZone, note.ExpiresIn)

	decodedKeyHash, err := base58.Decode(note.KeyHash)
	if err != nil {
		return &internal.Note{}, fmt.Errorf("note creation failed: %w", err)
	}

	newNote := &internal.NoteModel{
		Content:           note.Content,
		CreatedAt:         note.CreatedAt,
		ExpiresAt:         expiresAt,
		ExpiresAtTimeZone: tz.String(),
		KeyHash:           decodedKeyHash,
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
	return note, nil
}

func (s *NoteService) GetNote(ctx context.Context, id string, keyHash string) (*internal.Note, error) {
	v := validator.New()
	v.Check(len(id) == 12, "id should consist of 12 letters and/or digits and hyphens.")
	v.Check(validator.ValidateHyphenatedB58String(id), "id should consist of latin letters and/or digits and hyphens")
	v.Check(len(keyHash) == 44, "key hash should consist of 44 letters and/or digits")
	v.Check(validator.ValidateHyphenatedB58String(keyHash), "key hash should consist of latin letters and/or digits")
	if !v.CheckIsValid() {
		var validationErrors []error
		for _, err := range v.GetErrors() {
			validationErrors = append(validationErrors, err)
		}
		err := errors.Join(validationErrors...)
		return &internal.Note{}, err
	}

	gotError := false
	decodedID, err := s.idEncDec.Decode(id)
	if err != nil {
		gotError = true
	}

	decodedKeyHash, err := base58.Decode(keyHash)
	if err != nil {
		gotError = true
	}

	noteDB, err := s.repo.Get(ctx, decodedID)
	if err != nil {
		gotError = true
	}

	var noteKeyHash []byte
	var noteTimeZone string
	if noteDB != nil {
		noteKeyHash = noteDB.KeyHash
		noteTimeZone = noteDB.ExpiresAtTimeZone
	} else {
		noteKeyHash = decodedKeyHash
	}
	isAuthorized := compareKeyHashes(decodedKeyHash, noteKeyHash)

	tz, err := stringToTimeZone(noteTimeZone)
	if err != nil {
		return &internal.Note{}, errors.New("note has invalid time zone")
	}

	err = s.repo.Delete(ctx, decodedID)
	if err != nil {
		gotError = true
	}

	if !isAuthorized || gotError {
		return &internal.Note{}, ErrDoesNotExist
	}
	return &internal.Note{
		ID:                id,
		Content:           noteDB.Content,
		CreatedAt:         noteDB.CreatedAt,
		ExpiresAt:         noteDB.ExpiresAt,
		ExpiresAtTimeZone: tz,
		KeyHash:           keyHash,
	}, nil
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
	enc := base58.Encode(bin)
	enc = ed.padEncodedID(enc)
	return insertStringSeparator(enc, "-", 4)
}

func (ed *B58IDEncDec) padEncodedID(id string) string {
	if width := 10 - len(id); width > 0 {
		return padStringWith(id, "1", 10)
	}
	return id
}

func (ed *B58IDEncDec) Decode(str string) (uint64, error) {
	str, err := removeStringSeparator(str, "-", 4)
	if err != nil {
		return 0, err
	}
	dec, err := base58.Decode(str)
	if err != nil {
		return 0, fmt.Errorf("id decoding error: %w", err)
	}
	dec = dec[len(dec)-8:]
	num := binary.BigEndian.Uint64(dec)
	return num, nil
}

func compareKeyHashes(x, y []byte) bool {
	return subtle.ConstantTimeCompare(x, y) == 1
}

func padStringWith(s, padding string, totalWidth int) string {
	stringWidth := utf8.RuneCountInString(s)
	paddingWidth := utf8.RuneCountInString(padding)
	fillWidth := totalWidth - stringWidth
	if fillWidth > 0 {
		fill := strings.Repeat(padding, fillWidth/paddingWidth) + padding[:fillWidth%paddingWidth]
		s = fill + s
	}
	return s
}

func insertStringSeparator(s, sep string, segmentWidth int) string {
	approxSepCount := len(s) / segmentWidth
	if len(s)%segmentWidth == 0 {
		approxSepCount--
	}
	newWidth := len(s) + approxSepCount*len(sep)
	strBytes := make([]byte, newWidth)

	writeI := 0
	for byteI, charI := 0, 0; byteI < len(s)-1; byteI++ {
		strBytes[writeI] = s[byteI]
		writeI++

		if utf8.RuneStart(s[byteI+1]) {
			charI++
		}
		if charI == segmentWidth {
			for i := 0; i < len(sep); i++ {
				strBytes[writeI] = sep[i]
				writeI++
			}
			charI = 0
		}
	}
	strBytes[writeI] = s[len(s)-1]
	return string(strBytes[:writeI+1])
}

func removeStringSeparator(s, sep string, segmentWidth int) (string, error) {
	strBytes := make([]byte, len(s))

	writeI := 0
	for byteI, charI := 0, 0; byteI < len(s); byteI++ {
		if byteI == len(s)-1 || utf8.RuneStart(s[byteI+1]) {
			charI++
		}

		strBytes[writeI] = s[byteI]
		writeI++

		if charI == segmentWidth && byteI != len(s)-1 {
			charI = 0
			substr := s[byteI+1 : byteI+1+len(sep)]
			if substr != sep {
				return "", fmt.Errorf(
					"removeStringSeparator: expected %s at s[%d], but got %s", sep, byteI+1, substr,
				)
			}
			byteI += len(sep)
		}
	}
	return string(strBytes[:writeI]), nil
}
