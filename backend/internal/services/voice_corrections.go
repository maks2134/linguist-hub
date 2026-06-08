package services

import (
	"encoding/json"
	"strings"

	"github.com/makskozlov/linguist-hub/internal/dto"
)

type correctionsPayload struct {
	Corrections []dto.VoiceCorrection `json:"corrections"`
}

func splitAssistantAndCorrections(response string) (string, []dto.VoiceCorrection) {
	response = strings.TrimSpace(response)
	startIdx := strings.LastIndex(response, "{")
	if startIdx == -1 {
		return response, nil
	}
	jsonStr := response[startIdx:]
	var payload correctionsPayload
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil || len(payload.Corrections) == 0 {
		// try score payload for lessons
		return response, nil
	}
	text := strings.TrimSpace(response[:startIdx])
	return text, payload.Corrections
}

type lessonScorePayload struct {
	Score struct {
		Fluency       float32 `json:"fluency"`
		Grammar       float32 `json:"grammar"`
		Vocabulary    float32 `json:"vocabulary"`
		Pronunciation float32 `json:"pronunciation"`
		Notes         string  `json:"notes"`
	} `json:"score"`
}

func splitAssistantAndLessonScore(response string) (string, *lessonScorePayload) {
	response = strings.TrimSpace(response)
	startIdx := strings.LastIndex(response, "{")
	if startIdx == -1 {
		return response, nil
	}
	jsonStr := response[startIdx:]
	var payload lessonScorePayload
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil {
		return response, nil
	}
	if payload.Score.Fluency == 0 && payload.Score.Grammar == 0 && payload.Score.Vocabulary == 0 {
		return response, nil
	}
	text := strings.TrimSpace(response[:startIdx])
	return text, &payload
}
