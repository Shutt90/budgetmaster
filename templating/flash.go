package template

type flash struct {
	Message  string
	Negative bool
}

func NewFlash(message string, negative bool) flash {
	return flash{
		Message:  message,
		Negative: negative,
	}
}
