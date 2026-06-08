package audio

import "testing"

func TestValidate_supported(t *testing.T) {
	if err := Validate("audio/webm;codecs=opus", 1024); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidate_tooLarge(t *testing.T) {
	if err := Validate("audio/webm", MaxTurnBytes+1); err == nil {
		t.Fatal("expected error for large payload")
	}
}

func TestValidate_unsupportedMIME(t *testing.T) {
	if err := Validate("video/mp4", 100); err == nil {
		t.Fatal("expected unsupported mime error")
	}
}

func TestEncodeDecodeBase64(t *testing.T) {
	raw := []byte{1, 2, 3}
	enc := EncodeBase64(raw)
	dec, err := DecodeBase64(enc)
	if err != nil || string(dec) != string(raw) {
		t.Fatalf("roundtrip failed: %v %v", dec, err)
	}
}
