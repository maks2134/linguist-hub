package services

import (
	"strings"

	"github.com/makskozlov/linguist-hub/internal/models"
)

func normalizeTopic(topic string) string {
	return strings.TrimSpace(topic)
}

// MarkTopicCompleted records a topic as completed in the matching module.
// Returns true if the topic was newly marked.
func MarkTopicCompleted(modules *models.Modules, moduleID int, topic string) bool {
	topic = normalizeTopic(topic)
	if topic == "" {
		return false
	}
	for i := range *modules {
		if (*modules)[i].ID != moduleID {
			continue
		}
		for _, ct := range (*modules)[i].CompletedTopics {
			if normalizeTopic(ct) == topic {
				return false
			}
		}
		(*modules)[i].CompletedTopics = append((*modules)[i].CompletedTopics, topic)
		return true
	}
	return false
}

func IsModuleComplete(module models.Module) bool {
	if len(module.Topics) == 0 {
		return true
	}
	completed := make(map[string]bool, len(module.CompletedTopics))
	for _, ct := range module.CompletedTopics {
		completed[normalizeTopic(ct)] = true
	}
	for _, t := range module.Topics {
		if !completed[normalizeTopic(t)] {
			return false
		}
	}
	return true
}

// ModuleLocked returns true when module index > 0 and the previous module is incomplete.
func ModuleLocked(modules models.Modules, index int) bool {
	if index <= 0 {
		return false
	}
	return !IsModuleComplete(modules[index-1])
}

func IsModuleLocked(modules models.Modules, moduleID int) bool {
	for i, m := range modules {
		if m.ID == moduleID {
			return ModuleLocked(modules, i)
		}
	}
	return true
}
