package email

type Email struct {
	Text    string `json:"text"`
	Sender  string `json:"sender"`
	Contact string `json:"contact"`
}