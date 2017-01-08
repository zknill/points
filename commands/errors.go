package points

type argError struct {
	message string
	err     error
}

func (e argError) Error() string {
	return e.message
}
