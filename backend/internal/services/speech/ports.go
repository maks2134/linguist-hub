package speech

import (
	"context"
	"errors"
)

// ErrClientSide indicates STT/TTS is handled in the browser.
var ErrClientSide = errors.New("speech: handled on client")

type AudioBlob struct {
	Data     []byte
	MIMEType string
	Text     string // set when client already transcribed (MVP)
}

type Transcript struct {
	Text       string
	Locale     string
	Confidence float32
}

type SynthesisRequest struct {
	Text   string
	Locale string
	Voice  string
}

type Recognizer interface {
	Recognize(ctx context.Context, in AudioBlob, locale string) (Transcript, error)
}

type Synthesizer interface {
	Synthesize(ctx context.Context, req SynthesisRequest) (AudioBlob, error)
}
