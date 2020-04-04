package http

type Error struct {
	Cause      string `json:"error"`
	Stacktrace string `json:"stacktrace,omitempty"`
}

func WrapError(err error) Error {
	return Error{Cause: err.Error()}
}
