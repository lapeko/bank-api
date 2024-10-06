package api

type response struct {
	Ok   bool        `json:"ok"`
	Err  interface{} `json:"err"`
	Body interface{} `json:"body"`
}

func genFailBody(input interface{}) response {
	err := input
	if realErr, ok := err.(error); ok {
		err = realErr.Error()
	}
	return response{
		Ok:   false,
		Err:  err,
		Body: nil,
	}
}
func genOkBody(body interface{}) response {
	return response{
		Ok:   true,
		Err:  nil,
		Body: body,
	}
}
