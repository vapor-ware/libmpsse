package libmpsse


// MpsseError is the error that is returned when an MPSSE failure is
// detected. The error message it provides is the message retrieved from
// the ErrorString function.
type MpsseError struct {
	Message string
}

func (e *MpsseError) Error() string {
	return e.Message
}