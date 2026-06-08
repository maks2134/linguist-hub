package speech

import "context"

// GeminiTranscriber transcribes audio via Gemini (used when FEATURE_SERVER_STT=true).
type GeminiTranscriber interface {
	TranscribeBlob(ctx context.Context, data []byte, mimeType, locale string) (string, error)
}

type GeminiRecognizer struct {
	transcriber GeminiTranscriber
}

func NewGeminiRecognizer(t GeminiTranscriber) *GeminiRecognizer {
	return &GeminiRecognizer{transcriber: t}
}

func (g *GeminiRecognizer) Recognize(ctx context.Context, in AudioBlob, locale string) (Transcript, error) {
	if in.Text != "" {
		return Transcript{Text: in.Text, Locale: locale, Confidence: 1}, nil
	}
	if len(in.Data) == 0 {
		return Transcript{}, ErrClientSide
	}
	text, err := g.transcriber.TranscribeBlob(ctx, in.Data, in.MIMEType, locale)
	if err != nil {
		return Transcript{}, err
	}
	return Transcript{Text: text, Locale: locale, Confidence: 0.9}, nil
}
