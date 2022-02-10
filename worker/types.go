package worker

type Job struct {
	Url   string
	Hash  []byte
	Error error
}

type HttpError string

func (h HttpError) Error() string {
	return string(h)
}



