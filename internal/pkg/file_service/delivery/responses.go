package delivery

type ResponseURLBody struct {
	SlURL []string `json:"urls"`
}

type ResponseURLs struct {
	Status int             `json:"status"`
	Body   ResponseURLBody `json:"body"`
}

func NewResponseURLs(status int, slURL []string) *ResponseURLs {
	return &ResponseURLs{Status: status, Body: ResponseURLBody{SlURL: slURL}}
}
