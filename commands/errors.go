package points

type argError struct {
	message string
	err     error
}

// Error implements the error interface
func (e argError) Error() string {
	return e.message
}
