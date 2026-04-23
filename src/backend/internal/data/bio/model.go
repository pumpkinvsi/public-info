package bio
import "src/backend/internal/shared/model"

type Bio struct {
	Text model.LocalizedString `json:"text"`
}