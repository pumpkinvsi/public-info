package model

type LocalizedString struct {
	Rus string `json:"rus"`
	Eng string `json:"eng"`
}

type Level struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Contact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Bio struct {
	Text LocalizedString `json:"text"`
}

type Skill struct {
	Name  string `json:"name"`
	Level Level  `json:"level"`
}

type Project struct {
	ID          int             `json:"id"`
	Name        LocalizedString `json:"name"`
	Description LocalizedString `json:"description"`
	Skills      []Skill         `json:"skills"`
	Note        LocalizedString `json:"note"`
}

type Technology struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectGroup struct {
	Technology Technology `json:"technology"`
	Projects   []Project  `json:"projects"`
}

type Email struct {
	Text    string `json:"text"`
	Sender  string `json:"sender"`
	Contact string `json:"contact"`
}

type Contacts struct {
	Contacts []Contact `json:"contacts"`
}