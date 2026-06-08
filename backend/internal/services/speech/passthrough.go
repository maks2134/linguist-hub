package speech

import "context"

// PassthroughRecognizer accepts client-side transcripts (MVP).
type PassthroughRecognizer struct{}

func NewPassthroughRecognizer() *PassthroughRecognizer {
	return &PassthroughRecognizer{}
}

func (p *PassthroughRecognizer) Recognize(_ context.Context, in AudioBlob, locale string) (Transcript, error) {
	if in.Text != "" {
		return Transcript{Text: in.Text, Locale: locale, Confidence: 1}, nil
	}
	if len(in.Data) > 0 {
		return Transcript{}, ErrClientSide
	}
	return Transcript{}, ErrClientSide
}

// NoopSynthesizer defers TTS to the browser.
type NoopSynthesizer struct{}

func NewNoopSynthesizer() *NoopSynthesizer {
	return &NoopSynthesizer{}
}

func (n *NoopSynthesizer) Synthesize(_ context.Context, _ SynthesisRequest) (AudioBlob, error) {
	return AudioBlob{}, ErrClientSide
}
