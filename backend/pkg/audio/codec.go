package audio

import (
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	MaxChunkBytes = 256 * 1024
	MaxTurnBytes  = 1024 * 1024
)

var (
	ErrUnsupportedMIME = errors.New("audio: unsupported mime type")
	ErrTooLarge        = errors.New("audio: payload too large")
)

var supportedMIME = map[string]struct{}{
	"audio/webm":             {},
	"audio/webm;codecs=opus": {},
	"audio/ogg":              {},
	"audio/ogg;codecs=opus":  {},
	"audio/wav":              {},
	"audio/mpeg":             {},
	"audio/mp3":              {},
}

// Validate checks MIME and size for uploaded audio.
func Validate(mime string, size int) error {
	if size <= 0 {
		return fmt.Errorf("%w: empty payload", ErrTooLarge)
	}
	if size > MaxTurnBytes {
		return fmt.Errorf("%w: %d bytes", ErrTooLarge, size)
	}
	if !IsSupportedMIME(mime) {
		return fmt.Errorf("%w: %s", ErrUnsupportedMIME, mime)
	}
	return nil
}

// IsSupportedMIME reports whether the MIME type is allowed.
func IsSupportedMIME(mime string) bool {
	if mime == "" {
		return false
	}
	_, ok := supportedMIME[mime]
	return ok
}

// EncodeBase64 encodes raw audio bytes.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a base64 audio payload.
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
