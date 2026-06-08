package dto

type VoiceCorrection struct {
	Original  string `json:"original"`
	Suggested string `json:"suggested"`
	Reason    string `json:"reason"`
}

type VoiceTutorRequest struct {
	Message string `json:"message"`
	Locale  string `json:"locale,omitempty"`
	Topic   string `json:"topic,omitempty"`
}

type VoiceTutorResponse struct {
	UserText    string            `json:"user_text"`
	Assistant   string            `json:"assistant"`
	Corrections []VoiceCorrection `json:"corrections,omitempty"`
	Locale      string            `json:"locale"`
}

type VoiceLocalesResponse struct {
	Primary  string   `json:"primary"`
	Fallback []string `json:"fallback,omitempty"`
	Voices   []string `json:"voices,omitempty"`
}

type LessonStartRequest struct {
	ModuleID int    `json:"module_id"`
	Topic    string `json:"topic"`
	Locale   string `json:"locale,omitempty"`
}

type LessonTurnRequest struct {
	LessonID uint   `json:"lesson_id"`
	Message  string `json:"message"`
	Locale   string `json:"locale,omitempty"`
}

type LessonScoreDTO struct {
	Fluency       float32 `json:"fluency"`
	Grammar       float32 `json:"grammar"`
	Vocabulary    float32 `json:"vocabulary"`
	Pronunciation float32 `json:"pronunciation"`
	Notes         string  `json:"notes"`
}

type LessonResponse struct {
	LessonID      uint            `json:"lesson_id"`
	Stage         string          `json:"stage"`
	StageIndex    int             `json:"stage_index"`
	StageTotal    int             `json:"stage_total"`
	Prompt        string          `json:"prompt"`
	TargetPhrases []string        `json:"target_phrases"`
	HitPhrases    []string        `json:"hit_phrases"`
	Score         *LessonScoreDTO `json:"score,omitempty"`
	Locale        string          `json:"locale,omitempty"`
}

type STTResponse struct {
	Text       string  `json:"text"`
	Confidence float32 `json:"confidence"`
}
