package contacts

type Contact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Contacts struct {
	Contacts []Contact `json:"contacts"`
}