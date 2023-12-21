package obs

func ErrToStatus(err error) string {
	if err != nil {
		return "ERROR"
	}
	return "OK"
}
