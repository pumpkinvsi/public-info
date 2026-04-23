package email

import( 
	"github.com/microcosm-cc/bluemonday"
)

type sanitizer struct {
	strictPolicy *bluemonday.Policy
}

func New() *sanitizer {
	return &sanitizer{
		strictPolicy: bluemonday.StrictPolicy(),
	}
}

func (s *sanitizer) Sanitize(body string) string { 
	return s.strictPolicy.Sanitize(body)
}