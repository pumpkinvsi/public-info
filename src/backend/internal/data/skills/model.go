package skills

type Level struct {
	ID    int    `json:"id"`
	Level int    `json:"level"`
	Text  string `json:"text"`
}

type Skill struct {
	Name  string `json:"name"`
	Level Level  `json:"level"`
}
