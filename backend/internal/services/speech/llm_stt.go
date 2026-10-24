package speech

import "context"

// LLMTranscriber is the subset of LLMService used for optional server STT.
type LLMTranscriber interface {
	TranscribeBlob(ctx context.Context, data []byte, mimeType, locale string) (string, error)
}

type LLMRecognizer struct {
	transcriber LLMTranscriber
}

func NewLLMRecognizer(t LLMTranscriber) *LLMRecognizer {
	return &LLMRecognizer{transcriber: t}
}

func (g *LLMRecognizer) Recognize(ctx context.Context, in AudioBlob, locale string) (Transcript, error) {
	text, err := g.transcriber.TranscribeBlob(ctx, in.Data, in.MIMEType, locale)
	if err != nil {
		return Transcript{}, err
	}
	return Transcript{Text: text}, nil
}
