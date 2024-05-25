package menuify

const (
	ERR_CANCELLED = Error("calibrator: cancelled")
)

type Error string
func (e Error) Error() string {
	return string(e)
}